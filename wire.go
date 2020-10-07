//+build wireinject
package main

import (
	"net/http"

	"github.com/google/wire"
	"github.com/movitz-s/ssh-spawner/shells"

	docker "github.com/docker/docker/client"
)

func InitializeShellService() shells.ShellService {
	panic(wire.Build(NewDockerClient, shells.NewDockerShellService))
}

func NewDockerClient() *docker.Client {
	client, err := docker.NewClient("tcp://127.0.0.1:2375", "", &http.Client{Transport: http.DefaultTransport}, map[string]string{})
	if err != nil {
		panic(err)
	}
	return client
}
