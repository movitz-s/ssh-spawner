package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	docker "github.com/docker/docker/client"
	"golang.org/x/crypto/ssh"
)

// GitCommit is assigned at build time
var GitCommit string

func main() {
	if GitCommit == "" {
		GitCommit = "<unknown>"
	}
	fmt.Printf("SSH Loadbalancer\nGit commit %s\n", GitCommit)

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	config.AddHostKey(loadPrivateKey())

	client, err := docker.NewClient("tcp://127.0.0.1:2375", "", &http.Client{Transport: http.DefaultTransport}, map[string]string{})
	if err != nil {
		panic(err)
	}

	ds := NewSimpleDockerService(client)

	server := NewSSHServer(config, ds, "localhost", 22)
	server.Start()

}

func loadPrivateKey() (private ssh.Signer) {
	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		panic(fmt.Sprintf("Error while loading private key: %v\n", err))
	}

	private, err = ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing private key: %v\n", err))
	}
	return
}
