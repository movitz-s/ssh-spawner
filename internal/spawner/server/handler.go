package server

import (
	"fmt"
	"io"

	"github.com/manifoldco/promptui"
	"github.com/movitz-s/ssh-spawner/internal/spawner/shells"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func (server *Server) handle(channel ssh.Channel) {

	defer channel.Close()
	server.displayBanner(channel)
	imageID, err := server.pickImageID(channel)

	if err != nil {
		fmt.Printf("Could not get image id\n")
		fmt.Println(err)
		channel.Write([]byte("Something went wrong. Try again later.\n\r"))
		return
	}

	shell, err := server.ss.GetShell(imageID)

	if err != nil {
		fmt.Printf("Could not get a shell")
		fmt.Print(err)
		channel.Write([]byte("Could not allocate a shell for you. Try again later.\n\r"))
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

func (server *Server) displayBanner(channel ssh.Channel) {
	if server.config.SSH.Banner != "" {
		channel.Write([]byte(server.config.SSH.Banner + "\n\r"))
	}
}

func (server *Server) pickImageID(channel ssh.Channel) (shells.ImageID, error) {
	images := server.config.Images

	if len(images) == 0 {
		return shells.ImageID(""), errors.New("Server misconfigured, no images present")
	}

	if len(images) == 1 {
		return images[0].ImageID, nil
	}

	var displayNames []string
	for _, image := range images {
		displayNames = append(displayNames, image.DisplayName)
	}

	promt := promptui.Select{
		Stdin:    channel,
		Stdout:   channel,
		HideHelp: true,
		Label:    "Select challange",
		Items:    displayNames,
	}

	index, _, err := promt.Run()

	if err != nil {
		return shells.ImageID(""), errors.Wrap(err, "Cloud not pick image id")
	}

	if index < 0 || index >= len(images) {
		return shells.ImageID(""), errors.New("Received invalid index from promt")
	}

	return images[index].ImageID, nil
}
