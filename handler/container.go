package handler

import (
	"container-manager/types"
	"github.com/sirupsen/logrus"
	"net/http"
)

// ContainerRequest is the request object for the ContainerService.Create method.
type ContainerRequest struct {
	types.Container
}

// ContainerResponse is the response object for the ContainerService.Create method.
// job_id: The ID of the job that was created
type ContainerResponse struct {
	JobID string `json:"job_id"`
}

// ContainerService is the service that handles container creation.
type ContainerService struct{}

// Create creates a new container.
func (cs *ContainerService) Create(r *http.Request, req *ContainerRequest, res *ContainerResponse) error {
	logrus.WithFields(logrus.Fields{
		"image":     req.Image,
		"arguments": req.Arguments,
		"resources": req.Resources,
		"env":       req.Env,
	}).Debugf("Running job")
	// TODO: Implement job creation
	return nil
}
