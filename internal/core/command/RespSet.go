package command

import (
	"fmt"

	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type RespSet struct {
	Key   string
	Value serialization.RespItem
}

func (s RespSet) parse(args []serialization.RespItem) (RespCommand, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("too few arguments")
	}
	keyStr, ok := args[0].(serialization.RespString)
	if !ok {
		return nil, fmt.Errorf("key needs to be string")
	}
	s.Key = keyStr.Value
	s.Value = args[1]
	return s, nil
}

func (s RespSet) serialize() serialization.RespArray {
	return serialization.RespArray{
		Values: []serialization.RespItem{
			serialization.RespString{Value: "SET"},
			serialization.RespString{Value: s.Key},
			s.Value,
		},
	}
}
