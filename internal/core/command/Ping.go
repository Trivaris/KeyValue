package command

import (
	"fmt"

	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type Ping struct {
	Name string
	Port uint16
}

func (p Ping) parse(args []serialization.Entity) (Executable, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("too few arguments")
	}
	nameStr, ok := args[0].(serialization.BulkString)
	if !ok {
		return nil, fmt.Errorf("name needs to be a string")
	}
	port, ok := args[1].(serialization.Integer)
	if !ok {
		return nil, fmt.Errorf("port needs to be an int")
	}
	p.Name = nameStr.Value
	p.Port = uint16(port.Value)

	return p, nil
}

func (p Ping) serialize() serialization.Array {
	return serialization.Array{
		Values: []serialization.Entity{
			serialization.BulkString{Value: "PING"},
			serialization.BulkString{Value: p.Name},
			serialization.Integer{Value: int64(p.Port)},
		},
	}
}
