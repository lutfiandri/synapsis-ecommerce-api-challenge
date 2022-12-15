package model

import "gorm.io/gorm"

type Checkout struct {
	gorm.Model
	UserID     uint `json:"userID" binding:"required"`
	TotalPrice int  `json:"totalPrice"`
	Paid       bool `json:"paid" gorm:"default:false"`

	User      User
	CartItems []CartItem
}

func (Checkout) TableName() string {
	return "checkouts"
}
