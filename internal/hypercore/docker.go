package hypercore

import (
	"context"
	"fmt"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

type DockerClient struct {
	*client.Client
}

func NewDockerClient() (DockerClient, error) {
	// API version negotiation is enabled by default in the moby client.
	client, err := client.New()
	if err != nil {
		return DockerClient{}, err
	}

	return DockerClient{client}, nil
}

func (c DockerClient) Start(ctx context.Context, imageRef string) (string, error) {
	pullResp, err := c.ImagePull(ctx, imageRef, client.ImagePullOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to pull image: %w", err)
	}

	pullResp.Close()

	containerResp, err := c.ContainerCreate(ctx, client.ContainerCreateOptions{
		Config: &container.Config{
			Image: imageRef,
		},
		HostConfig: &container.HostConfig{
			AutoRemove: true,
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create container: %w", err)
	}

	if _, err = c.ContainerStart(ctx, containerResp.ID, client.ContainerStartOptions{}); err != nil {
		return "", fmt.Errorf("failed to start container: %w", err)
	}

	return "", nil
}

func (c DockerClient) Stop(ctx context.Context, containerID string) error {
	timeout := 15
	_, err := c.ContainerStop(ctx, containerID, client.ContainerStopOptions{Timeout: &timeout})

	if err != nil {
		return fmt.Errorf("failed to remove container %s: %w", containerID, err)
	}

	return nil
}
