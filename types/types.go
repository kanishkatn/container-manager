package types

import "fmt"

// Container is the object that represents a container to be run.
// image: The container image to run
// arguments: The arguments to pass to the container
// resources: The resources to allocate for the job
// env: The environment variables to set for the job
type Container struct {
	Image     string            `json:"image"`
	Arguments []string          `json:"arguments"`
	Env       map[string]string `json:"env"`
}

func (c Container) Validate() error {
	if c.Image == "" {
		return fmt.Errorf("image is required")
	}
	return nil
}

type JobStatus string

const (
	JobStatusPending  JobStatus = "pending"
	JobStatusComplete JobStatus = "complete"
	JobStatusFailed   JobStatus = "failed"
)

func (js JobStatus) String() string {
	return string(js)
}

type P2PMessageType string

const (
	P2PMessageTypeDeployContainer P2PMessageType = "deploy_container"
)

func (pm P2PMessageType) String() string {
	return string(pm)
}
