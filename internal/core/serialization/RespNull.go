package serialization

import (
	"fmt"
)

type RespNull struct{}

func (rNull RespNull) Serialize() []byte {
	msg := fmt.Sprintf("$-1%s", crlf)
	fmt.Printf("null serialized: %q\n", msg)
	return []byte(msg)
}

func (rNull RespNull) deserialize(buf []byte, start int) (RespItem, int, error) {
	return nil, 0, fmt.Errorf("should not try to deserialize null, it is handled as a bulk string of len -1")
}
