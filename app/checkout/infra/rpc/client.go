package rpc

import (
	"github.com/cloudwego/biz-demo/gomall/app/checkout/conf"
	"github.com/cloudwego/biz-demo/gomall/common/clientsuite"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart/cartservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order/orderservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/payment/paymentservice"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product/productcatalogservice"
	"github.com/cloudwego/kitex/client"
	"sync"
)

var (
	ServiceName   = conf.GetConf().Kitex.Service
	RegistryAddr  = conf.GetConf().Registry.RegistryAddress[0]
	CartClient    cartservice.Client
	ProductClient productcatalogservice.Client
	PaymentClient paymentservice.Client
	OrderClient   orderservice.Client
	once          sync.Once
	err           error
)

func InitClient() {
	once.Do(func() {
		initCartClient()
		initProductClient()
		initPaymentClient()
		initOrderClient()
	})
}
func initCartClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
	}
	//r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	//if err != nil {
	//	panic(err)
	//}
	//opts = append(opts, client.WithResolver(r))
	//opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GetConf().Kitex.Service}),
	//	client.WithTransportProtocol(transport.GRPC),
	//	client.WithMetaHandler(transmeta.ClientHTTP2Handler),
	//)
	CartClient, err = cartservice.NewClient("cart", opts...)
	if err != nil {
		panic(err)
	}
}
func initProductClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
	}
	//r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	//if err != nil {
	//	panic(err)
	//}
	//opts = append(opts, client.WithResolver(r))
	//opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GetConf().Kitex.Service}),
	//	client.WithTransportProtocol(transport.GRPC),
	//	client.WithMetaHandler(transmeta.ClientHTTP2Handler),
	//)
	ProductClient, err = productcatalogservice.NewClient("product", opts...)
	if err != nil {
		panic(err)
	}
}
func initPaymentClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
	}
	//r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	//if err != nil {
	//	panic(err)
	//}
	//opts = append(opts, client.WithResolver(r))
	//opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GetConf().Kitex.Service}),
	//	client.WithTransportProtocol(transport.GRPC),
	//	client.WithMetaHandler(transmeta.ClientHTTP2Handler),
	//)
	PaymentClient, err = paymentservice.NewClient("payment", opts...)
	if err != nil {
		panic(err)
	}
}
func initOrderClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonClientSuite{
			CurrentServiceName: ServiceName,
			RegistryAddr:       RegistryAddr,
		}),
	}

	//r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	//if err != nil {
	//	panic(err)
	//}
	//opts = append(opts, client.WithResolver(r))
	//opts = append(opts, client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: conf.GetConf().Kitex.Service}),
	//	client.WithTransportProtocol(transport.GRPC),
	//	client.WithMetaHandler(transmeta.ClientHTTP2Handler),
	//)
	OrderClient, err = orderservice.NewClient("order", opts...)
	if err != nil {
		panic(err)
	}
}
