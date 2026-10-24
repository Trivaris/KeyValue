package serialization

import (
	"bytes"
	"fmt"
	"strconv"
)

type Integer struct{ Value int64 }

func (i Integer) Serialize() []byte {
	msg := fmt.Sprintf(":%d%s", i.Value, crlf)
	fmt.Printf("int serialized: %q\n", msg)
	return []byte(msg)
}

func (i Integer) deserialize(buf []byte, start int) (Entity, int, error) {
	raw, _, found := bytes.Cut(buf[start+1:], crlf)
	if !found {
		return nil, 0, fmt.Errorf("could not partition while deserializing int: %q", string(buf[start:]))
	}

	n, err := strconv.Atoi(string(raw))
	if err != nil {
		return nil, 0, err
	}

	i.Value = int64(n)

	fmt.Printf("int deserialized: %d\n", i.Value)
	consumed := len(":") + len(raw) + len(crlf)
	return i, consumed, nil
}
