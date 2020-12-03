package main

import (
	"context"
	"fmt"
	"github.com/apache/thrift/lib/go/thrift"
	"goday-onebyone/simple_server/use_thrift/gen-go/compute"
	"net"
	"os"
)

func main() {
	// 先建立和服务器的连接的socket，再通过socket建立Transport
	socket, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "9999"))
	if err != nil {
		fmt.Println("Error opening socket:", err)
		os.Exit(1)
	}
	transport := thrift.NewTFramedTransport(socket)

	// 创建二进制协议
	protocol := thrift.NewTBinaryProtocolTransport(transport)
	// 打开Transport，与服务器进行连接
	if err := transport.Open(); err != nil {
		fmt.Fprintln(os.Stderr, "Error opening socket to "+"localhost"+":"+"9999", err)
		os.Exit(1)
	}
	defer transport.Close()

	// 接口需要context，以便在长操作时用户可以取消RPC调用
	ctx := context.Background()

	// 使用divmod服务
	divmodProtocol := thrift.NewTMultiplexedProtocol(protocol, "divmod")
	// 创建代理客户端，使用TMultiplexedProtocol访问对应的服务
	c := thrift.NewTStandardClient(divmodProtocol, divmodProtocol)

	client := compute.NewDivModClient(c)
	res, err := client.DoDivMod(ctx, 100, 3)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(res)

	// 使用mulrange服务
	// 步骤与上面的相同
	mulProtocol := thrift.NewTMultiplexedProtocol(protocol, "mulrange")
	c = thrift.NewTStandardClient(mulProtocol, mulProtocol)
	client2 := compute.NewMulRangeClient(c)
	num, err := client2.BigRange(ctx, 100)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(num)
}
