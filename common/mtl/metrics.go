package mtl

import (
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net"
	"net/http"
)

var Register *prometheus.Registry

func InitMetric(serviceName, metricsPort, registryAddr string) (registry.Registry, *registry.Info) {
	Register = prometheus.NewRegistry()
	//创建了一个新的 Prometheus 注册表实例。注册表用于管理所有的指标。
	//通过这个注册表，可以注册和查询所有的监控指标。
	Register.MustRegister(collectors.NewGoCollector())
	//将一个新的 GoCollector 注册到 Prometheus 注册表中。GoCollector 会收集 Go 程序的运行时信息，
	//例如 Go 内存使用、垃圾回收等信息。
	//通过注册这个收集器，Prometheus 将能够获取 Go 运行时的监控指标。
	Register.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	//册了 ProcessCollector，它会收集当前进程（程序）的相关信息，比如 CPU、内存使用等。
	//ProcessCollectorOpts{} 是 ProcessCollector 的配置选项，这里使用的是默认配置。
	r, _ := consul.NewConsulRegister(registryAddr)
	addr, _ := net.ResolveTCPAddr("tcp", metricsPort)
	//将 metricsPort 转换为一个 TCP 地址。metricsPort 是 Prometheus 指标的暴露端口，
	//表示 Prometheus 从该端口获取指标数据。
	registryInfo := &registry.Info{
		ServiceName: "prometheus",
		Addr:        addr,
		Weight:      1,
		//将 metricsPort 转换为一个 TCP 地址。metricsPort 是 Prometheus 指标的暴露端口，
		//表示 Prometheus 从该端口获取指标数据。
		Tags: map[string]string{"service": serviceName},
		//标签，用于给该服务打上标记。这里用 serviceName 标记该服务的名称。
	}

	_ = r.Register(registryInfo)
	server.RegisterShutdownHook(func() {
		r.Deregister(registryInfo)
	})
	//钩子会取消服务在 Consul 上的注册

	http.Handle("/metrics", promhttp.HandlerFor(Register, promhttp.HandlerOpts{}))
	//返回 Prometheus 指标数据。promhttp.HandlerFor(Register, promhttp.HandlerOpts{}) 会根据之前注册的
	//Register（Prometheus 注册表） 生成一个 HTTP handler，供 Prometheus 拉取指标。
	go http.ListenAndServe(metricsPort, nil)
	//启动 HTTP 服务，监听 metricsPort 端口。这里使用了 go 关键字来启动一个新的 goroutine，
	//使得这个 HTTP 服务器可以在后台运行，不会阻塞主线程。
	return r, registryInfo
}
