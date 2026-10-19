package command

import (
	"fmt"

	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type RespGet struct {
	Keys []string
}

func (g RespGet) parse(args []serialization.RespItem) (RespCommand, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("too few arguments")
	}
	values := make([]string, len(args))
	for i, item := range args {
		str, ok := item.(serialization.RespString)
		if !ok {
			return nil, fmt.Errorf("keys may only exist as strings")
		}
		values[i] = str.Value
	}
	g.Keys = values
	return g, nil
}

func (g RespGet) serialize() serialization.RespArray {
	values := make([]serialization.RespItem, 1+len(g.Keys))
	values[0] = serialization.RespString{Value: "GET"}

	for i, key := range g.Keys {
		values[i+1] = serialization.RespString{Value: key}
	}

	return serialization.RespArray{Values: values}
}
