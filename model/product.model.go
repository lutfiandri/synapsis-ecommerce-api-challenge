package model

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Price       int    `json:"price" binding:"required"`
}

func (Product) TableName() string {
	return "products"
}
