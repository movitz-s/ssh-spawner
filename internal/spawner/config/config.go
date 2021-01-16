package config

import (
	"github.com/movitz-s/ssh-spawner/internal/spawner/shells"
)

// Config holds the configuration for the whole application
type Config struct {
	DockerURL string    `yaml:"docker_url"`
	Images    []Image   `yaml:"images"`
	SSH       SSHConfig `yaml:"ssh"`
}

type Image struct {
	ImageID     shells.ImageID `yaml:"image_id"`
	DisplayName string         `yaml:"display_name"`
}

type SSHConfig struct {
	Port   int    `yaml:"port"`
	Host   string `yaml:"host"`
	Banner string `yaml:"banner"`
}
