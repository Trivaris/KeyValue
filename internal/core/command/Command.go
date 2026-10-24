package command

import (
	"fmt"
	"net"

	"git.trivaris.org/trivaris/keyvalue/internal/core"
	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type Executable interface {
	parse([]serialization.Entity) (Executable, error)
	serialize() serialization.Array
}

func ParseRaw(raw serialization.Array) (Executable, error) {
	if len(raw.Values) < 1 {
		return nil, fmt.Errorf("too few arguments")
	}
	commandType, ok := raw.Values[0].(serialization.BulkString)
	if !ok {
		return nil, fmt.Errorf("no command keyword found: %s", raw.Values[0])
	}
	return Parse(commandType.Value, raw.Values[1:])
}

func Parse(command string, args []serialization.Entity) (Executable, error) {
	switch command {
	case "GET":
		return Get{}.parse(args)
	case "SET":
		return Set{}.parse(args)
	case "CLUSTER":
		{
			if len(args) <= 0 {
				return nil, fmt.Errorf("no subcommand provided for CLUSTER command")
			}

			subcmd, ok := args[0].(serialization.BulkString)

			if !ok {
				return nil, fmt.Errorf("could not match command")
			}

			switch subcmd.Value {
			case "MEET":
				return ClusterMeet{}.parse(args[1:])
			default:
				return nil, fmt.Errorf("could not match command")
			}
		}
	default:
		return nil, fmt.Errorf("could not match command")
	}
}

func Execute(c Executable, conn *net.TCPConn) (serialization.Entity, error) {
	request := c.serialize()
	return core.SendRequest(request.Serialize(), conn)
}
