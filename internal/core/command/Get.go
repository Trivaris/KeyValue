package command

import (
	"fmt"

	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type Get struct {
	Keys []string
}

func (g Get) parse(args []serialization.Entity) (Executable, error) {
	if len(args) < 1 {
		return nil, fmt.Errorf("too few arguments")
	}
	values := make([]string, len(args))
	for i, item := range args {
		str, ok := item.(serialization.BulkString)
		if !ok {
			return nil, fmt.Errorf("keys may only exist as strings")
		}
		values[i] = str.Value
	}
	g.Keys = values
	return g, nil
}

func (g Get) serialize() serialization.Array {
	values := make([]serialization.Entity, 1+len(g.Keys))
	values[0] = serialization.BulkString{Value: "GET"}

	for i, key := range g.Keys {
		values[i+1] = serialization.BulkString{Value: key}
	}

	return serialization.Array{Values: values}
}
