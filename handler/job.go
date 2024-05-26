package handler

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

type JobRequest struct {
	TargetMachine  string            `json:"target_machine"`
	ContainerImage string            `json:"container_image"`
	Arguments      []string          `json:"arguments"`
	Resources      map[string]string `json:"resources"`
	Env            map[string]string `json:"env"`
}

type JobResponse struct {
	JobID   string `json:"job_id"`
	Message string `json:"message"`
}

type JobService struct{}

func (js *JobService) CreateJob(r *http.Request, req *JobRequest, res *JobResponse) error {
	logrus.WithFields(logrus.Fields{
		"target_machine":  req.TargetMachine,
		"container_image": req.ContainerImage,
		"arguments":       req.Arguments,
		"resources":       req.Resources,
		"env":             req.Env,
	}).Debugf("Running job")
	// TODO: Implement job creation
	return nil
}
