package model

import "gorm.io/gorm"

type OrderItem struct {
	gorm.Model
	ProductId    uint32  `gorm:"type:int(11)"`
	OrderIdRefer string  `gorm:"type:varchar(100);index"` //指示 GORM 在该字段上创建一个索引，以提高查询性能
	Quantity     uint32  `gorm:"type:int(11)"`
	Cost         float32 `gorm:"type:decimal(10,2)"`
}

func (OrderItem) TableName() string {
	return "order_item"
}
