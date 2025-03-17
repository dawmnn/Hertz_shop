package clientsuite

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

type CommonClientSuite struct {
	CurrentServiceName string
	RegistryAddr       string
}

func (s CommonClientSuite) Options() []client.Option {
	opts := []client.Option{
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC),
		//是一个配置选项，用于设置客户端与服务器之间通信所使用的传输协议。
		//transport.GRPC 表示使用 gRPC 协议进行通信。
		//gRPC 是一种高性能的远程过程调用（RPC）协议，通常用于客户端与服务器之间的高效通信。
		client.WithSuite(tracing.NewClientSuite()),
	}
	r, err := consul.NewConsulResolver(s.RegistryAddr)

	if err != nil {
		panic(err)
	}
	opts = append(opts, client.WithResolver(r))
	return opts
}
