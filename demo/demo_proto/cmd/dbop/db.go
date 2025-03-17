package main

import (
	"github.com/cloudwego/biz-demo/gomall/demo/demo_proto/biz/dal"
	"github.com/cloudwego/biz-demo/gomall/demo/demo_proto/biz/dal/mysql"
	"github.com/cloudwego/biz-demo/gomall/demo/demo_proto/biz/model"

	//"github.com/gogo/protobuf/protoc-gen-gogo/testdata/imports/fmt"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	dal.Init()

	//mysql.DB.Create(&model.User{Email: "demo@exmalple.com", Password: "12345"})
	//mysql.DB.Model(&model.User{}).Where("email=?", "demo@exmalple.com").Update("password", "11111")
	//var row model.User
	//mysql.DB.Model(&model.User{}).Where("email=?", "demo@exmalple.com").First(&row)
	//fmt.Printf("row:%+v\n", row)
	//mysql.DB.Where("email=?", "demo@exmalple.com").Delete(&model.User{})

	mysql.DB.Unscoped().Where("email=?", "demo@exmalple.com").Delete(&model.User{})
}
