package command

import (
	"fmt"
	"net"

	"git.trivaris.org/trivaris/keyvalue/internal/core"
	"git.trivaris.org/trivaris/keyvalue/internal/core/serialization"
)

type RespCommand interface {
	parse([]serialization.RespItem) (RespCommand, error)
	serialize() serialization.RespArray
}

func ParseRaw(raw serialization.RespArray) (RespCommand, error) {
	if len(raw.Values) < 1 {
		return nil, fmt.Errorf("too few arguments")
	}
	commandType, ok := raw.Values[0].(serialization.RespString)
	if !ok {
		return nil, fmt.Errorf("no command keyword found: %s", raw.Values[0])
	}
	return Parse(commandType.Value, raw.Values[1:])
}

func Parse(command string, args []serialization.RespItem) (RespCommand, error) {
	switch command {
	case "GET":
		return RespGet{}.parse(args)
	case "SET":
		return RespSet{}.parse(args)
	case "ADDSERVER":
		return RespAddServer{}.parse(args)
	case "PING":
		return RespPing{}.parse(args)
	default:
		return nil, fmt.Errorf("could not match command")
	}
}

func Execute(c RespCommand, conn *net.TCPConn) (serialization.RespItem, error) {
	request := c.serialize()
	return core.SendRequest(request.Serialize(), conn)
}
