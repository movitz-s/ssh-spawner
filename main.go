package main

import (
	"fmt"
	"io/ioutil"

	"github.com/movitz-s/ssh-spawner/remote"

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

	ss := InitializeShellService()

	server := remote.NewServer(config, ss, "localhost", 22)
	server.Start()

}

func loadPrivateKey() ssh.Signer {
	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		panic(fmt.Sprintf("Error while loading private key: %v\n", err))
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing private key: %v\n", err))
	}
	return private
}
