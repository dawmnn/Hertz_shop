package service

import (
	"context"
	"fmt"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/infra/mq"
	"github.com/cloudwego/biz-demo/gomall/app/checkout/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/cart"
	checkout "github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/checkout"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/email"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/order"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/payment"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/kitex/pkg/kerrors"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/protobuf/proto"
)

type CheckoutService struct {
	ctx context.Context
} // NewCheckoutService new CheckoutService
func NewCheckoutService(ctx context.Context) *CheckoutService {
	return &CheckoutService{ctx: ctx}
}

// Run create note info
func (s *CheckoutService) Run(req *checkout.CheckoutReq) (resp *checkout.CheckoutResp, err error) {
	// Finish your business logic.
	cartResult, err := rpc.CartClient.GetCart(s.ctx, &cart.GetCartReq{UserId: req.UserId})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(5005001, err.Error())
	}
	if cartResult == nil || cartResult.Items == nil {
		return nil, kerrors.NewGRPCBizStatusError(5004001, "cart is empty")
	}
	var (
		total float32
		oi    []*order.OrderItem
	)
	for _, cartItem := range cartResult.Items {
		productResp, resultErr := rpc.ProductClient.GetProduct(s.ctx, &product.GetProductReq{
			Id: cartItem.ProductId,
		})
		if resultErr != nil {
			return nil, resultErr
		}
		if productResp.Product == nil {
			continue
		}
		p := productResp.Product.Price
		cost := p * float32(cartItem.Quantity)
		total += cost
		oi = append(oi, &order.OrderItem{
			Item: &cart.CartItem{
				ProductId: cartItem.ProductId,
				Quantity:  cartItem.Quantity,
			},
			Cost: cost,
		})
	}

	var orderId string
	orderResp, err := rpc.OrderClient.PlaceOrder(s.ctx, &order.PlaceOrderReq{
		UserId: req.UserId,
		Email:  req.Email,
		Address: &order.Address{
			StreetAddress: req.Address.StreetAddress,
			City:          req.Address.City,
			State:         req.Address.State,
			Country:       req.Address.Country,
			ZipCode:       req.Address.ZipCode,
		},
		Items: oi,
	})
	if err != nil {
		return nil, kerrors.NewGRPCBizStatusError(5004002, err.Error())
	}
	if orderResp.Order != nil && orderResp != nil {
		orderId = orderResp.Order.OrderId
	}

	payReq := &payment.ChargeReq{
		UserId:  req.UserId,
		OrderId: orderId,
		Amount:  total,
		CreditCard: &payment.CreditCardInfo{
			CreditCardNumber:          req.CreditCard.CreditCardNumber,
			CreditCardCvv:             req.CreditCard.CreditCardCvv,
			CreditCardExpirationMonth: req.CreditCard.CreditCardExpirationMonth,
			CreditCardExpirationYear:  req.CreditCard.CreditCardExpirationYear,
		},
	}
	_, err = rpc.CartClient.EmptyCart(s.ctx, &cart.EmptyCartReq{UserId: req.UserId})
	if err != nil {
		klog.Error(err.Error())
	}
	paymentResult, err := rpc.PaymentClient.Charge(s.ctx, payReq)
	if err != nil {
		err = fmt.Errorf("Charge.err:%v", err)
		return
	}

	data, _ := proto.Marshal(&email.EmailReq{
		From:        "from@example.com",
		To:          req.Email,
		ContentType: "text/plain",
		Subject:     "You have just created an order in the CloudWeGo shop",
		Content:     "You have just created an order in the CloudWeGo shop",
	})

	//将结构体 email.EmailReq 编码成二进制数据格式，以便可以通过网络传输或存储。
	msg := &nats.Msg{Subject: "email", Data: data, Header: make(nats.Header)}
	//otel.GetTextMapPropagator()：获取 OpenTelemetry 的默认文本映射传播器（TextMapPropagator）。
	//这个传播器是用于从一个上下文（ctx）提取或将追踪信息注入到可传输的格式
	//（如 HTTP 头部、消息队列消息头等）中的工具。
	//propagation.HeaderCarrier(msg.Header)：HeaderCarrier 是 OpenTelemetry 的一个适配器，
	//用于将 nats.Header 作为容器，允许将追踪信息注入到消息头部。
	//它实现了 OpenTelemetry 的 TextMapCarrier 接口，使得 msg.Header 可以像 HTTP 头部一样处理。
	otel.GetTextMapPropagator().Inject(s.ctx, propagation.HeaderCarrier(msg.Header))
	//Inject 方法会从 s.ctx 中提取这些信息，并注入到 msg.Header 中，确保追踪信息得以传播。
	_ = mq.Nc.PublishMsg(msg)
	//用于将消息发送到 NATS 服务器

	klog.Info(paymentResult)
	klog.Info(orderResp)
	resp = &checkout.CheckoutResp{
		OrderId:     orderId,
		Transaction: paymentResult.TransactionId,
	}
	return
}
