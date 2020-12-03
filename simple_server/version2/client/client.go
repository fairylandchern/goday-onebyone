package main

import (
	"fmt"
	"goday-onebyone/simple_server/version2/protocol"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
)

//var (
//	rwlock = sync.RWMutex{}
//)

type Client struct {
	seqId     int32
	queuelen  int32
	tmout     int32
	conn      net.Conn
	queueChan chan []byte
	lock      sync.RWMutex
	addr      string
}

func main() {
	addr := "localhost:9999"
	cli := NewClient(addr, 30)
	err := cli.StartServer()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 50; i++ {
		cli.SendData([]byte(fmt.Sprintf("helloworld:%v", i)))
	}
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func NewClient(addr string, tmout int32) *Client {
	return &Client{addr: addr, tmout: tmout, queuelen: 100, queueChan: make(chan []byte, 100)}
}

func (c *Client) StartServer() error {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		log.Println("err dial server:", err)
		return err
	}
	c.conn = conn
	go c.recv()
	go c.send()
	return nil
}

func (c *Client) SendData(data []byte) {
	c.queueChan <- data
}

func (c *Client) send() {
	for {
		data, ok := <-c.queueChan
		if ok {
			count, err := c.conn.Write(protocol.MarshalData(data))
			if err != nil {
				log.Fatal("err write data to server:", err)
				return
			}
			log.Println("write data to server success:", count, " data:", data)
		}
	}
}

func (c *Client) recv() {
	var data []byte
	buf := make([]byte, 4096)
	for {
		count, err := c.conn.Read(buf)
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
