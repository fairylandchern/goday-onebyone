package gpool

import (
	"runtime"
	"testing"
)

func init() {
	println("using MAXPROC")
	numCPUs := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPUs)
}

func TestNewPool(t *testing.T) {
	p := NewPool(100, 100)
	for i := 0; i < 1000; i++ {
		p.PushJob(func() {
			t.Log("hello,world:", i)
		})
	}
	p.Release()
}
