package serialization

import (
	"bytes"
	"fmt"
	"strconv"
)

type RespInteger struct{ Value int64 }

func (rInt RespInteger) Serialize() []byte {
	msg := fmt.Sprintf(":%d%s", rInt.Value, crlf)
	fmt.Printf("int serialized: %q\n", msg)
	return []byte(msg)
}

func (rInt RespInteger) deserialize(buf []byte, start int) (RespItem, int, error) {
	raw, _, found := bytes.Cut(buf[start+1:], crlf)
	if !found {
		return nil, 0, fmt.Errorf("could not partition while deserializing int: %q", string(buf[start:]))
	}

	n, err := strconv.Atoi(string(raw))
	if err != nil {
		return nil, 0, err
	}

	rInt.Value = int64(n)

	fmt.Printf("int deserialized: %d\n", rInt.Value)
	consumed := len(":") + len(raw) + len(crlf)
	return rInt, consumed, nil
}
