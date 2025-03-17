package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/bytedance/gopkg/cloud/metainfo"
	"github.com/cloudwego/biz-demo/gomall/demo/demo_proto/kitex_gen/pbapi"
	"github.com/cloudwego/biz-demo/gomall/demo/demo_proto/middleware"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"

	//"github.com/cloudwego/biz-demo/gomall/demo/demo_thrift/kitex_gen/api/echo"
	"github.com/cloudwego/biz-demo/gomall/demo/demo_proto/kitex_gen/pbapi/echoservice"
	"github.com/cloudwego/kitex/client"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
)

func main() {
	r, err := consul.NewConsulResolver("127.0.0.1:8500")
	//这一行代码初始化了一个 Consul 注册中心解析器，用于在 Consul 上查找服务。
	//它连接到 127.0.0.1:8500 这个 Consul 服务器。
	if err != nil {
		log.Fatal(err)
	}
	c, err := echoservice.NewClient("demo_proto", client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithMiddleware(middleware.Middleware),
	)
	//这里创建了一个新的 Kitex 客户端，指定了服务的名称 demo_proto（这个名字应与服务端提供的名称匹配）。
	//client.WithResolver(r) 表明这个客户端使用之前创建的 Consul 解析器来发现服务。
	if err != nil {
		log.Fatal(err)
	}
	ctx := metainfo.WithPersistentValue(context.Background(), "CLIENT_NAME", "demo_proto_client")
	res, err := c.Echo(ctx, &pbapi.Request{Message: "hello"})
	//这一行代码向服务端发送一个 Echo 请求，包含一个消息 "hello"。context.TODO()
	//表示使用一个空的上下文（用于没有特别要求的情况）。
	//&pbapi.Request{Message: "hello"} 是请求体，包含一个字段 Message，其值为 "hello"。
	var bizErr *kerrors.GRPCBizStatusError
	if err != nil {
		ok := errors.As(err, &bizErr)
		if ok {
			fmt.Printf("%#v", bizErr)
		}
		log.Fatal(err)
	}
	fmt.Printf("%v", res)
}
