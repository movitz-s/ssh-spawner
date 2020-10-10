package remote

import (
	"fmt"
	"io"
	"net"

	"github.com/movitz-s/ssh-spawner/shells"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// Server listens to SSH reqs and delegates to ShellService
type Server struct {
	config *ssh.ServerConfig
	ss     shells.ShellService
	host   string
	port   int
}

// NewServer constructs a new Server
func NewServer(config *ssh.ServerConfig, ss shells.ShellService, host string, port int) *Server {
	return &Server{config, ss, host, port}
}

// Start initializes a tcp connection and delegate requests
func (server *Server) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.host, server.port))

	if err != nil {
		return errors.Wrap(err, "Could not start SSH server")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return errors.Wrap(err, "Could not listen on SSH server")
		}

		go server.bootstrap(conn)
	}
}

func (server *Server) bootstrap(conn net.Conn) {
	_, newChannels, _, err := ssh.NewServerConn(conn, server.config)
	if err != nil {
		fmt.Printf("Handshake failed with client: %v\n", err)
		return
	}

	for newChannel := range newChannels {
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			fmt.Printf("Could not accept channel: %v\n", err)
			continue
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				switch req.Type {
				case "shell", "pty-req":
					req.Reply(true, nil)
				default:
					req.Reply(false, nil)
				}
			}
		}(requests)

		server.handle(channel)

	}

}

func (server *Server) handle(channel ssh.Channel) {

	defer func() {
		channel.Close()
	}()

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
