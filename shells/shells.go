package shells

import (
	"io"
)

type Shell interface {
	io.Reader
	io.Writer
	io.Closer
}

// ShellService describes a service that serves shells
type ShellService interface {
	GetShell() (shell Shell, err error)
}
