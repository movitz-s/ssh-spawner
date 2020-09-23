package main

import (
	"fmt"
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

type SSHServer struct {
	config *ssh.ServerConfig
	ds     *DockerService
	host   string
	port   int
}

// NewSSHServer constructs a new SSHServer
func NewSSHServer(config *ssh.ServerConfig, host string, port int) *SSHServer {
	return &SSHServer{
		config,
		NewDockerService(),
		host,
		port,
	}
}

// Start initializes a tcp connection and delegate requests
func (server *SSHServer) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.host, server.port))
	if err != nil {
		panic(fmt.Sprintf("Could not listen: %v", err))
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(fmt.Sprintf("Failed to accept incoming connection: %v\n", err))
		}

		go server.bootstrap(conn)
	}
}

func (server *SSHServer) bootstrap(conn net.Conn) {
	_, newChannels, _, err := ssh.NewServerConn(conn, server.config)
	if err != nil {
		fmt.Println("Failed to handshake with new client: %v", err)
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
			return
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

func (server *SSHServer) handle(channel ssh.Channel) {

	defer func() {
		channel.Close()
	}()

	hijack, err := server.ds.HijackShell()

	if err != nil {
		fmt.Printf("Could not get a hijack connection: %+v\n", err)
		return
	}

	go func() {
		_, err := io.Copy(hijack.Conn, channel)
		if err != nil {
			fmt.Printf("Error while copying from hijack to client: %v\n", err)
		}
	}()
	_, err = io.Copy(channel, hijack.Conn)
	if err != nil {
		fmt.Printf("Error while copying from client to hijack: %v\n", err)
	}
}
