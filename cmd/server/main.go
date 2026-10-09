package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"strconv"
)

var (
	values sync.Map
)

func main() {
	values.Store("key", "value")
	startListen()
}

func startListen() {
	println("Server init...")
	server, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer server.Close()

	println("Waiting for connections...")
	for {
		conn, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Read error:", err)
		}
		return
	}

	response := []byte("Key Not found or malformed\n")
	key, err := handleInput(n, buf)
	if err == nil {
		val, ok := values.Load(key)
		if ok {
			response = fmt.Appendf(nil, "%v\n", val)
		}
	}

	if _, err := conn.Write(response); err != nil {
		fmt.Println("Write error:", err)
		return
	}
}

func handleInput(n int, buf []byte) (any, error) {
	fmt.Println("Handling Input...")
	switch buf[0] {
	case []byte("+")[0]: {
		end := n - len([]byte("\r\n")) 
		str := buf[1:end]
		return string(str), nil;
	}
	// case []byte("-")[0]: {
	// 	end := n - len([]byte("\r\n")) 
	// 	str := buf[1:end]
	// 	return fmt.Errorf("%q", str), nil
	// }
	case []byte(":")[0]: 
		end := n - len([]byte("\r\n"))
		n, err := strconv.Atoi(string(buf[1:end]))
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		return n, nil
	default: return nil, fmt.Errorf("Invalid Type")
	}
}