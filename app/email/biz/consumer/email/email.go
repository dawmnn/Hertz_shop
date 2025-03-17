package email

import (
	"context"
	"github.com/cloudwego/biz-demo/gomall/app/email/infra/mq"
	"github.com/cloudwego/biz-demo/gomall/app/email/infra/notify"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/email"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
)

func ConsumerInit() {
	//从 OpenTelemetry 包中获取一个 Tracer 对象。
	//Tracer 是用于创建、管理和记录追踪信息的核心工具

	//"shop-nats-consumer" 是追踪的名称，可以看作是标识符。不同的名称可以用来标识不同的系统、服务或者消费者。
	//在这个例子中，它表示 shop-nats-consumer 这个服务或消费者。
	tracer := otel.Tracer("shop-nats-consumer")
	//表示在 NATS 上订阅一个主题为 "email" 的消息。当有新的消息发送到 "email" 主题时，
	//回调函数 func(msg *nats.Msg) 会被调用，并且 msg 参数包含了收到的消息。
	sub, err := mq.Nc.Subscribe("email", func(msg *nats.Msg) {
		var req email.EmailReq
		err := proto.Unmarshal(msg.Data, &req)
		if err != nil {
			klog.Error(err)
			return
		}

		ctx := context.Background()
		//otel.GetTextMapPropagator() 获取 OpenTelemetry 的默认文本映射传播器。
		//这个传播器用于从外部传入的请求头（如 HTTP Header 或消息头）中提取追踪信息。
		ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(msg.Header))
		//这个方法从 msg.Header 中提取追踪上下文，
		//并将它合并到当前的 ctx 中。这样，你就能够在追踪中保持父级信息，进行追踪传播。
		_, span := tracer.Start(ctx, "shop-email-consumer")

		defer span.End()

		noopEmail := notify.NewNoopEmail()
		_ = noopEmail.Send(&req)
	})
	if err != nil {
		panic(err)
	}
	server.RegisterShutdownHook(func() {
		//server.RegisterShutdownHook 注册了一个关闭钩子。当程序关闭时，会执行传入的函数。
		//这个函数的作用是取消对 "email" 主题的订阅 (sub.Unsubscribe())
		//并且关闭 NATS 客户端连接 (mq.Nc.Close())。
		sub.Unsubscribe()
		mq.Nc.Close()
	})
}
