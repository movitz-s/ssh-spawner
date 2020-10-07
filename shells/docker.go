package shells

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
)

// DockerShellService only creates containers, nothing fancy
type DockerShellService struct {
	targetImageID string
	dockerClient  *docker.Client
}

// NewDockerShellService constructs a new ShellService with a docker backend
func NewDockerShellService(client *docker.Client) ShellService {
	return DockerShellService{"debian", client}
}

// GetShell starts a container and retreives a shell from the container
func (dm DockerShellService) GetShell() (Shell, error) {

	resp, err := dm.dockerClient.ContainerCreate(
		context.TODO(),
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

	hijack, err := dm.dockerClient.ContainerAttach(context.TODO(), resp.ID, types.ContainerAttachOptions{
		Stderr: true,
		Stdin:  true,
		Stdout: true,
		Stream: true,
	})

	if err != nil {
		return nil, err
	}

	err = dm.dockerClient.ContainerStart(context.TODO(), resp.ID, types.ContainerStartOptions{})
	return hijack.Conn, err

}
