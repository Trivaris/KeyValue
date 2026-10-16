package main

import(
	"git.trivaris.org/trivaris/keyvalue/cmd/internal/resp"
)

func main() {
	dataMan 	:= DataManagerImpl{
		values: map[string]resp.RespItem{
			"foo":resp.RespString{Value: "bar"},
			"greeting":resp.RespString{Value: "Hello world!"},
		},
	}
	logicMan 	:= LogicManagerImpl{ dataMan: &dataMan }
	netManager 	:= NetManagerImpl{ logicMan: logicMan }
	netManager.Listen()
}
