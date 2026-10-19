package serialization

import (
	"bytes"
	"fmt"
	"strconv"
)

type RespString struct{ Value string }

func (rString RespString) Serialize() []byte {
	msg := fmt.Sprintf("$%d%s%s%s", len(rString.Value), crlf, rString.Value, crlf)
	fmt.Printf("string serialized: %q\n", msg)
	return []byte(msg)
}

func (rString RespString) deserialize(buf []byte, start int) (RespItem, int, error) {
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
		return RespNull{}, 0, nil
	case 0:
		{
			rString.Value = ""
			return rString, 0, nil
		}
	}

	rString.Value = string(rest[:stringLength])

	fmt.Printf("string deserialized: %q\n", rString.Value)
	consumed := len("$") + len(stringLengthRaw) + len(crlf) + stringLength + len(crlf)
	return rString, consumed, nil
}

func IsSuccessString(item RespItem) bool {
	str, ok := item.(RespString)
	return ok && str.Value == "Success!"
}

func SuccessString() RespItem {
	return RespString{Value: "Success!"}
}
