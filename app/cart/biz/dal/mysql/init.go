package mysql

import (
	"fmt"
	"github.com/cloudwego/biz-demo/gomall/app/cart/biz/model"
	"github.com/cloudwego/biz-demo/gomall/app/cart/conf"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"os"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	dsn := fmt.Sprintf(conf.GetConf().MySQL.DSN, os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"))
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:            true,
			SkipDefaultTransaction: true,
		},
	)
	//将一个不包含度量指标的追踪插件应用到数据库连接（DB）上
	if err = DB.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
		panic(err)
	}
	DB.AutoMigrate(&model.Cart{})
	if err != nil {
		panic(err)
	}
}
