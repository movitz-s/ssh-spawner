package shells

import (
	"io"
)

type Shell interface {
	io.Reader
	io.Writer
	io.Closer
}

type ImageID string

// ShellService describes a service that serves shells
type ShellService interface {
	GetShell(imageID ImageID) (shell Shell, err error)
}
