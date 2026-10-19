package main

import (
	"fmt"
	"net"

	"git.trivaris.org/trivaris/keyvalue/internal/resp"
)

type client struct {
}

func (c *client) sendRequest(request []byte, conn *net.TCPConn) (resp.RespItem, error) {
	_, err := conn.Write(request)
	if err != nil {
		fmt.Print(err)
		return nil, fmt.Errorf("Write Error: %s\n", err)
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Read Error: %s\n", err)
	}

	item, _, err := resp.Deserialize(buf[:n], 0)
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Deserialization Error: %s\n", err)
	}

	return item, nil
}

func (c *client) createConn(address string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
