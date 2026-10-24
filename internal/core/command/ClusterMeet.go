package command

import (
	"fmt"

	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type ClusterMeet struct {
	Name    string
	Address string
}

func (cm ClusterMeet) parse(args []serialization.Entity) (Executable, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("too few arguments")
	}
	nameStr, ok := args[0].(serialization.BulkString)
	if !ok {
		return nil, fmt.Errorf("name needs to be string")
	}
	addressStr, ok := args[1].(serialization.BulkString)
	if !ok {
		return nil, fmt.Errorf("address needs to be string")
	}
	cm.Name = nameStr.Value
	cm.Address = addressStr.Value
	return cm, nil
}

func (cm ClusterMeet) serialize() serialization.Array {
	return serialization.Array{
		Values: []serialization.Entity{
			serialization.BulkString{Value: "CLUSTER"},
			serialization.BulkString{Value: "MEET"},
			serialization.BulkString{Value: cm.Name},
			serialization.BulkString{Value: cm.Address},
		},
	}
}
