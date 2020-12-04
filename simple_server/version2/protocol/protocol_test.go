package protocol

import (
	"reflect"
	"sort"
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
	count, _, err := UnmarshalData(data)
	if err != nil {
		t.Fatal("err unmarshal:", err)
		return
	}
	readData := data[4:count]
	t.Log("data count:", count, " data:", string(data), " read data:", string(readData))
}

// 测试如何使用sort.search方法
func TestSortSearch(t *testing.T) {
	x := 2
	arr := []int{1, 2, 3, 5, 8, 10, 20, 30}
	i := sort.Search(len(arr), func(i int) bool {
		return arr[i] >= x
	})
	t.Log("index of x:", i)
}
