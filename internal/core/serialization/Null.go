package serialization

import (
	"fmt"
)

type Null struct{}

func (n Null) Serialize() []byte {
	msg := fmt.Sprintf("$-1%s", crlf)
	fmt.Printf("null serialized: %q\n", msg)
	return []byte(msg)
}

func (n Null) deserialize(buf []byte, start int) (Entity, int, error) {
	return nil, 0, fmt.Errorf("should not try to deserialize null, it is handled as a bulk string of len -1")
}
