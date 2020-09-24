package containers

import (
	"io"
)

type Shell interface {
	io.Reader
	io.Writer
	io.Closer
}

// ContainerService describes a service that serves shells
type ContainerService interface {
	GetShell() (shell Shell, err error)
}
