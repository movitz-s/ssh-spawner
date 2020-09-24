package containers

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
)

// SimpleDockerService only creates containers, nothing fancy
type SimpleDockerService struct {
	targetImageID string
	dockerClient  *docker.Client
}

// NewSimpleDockerService constructs a new ContainerService with a docker backend
func NewSimpleDockerService(client *docker.Client) ContainerService {
	return SimpleDockerService{"debian", client}
}

// GetShell starts a container and retreives a shell from the container
func (dm SimpleDockerService) GetShell() (Shell, error) {

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
		return nil, err
	}

	hijack, err := dm.dockerClient.ContainerAttach(context.Background(), resp.ID, types.ContainerAttachOptions{
		Stderr: true,
		Stdin:  true,
		Stdout: true,
		Stream: true,
	})

	if err != nil {
		return nil, err
	}

	err = dm.dockerClient.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{})
	return hijack.Conn, err

}
