package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
)

// DockerService manages everything docker related
type DockerService struct {
	targetImageID string
	dockerClient  *docker.Client
}

// NewDockerService constructs a new DockerService
func NewDockerService() *DockerService {
	client, err := docker.NewClient("tcp://127.0.0.1:2375", "", &http.Client{Transport: http.DefaultTransport}, map[string]string{})
	if err != nil {
		fmt.Println(err)
	}
	return &DockerService{
		"debian",
		client,
	}
}

// HijackShell starts a container and retreives a hijacked shell from the container
func (dm DockerService) HijackShell() (hijack types.HijackedResponse, err error) {

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
