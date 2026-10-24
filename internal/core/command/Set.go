package command

import (
	"fmt"

	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type Set struct {
	Key   string
	Value serialization.Entity
}

func (s Set) parse(args []serialization.Entity) (Executable, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("too few arguments")
	}
	keyStr, ok := args[0].(serialization.BulkString)
	if !ok {
		return nil, fmt.Errorf("key needs to be string")
	}
	s.Key = keyStr.Value
	s.Value = args[1]
	return s, nil
}

func (s Set) serialize() serialization.Array {
	return serialization.Array{
		Values: []serialization.Entity{
			serialization.BulkString{Value: "SET"},
			serialization.BulkString{Value: s.Key},
			s.Value,
		},
	}
}
