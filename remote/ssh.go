package remote

import (
	"fmt"
	"net"

	"github.com/movitz-s/ssh-spawner/config"
	"github.com/movitz-s/ssh-spawner/shells"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

// Server listens to SSH reqs and delegates to ShellService
type Server struct {
	serverConfig *ssh.ServerConfig
	ss           shells.ShellService
	config       *config.Config
}

// NewServer constructs a new Server
func NewServer(serverConfig *ssh.ServerConfig, ss shells.ShellService, config *config.Config) *Server {
	return &Server{serverConfig, ss, config}
}

// Start initializes a tcp connection and delegate requests
func (server *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", server.config.SSH.Host, server.config.SSH.Port)
	fmt.Println("Listening on ", addr)
	listener, err := net.Listen("tcp", addr)

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
	_, newChannels, _, err := ssh.NewServerConn(conn, server.serverConfig)
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
				ok := req.Type == "shell" || req.Type == "pty-req"
				req.Reply(ok, nil)
			}
		}(requests)

		server.handle(channel)

	}

}
