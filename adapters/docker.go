package adapters

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"

	"github.com/IMEsec-USP/container-manager/applications"
)

// DockerAdapter is the adapter that communicates with dockerd.
type DockerAdapter struct {
	client *client.Client
}

type ContainerConfig struct {
	id               string
	config           *container.Config
	hostConfig       *container.HostConfig
	networkingConfig *network.NetworkingConfig
}

// NewDockerAdapter returns a DockerAdapter based on environment vars.
func NewDockerAdapter() (*DockerAdapter, error) {
	dockerClient, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &DockerAdapter{client: dockerClient}, nil
}

// PullImage pulls a docker image of the provided app.
func (d *DockerAdapter) PullImage(ctx context.Context, app applications.Application) (io.Reader, error) {
	return d.client.ImagePull(ctx, app.ImageReference(), types.ImagePullOptions{})
}

func (d *DockerAdapter) GetContainerConfig(ctx context.Context, app applications.Application) (*ContainerConfig, error) {
	containerJSON, err := d.client.ContainerInspect(ctx, app.ContainerName())
	if err != nil {
		return nil, err
	}

	return &ContainerConfig{
		id:               containerJSON.ID,
		config:           containerJSON.Config,
		hostConfig:       containerJSON.HostConfig,
		networkingConfig: &network.NetworkingConfig{EndpointsConfig: containerJSON.NetworkSettings.Networks},
	}, nil
}

func (d *DockerAdapter) RemoveContainer(ctx context.Context, config *ContainerConfig) error {
	timeout := 10 * time.Second
	err := d.client.ContainerStop(ctx, config.id, &timeout)
	if err != nil {
		return err
	}
	return d.client.ContainerRemove(ctx, config.id, types.ContainerRemoveOptions{Force: true})
}

func (d *DockerAdapter) RunImage(ctx context.Context, app applications.Application, config *ContainerConfig) error {
	created, err := d.client.ContainerCreate(ctx, config.config, config.hostConfig, nil, app.ContainerName())
	if err != nil {
		return err
	}
	return d.client.ContainerStart(ctx, created.ID, types.ContainerStartOptions{})
}
