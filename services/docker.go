package services

import (
	"container-manager/types"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"time"

	"github.com/docker/docker/api/types/image"

	dockerContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// DockerService service interface to deploy and get container status
type DockerService interface {
	DeployContainer(container types.Container) (string, error)
	GetContainerStatus(containerID string) (string, error)
}

// DockerServiceHandler is the implementation of the DockerService interface
// client: The Docker client
type DockerServiceHandler struct {
	client *client.Client
}

// NewDockerService creates a new DockerServiceHandler instance
func NewDockerService() (*DockerServiceHandler, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	return &DockerServiceHandler{
		client: cli,
	}, nil
}

// DeployContainer deploys a container using Docker
func (ds *DockerServiceHandler) DeployContainer(container types.Container) (string, error) {
	logrus.WithField("container", container).Debug("Deploying container")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var envVars []string
	for key, value := range container.Env {
		envVars = append(envVars, key+"="+value)
	}

	reader, err := ds.client.ImagePull(ctx, container.Image, image.PullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}
	defer reader.Close()

	io.Copy(io.Discard, reader)

	resp, err := ds.client.ContainerCreate(ctx, &dockerContainer.Config{
		Image: container.Image,
		Cmd:   container.Arguments,
		Env:   envVars,
	}, nil, nil, nil, "")
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	if err := ds.client.ContainerStart(ctx, resp.ID, dockerContainer.StartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return resp.ID, nil
}

// GetContainerStatus gets the status of a container by container ID
func (ds *DockerServiceHandler) GetContainerStatus(containerID string) (string, error) {
	logrus.WithField("container_id", containerID).Debug("Getting container status")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	containerJSON, err := ds.client.ContainerInspect(ctx, containerID)
	if err != nil {
		return "", fmt.Errorf("failed to inspect container: %w", err)
	}

	return containerJSON.State.Status, nil
}
