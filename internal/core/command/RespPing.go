package command

import (
	"fmt"

	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type RespPing struct {
	Name string
	Port uint16
}

func (p RespPing) parse(args []serialization.RespItem) (RespCommand, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("too few arguments")
	}
	nameStr, ok := args[0].(serialization.RespString)
	if !ok {
		return nil, fmt.Errorf("name needs to be a string")
	}
	port, ok := args[1].(serialization.RespInteger)
	if !ok {
		return nil, fmt.Errorf("port needs to be an int")
	}
	p.Name = nameStr.Value
	p.Port = uint16(port.Value)

	return p, nil
}

func (p RespPing) serialize() serialization.RespArray {
	return serialization.RespArray{
		Values: []serialization.RespItem{
			serialization.RespString{Value: "PING"},
			serialization.RespString{Value: p.Name},
			serialization.RespInteger{Value: int64(p.Port)},
		},
	}
}
