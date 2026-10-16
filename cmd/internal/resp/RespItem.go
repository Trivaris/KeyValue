package resp

import (
	"bytes"
	"fmt"
	"strconv"
)

const (
	SERIALIZATION_ERROR string = "[SERIALIZATION ERROR]"
	DESERIALIZATION_ERROR string = "[DESERIALIZATION ERROR]"
)

var (
	crlf = []byte("\r\n")
	crlfLen = len(crlf)
)

type RespItem interface {
	Serialize() []byte
	deserialize([]byte, int) (RespItem, int, error)
}

type RespNull 		struct { }
type RespInteger 	struct { Value int64 }
type RespString 	struct { Value string }
type RespArray 		struct { Values []RespItem }
type RespError 		struct { ErrorType string; ErrorMessage string }

func (rNull RespNull) Serialize() []byte {
	return []byte("$-1\r\n")
}

func (rNull RespNull) deserialize(buf []byte, start int) (RespItem, int, error) {
	return nil, 0, fmt.Errorf("%s Should Not Try To Deserialize Null\n", DESERIALIZATION_ERROR)
}

func (rInt RespInteger) Serialize() []byte {
	msg := fmt.Sprintf(":%d\r\n", rInt.Value)
	return []byte(msg)
}

func (rInt RespInteger) deserialize(buf []byte, start int) (RespItem, int, error) {
	raw, _, found := bytes.Cut(buf[start + 1:], crlf)
	if !found {
		return nil, 0, fmt.Errorf("%s Could Not Partition While Deserializing Int: %q\n", DESERIALIZATION_ERROR, string(buf[start:]))
	}

	n, err := strconv.Atoi(string(raw))
	if err != nil {
		return nil, 0, err
	}
	
	rInt.Value = int64(n)

	consumed := len(":") + len(raw) + crlfLen
	return rInt, consumed, nil
}

func (rString RespString) Serialize() []byte {
	bytes := []byte(rString.Value)
	bytes = append(fmt.Appendf(nil, "$%d\r\n", len(bytes)),bytes...)
	bytes = append(bytes, []byte(crlf)...)
	return bytes
}

func (rString RespString) deserialize(buf []byte, start int) (RespItem, int, error) {
	stringLengthRaw, rest, found := bytes.Cut(buf[start + 1:], crlf)
	if !found {
		return nil, 0, fmt.Errorf("%s Could Not Partition While Deserializing String: %q\n", DESERIALIZATION_ERROR, string(buf[start:]))
	}

	stringLength, err := strconv.Atoi(string(stringLengthRaw))
	if err != nil {
		return nil, 0, err
	}
	if (stringLength > len(rest)) {
		return nil, 0, fmt.Errorf("%s Parsed Length Did Not Match Actual Length Of Serialized String: %q\n", DESERIALIZATION_ERROR, string(buf[start:stringLength]))
	}

	switch stringLength {
		case -1: return RespNull{}, 0, nil
		case 0: {
			rString.Value = ""
			return rString, 0, nil
		}
	}

	rString.Value = string(rest[:stringLength])
	
	consumed := len("$") + len(stringLengthRaw) + crlfLen + stringLength + crlfLen
	return rString, consumed, nil
}

func (rError RespError) Serialize() []byte {
	msg := fmt.Sprintf("-%s %s\r\n", rError.ErrorType, rError.ErrorMessage)
	return []byte(msg)
}

func (rError RespError) deserialize(buf []byte, start int) (RespItem, int, error) {
	errorTypeRaw, rest, found := bytes.Cut(buf[start + 1:], []byte(" "))
	if !found { return nil, 0, fmt.Errorf("%s Could Not Parse Error Type while deserializing Error: %q\n", DESERIALIZATION_ERROR, string(buf[start:])) }
	
	errorMessageRaw, _, found := bytes.Cut(rest, crlf)
	if !found { return nil, 0, fmt.Errorf("%s Could Not Parse Error Message while deserializing Error: %q\n", DESERIALIZATION_ERROR, string(buf[start:])) }

	rError.ErrorType = string(errorTypeRaw)
	rError.ErrorMessage = string(errorMessageRaw)

	consumed := len(errorTypeRaw) + len(" ") + len(errorMessageRaw)
	return rError, consumed, nil
}

func (rArray RespArray) Serialize() []byte {
	buf := make([]byte, 0)
	buf = append(buf, fmt.Appendf(nil, "*%d\r\n", len(rArray.Values))...)
	for _, item := range rArray.Values {
		buf = append(buf, item.Serialize()...)
	}
	return buf
}

func (rArray RespArray) deserialize(buf []byte, start int) (RespItem, int, error) {
	arrLengthRaw, _, found := bytes.Cut(buf[start + 1:], crlf)
	if !found {
		return nil, 0, fmt.Errorf("%s Could Not Partition While Deserializing Array: %q\n", DESERIALIZATION_ERROR, string(buf[start:]))
	}

	arrLength, err := strconv.Atoi(string(arrLengthRaw))
	if err != nil {
		return nil, 0, err
	}

	values := make([]RespItem, arrLength)
	offset := start + 1 + len(arrLengthRaw) + crlfLen
	for i := range(values) {
		item, consumed, err := Deserialize(buf, offset)
		if err != nil {
			return nil, 0, err
		}
		offset += consumed
		values[i] = item
	}
	
	rArray.Values = values

	consumed := offset - start
	return rArray, consumed, nil
}

func Deserialize(buf []byte, start int) (RespItem, int, error) {
	if start >= len(buf) { return nil, 0, fmt.Errorf("Index out of bounds\n") }
	switch buf[start] {
		case '+': { return nil, 0, fmt.Errorf("%s Deserialization Of Simple Strings Is Not Supported\n", DESERIALIZATION_ERROR) 		}
		case '-': { return RespError{}.deserialize(buf, start) 	}
		case ':': { return RespInteger{}.deserialize(buf, start)	}
		case '$': { return RespString{}.deserialize(buf, start)	}
		case '*': { return RespArray{}.deserialize(buf, start)	}
		default:  { return nil, 0, fmt.Errorf("%s Could Not Determine Type Of Serialized Data \n", DESERIALIZATION_ERROR) }
	}
}
