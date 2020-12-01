package protocol

import (
	"encoding/binary"
	"errors"
)

// transport message format

// simple format,head+payload

func UnmarshalData(data []byte) (uint32, error) {
	if len(data) < 4 {
		return 0, errors.New("err not have enough length")
	}
	lenth := binary.BigEndian.Uint32(data)
	if lenth == 0 {
		return 0, errors.New("err read data empty")
	}
	if lenth > uint32(len(data)) {
		return 0, errors.New("err not have enough data")
	}
	return lenth, nil
}

func MarshalData(data []byte) []byte {
	lenth := make([]byte, 4)
	binary.BigEndian.PutUint32(lenth, uint32(len(data))+4)
	data = append(lenth, data...)
	return data
}
