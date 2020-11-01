package config

import (
	"github.com/movitz-s/ssh-spawner/shells"
)

// Config holds the configuration for the whole application
type Config struct {
	Images []Image
	SSH    SSHConfig
}

type Image struct {
	ImageID     shells.ImageID
	DisplayName string
}

type SSHConfig struct {
	Port   int
	Host   string
	Banner string
}
