package protocol

import (
	"reflect"
	"testing"
)

func TestArrLen(t *testing.T) {
	a := "abcdefg"
	a = a[len(a):]
	t.Log("a:", reflect.TypeOf(a))
}

func TestUnmarshalData(t *testing.T) {
	str := "hello,world"
	data := MarshalData([]byte(str))
	count, err := UnmarshalData(data)
	if err != nil {
		t.Fatal("err unmarshal:", err)
		return
	}
	readData := data[4:count]
	t.Log("data count:", count, " data:", string(data), " read data:", string(readData))
}
