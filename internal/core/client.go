package core

import (
	"fmt"
	"net"

	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

func SendRequest(request []byte, conn *net.TCPConn) (serialization.RespItem, error) {
	_, err := conn.Write(request)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	item, _, err := serialization.Deserialize(buf[:n], 0)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return item, nil
}

func CreateConn(address string) (*net.TCPConn, error) {
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
