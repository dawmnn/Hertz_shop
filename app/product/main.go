package main

import (
	"context"
	"github.com/cloudwego/biz-demo/gomall/app/product/biz/dal"
	"github.com/cloudwego/biz-demo/gomall/common/mtl"
	"github.com/cloudwego/biz-demo/gomall/common/serversuite"
	"github.com/joho/godotenv"
	"net"
	"time"

	"github.com/cloudwego/biz-demo/gomall/app/product/conf"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	//ServiceName  = conf.GetConf().Kitex.Service
	ServiceName = "product"
	//RegistryAddr = conf.GetConf().Registry.RegistryAddress[0]
	RegistryAddr = "127.0.0.1:8500"
)

func main() {
	_ = godotenv.Load()
	mtl.InitMetric(ServiceName, conf.GetConf().Kitex.MetricsPort, RegistryAddr)
	p := mtl.InitTracing(ServiceName)
	defer p.Shutdown(context.Background())
	dal.Init()
	opts := kitexInit()

	svr := productcatalogservice.NewServer(new(ProductCatalogServiceImpl), opts...)

	err := svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		klog.Error(err)
	}
	opts = append(opts, server.WithServiceAddr(addr), server.WithSuite(serversuite.CommonServiceSuite{
		CurrentServiceName: ServiceName,
		RegistryAddr:       RegistryAddr,
	}))

	//r, err := consul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0])
	//if err != nil {
	//	klog.Fatal(err)
	//}
	//// service info
	//opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
	//	ServiceName: conf.GetConf().Kitex.Service,
	//}), server.WithRegistry(r))

	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	server.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	return
}
