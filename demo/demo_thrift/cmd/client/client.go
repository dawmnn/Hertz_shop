package main

import (
	"context"
	"fmt"
	"github.com/cloudwego/biz-demo/gomall/demo/demo_thrift/kitex_gen/api"
	"github.com/cloudwego/biz-demo/gomall/demo/demo_thrift/kitex_gen/api/echo"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
)

func main() {
	cli, err := echo.NewClient("demo_thrift",
		client.WithHostPorts("localhost:8888"),
		client.WithMetaHandler(transmeta.ClientTTHeaderHandler),
		//这行配置设置了一个元数据处理器 MetaHandler，transmeta.ClientTTHeaderHandler
		//是一个处理客户端请求时与传输协议相关的元数据的函数或对象。
		client.WithTransportProtocol(transport.TTHeader),
		//这行代码设置了客户端使用的传输协议，具体是 transport.TTHeader。
		//传输协议 (transport.TTHeader) 是一种在客户端与服务器之间交换数据时使用的协议
		//，它可能用于指定底层的消息格式、编码方式等。
		//WithTransportProtocol 配置选项允许你选择不同的协议，
		//可能是为了适配不同的通信模式或优化。
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: "demo_thrift_client",
		}),
		//这行代码配置客户端的基础信息，使用 rpcinfo.EndpointBasicInfo 来提供有关客户端的信息。
		//ServiceName: "demo_thrift_client"：ServiceName 表示客户端的服务名称，
		//"demo_thrift_client" 是该客户端的服务标识。
		//EndpointBasicInfo 用于存储与客户端相关的基本信息，
		//如服务名称、服务版本、客户端 IP 等，用于在 RPC 请求中传递客户端的基本身份信息。
	)
	//NewClient 函数会根据传入的参数初始化客户端，"demo_thrift" 是服务器端的服务名称，
	//client.WithHostPorts("localhost:8888") 表示将客户端连接到本地的 8888 端口。
	if err != nil {
		panic(err)
	}
	res, err := cli.Echo(context.Background(), &api.Request{
		Message: "hello",
	})
	//context.Background() 来传递请求上下文。context 在 Go 中用于管理请求的生命周期，
	//通常用于取消请求、设置超时、传递元数据等。Background 是最基础的上下文，适用于没有父上下文的情况。
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%v", res)
}
