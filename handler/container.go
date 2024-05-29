package handler

import (
	"container-manager/services"
	"container-manager/types"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// ContainerCreateRequest is the request object for the ContainerService.Create method.
type ContainerCreateRequest struct {
	types.Container
}

// ContainerCreateResponse is the response object for the ContainerService.Create method.
// JobID: The ID of the job that was created
// Message: A response message
type ContainerCreateResponse struct {
	JobID   string `json:"job_id"`
	Message string `json:"message"`
}

// ContainerStatusRequest is the request object for the ContainerService.Status method.
type ContainerStatusRequest struct {
	JobID string `json:"job_id"`
}

// ContainerStatusResponse is the response object for the ContainerService.Status method.
// JobID: The ID of the job
// Message: The status of the job
type ContainerStatusResponse struct {
	JobID  string `json:"job_id"`
	Status string `json:"status"`
}

// ContainerService is the service that handles container creation.
type ContainerService struct {
	jobQueue   services.Queue
	p2pService services.P2PService
}

// NewContainerService creates a new container service.
func NewContainerService(jobQueue services.Queue, p2pService services.P2PService) *ContainerService {
	return &ContainerService{
		jobQueue:   jobQueue,
		p2pService: p2pService,
	}
}

// Create creates a new container.
func (cs *ContainerService) Create(r *http.Request, req *ContainerCreateRequest, res *ContainerCreateResponse) error {
	if req == nil {
		return fmt.Errorf("invalid request")
	}

	logrus.WithFields(logrus.Fields{
		"image":     req.Image,
		"arguments": req.Arguments,
		"env":       req.Env,
	}).Debugf("queueing and broadcasting job")
	if err := req.Container.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	jobID, err := uuid.NewUUID()
	if err != nil {
		return fmt.Errorf("failed to generate job ID: %w", err)
	}

	// Enqueue the job
	if err := cs.jobQueue.Enqueue(jobID.String(), req.Container); err != nil {
		return fmt.Errorf("failed to enqueue job: %w", err)
	}

	// forward the job to the p2p network
	containerData, err := json.Marshal(req.Container)
	if err != nil {
		return fmt.Errorf("failed to marshal container data: %w", err)
	}
	msg := services.Message{
		Type:  types.P2PMessageTypeDeployContainer,
		JobID: jobID.String(),
		Data:  containerData,
	}
	if err := cs.p2pService.Broadcast(msg); err != nil {
		return fmt.Errorf("failed to send job to p2p network: %w", err)
	}

	res.JobID = jobID.String()
	res.Message = "Job created successfully"

	return nil
}

// Status returns the status of a job.
func (cs *ContainerService) Status(r *http.Request, req *ContainerStatusRequest, res *ContainerStatusResponse) error {
	if req == nil {
		return fmt.Errorf("invalid request")
	}

	logrus.WithField("job_id", req.JobID).Debug("getting job status")

	status, ok := cs.jobQueue.GetStatus(req.JobID)
	if !ok {
		return fmt.Errorf("job not found")
	}

	res.JobID = req.JobID
	res.Status = status.String()

	return nil
}
