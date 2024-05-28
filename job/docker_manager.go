package job

import (
	"container-manager/types"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/docker/docker/api/types/image"

	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerManager interface to deploy and get container status
type DockerManager interface {
	DeployContainer(container types.Container) (string, error)
	GetContainerStatus(containerID string) (string, error)
}

// dockerManager is the implementation of the DockerManager interface
// client: The Docker client
type dockerManager struct {
	client *client.Client
}

// newDockerManager creates a new DockerManager instance
func newDockerManager() (*dockerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	return &dockerManager{
		client: cli,
	}, nil
}

// DeployContainer deploys a container using Docker
func (d *dockerManager) DeployContainer(container types.Container) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var envVars []string
	for key, value := range container.Env {
		envVars = append(envVars, key+"="+value)
	}

	reader, err := d.client.ImagePull(ctx, container.Image, image.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	io.Copy(io.Discard, reader)

	resp, err := d.client.ContainerCreate(ctx, &dockerContainer.Config{
		Image: container.Image,
		Cmd:   container.Arguments,
		Env:   envVars,
	}, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	if err := d.client.ContainerStart(ctx, resp.ID, dockerContainer.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return resp.ID, nil
}

// GetContainerStatus gets the status of a container by container ID
func (d *dockerManager) GetContainerStatus(containerID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	containerJSON, err := d.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	return containerJSON.State.Status, nil
}
