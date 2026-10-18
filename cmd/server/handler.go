package main

import (
	"fmt"
	"net"

	"git.trivaris.org/trivaris/keyvalue/cmd/internal/resp"
)

const INTERNAL_ERROR = "[INTERNAL ERROR]"
const INTERNAL_INFO = "[INTERNAL INFO]"

type handler struct {
	data     *store
	cluster *cluster
	client *client
}

func (h *handler) handleRequest(buf []byte, conn *net.TCPConn) []byte {
	item, _, err := resp.Deserialize(buf, 0)

	if err != nil {
		fmt.Println(err)
		return []byte(err.Error())
	}

	arr, ok := item.(resp.RespArray)
	if !ok {
		msg := fmt.Sprintf("%s No Command Array Found: %s\n", INTERNAL_ERROR, item.Serialize())
		fmt.Println(msg)
		return []byte(msg)
	}

	command, ok := arr.Values[0].(resp.RespString)
	if !ok {
		msg := fmt.Sprintf("%s No Command Keyword Found: %s\n", INTERNAL_ERROR, arr.Values[0])
		fmt.Println(msg)
		return []byte(msg)
	}

	switch command.Value {
	case "GET": 		return h.handleGet(arr)
	case "SET": 		return h.handleSet(arr)
	case "ADDSERVER": 	return h.handleAddServer(arr)
	case "PING":		return h.handlePing(arr, conn)
	default:
		{
			msg := fmt.Sprintf("%s %s Request Not Yet Implemented: %s\n", INTERNAL_ERROR, command.Value, arr.Values)
			fmt.Println(msg)
			return []byte(msg)
		}
	}
}

func (h *handler) handleGet(arr resp.RespArray) []byte {
	values := make([]resp.RespItem, len(arr.Values)-1)

	for i, rawKey := range arr.Values[1:] {
		rawKey, ok := rawKey.(resp.RespString)
		if !ok {
			msg := fmt.Sprintf("%s GET Request Array Does Not Consist Only Of Keys Of Type String: %s\n", INTERNAL_ERROR, arr.Values[0])
			fmt.Println(msg)
			return []byte(msg)
		}

		key := rawKey.Value

		item, ok := h.data.getValue(key)
		if !ok {
			item = resp.RespNull{}
		}

		values[i] = item
	}
	response := resp.RespArray{Values: values}
	fmt.Printf("%s GET Request Successfully completed: %#v\n", INTERNAL_INFO, response)
	return response.Serialize()
}

func (h *handler) handleSet(arr resp.RespArray) []byte {
	if len(arr.Values) != 3 {
		msg := fmt.Sprintf("%s SET Request Array Does Not Consist Of Exactly 2 arguments: %s\n", INTERNAL_ERROR, arr.Values)
		fmt.Println(msg)
		return []byte(msg)
	}

	key, ok := arr.Values[1].(resp.RespString)
	if !ok {
		msg := fmt.Sprintf("%s SET Request Key Is Not Of Type String: %s\n", INTERNAL_ERROR, arr.Values)
		fmt.Println(msg)
		return []byte(msg)
	}

	value := arr.Values[2]
	h.data.setValue(key.Value, value)
	fmt.Printf("%s SET Request Successfully completed. Key: %q, Value: %#v\n", INTERNAL_INFO, key, value)
	return resp.SuccessString().Serialize()
}

func (h *handler) handleAddServer(arr resp.RespArray) []byte {
	if len(arr.Values) != 3 {
		msg := fmt.Sprintf("%s ADDSERVER Request Array Does Not Consist Of Exactly 2 arguments: %s\n", INTERNAL_ERROR, arr.Values)
		fmt.Println(msg)
		return []byte(msg)
	}

	name, ok := arr.Values[1].(resp.RespString)
	if !ok {
		msg := fmt.Sprintf("%s ADDSERVER Request Name Is Not Of Type String: %s\n", INTERNAL_ERROR, arr.Values)
		fmt.Println(msg)
		return []byte(msg)
	}

	address, ok := arr.Values[2].(resp.RespString)
	if !ok {
		msg := fmt.Sprintf("%s ADDSERVER Request Name Is Not Of Type String: %s\n", INTERNAL_ERROR, arr.Values)
		fmt.Println(msg)
		return []byte(msg)
	}

	fmt.Printf("%s Adding %s at address %s\n", DATA_INFO, name.Value, address.Value)

	dialConn, err := h.client.createConn(address.Value)
	if err != nil {
		fmt.Println(err)
		return []byte(err.Error())
	}

	request := resp.RespArray{Values: []resp.RespItem{
		resp.RespString{Value: "PING"},
		resp.RespString{Value: NODE_NAME},
		resp.RespInteger{Value: int64(LISTEN_PORT)},
	}}.Serialize()

	item, err := h.client.sendRequest(request, dialConn)
	if err != nil {
		fmt.Println(err)
		return []byte(err.Error())
	}

	if !resp.IsSuccessString(item) {
		return []byte("Unsuccessfult thing")
	}

	h.cluster.addShardConn(name.Value, dialConn)

	return resp.SuccessString().Serialize()
}

func (h *handler) handlePing(arr resp.RespArray, conn *net.TCPConn) []byte {
	if len(arr.Values) != 3 {
		msg := fmt.Sprintf("%s PING Request Array Does Not Consist Of Exactly 2 arguments: %s\n", INTERNAL_ERROR, arr.Values)
		fmt.Println(msg)
		return []byte(msg)
	}

	name, ok := arr.Values[1].(resp.RespString)
	if !ok {
		msg := fmt.Sprintf("%s PING Request Name Is Not Of Type String: %s\n", INTERNAL_ERROR, arr.Values)
		fmt.Println(msg)
		return []byte(msg)
	}

	port, ok := arr.Values[2].(resp.RespInteger)
	if !ok {
		msg := fmt.Sprintf("%s PING Request Port Is Not Of Type Integer: %s\n", INTERNAL_ERROR, arr.Values)
		fmt.Println(msg)
		return []byte(msg)
	}

	tcpAddr, ok := conn.RemoteAddr().(*net.TCPAddr)
	if !ok {
		return []byte("Could not convert to tcp addr")
	}

	dialConn, err := h.client.createConn(fmt.Sprintf("%s:%d", tcpAddr.IP, port.Value))
	if err != nil {
		fmt.Println(err)
		return []byte(err.Error())
	}

	h.cluster.addShardConn(name.Value, dialConn)

	return resp.SuccessString().Serialize()
}