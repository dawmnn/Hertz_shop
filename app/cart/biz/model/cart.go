package model

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type Cart struct {
	gorm.Model
	UserId    uint32 `gorm:"type:int(11);not null;index:idx_user_id"`
	ProductId uint32 `gorm:"type:int(11);not null;"`
	Qty       uint32 `gorm:"type:int(11);not null;"`
}

func (Cart) TableName() string {
	return "cart"
}

func AddItem(ctx context.Context, db *gorm.DB, cart *Cart) error {
	var row Cart
	err := db.WithContext(ctx).Model(&Cart{}).Where(&Cart{UserId: cart.UserId, ProductId: cart.ProductId}).First(&row).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if row.ID > 0 {
		return db.WithContext(ctx).Model(&Cart{}).
			Where(&Cart{UserId: cart.UserId, ProductId: cart.ProductId}).
			UpdateColumn("qty", gorm.Expr("qty+?", cart.Qty)).Error
	}

	return db.WithContext(ctx).Create(cart).Error
}

func EmptyCart(ctx context.Context, db *gorm.DB, userId uint32) error {
	if userId == 0 {
		return errors.New("user id is required")
	}
	return db.WithContext(ctx).Delete(&Cart{}, "user_id=?", userId).Error
}

func GetCartByUserId(ctx context.Context, db *gorm.DB, userId uint32) (cartList []*Cart, err error) {
	err = db.Debug().WithContext(ctx).Model(&Cart{}).Find(&cartList, "user_id = ?", userId).Error
	return cartList, err
}
