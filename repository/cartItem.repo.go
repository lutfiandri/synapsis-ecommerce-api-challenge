package repository

import (
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	Create(*model.CartItem) error
	FindOneByID(*string) (model.CartItem, error)
}

type cartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository(database *gorm.DB) CartItemRepository {
	database.AutoMigrate(&model.CartItem{})
	return &cartItemRepository{
		db: database,
	}
}

func (r *cartItemRepository) Create(cartItem *model.CartItem) error {
	err := r.db.Create(&cartItem).Error
	return err
}

func (r *cartItemRepository) FindOneByID(id *string) (model.CartItem, error) {
	var cartItem model.CartItem
	err := r.db.Preload("User").Preload("Product").First(&cartItem, "id = ?", id).Error
	return cartItem, err
}
