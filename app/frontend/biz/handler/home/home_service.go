package home

import (
	"context"

	"github.com/cloudwego/biz-demo/gomall/app/frontend/biz/service"
	"github.com/cloudwego/biz-demo/gomall/app/frontend/biz/utils"
	common "github.com/cloudwego/biz-demo/gomall/app/frontend/hertz_gen/frontend/common"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

//Home 是一个处理 HTTP 请求的函数。
//@router / [GET] 是注释，它使用了文档生成工具（例如 Swag）来自动生成 API 文档，
//表示这个函数处理 GET 请求并且对应的 URL 路径是 /

// Home .
// @router / [GET]
func Home(ctx context.Context, c *app.RequestContext) {
	var err error
	var req common.Empty
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	//req：定义了一个 home.Empty 类型的变量，表示请求体的内容。这里 Empty 可能是一个空的结构体，
	//意味着不需要传入任何参数。
	//c.BindAndValidate(&req)：从请求中绑定并验证数据。它将请求中的数据绑定到 req 变量中，
	//并验证数据是否符合要求。如果绑定或验证失败，将返回错误。
	//如果有错误，使用 utils.SendErrResponse 发送错误响应，返回给客户端。

	//resp := &home.Empty{}
	resp, err := service.NewHomeService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	//创建一个空的 home.Empty 类型的响应对象 resp。
	//service.NewHomeService(ctx, c).Run(&req)：
	//调用 NewHomeService 创建一个新的服务实例，并执行其 Run 方法，
	//传入请求数据 req。这通常会处理业务逻辑，并返回结果（存储在 resp 中）。
	//如果 Run 方法返回错误，发送错误响应并返回。
	//resp["user_id"] = 22
	c.HTML(consts.StatusOK, "home", utils.WarpResponse(ctx, c, resp))
	//c.HTML(consts.StatusOK, "home", resp)
	//utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
