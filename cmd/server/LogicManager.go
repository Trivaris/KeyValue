package main

import (
	"fmt"
	"net"
	"git.trivaris.org/trivaris/keyvalue/cmd/internal/resp"
)

const INTERNAL_ERROR = "[INTERNAL ERROR]"
const INTERNAL_INFO = "[INTERNAL INFO]"

type logicManager struct {
	data *dataManager
	sharding *shardingManager
}

func (this *logicManager) handleRequest(n int, buf []byte, conn *net.TCPConn) ([]byte, bool) {
	item, _, err := resp.Deserialize(buf[:n], 0)

	if err != nil {
		fmt.Println(err)
		return []byte(err.Error()), true
	}

	arr, ok := item.(resp.RespArray)
	if !ok {
		msg := fmt.Sprintf("%s No Command Array Found: %s\n", INTERNAL_ERROR, item.Serialize())
		fmt.Println(msg)
		return []byte(msg), true
	}

	command, ok := arr.Values[0].(resp.RespString)
	if !ok {
		msg := fmt.Sprintf("%s No Command Keyword Found: %s\n", INTERNAL_ERROR, arr.Values[0])
		fmt.Println(msg)
		return []byte(msg), true
	}

	switch command.Value {
		case "GET": {
			values := make([]resp.RespItem, len(arr.Values) - 1)

			for i, rawKey := range(arr.Values[1:]) {
				rawKey, ok := rawKey.(resp.RespString)
				if (!ok) {
					msg := fmt.Sprintf("%s GET Request Array Does Not Consist Only Of Keys Of Type String: %s\n", INTERNAL_ERROR, arr.Values[0])
					fmt.Println(msg)
					return []byte(msg), true
				}

				key := rawKey.Value

				item, ok = this.data.getValue(key)
				if !ok {
					item = resp.RespNull{}
				}

				values[i] = item
			}
			response := resp.RespArray{Values: values}
			fmt.Printf("%s GET Request Successfully completed: %#v\n", INTERNAL_INFO, response)
			return response.Serialize(), true
		}
		case "SET": {
			if len(arr.Values) != 3 {
				msg := fmt.Sprintf("%s SET Request Array Does Not Consist Of Exactly 2 arguments: %s\n", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg), true
			}
			
			key, ok := arr.Values[1].(resp.RespString)
			if !ok {
				msg := fmt.Sprintf("%s SET Request Key Is Not Of Type String: %s\n", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg), true
			}
			
			value := arr.Values[2]
			this.data.setValue(key.Value, value)
			fmt.Printf("%s SET Request Successfully completed. Key: %q, Value: %#v\n", INTERNAL_INFO, key, value)
			return resp.RespString{Value: "Success!"}.Serialize(), true
		}
		case "ADDSERVER": {
			if len(arr.Values) != 3 {
				msg := fmt.Sprintf("%s ADDSERVER Request Array Does Not Consist Of Exactly 2 arguments: %s\n", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg), true
			}

			name, ok := arr.Values[1].(resp.RespString)
			if !ok {
				msg := fmt.Sprintf("%s ADDSERVER Request Name Is Not Of Type String: %s\n", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg), true
			}

			address, ok := arr.Values[2].(resp.RespString)
			if !ok {
				msg := fmt.Sprintf("%s ADDSERVER Request Name Is Not Of Type String: %s\n", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg), true
			}

			fmt.Printf("%s Adding %s at address %s\n", DATA_INFO, name.Value, address.Value)

			tcpAddr, err := net.ResolveTCPAddr("tcp", address.Value)
			if err != nil {
				fmt.Println(err)
				return []byte(err.Error()), true
			}
			
			dialConn, err := net.DialTCP("tcp", nil, tcpAddr)
			if err != nil {
				fmt.Println(err)
				return []byte(err.Error()), true
			}
			
			this.sharding.addServerConn(name.Value, dialConn)
			request := resp.RespArray{ Values: []resp.RespItem{
				resp.RespString{Value: "PING"},
				resp.RespString{Value: NODE_NAME},
				resp.RespInteger{Value: int64(LISTEN_PORT)},
			}}.Serialize()

			_, err = dialConn.Write(request)
			if err != nil {
				fmt.Println(err)
				return []byte(err.Error()), true
			}

			buf := make([]byte, 1024)
			n, err := dialConn.Read(buf)
			if err != nil {
				fmt.Println(err)
				return []byte(err.Error()), true
			}

			response, _, err := resp.Deserialize(buf[:n], 0)
			if err != nil {
				fmt.Println(err)
				return []byte(err.Error()), true
			}

			item, ok := response.(resp.RespString)
			if !ok || item.Value != "Success!" {
				fmt.Println("Did not get a success response")
			}
			
			return resp.RespString{Value: "Success!"}.Serialize(), true
		}
		case "PING": {
			if len(arr.Values) != 3 {
				msg := fmt.Sprintf("%s PING Request Array Does Not Consist Of Exactly 2 arguments: %s\n", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg), true
			}

			name, ok := arr.Values[1].(resp.RespString)
			if !ok {
				msg := fmt.Sprintf("%s PING Request Name Is Not Of Type String: %s\n", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg), true
			}

			port, ok := arr.Values[2].(resp.RespInteger)
			if !ok {
				msg := fmt.Sprintf("%s PING Request Port Is Not Of Type Integer: %s\n", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg), true
			}

			tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
			if !ok {
				msg := fmt.Sprintf("%s Could not convert Addr to TCPAddr: %s\n", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg), true
			}

			targetAddr := &net.TCPAddr{
				IP:   tcpAddr.IP,
				Port: int(port.Value),
			}

			dialConn, err := net.DialTCP("tcp", nil, targetAddr)
			if err != nil {
				fmt.Println(err)
				return []byte(err.Error()), true
			}

			this.sharding.addServerConn(name.Value, dialConn)
			
			return resp.RespString{Value: "Success!"}.Serialize(), true
		}
		default: {
			msg := fmt.Sprintf("%s %s Request Not Yet Implemented: %s\n", INTERNAL_ERROR, command.Value, arr.Values)
			fmt.Println(msg)
			return []byte(msg), true
		}
	}
}
