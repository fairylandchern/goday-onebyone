package main

import (
	"fmt"
	"goday-onebyone/simple_server/version2/protocol"
	"log"
	"net"
	"os"
	"os/signal"
)

//var (
//	rwlock = sync.RWMutex{}
//)

func main() {
	conn, err := net.Dial("tcp", ":9999")
	if err != nil {
		log.Fatal("err dial server:", err)
		return
	}
	go readData(conn)
	for i := 0; i < 50; i++ {
		go func(idx int) {
			//rwlock.Lock()
			//defer rwlock.Unlock()
			str := fmt.Sprintf("hello,world:%v", idx)
			count, err := conn.Write(protocol.MarshalData([]byte(str)))
			if err != nil {
				log.Fatal("err write data to server:", err)
				return
			}
			log.Println("write data to server success:", count, " data:", str)
		}(i)
	}
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func readData(conn net.Conn) {
	// read or write conn data here
	var data []byte
	buf := make([]byte, 4096)
	for {
		count, err := conn.Read(buf)
		if err != nil {
			log.Fatal("err read data:", err)
			return
		}
		data = append(data, buf[:count]...)
		// process protocol here,to parse the syntax exactly
		for {
			lenth, err := protocol.UnmarshalData(data)
			if err != nil {
				log.Println("err data not enough:", len(data))
				break
			}
			if lenth == 0 {
				data = data[4:]
				break
			}
			needData := data[4:lenth]
			// can have some extra process function here to understand the syntax exactly
			log.Println("client:understand syntax here", string(needData))
			data = data[lenth:]
		}
	}
}
