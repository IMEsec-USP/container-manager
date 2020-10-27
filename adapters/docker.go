package adapters

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/IMEsec-USP/container-manager/applications"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/rs/zerolog"
)

// DockerAdapter is the adapter that communicates with dockerd.
type DockerAdapter struct {
	client *client.Client
}

type ContainerConfig struct {
	containerName    string
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

func (d *DockerAdapter) GetContainerConfig(ctx context.Context, app applications.Application, logger zerolog.Logger) (*ContainerConfig, error) {
	containers, err := d.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}

	flattentedContainerNames := flattenNames(containers)
	logger.Info().Strs("allContainerNames", flattentedContainerNames).Msg("searching for the app in containers")

	// TODO we should probably restart -all- containers that match the regexes, so that
	// we deal with having replicas correctly.
	containerName := findContainerNameOfApp(containers, app)
	if containerName == "" {
		return nil, errors.New("no container currently up matches the regexes of the app")
	}

	containerJSON, err := d.client.ContainerInspect(ctx, containerName)
	if err != nil {
		return nil, err
	}

	return &ContainerConfig{
		containerName:    containerName,
		id:               containerJSON.ID,
		config:           containerJSON.Config,
		hostConfig:       containerJSON.HostConfig,
		networkingConfig: &network.NetworkingConfig{EndpointsConfig: containerJSON.NetworkSettings.Networks},
	}, nil
}

func flattenNames(containers []types.Container) []string {
	names := make([]string, 0, len(containers))
	for _, container := range containers {
		for _, name := range container.Names {
			names = append(names, name)
		}
	}
	return names
}

func findContainerNameOfApp(containers []types.Container, app applications.Application) string {
	for _, container := range containers {
		for _, name := range container.Names {
			if app.Matches(name) {
				return name
			}
		}
	}
	return ""
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
	created, err := d.client.ContainerCreate(ctx, config.config, config.hostConfig, config.networkingConfig, config.containerName)
	if err != nil {
		return err
	}
	return d.client.ContainerStart(ctx, created.ID, types.ContainerStartOptions{})
}
