package serialization

import (
	"bytes"
	"fmt"
	"strconv"
)

type RespArray struct{ Values []RespItem }

func (rArray RespArray) Serialize() []byte {
	buf := make([]byte, 0)
	buf = append(buf, fmt.Appendf(nil, "*%d\r\n", len(rArray.Values))...)
	for _, item := range rArray.Values {
		buf = append(buf, item.Serialize()...)
	}
	fmt.Printf("array serialized\n")
	return buf
}

func (rArray RespArray) deserialize(buf []byte, start int) (RespItem, int, error) {
	arrLengthRaw, _, found := bytes.Cut(buf[start+1:], crlf)
	if !found {
		return nil, 0, fmt.Errorf("could not partition while deserializing array: %q", string(buf[start:]))
	}

	arrLength, err := strconv.Atoi(string(arrLengthRaw))
	if err != nil {
		return nil, 0, err
	}

	values := make([]RespItem, arrLength)
	offset := start + 1 + len(arrLengthRaw) + len(crlf)
	for i := range values {
		item, consumed, err := Deserialize(buf, offset)
		if err != nil {
			return nil, 0, err
		}
		offset += consumed
		values[i] = item
	}

	rArray.Values = values

	fmt.Printf("array deserialized\n")
	consumed := offset - start
	return rArray, consumed, nil
}
