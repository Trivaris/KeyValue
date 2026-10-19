package command

import (
	"fmt"
	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type RespAddServer struct {
	Name    string
	Address string
}

func (as RespAddServer) parse(args []serialization.RespItem) (RespCommand, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("too few arguments")
	}
	nameStr, ok := args[0].(serialization.RespString)
	if !ok {
		return nil, fmt.Errorf("name needs to be string")
	}
	addressStr, ok := args[1].(serialization.RespString)
	if !ok {
		return nil, fmt.Errorf("address needs to be string")
	}
	as.Name = nameStr.Value
	as.Address = addressStr.Value
	return as, nil
}

func (as RespAddServer) serialize() serialization.RespArray {
	return serialization.RespArray{
		Values: []serialization.RespItem{
			serialization.RespString{Value: "ADDSERVER"},
			serialization.RespString{Value: as.Name},
			serialization.RespString{Value: as.Address},
		},
	}
}
