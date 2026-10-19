package serialization

import (
	"bytes"
	"fmt"
)

type RespError struct{ Error error }

func (rError RespError) Serialize() []byte {
	msg := fmt.Sprintf("-%s\r\n", rError.Error)
	fmt.Printf("error serialized: %q\n", msg)
	return []byte(msg)
}

func (rError RespError) deserialize(buf []byte, start int) (RespItem, int, error) {
	errorMessageRaw, _, found := bytes.Cut(buf[start+1:], crlf)
	if !found {
		return nil, 0, fmt.Errorf("could not parse error message while deserializing error: %q", string(buf[start:]))
	}

	rError.Error = fmt.Errorf("%s", string(errorMessageRaw))

	fmt.Printf("error deserialized: %q\n", rError.Error)
	consumed := len("-") + len(errorMessageRaw) + len(crlf)
	return rError, consumed, nil
}
