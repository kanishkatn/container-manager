package types

// Container is the object that represents a container to be run.
// target_machine: The machine to run the job on
// container_image: The container image to run
// arguments: The arguments to pass to the container
// resources: The resources to allocate for the job
// env: The environment variables to set for the job
type Container struct {
	TargetMachine  string            `json:"target_machine"`
	ContainerImage string            `json:"container_image"`
	Arguments      []string          `json:"arguments"`
	Resources      map[string]string `json:"resources"`
	Env            map[string]string `json:"env"`
}

type JobStatus string

const (
	JobStatusPending  JobStatus = "pending"
	JobStatusComplete JobStatus = "complete"
	JobStatusFailed   JobStatus = "failed"
)
