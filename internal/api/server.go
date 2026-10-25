package api

import (
	"errors"
	"fmt"
	"io"
	"net"
	"net/netip"
)

type Server struct {
	handler       *Handler
	cbHandler *ClusterBusHandler
	listenAddress netip.AddrPort
	busAddress    netip.AddrPort
}

func NewServer(handler *Handler, cbHandler *ClusterBusHandler, addr netip.Addr, port uint16, cport uint16) *Server {
	return &Server{
		handler:       handler,
		cbHandler: cbHandler,
		listenAddress: netip.AddrPortFrom(addr, port),
		busAddress:    netip.AddrPortFrom(addr, cport),
	}
}

func (s *Server) Listen() {
	server, err := net.ListenTCP("tcp", net.TCPAddrFromAddrPort(s.listenAddress))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = server.Close() }()

	fmt.Printf("server listening on %s\n", s.listenAddress.String())

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) ListenClusterBus() {
	server, err := net.ListenTCP("tcp", net.TCPAddrFromAddrPort(s.busAddress))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() { _ = server.Close() }()

	fmt.Printf("cluster bus server listening on %s\n", s.listenAddress.String())

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go s.handleClusterBusConnection(conn)
	}
}

func (s *Server) handleConnection(conn *net.TCPConn) {
	defer func() { _ = conn.Close() }()
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)

	if n > 0 {
		response := s.handler.handleRequest(buf[:n]).Serialize()
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

func (s *Server) handleClusterBusConnection(conn *net.TCPConn) {
	defer func() { _ = conn.Close() }()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)

		if n > 0 {
			response := s.cbHandler.handleRequest(buf[:n])
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
