package rpc

import (
	"context"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/conf"
	frontendUtils "github.com/cloudwego/biz-demo/gomall/app/frontend/utils"
	"github.com/cloudwego/biz-demo/gomall/common/clientsuite"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/checkout/checkoutservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/user/userservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	consulclient "github.com/kitex-contrib/config-consul/client"
	"github.com/kitex-contrib/config-consul/consul"
	"sync"
)

var (
	UserClient     userservice.Client
	ProductClient  productcatalogservice.Client
	CartClient     cartservice.Client
	OrderClient    orderservice.Client
	CheckoutClient checkoutservice.Client
	once           sync.Once

	ServiceName  = frontendUtils.ServiceName
	RegistryAddr = conf.GetConf().Hertz.RegistryAddr
	err          error
)

func InitClient() {
	//initUserClient()
	//initProductClient()
	once.Do(func() {
		initUserClient()
		initProductClient()
		initCartClient()
		initOrderClient()
		initCheckoutClient()
	})
}

func initUserClient() {
	//var opts []client.Option
	//r, err := consul.NewConsulResolver(conf.GetConf().Hertz.RegistryAddr)
	//
	//frontendUtils.MustHandleError(err)
	//opts = append(opts, client.WithResolver(r))
	UserClient, err = userservice.NewClient("user", client.WithSuite(clientsuite.CommonClientSuite{
		CurrentServiceName: ServiceName,
		RegistryAddr:       RegistryAddr,
	}))

	frontendUtils.MustHandleError(err)
}

func initProductClient() {
	//var opts []client.Option
	//r, err := consul.NewConsulResolver(conf.GetConf().Hertz.RegistryAddr)
	//
	//frontendUtils.MustHandleError(err)
	//opts = append(opts, client.WithResolver(r))

	//断路器的目的是避免在服务不可用时反复请求失败的服务，从而避免系统的雪崩效应。
	consulClient, err := consul.NewClient(consul.Options{})
	cbs := circuitbreak.NewCBSuite(func(ri rpcinfo.RPCInfo) string {
		return circuitbreak.RPCInfo2Key(ri)
	})
	//NewCBSuite 是一个创建新的 CBSuite（断路器套件）的方法。
	//该方法接受一个回调函数 func(ri rpcinfo.RPCInfo) string，
	//用于生成一个唯一的键，通常这个键与某个 RPC（远程过程调用）请求相关联。
	//
	//rpcinfo.RPCInfo 是一个结构体，通常包含关于当前 RPC 调用的信息，
	//比如目标服务名、请求 ID、调用的具体方法等。
	//
	//circuitbreak.RPCInfo2Key(ri) 是一个将 RPCInfo 转换为一个字符串键的工具方法。
	//这个键用于标识服务实例，帮助断路器区分不同的服务或方法，进行独立的监控和管理。
	cbs.UpdateServiceCBConfig("frontend/product/GetProduct",
		circuitbreak.CBConfig{Enable: true, // 启用断路器
			ErrRate:   0.5, // 错误率阈值，表示当错误率超过50%时，断路器触发
			MinSample: 2,   // 最小样本量，表示至少需要有2个请求进行统计，才会触发断路器判断
		},
	)
	//pdateServiceCBConfig 方法用于更新某个服务的断路器配置。在这里，它接收两个参数：
	//第一个参数 "frontend/product/GetProduct" 是服务的名称或路径，
	//表示这段代码是针对 frontend/product/GetProduct 这个方法的断路器配置。
	//通常，这个值会对应到某个具体的服务接口。
	//第二个参数是 circuitbreak.CBConfig 类型的配置对象，用来设置断路器的各种参数。
	ProductClient, err = productcatalogservice.NewClient("product", client.WithSuite(clientsuite.CommonClientSuite{
		CurrentServiceName: ServiceName,
		RegistryAddr:       RegistryAddr, //配置了后备策略（Fallback Policy）。如果请求失败，后备策略决定客户端会如何应对。
	}), client.WithCircuitBreaker(cbs), client.WithFallback(
		fallback.NewFallbackPolicy(
			fallback.UnwrapHelper(
				func(ctx context.Context, req, resp interface{}, err error) (fbResp interface{}, fbErr error) {
					if err == nil {
						return resp, nil
					}
					methodName := rpcinfo.GetRPCInfo(ctx).To().Method()
					//获取当前请求的 methodName，也就是 RPC 调用的方法名。
					if methodName != "ListProducts" {
						return resp, err
					}
					return &product.ListProductsResp{
						Products: []*product.Product{
							{
								Price:       6.6,
								Id:          3,
								Picture:     "/static/image/t-shirt-1.jpeg",
								Name:        "T-Shirt",
								Description: "CloudWeGo T-Shirt",
							},
						},
					}, nil
				}),
		),
	),
		client.WithSuite(consulclient.NewSuite("product", ServiceName, consulClient)),
	)

	frontendUtils.MustHandleError(err)
}

func initCartClient() {
	CartClient, err = cartservice.NewClient("cart", client.WithSuite(clientsuite.CommonClientSuite{
		CurrentServiceName: ServiceName,
		RegistryAddr:       RegistryAddr,
	}))

	frontendUtils.MustHandleError(err)
}

func initCheckoutClient() {

	CheckoutClient, err = checkoutservice.NewClient("checkout", client.WithSuite(clientsuite.CommonClientSuite{
		CurrentServiceName: ServiceName,
		RegistryAddr:       RegistryAddr,
	}))

	frontendUtils.MustHandleError(err)
}

func initOrderClient() {

	OrderClient, err = orderservice.NewClient("order", client.WithSuite(clientsuite.CommonClientSuite{
		CurrentServiceName: ServiceName,
		RegistryAddr:       RegistryAddr,
	}))

	frontendUtils.MustHandleError(err)
}
