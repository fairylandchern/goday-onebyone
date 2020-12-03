package main

import (
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"goday-onebyone/simple_server/use_thrift/gen-go/compute"
	"net"
	"os"
)

func main() {
	serverTransport, err := thrift.NewTServerSocket(net.JoinHostPort("127.0.0.1", "9999"))
	if err != nil {
		fmt.Println("Error!", err)
		os.Exit(1)
	}

	// 创建二进制协议
	protocolFactory := thrift.NewTBinaryProtocolFactoryDefault()
	transportFactory := thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())

	// 创建Processor，用一个端口处理多个服务
	divmodProcessor := compute.NewDivModProcessor(new(divmodThrift))
	mulrangeProcessor := compute.NewMulRangeProcessor(new(mulrangeThrift))

	multiProcessor := thrift.NewTMultiplexedProcessor()
	// 给每个service起一个名字
	multiProcessor.RegisterProcessor("divmod", divmodProcessor)
	multiProcessor.RegisterProcessor("mulrange", mulrangeProcessor)

	// 启动服务器
	server := thrift.NewTSimpleServer4(multiProcessor, serverTransport, transportFactory, protocolFactory)
	server.Serve()
	// 退出时停止服务器
	defer server.Stop()
}
