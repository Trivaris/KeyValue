package main

import (
	"fmt"
	"git.trivaris.org/trivaris/keyvalue/cmd/internal/resp"
)

const INTERNAL_ERROR = "[INTERNAL ERROR]"

type LogicManager interface {
	handleRequest(int, []byte) []byte
}

type LogicManagerImpl struct {
	dataMan DataManager
}

func (logicMan LogicManagerImpl) handleRequest(n int, buf []byte) []byte {
	fmt.Println("Handling Input...")

	item, _, err := resp.Deserialize(buf[:n], 0)

	if err != nil {
		fmt.Println(err)
		return []byte(err.Error())
	}

	arr, ok := item.(resp.RespArray)
	if !ok {
		msg := fmt.Sprintf("%s No Command Array Found: %s", INTERNAL_ERROR, item.Serialize())
		fmt.Println(msg)
		return []byte(msg)
	}

	command, ok := arr.Values[0].(resp.RespString)
	if !ok {
		msg := fmt.Sprintf("%s No Command Keyword Found: %s", INTERNAL_ERROR, arr.Values[0])
		fmt.Println(msg)
		return []byte(msg)
	}

	switch command.Value {
		case "GET": {
			values := make([]resp.RespItem, len(arr.Values) - 1)

			for i, rawKey := range(arr.Values[1:]) {
				rawKey, ok := rawKey.(resp.RespString)
				if (!ok) {
					msg := fmt.Sprintf("%s GET Request Array Does Not Consist Only Of Keys Of Type String: %s", INTERNAL_ERROR, arr.Values[0])
					fmt.Println(msg)
					return []byte(msg)
				}

				key := rawKey.Value
				values[i] = logicMan.dataMan.getValue(key)
			}
			response := resp.RespArray{Values: values}
			return response.Serialize()
		}
		case "SET": {
			if len(arr.Values) != 3 {
				msg := fmt.Sprintf("%s SET Request Array Does Not Consist Of Exactly 2 arguments: %s", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg)
			}
			
			key, ok := arr.Values[1].(resp.RespString)
			if !ok {
				msg := fmt.Sprintf("%s SET Request Key Is Not Of Type String: %s", INTERNAL_ERROR, arr.Values)
				fmt.Println(msg)
				return []byte(msg)
			}
			
			logicMan.dataMan.setValue(key.Value, arr.Values[2])
			return resp.RespString{Value: "Success!"}.Serialize()
		}
		default: {
			msg := fmt.Sprintf("%s %s Request Not Yet Implemented: %s", INTERNAL_ERROR, command.Value, arr.Values)
			fmt.Println(msg)
			return []byte(msg)
		}
	}
}
