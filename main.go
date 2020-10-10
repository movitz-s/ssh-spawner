package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	docker "github.com/docker/docker/client"
	"github.com/movitz-s/ssh-spawner/remote"
	"github.com/pkg/errors"

	"golang.org/x/crypto/ssh"
)

// GitCommit is assigned at build time
var GitCommit = "<unknown>"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Printf("SSH Spawner\nGit commit %s\n", GitCommit)

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	key, err := loadPrivateKey()
	if err != nil {
		return err
	}
	config.AddHostKey(key)

	ss, err := initializeShellService()
	if err != nil {
		return err
	}

	server := remote.NewServer(config, ss, "localhost", 22)
	err = server.Start()
	return err
}

func loadPrivateKey() (ssh.Signer, error) {
	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		return nil, errors.Wrap(err, "Could not load SSH private key file")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	return private, errors.Wrap(err, "Could not parse SSH private key")
}

func newDockerClient() (*docker.Client, error) {
	client, err := docker.NewClient("tcp://127.0.0.1:2375", "", &http.Client{Transport: http.DefaultTransport}, map[string]string{})
	return client, errors.Wrap(err, "Could not create docker client")
}
