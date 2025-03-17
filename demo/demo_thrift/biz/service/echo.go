package service

import (
	"context"
	"fmt"
	api "github.com/cloudwego/biz-demo/gomall/demo/demo_thrift/kitex_gen/api"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
)

type EchoService struct {
	ctx context.Context
} // NewEchoService new EchoService
func NewEchoService(ctx context.Context) *EchoService {
	return &EchoService{ctx: ctx}
}

// Run create note info
func (s *EchoService) Run(req *api.Request) (resp *api.Response, err error) {
	// Finish your business logic.
	info := rpcinfo.GetRPCInfo(s.ctx)
	//调用了 rpcinfo.GetRPCInfo 函数，获取关于当前 RPC 调用的信息。
	//s.ctx 是 EchoService 的上下文（Context）。
	//Context 通常用于跨 API 边界传递请求范围的信息，比如元数据、超时设置等。
	fmt.Println(info.From().ServiceName())

	return &api.Response{Message: req.Message}, nil
}
