package spawner

import (
	"io/ioutil"
	"net/http"

	docker "github.com/docker/docker/client"
	"github.com/movitz-s/ssh-spawner/internal/spawner/config"
	"github.com/movitz-s/ssh-spawner/internal/spawner/shells"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

func loadPrivateKey() (ssh.Signer, error) {
	privateBytes, err := ioutil.ReadFile("id_rsa")
	if err != nil {
		return nil, errors.Wrap(err, "Could not load SSH private key file")
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	return private, errors.Wrap(err, "Could not parse SSH private key")
}

func newSSHConfig(key ssh.Signer) *ssh.ServerConfig {
	challCallback := func(conn ssh.ConnMetadata, client ssh.KeyboardInteractiveChallenge) (*ssh.Permissions, error) {
		answers, err := client(conn.User(), "hej du", []string{"hehe ( :"}, []bool{true})
		if err != nil {
			return nil, err
		}
		if answers[0] != "nja" {
			return nil, errors.New("nope")
		}
		return nil, nil
	}

	var config ssh.ServerConfig
	config.KeyboardInteractiveCallback = challCallback
	config.AddHostKey(key)
	return &config
}

func newDockerClient() (*docker.Client, func(), error) {
	client, err := docker.NewClient("tcp://127.0.0.1:2375", "", &http.Client{Transport: http.DefaultTransport}, map[string]string{})

	cleanup := func() {
		client.Close()
	}

	return client, cleanup, errors.Wrap(err, "Could not create docker client")
}

func newConfig() *config.Config {
	return &config.Config{
		Images: []config.Image{
			{
				DisplayName: "Debian",
				ImageID:     shells.ImageID("debian"),
			},
			{
				DisplayName: "Hackerbox",
				ImageID:     shells.ImageID("debian"),
			},
		},
		SSH: config.SSHConfig{
			Port: 22,
			Host: "0.0.0.0",
		},
	}
}
