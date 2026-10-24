package serialization

import (
	"fmt"
)

var (
	crlf = []byte("\r\n")
)

type Entity interface {
	Serialize() []byte
	deserialize([]byte, int) (Entity, int, error)
}

func Deserialize(buf []byte, start int) (Entity, int, error) {
	if start >= len(buf) {
		return nil, 0, fmt.Errorf("index out of bounds")
	}
	switch buf[start] {
	case '+':
		return nil, 0, fmt.Errorf("deserialization of simple strings is not supported")
	case '-':
		return Error{}.deserialize(buf, start)
	case ':':
		return Integer{}.deserialize(buf, start)
	case '$':
		return BulkString{}.deserialize(buf, start)
	case '*':
		return Array{}.deserialize(buf, start)
	default:
		return nil, 0, fmt.Errorf("could not determine type of serialized data")
	}
}
