package model

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type Product struct {
	Base
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Picture     string  `json:"picture"`
	Price       float32 `json:"price"`

	Categories []Category `json:"categories" gorm:"many2many:product_category"`
}

//这是一个 GORM 标签，它用于定义数据库中的关系。具体来说，many2many 表示这是一个多对多关系，
//意味着 Product 和 Category 之间存在多对多的关系。
//product_category 是用来表示这个多对多关系的中间表的名称。也就是说，
//GORM 会创建一个名为 product_category 的中间表，用来关联 Product 和 Category 表中的记录。

func (p Product) TableName() string {
	return "product"
}

type ProductQuery struct {
	ctx context.Context
	db  *gorm.DB
}

func (p ProductQuery) GetById(productId int) (product Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).First(&product, productId).Error
	return
}

func (p ProductQuery) SearchProducts(q string) (Products []*Product, err error) {
	err = p.db.WithContext(p.ctx).Model(&Product{}).Find(&Products,
		"name like ? or description like ?", "%"+q+"%", "%"+q+"%").Error
	return
}

func NewProductQuery(ctx context.Context, db *gorm.DB) *ProductQuery {
	return &ProductQuery{
		ctx: ctx,
		db:  db,
	}
}

//服务优化

type CachedProductQuery struct {
	productQuery ProductQuery
	cacheClient  *redis.Client
	prefix       string
}

func (c CachedProductQuery) GetById(productId int) (product Product, err error) {
	cachedKey := fmt.Sprintf("%s_%s_%d", c.prefix, "product_by_id", productId)
	cachedResult := c.cacheClient.Get(c.productQuery.ctx, cachedKey)

	err = func() error {
		if err := cachedResult.Err(); err != nil {
			return err
		}
		cachedResultByte, err := cachedResult.Bytes()
		if err != nil {
			return err
		}
		err = json.Unmarshal(cachedResultByte, &product)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		product, err = c.productQuery.GetById(productId)
		if err != nil {
			return Product{}, err
		}
		encoded, err := json.Marshal(product)
		if err != nil {
			return product, nil
		}
		_ = c.cacheClient.Set(c.productQuery.ctx, cachedKey, encoded, time.Hour)
	}
	return
}
func (c CachedProductQuery) SearchProducts(q string) (product []*Product, err error) {
	return c.productQuery.SearchProducts(q)
}

func NewCachedProductQuery(ctx context.Context, db *gorm.DB, cacheClient *redis.Client) *CachedProductQuery {
	return &CachedProductQuery{
		productQuery: *NewProductQuery(ctx, db),
		cacheClient:  cacheClient,
		prefix:       "shop",
	}
}

//读写分离

type ProductMutation struct {
	ctx context.Context
	db  *gorm.DB
}
