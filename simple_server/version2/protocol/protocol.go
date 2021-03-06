package protocol

import (
	"encoding/binary"
	"errors"
	"log"
)

// transport message format

// simple format,head+payload

func UnmarshalData(data []byte) (uint32, uint32, error) {
	if len(data) < 4 {
		return 0, 0, errors.New("err not have enough length")
	}
	lenth := binary.BigEndian.Uint32(data[0:4])
	log.Println("lenth:", lenth, " data len:", len(data))
	if lenth == 0 {
		return 0, 0, errors.New("err read data empty")
	}
	if lenth > uint32(len(data)) {
		return 0, 0, errors.New("err not have enough data")
	}
	seqID := binary.BigEndian.Uint32(data[4:8])
	return lenth, seqID, nil
}

func MarshalData(data []byte) []byte {
	lenth := make([]byte, 4)
	binary.BigEndian.PutUint32(lenth, uint32(len(data))+4)
	data = append(lenth, data...)
	return data
}
