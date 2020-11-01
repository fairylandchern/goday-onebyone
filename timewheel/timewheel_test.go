package timewheel

import (
	"log"
	"testing"
)

func TestUsingWheelChannel(t *testing.T) {
	ring := initRing(100)
	go ring.Start()
	go ring.Process()
	ring.RegisterTask(Task{key: "test1", process: func() { log.Print("hello,world1") }, tmWait: 20})
	ring.RegisterTask(Task{key: "test2", process: func() { log.Print("hello,world2") }, tmWait: 120})
	ring.DeleteKeyTask("test1")
	select {}
}


//can use in later to split tm precisely
type TimeTp int

const (
	Unknown = iota
	Hour
	Minute
	Sec
	MillSec
	Nano
)

