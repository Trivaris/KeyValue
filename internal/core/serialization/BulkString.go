package serialization

import (
	"bytes"
	"fmt"
	"strconv"
)

type BulkString struct{ Value string }

func (b BulkString) Serialize() []byte {
	msg := fmt.Sprintf("$%d%s%s%s", len(b.Value), crlf, b.Value, crlf)
	fmt.Printf("string serialized: %q\n", msg)
	return []byte(msg)
}

func (b BulkString) deserialize(buf []byte, start int) (Entity, int, error) {
	stringLengthRaw, rest, found := bytes.Cut(buf[start+1:], crlf)
	if !found {
		return nil, 0, fmt.Errorf("could not partition while deserializing string: %q", string(buf[start:]))
	}

	stringLength, err := strconv.Atoi(string(stringLengthRaw))
	if err != nil {
		return nil, 0, err
	}
	if stringLength > len(rest) {
		return nil, 0, fmt.Errorf("parsed length did not match actual length of serialized string: %q", string(buf[start:stringLength]))
	}

	switch stringLength {
	case -1:
		return Null{}, 0, nil
	case 0:
		{
			b.Value = ""
			return b, 0, nil
		}
	}

	b.Value = string(rest[:stringLength])

	fmt.Printf("string deserialized: %q\n", b.Value)
	consumed := len("$") + len(stringLengthRaw) + len(crlf) + stringLength + len(crlf)
	return b, consumed, nil
}

func IsSuccessString(item Entity) bool {
	str, ok := item.(BulkString)
	return ok && str.Value == "Success!"
}

func SuccessString() Entity {
	return BulkString{Value: "Success!"}
}
