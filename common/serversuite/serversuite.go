package serversuite

import (
	"github.com/cloudwego/biz-demo/gomall/common/mtl"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	prometheus "github.com/kitex-contrib/monitor-prometheus"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	consul "github.com/kitex-contrib/registry-consul"
)

type CommonServiceSuite struct {
	CurrentServiceName string
	RegistryAddr       string
}

func (s CommonServiceSuite) Options() []server.Option {
	opts := []server.Option{
		server.WithMetaHandler(transmeta.ServerHTTP2Handler),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
		server.WithTracer(prometheus.NewServerTracer("",
			"",
			prometheus.WithDisableServer(true),
			//表示禁用 Prometheus 服务器端的追踪功能，通常是指不采集服务器端的追踪数据
			prometheus.WithRegistry(mtl.Register)),
		),
		server.WithSuite(tracing.NewServerSuite()),
	}
	r, err := consul.NewConsulRegister(s.RegistryAddr)
	if err != nil {
		klog.Fatal(err)
	}
	opts = append(opts, server.WithRegistry(r))
	return opts
}
