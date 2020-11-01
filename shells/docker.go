package shells

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	docker "github.com/docker/docker/client"
	"github.com/pkg/errors"
)

// DockerShellService only creates containers, nothing fancy
type DockerShellService struct {
	dockerClient *docker.Client
}

// NewDockerShellService constructs a new ShellService with a docker backend
func NewDockerShellService(client *docker.Client) ShellService {
	return DockerShellService{client}
}

// GetShell starts a container and retreives a shell from the container
func (dm DockerShellService) GetShell(imageID ImageID) (Shell, error) {

	resp, err := dm.dockerClient.ContainerCreate(
		context.TODO(),
		&container.Config{
			Image:           string(imageID),
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
		return nil, errors.Wrap(err, "Could not create container")
	}

	hijack, err := dm.dockerClient.ContainerAttach(context.TODO(), resp.ID, types.ContainerAttachOptions{
		Stderr: true,
		Stdin:  true,
		Stdout: true,
		Stream: true,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Could not attach container")
	}

	err = dm.dockerClient.ContainerStart(context.TODO(), resp.ID, types.ContainerStartOptions{})
	return hijack.Conn, errors.Wrap(err, "Could not start container")

}
