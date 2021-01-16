//+build wireinject

package spawner

import (
	"github.com/google/wire"
	"github.com/movitz-s/ssh-spawner/internal/spawner/server"
	"github.com/movitz-s/ssh-spawner/internal/spawner/shells"
)

func initializeShellService() (shells.ShellService, func(), error) {
	panic(wire.Build(newDockerClient, shells.NewDockerShellService))
}

func InitializeSSHServer() (*server.Server, func(), error) {
	panic(wire.Build(loadPrivateKey, newSSHConfig, server.NewServer, initializeShellService, newConfig))
}
