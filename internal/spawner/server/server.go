package server

import (
	"fmt"
	"net"
	"sync"

	"github.com/movitz-s/ssh-spawner/internal/spawner/config"
	"github.com/movitz-s/ssh-spawner/internal/spawner/shells"
	"github.com/pkg/errors"
	"golang.org/x/crypto/ssh"
)

type Session struct {
	sshChannel ssh.Channel
	shell      shells.Shell
}

// Server listens to SSH reqs and delegates to ShellService
type Server struct {
	serverConfig *ssh.ServerConfig
	ss           shells.ShellService
	config       config.Config
	ips          map[string]bool
	ipsLock      sync.Mutex
}

// NewServer constructs a new Server
func NewServer(serverConfig *ssh.ServerConfig, ss shells.ShellService, config config.Config) *Server {
	return &Server{
		serverConfig: serverConfig,
		ss:           ss,
		config:       config,
		ips:          make(map[string]bool),
	}
}

// Start initializes a tcp connection and delegate requests
func (server *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", server.config.SSH.Host, server.config.SSH.Port)
	fmt.Println("Listening on", addr)
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return errors.Wrap(err, "Could not start SSH server")
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return errors.Wrap(err, "Could not listen on SSH server")
		}

		go func() {
			server.bootstrap(conn)
		}()
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

		// IP Check
		ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
		if !server.ipCheck(ip) {
			channel.Write([]byte("Your IP has already made a connection to this server\n"))
			channel.Close()
			continue
		}
		defer server.ipCheckClear(ip)

		server.handle(channel)
	}

}

func (server *Server) ipCheck(ip string) bool {
	server.ipsLock.Lock()
	defer server.ipsLock.Unlock()
	if server.ips[ip] {
		return false
	}
	server.ips[ip] = true
	return true
}

func (server *Server) ipCheckClear(ip string) {
	server.ipsLock.Lock()
	delete(server.ips, ip)
	server.ipsLock.Unlock()
}
