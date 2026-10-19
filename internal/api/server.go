package api

import (
	"errors"
	"fmt"
	"io"
	"net"

	"git.trivaris.org/trivaris/keyvalue/internal/config"
)

type Server struct {
	handler *Handler
	config  *config.Config
}

func NewServer(handler *Handler, config *config.Config) *Server {
	return &Server{
		handler: handler,
		config:  config,
	}
}

func (s *Server) Listen() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", s.config.ListenAddress, s.config.ListenPort))
	if err != nil {
		fmt.Println(err)
		return
	}

	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = server.Close() }()

	fmt.Printf("server listening on %s:%d\n", s.config.ListenAddress, s.config.ListenPort)

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn *net.TCPConn) {
	defer func() { _ = conn.Close() }()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		
		if n > 0 {
			response := s.handler.handleRequest(buf[:n], conn).Serialize()
			if _, err := conn.Write(response); err != nil {
				fmt.Println(err)
				return
			}
		}
		
		if err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Println(err)
			}
			return
		}
	}
}
