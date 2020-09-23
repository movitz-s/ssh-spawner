package main

import (
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// GitCommit
var GitCommit string

func main() {
	if GitCommit == "" {
		GitCommit = "Unknown"
	}
	fmt.Printf("SSH Loadbalancer\nGit commit %s\n", GitCommit)

	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		panic(fmt.Sprintf("Error while loading private key: %v\n", err))
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		panic(fmt.Sprintf("Error while parsing private key: %v\n", err))
	}

	config := &ssh.ServerConfig{
		NoClientAuth: true,
	}

	config.AddHostKey(private)

	server := NewSSHServer(config, "localhost", 22)
	server.Start()

}
