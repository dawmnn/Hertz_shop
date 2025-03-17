package model

import (
	"context"
	"gorm.io/gorm"
)

type Category struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`

	Products []Product `json:"products" gorm:"many2many:product_category"`
}

func (c Category) TableName() string {
	return "category"
}

type CategoryQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func (c CategoryQuery) GetProductsByCategoryName(name string) (categories []Category, err error) {
	err = c.db.WithContext(c.ctx).Model(&Category{}).Where(&Category{Name: name}).
		Preload("Products").Find(&categories).Error
	//Preload("Products")：这是 GORM 的预加载方法，
	//它会自动加载与 Category 关联的 Products（假设 Category 和 Product 存在一对多或多对多的关系）。
	//预加载会在查询 Category 的同时，将与之关联的 Products 一并加载进来，避免后续的额外查询。
	return
}
func NewCategoryQuery(ctx context.Context, db *gorm.DB) *CategoryQuery {
	return &CategoryQuery{
		ctx: ctx,
		db:  db,
	}
}
