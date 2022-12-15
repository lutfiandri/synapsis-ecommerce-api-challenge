package model

import (
	"gorm.io/gorm"
)

type CartItem struct {
	gorm.Model
	UserID     uint  `json:"userID" binding:"required"`
	ProductID  uint  `json:"productID" binding:"required"`
	CheckoutID *uint `json:"checkoutID"`
	Quantity   int   `json:"quantity" gorm:"default:1"`

	User    User
	Product Product
}

func (CartItem) TableName() string {
	return "cart_items"
}
