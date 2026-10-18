package main

import (
	"fmt"
	"net"
)

const NETWORK_INFO string = "[NETWORK INFO]"
const NETWORK_ERROR string = "[NETWORK ERROR]"
const LISTEN_ADDR string = "0.0.0.0"
const LISTEN_PORT uint16 = 8080

type server struct {
	handler *handler
}

func (s *server) listen() {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", LISTEN_ADDR, LISTEN_PORT))
	if err != nil {
		fmt.Println(err)
		return
	}

	server, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer server.Close()

	fmt.Printf("%s Server Listening On %q on port %d\n", NETWORK_INFO, LISTEN_ADDR, LISTEN_PORT)

	for {
		conn, err := server.AcceptTCP()
		if err != nil {
			fmt.Printf("%s %s\n", NETWORK_ERROR, err.Error())
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *server) handleConnection(conn *net.TCPConn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("%s Read Error: %s\n", NETWORK_ERROR, err.Error())
			return
		}

		fmt.Printf("%s Received Request\n", NETWORK_INFO)

		response := s.handler.handleRequest(buf[:n], conn)
		response = append(response, []byte("\r\n")...)

		if _, err := conn.Write(response); err != nil {
			fmt.Printf("%s Write Error: %s\n", NETWORK_ERROR, err.Error())
			return
		}
	}
}
