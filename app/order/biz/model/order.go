package model

import (
	"context"
	"gorm.io/gorm"
)

type Consignee struct {
	Email         string
	StreetAddress string
	City          string
	State         string
	County        string
	ZipCode       string
}

type Order struct {
	gorm.Model
	OrderId      string    `gorm:"type:varchar(100);uniqueIndex"`
	UserId       uint32    `gorm:"type:int(11)"`
	UserCurrency string    `gorm:"type:varchar(10)"`
	Consignee    Consignee `gorm:"embedded"`
	//让嵌套结构体的字段直接展平到当前结构体中，使得它们成为数据库表的一部分，而不是单独的关联表。
	OrderItems []OrderItem `gorm:"foreignKey:OrderIdRefer;references:OrderId"` //一对多
	//foreignKey:OrderId：OrderItem 表的 OrderId 字段是外键，指向 Order 表。
	//references:OrderId：Order 表的 OrderId 字段是被引用的字段，通常是父表的主键。
}

func (Order) TableName() string {
	return "order"
}

func ListOrder(ctx context.Context, db *gorm.DB, UserId uint32) ([]*Order, error) {
	var orders []*Order
	//查找关联表
	err := db.WithContext(ctx).Where("user_id=?", UserId).Preload("OrderItems").Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
