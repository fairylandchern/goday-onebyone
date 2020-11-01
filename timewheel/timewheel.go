package timewheel

import (
	"container/list"
	"log"
	"sync"
	"time"
)

type RingBuffer struct {
	slot []*list.List
	size int64
	//recv channel
	taskRecv   chan Task
	runflag    chan bool
	reverseMap map[string]int64
	taskSend   chan Task
	sync.RWMutex
	count int64
}

type Task struct {
	key     string
	process ProcessFunc
	tmWait  int64
	round   int64
}

type ProcessFunc func()

func initRing(size int64) *RingBuffer {
	ring := new(RingBuffer)
	ring.size = size
	ring.taskRecv = make(chan Task)
	ring.slot = make([]*list.List, size)
	for i := int64(0); i < size; i++ {
		ring.slot[i] = list.New()
	}
	ring.runflag = make(chan bool, 1)
	ring.reverseMap = make(map[string]int64)
	return ring
}

func (r *RingBuffer) putTaskToRing(task Task) {
	part, round := r.partAndRoundTrip(task.tmWait)
	task.round = round
	log.Print("round:", round, "part:", part)
	r.reverseMap[task.key] = part
	if r.slot[part] == nil {
		r.slot[part] = list.New()
	}
	r.slot[part].PushBack(task)
}

func (r *RingBuffer) partAndRoundTrip(count int64) (part, round int64) {
	r.Lock()
	round = (r.count+count) / r.size
	part = (r.count+count) % r.size
	r.Unlock()
	return
}

func (r *RingBuffer) Process() {
	r.count = 0
	for {
		r.Lock()
		r.count %= r.size
		r.ProcessTask(r.count)
		r.count++
		r.Unlock()
		time.Sleep(time.Second)
	}
}

func (r *RingBuffer) ProcessTask(idx int64) {
	l := r.slot[idx]
	lenth := l.Len()
	for i := 0; i < lenth; i++ {
		h := l.Front()
		task := h.Value.(Task)
		l.Remove(h)
		if task.round == 0 {
			task.process()
			continue
		}
		task.round -= 1
		l.PushBack(task)
		h = h.Next()
	}
}

func (r *RingBuffer) Start() {
	log.Print("ring start now")
	for {
		log.Print("ring wait here now")
		select {
		case <-r.runflag:
			log.Printf("recv the signal,exit ring now")
		case data := <-r.taskRecv:
			r.putTaskToRing(data)
		}
	}
}

func (r *RingBuffer) End() {
	r.runflag <- true
}

func (r *RingBuffer) RegisterTask(task Task) {
	r.taskRecv <- task
}

func (r *RingBuffer) DeleteKeyTask(key string) {
	idx:=r.reverseMap[key]
	l:=r.slot[idx]
	lenth:=l.Len()
	h:=l.Front()
	for i:=0;i<lenth;i++{
		t:=h.Value.(Task)
		if t.key==key{
			l.Remove(h)
		}
		h=h.Next()
		if h==nil{
			break
		}
	}
}