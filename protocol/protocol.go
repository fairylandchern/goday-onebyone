package protocol

import (
	"encoding/binary"
	"log"
)

// 自定义协议

// head lenth
const (
	HEADLENGTH = 24
)

// metadata and payload info
type Message struct {
	Head    uint32
	Seq     uint64
	Method  string
	Path    string
	Payload []byte
}

// type Head [12]byte

// marshal send data
// how to modify head in binary.putUint64
// concat data transfer from binary.littleEndian or bigEndian.
// total head length.
func MarshalData(message *Message) []byte {
	head := make([]byte, 24)
	// put data into header
	binary.LittleEndian.PutUint32(head[:3], message.Head)
	binary.LittleEndian.PutUint64(head[4:8+4], message.Seq)
	binary.LittleEndian.PutUint16(head[12:12+2], uint16(len(message.Path)))
	binary.LittleEndian.PutUint16(head[14:14+2], uint16(len(message.Method)))
	binary.LittleEndian.PutUint64(head[14+2:16+8], uint64(len(message.Payload)))
	// append data to head
	data := append(head, message.Path...)
	data = append(data, message.Method...)
	data = append(data, message.Payload...)
	return data
}

// unmarshal receive data
func UnmarshalData(data chan []byte) *Message {
	var bufferData []byte

	for {
		select {
		case recvData := <-data:
			bufferData = append(bufferData, recvData...)
			// unmarshal data
			head := bufferData[:HEADLENGTH]
			var m Message
			m.Head = binary.LittleEndian.Uint32(head[:3])
			m.Seq = binary.LittleEndian.Uint64(head[4 : 8+4])
			pathlen := binary.LittleEndian.Uint16(head[12 : 12+2])
			methodLen := binary.LittleEndian.Uint16(head[14 : 14+2])
			payloadLen := binary.LittleEndian.Uint64(head[14+2 : 16+8])
			oneTotalLength:=uint64(pathlen+methodLen)+payloadLen
			unitData:=bufferData[HEADLENGTH:oneTotalLength+HEADLENGTH]
			bufferData=bufferData[oneTotalLength+HEADLENGTH:]
			m.Path=string(unitData[:pathlen])
			m.Method=string(unitData[pathlen:pathlen+methodLen])
			// send unit message to upper layer


		default:
			log.Println("recv data err")
		}
	}
}

// unpack recv data  length here
//func unpackLength(head []byte) {
//	//var m Message
//	//m.Head=binary.LittleEndian.Uint32(head[:3])
//	//m.Seq=binary.LittleEndian.Uint64(head[4:8+4])
//	//pathlen:=binary.LittleEndian.Uint16(head[12:12+2])
//	//methodLen:=binary.LittleEndian.Uint16(head[14:14+2])
//	//payloadLen:=binary.LittleEndian.Uint64(head[14+2:16+8])
//}
//
//// unpack data and ret data to upper process
//func unpack(data []byte) {
//
//}
