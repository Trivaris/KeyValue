package serialization

import (
	"bytes"
	"fmt"
)

type Error struct{ Error error }

func (e Error) Serialize() []byte {
	msg := fmt.Sprintf("-%s\r\n", e.Error)
	fmt.Printf("error serialized: %q\n", msg)
	return []byte(msg)
}

func (e Error) deserialize(buf []byte, start int) (Entity, int, error) {
	errorMessageRaw, _, found := bytes.Cut(buf[start+1:], crlf)
	if !found {
		return nil, 0, fmt.Errorf("could not parse error message while deserializing error: %q", string(buf[start:]))
	}

	e.Error = fmt.Errorf("%s", string(errorMessageRaw))

	fmt.Printf("error deserialized: %q\n", e.Error)
	consumed := len("-") + len(errorMessageRaw) + len(crlf)
	return e, consumed, nil
}
