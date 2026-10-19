package serialization

import (
	"fmt"
)

var (
	crlf = []byte("\r\n")
)

type RespItem interface {
	Serialize() []byte
	deserialize([]byte, int) (RespItem, int, error)
}

func Deserialize(buf []byte, start int) (RespItem, int, error) {
	if start >= len(buf) {
		return nil, 0, fmt.Errorf("index out of bounds")
	}
	switch buf[start] {
	case '+':
		return nil, 0, fmt.Errorf("deserialization of simple strings is not supported")
	case '-':
		return RespError{}.deserialize(buf, start)
	case ':':
		return RespInteger{}.deserialize(buf, start)
	case '$':
		return RespString{}.deserialize(buf, start)
	case '*':
		return RespArray{}.deserialize(buf, start)
	default:
		return nil, 0, fmt.Errorf("could not determine type of serialized data")
	}
}
