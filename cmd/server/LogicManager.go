package main

import (
	"fmt"
	"git.trivaris.org/trivaris/keyvalue/cmd/internal/resp"
)

const INTERNAL_ERROR = "[INTERNAL ERROR]"
const INTERNAL_INFO = "[INTERNAL INFO]"

type LogicManager interface {
	handleRequest(int, []byte) []byte
}

type LogicManagerImpl struct {
	dataMan DataManager
}

func (logicMan LogicManagerImpl) handleRequest(n int, buf []byte) []byte {
	item, _, err := resp.Deserialize(buf[:n], 0)

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
		case "GET": {
			values := make([]resp.RespItem, len(arr.Values) - 1)

			for i, rawKey := range(arr.Values[1:]) {
				rawKey, ok := rawKey.(resp.RespString)
				if (!ok) {
					msg := fmt.Sprintf("%s GET Request Array Does Not Consist Only Of Keys Of Type String: %s\n", INTERNAL_ERROR, arr.Values[0])
					fmt.Println(msg)
					return []byte(msg)
				}

				key := rawKey.Value
				values[i] = logicMan.dataMan.getValue(key)
			}
			response := resp.RespArray{Values: values}
			fmt.Printf("%s GET Request Successfully completed: %#v\n", INTERNAL_INFO, response)
			return response.Serialize()
		}
		case "SET": {
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
			logicMan.dataMan.setValue(key.Value, value)
			fmt.Printf("%s SET Request Successfully completed. Key: %q, Value: %#v\n", INTERNAL_INFO, key, value)
			return resp.RespString{Value: "Success!"}.Serialize()
		}
		default: {
			msg := fmt.Sprintf("%s %s Request Not Yet Implemented: %s\n", INTERNAL_ERROR, command.Value, arr.Values)
			fmt.Println(msg)
			return []byte(msg)
		}
	}
}
