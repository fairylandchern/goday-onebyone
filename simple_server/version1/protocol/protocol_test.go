package protocol

import (
	"reflect"
	"testing"
)

func TestArrLen(t *testing.T) {
	a:="abcdefg"
	a=a[len(a):]
	t.Log("a:",reflect.TypeOf(a))
}
