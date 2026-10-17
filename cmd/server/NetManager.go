package main

import (
	"fmt"
	"net"
)

const NETWORK_INFO string = "[NETWORK INFO]"
const NETWORK_ERROR string = "[NETWORK ERROR]"
const LISTEN_ADDR string = "0.0.0.0:8080"

type netManager struct {
	logic logicManager
}

func (this netManager) Listen() {
	server, err := net.Listen("tcp", LISTEN_ADDR)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer server.Close()
	fmt.Printf("%s Server Listening On %q\n", NETWORK_INFO, LISTEN_ADDR)
	
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Printf("%s %s\n", NETWORK_ERROR, err.Error())
			continue
		}

		go this.handleConnection(conn)
	}
}

func (this netManager) handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Printf("%s Read Error: %s\n", NETWORK_ERROR, err.Error())
		return
	}

	fmt.Printf("%s Received Request\n", NETWORK_INFO)
	response := this.logic.handleRequest(n, buf)
	response = append(response, []byte("\r\n")...)
	
	if _, err := conn.Write(response); err != nil {
		fmt.Printf("%s Write Error: %s\n", NETWORK_ERROR, err.Error())
		return
	}
}
