package main

import (
	"context"
	"encoding/binary"
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
type RespCh chan []byte

type Client struct {
	seqId     int32
	queuelen  int32
	tmout     int32
	conn      net.Conn
	queueChan chan []byte
	lock      sync.RWMutex
	addr      string
	respMap   map[int32]RespCh
}

type ClientProtoco interface {
	Recv([]byte)
}

func main() {
	addr := "localhost:9999"
	cli := NewClient(addr, 30)
	err := cli.StartServer()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 50; i++ {
		err = cli.Call(context.Background(), "test", []byte(fmt.Sprintf("helloworld:%v", i)), nil)
		if err != nil {
			log.Println("err call server method:", err, " idx:", i)
		}
	}
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, os.Kill)
	<-sig
}

func NewClient(addr string, tmout int32) *Client {
	return &Client{addr: addr, tmout: tmout, queuelen: 100, queueChan: make(chan []byte, 100), respMap: make(map[int32]RespCh)}
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
			lenth, seqID, err := protocol.UnmarshalData(data)
			if err != nil {
				log.Println("err data not enough:", len(data))
				break
			}
			if lenth == 0 {
				data = data[4:]
				break
			}
			c.respMap[int32(seqID)] <- data[8:lenth]
			needData := data[4:lenth]
			// can have some extra process function here to understand the syntax exactly
			log.Println("client:understand syntax here", string(needData), " seqID:", seqID)
			data = data[lenth:]
		}
	}
}

// 两步走：发出去和收回来(异步收回来的话，方案：收到数据调用recv方法，数据赋值及管道相关操作放在该方法中，取消息的地方，根据管道数据和需要的数据进行校验即可）
// client抽象出来一个方法，专供上层方法调用，可以验证一个猜想异步回调相关
// 应该是需要用到消息的序列号和方法名了，先考虑只用到序列号的情况下应该如何处理
func (c *Client) Call(ctx context.Context, method string, data []byte, resp interface{}) error {
	c.seqId++
	c.respMap[c.seqId] = make(RespCh)
	// 讲seqid注入到msg中，并重新调整协议
	seqLen := make([]byte, 4)
	binary.BigEndian.PutUint32(seqLen, uint32(c.seqId))
	c.SendData(append(seqLen, data...))
	// 异步读取返回值
	return c.GetMsg(ctx, c.seqId, method, resp)
}

func (c *Client) GetMsg(ctx context.Context, seqId int32, method string, resp interface{}) error {
	data := <-c.respMap[seqId]
	log.Printf("resp seqId:%v,data:%v", seqId, string(data))
	return nil
}

func (c *Client) SendFullData(ctx context.Context, seqId int32, method string, data []byte) {

}
