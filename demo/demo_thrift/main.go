package main

import (
	"net"
	"time"

	"github.com/cloudwego/biz-demo/gomall/demo/demo_thrift/conf"
	"github.com/cloudwego/biz-demo/gomall/demo/demo_thrift/kitex_gen/api/echo"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	opts := kitexInit()
	//调用 kitexInit() 函数： 这一行调用了 kitexInit 函数并将返回的配置选项存储在 opts 变量中。
	//kitexInit 函数的作用是初始化服务器的相关配置和选项，并返回一个包含这些配置的切片。
	svr := echo.NewServer(new(EchoImpl), opts...)
	//echo.NewServer(new(EchoImpl), opts...) 创建了一个新的 Echo 服务器实例。
	//这里 EchoImpl 是实现了服务接口的结构体（假设它实现了某个接口 Echo）。
	//opts... 是从 kitexInit 函数返回的配置选项，将它们传递给 NewServer 来初始化服务器。
	//opts... 代表一个切片的“展开”操作，将其中的每个选项都传递给 NewServer。
	err := svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	//net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address) 解析配置中的服务器地址，
	//conf.GetConf().Kitex.Address 提供了一个地址（如 localhost:8080）。
	//它返回一个 TCPAddr 对象，并且 err 用于捕捉是否解析失败。
	//如果解析失败，调用 panic(err) 抛出错误并终止程序。
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))
	//server.WithServiceAddr(addr) 将解析出的服务地址（addr）作为选项添加到 opts 中。
	//这告诉服务器在哪个地址和端口上进行监听。
	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))
	//server.WithServerBasicInfo 配置服务器的一些基本信息，
	//其中 rpcinfo.EndpointBasicInfo 定义了服务器的基本信息，比如服务名称
	//在这里，conf.GetConf().Kitex.Service 提供了服务的名称。
	//这段代码将服务名称添加到选项中，以便服务器在运行时知道要提供哪个服务。
	// thrift meta handler
	opts = append(opts, server.WithMetaHandler(transmeta.ServerTTHeaderHandler))
	//server.WithMetaHandler(transmeta.ServerTTHeaderHandler) 配置了一个元数据处理器，
	//用于处理服务器的 Thrift 请求头。
	//这个处理器通常用于管理请求的元数据（如 HTTP 头部信息），确保通信协议的一致性。
	// klog
	logger := kitexlogrus.NewLogger()
	//创建了一个新的日志记录器，基于 logrus 库，这通常是一个结构化的日志系统。
	klog.SetLogger(logger)
	//设置日志系统为 logger，使得 klog 使用该日志系统进行日志输出。
	klog.SetLevel(conf.LogLevel())
	//根据配置设置日志的级别，conf.LogLevel() 返回一个日志级别（如 INFO, ERROR 等）。
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
	//zapcore.BufferedWriteSyncer 是一个用于日志异步写入的同步器。
	//zapcore.AddSync 会将日志输出同步到实际的写入设备（例如文件）。
	//lumberjack.Logger 是一个日志轮转工具，用于管理日志文件的大小和备份。
	//例如，MaxSize 控制单个日志文件的最大大小，MaxBackups 控制保留的最大备份数量，
	//MaxAge 控制日志文件的最大存活时间。
	//asyncWriter 配置了一个每分钟刷新一次的异步写入器，避免频繁的 I/O 操作。
	//klog.SetOutput(asyncWriter) 将日志的输出设置为异步写入器，从而将日志写入文件，并在后台异步刷新。
	server.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	//server.RegisterShutdownHook 注册了一个关闭钩子函数，在服务器关闭时调用。
	//在此，关闭时会调用 asyncWriter.Sync()，确保缓冲区的日志数据被及时刷新到磁盘。
	return
}
