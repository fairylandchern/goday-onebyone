package main

import (
	"goday-onebyone/simple_server/version1/protocol"
	"log"
	"net"
)

func main() {
	ln,err:=net.Listen("tcp",":9999")
	if err!=nil{
		log.Fatal("err bind port:",err)
		return
	}
	for{
		conn,err:=ln.Accept()
		if err!=nil{
			log.Fatal("err accept conn:",err)
			return
		}
		go func() {
			// read or write conn data here
			var data []byte
			buf:=make([]byte,4096)
			count,err:=conn.Read(buf)
			if err!=nil{
				log.Fatal("err read data:",err)
				return
			}
			data=append(data,buf[:count]...)
			// process protocol here,to parse the syntax exactly
			for{
				lenth,err:=protocol.UnmarshalData(data)
				if err!=nil{
					log.Println("err data not enough:",len(data))
					break
				}
				if lenth==0{
					data=data[4:]
					break
				}
				needData:=data[4:lenth]
				// can have some extra process function here to understand the syntax exactly
				log.Println("understand syntax here",string(needData))
				count,err=conn.Write(protocol.MarshalData(needData))
				if err!=nil{
					log.Println("err write data to client:",err)
				}
				log.Println("write data count to client:",count)
				data=data[lenth:]
			}
		}()
	}
}
