package protocol

import (
	"encoding/binary"
	"testing"
)

func TestMarshal(t *testing.T) {

}

func TestUnmarshal(t *testing.T) {

}


func TestBinaryLittle(t *testing.T) {
	a:=make([]byte,24)
	b:=uint32(1234456)
	binary.LittleEndian.PutUint32(a[:4],b)
	t.Log("endian uint32:",a[:4])
	t.Log("unmarshal uint32:",binary.LittleEndian.Uint32(a[:4]))
	// slice use
	c:="0123456789"
	t.Log("slice:",c[:4]," len:",len(c[:4])," after 4:",c[4:])
}