package main

import (
	"fmt"
	"goday-onebyone/simple_server/version1/protocol"
	"log"
	"net"
)

func main() {
	conn,err:=net.Dial("tcp",":9999")
	if err!=nil{
		log.Fatal("err dial server:",err)
		return
	}
	for i:=0;i<1000;i++{
		go func(idx int) {
			str:=fmt.Sprintf("hello,world:%v",idx)
			count,err:=conn.Write(protocol.MarshalData([]byte(str)))
			if err!=nil{
				log.Fatal("err write data to server:",err)
				return
			}
			log.Println("write data to server success:",count)
			data:=make([]byte,1024)
			count,err=conn.Read(data)
			if err!=nil{
				log.Println("err read data:",idx)
			}
			log.Println("read data count:",count," data:",string(data[4:count]))
		}(i)
	}
	select {

	}
}

