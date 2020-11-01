//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/movitz-s/ssh-spawner/remote"
	"github.com/movitz-s/ssh-spawner/shells"
)

func initializeShellService() (shells.ShellService, error) {
	panic(wire.Build(newDockerClient, shells.NewDockerShellService))
}

func initializeSSHServer() (*remote.Server, error) {
	panic(wire.Build(loadPrivateKey, newSSHConfig, remote.NewServer, initializeShellService, newConfig))
}
