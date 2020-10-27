package remote

import (
	"fmt"
	"io"

	"golang.org/x/crypto/ssh"
)

func (server *Server) handle(channel ssh.Channel) {

	defer channel.Close()

	shell, err := server.ss.GetShell()

	if err != nil {
		fmt.Printf("Could not get a shell: %+v\n", err)
		return
	}

	go func() {
		_, err := io.Copy(shell, channel)
		if err != nil {
			fmt.Printf("Error while copying from shell to client: %v\n", err)
		}
	}()

	_, err = io.Copy(channel, shell)
	if err != nil {
		fmt.Printf("Error while copying from client to shell: %v\n", err)
	}
}
