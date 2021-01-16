package main

import (
	"fmt"
	"os"

	"github.com/movitz-s/ssh-spawner/internal/spawner"
)

// GitCommit is assigned at build time
var GitCommit = "<unknown>"

func main() {
	err := realMain()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func realMain() error {
	fmt.Printf("SSH Spawner\nGit commit %s\n", GitCommit)

	server, cleanup, err := spawner.InitializeSSHServer()
	if err != nil {
		return err
	}
	defer cleanup()

	err = server.Start()
	return err
}
