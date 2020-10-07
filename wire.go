//+build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/movitz-s/ssh-spawner/shells"
)

func InitializeShellService() shells.ShellService {
	panic(wire.Build(NewDockerClient, shells.NewDockerShellService))
}
