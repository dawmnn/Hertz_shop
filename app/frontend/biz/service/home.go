package service

import (
	"context"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/common"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/infra/rpc"
	"github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen/product"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/kitex/pkg/klog"
)

type HomeService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewHomeService(Context context.Context, RequestContext *app.RequestContext) *HomeService {
	return &HomeService{RequestContext: RequestContext, Context: Context}
}

func (h *HomeService) Run(req *common.Empty) (map[string]any, error) {
	//
	//var resp = make(map[string]any)
	//items := []map[string]any{
	//	{"Name": "T-shirt-1", "Price": 100, "Picture": "/static/image/t-shirt-1.jpeg"},
	//	{"Name": "T-shirt-2", "Price": 110, "Picture": "/static/image/t-shirt-1.jpeg"},
	//	{"Name": "T-shirt-3", "Price": 120, "Picture": "/static/image/t-shirt-2.jpeg"},
	//	{"Name": "T-shirt-4", "Price": 130, "Picture": "/static/image/t-shirt-2.jpeg"},
	//	{"Name": "T-shirt-5", "Price": 140, "Picture": "/static/image/t-shirt-1.jpeg"},
	//	{"Name": "T-shirt-6", "Price": 105, "Picture": "/static/image/t-shirt-3.jpeg"},
	//}
	//resp["Title"] = "Hot Sales"
	//resp["Items"] = items
	//return resp, nil

	p, err := rpc.ProductClient.ListProducts(h.Context, &product.ListProductsReq{})
	if err != nil {
		klog.Error(err)
	}
	return utils.H{
		"Title": "Hot sale",
		"Items": p.Products,
	}, nil
}
