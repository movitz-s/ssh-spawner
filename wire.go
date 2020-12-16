//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/movitz-s/ssh-spawner/server"
	"github.com/movitz-s/ssh-spawner/shells"
)

func initializeShellService() (shells.ShellService, error) {
	panic(wire.Build(newDockerClient, shells.NewDockerShellService))
}

func initializeSSHServer() (*server.Server, error) {
	panic(wire.Build(loadPrivateKey, newSSHConfig, server.NewServer, initializeShellService, newConfig))
}
