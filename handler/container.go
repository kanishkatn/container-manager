package handler

import (
	"container-manager/job"
	"container-manager/p2p"
	"container-manager/types"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
)

// ContainerRequest is the request object for the ContainerService.Create method.
type ContainerRequest struct {
	types.Container
}

// ContainerResponse is the response object for the ContainerService.Create method.
// JobID: The ID of the job that was created
// Message: A response message
type ContainerResponse struct {
	JobID   string `json:"job_id"`
	Message string `json:"message"`
}

// ContainerService is the service that handles container creation.
type ContainerService struct {
	jobQueue   job.Queue
	p2pService p2p.P2PService
}

// NewContainerService creates a new container service.
func NewContainerService(jobQueue job.Queue, p2pService p2p.P2PService) *ContainerService {
	return &ContainerService{
		jobQueue:   jobQueue,
		p2pService: p2pService,
	}
}

// Create creates a new container.
func (cs *ContainerService) Create(r *http.Request, req *ContainerRequest, res *ContainerResponse) error {
	logrus.WithFields(logrus.Fields{
		"image":     req.Image,
		"arguments": req.Arguments,
		"resources": req.Resources,
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
	msg := p2p.Message{
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
func (cs *ContainerService) Status(r *http.Request, jobID string, res *ContainerResponse) error {
	status, ok := cs.jobQueue.GetStatus(jobID)
	if !ok {
		return fmt.Errorf("job not found")
	}

	res.JobID = jobID
	res.Message = status.String()

	return nil
}
