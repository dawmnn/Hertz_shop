.PHONY: gen-demo-proto
gen-demo-proto:
@cd demo/demo_proto && cwgo server -I ../../idl --type RPC --module github.com/cloudwego/biz-demo/gomall/demo/demo_proto --service demo_proto --idl ../../idl/echo.proto

.PHONY: gen-demo-thrift
gen-demo-thrift:
@cd demo/demo_thrift && cwgo server -I ../../idl --type RPC --module github.com/cloudwego/biz-demo/gomall/demo/demo_thrift --service demo_thrift --idl ../../idl/echo.thrift

.PHONY:gen-frontend
gen-frontend:
cwgo server --type HTTP --idl ../../idl/frontend/checkout_page.proto --service frontend -module github.com/cloudwego/biz-demo/gomall/app/frontend -I ../../idl

cwgo client --type RPC --service user --module github.com/cloudwego/biz-demo/gomall/rpc_gen --I ../idl --idl ../idl/user.proto

cwgo server --type RPC --service user --module github.com/cloudwego/biz-demo/gomall/app/user --pass "-use github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/user.proto



cwgo client --type RPC --service product --module github.com/cloudwego/biz-demo/gomall/rpc_gen --I ../idl --idl ../idl/product.proto

cwgo server --type RPC --service product --module github.com/cloudwego/biz-demo/gomall/app/product --pass "-use github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/product.proto




cwgo client --type RPC --service cart --module github.com/cloudwego/biz-demo/gomall/rpc_gen --I ../idl --idl ../idl/cart.proto

cwgo server --type RPC --service cart --module github.com/cloudwego/biz-demo/gomall/app/cart --pass "-use github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/cart.proto


cwgo client --type RPC --service payment --module github.com/cloudwego/biz-demo/gomall/rpc_gen --I ../idl --idl ../idl/payment.proto

cwgo server --type RPC --service payment --module github.com/cloudwego/biz-demo/gomall/app/payment --pass "-use github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/payment.proto


cwgo client --type RPC --service checkout --module github.com/cloudwego/biz-demo/gomall/rpc_gen --I ../idl --idl ../idl/checkout.proto

cwgo server --type RPC --service checkout --module github.com/cloudwego/biz-demo/gomall/app/checkout --pass "-use github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/checkout.proto


cwgo client --type RPC --service order --module github.com/cloudwego/biz-demo/gomall/rpc_gen --I ../idl --idl ../idl/order.proto

cwgo server --type RPC --service order --module github.com/cloudwego/biz-demo/gomall/app/order --pass "-use github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/order.proto


cwgo client --type RPC --service email --module github.com/cloudwego/biz-demo/gomall/rpc_gen --I ../idl --idl ../idl/email.proto

cwgo server --type RPC --service email --module github.com/cloudwego/biz-demo/gomall/app/email --pass "-use github.com/cloudwego/biz-demo/gomall/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/email.proto

docker run -v D:/golang_study/Hertz/app/frontend/conf:/opt/gomall/frontend/conf --network idl_default --env-file ./app/frontend/.env -p 8080:8080 dafc4fed59d
