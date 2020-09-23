package main

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
)

// DockerService describes a service that serves shells from docker
type DockerService interface {
	HijackShell() (types.HijackedResponse, error)
}

// SimpleDockerService only creates containers, nothing fancy
type SimpleDockerService struct {
	targetImageID string
	dockerClient  *docker.Client
}

// NewSimpleDockerService constructs a new DockerService
func NewSimpleDockerService(client *docker.Client) DockerService {
	return SimpleDockerService{"debian", client}
}

// HijackShell starts a container and retreives a hijacked shell from the container
func (dm SimpleDockerService) HijackShell() (hijack types.HijackedResponse, err error) {

	resp, err := dm.dockerClient.ContainerCreate(
		context.Background(),
		&container.Config{
			Image:           dm.targetImageID,
			NetworkDisabled: true,
			Cmd:             []string{"bash"},
			OpenStdin:       true,
			AttachStderr:    true,
			AttachStdin:     true,
			AttachStdout:    true,
			StdinOnce:       true,
			Tty:             true,
		},
		&container.HostConfig{
			AutoRemove: true,
		},
		&network.NetworkingConfig{},
		"",
	)

	if err != nil {
		return
	}

	hijack, err = dm.dockerClient.ContainerAttach(context.Background(), resp.ID, types.ContainerAttachOptions{
		Stderr: true,
		Stdin:  true,
		Stdout: true,
		Stream: true,
	})

	if err != nil {
		return
	}

	err = dm.dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	return
}
