package spawner

import (
	"io/ioutil"

	docker "github.com/docker/docker/client"
	"github.com/movitz-s/ssh-spawner/internal/spawner/config"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
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
	var config ssh.ServerConfig
	config.AddHostKey(key)
	config.NoClientAuth = true
	return &config
}

func newDockerClient(c config.Config) (*docker.Client, func(), error) {
	client, err := docker.NewClient(c.DockerURL, "", nil, map[string]string{})

	cleanup := func() {
		client.Close()
	}

	return client, cleanup, errors.Wrap(err, "Could not create docker client")
}

func newConfig() (c config.Config, err error) {
	bs, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		return
	}
	err = yaml.Unmarshal(bs, &c)
	return
}
