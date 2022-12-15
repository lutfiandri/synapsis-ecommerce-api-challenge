package repository

import (
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
	"gorm.io/gorm"
)

type CartItemRepository interface {
	Create(*model.CartItem) error
	FindOneByID(*string) (model.CartItem, error)
	FindAll() ([]model.CartItem, error)
	UpdateOneByID(*string, *model.CartItem) error
	DeleteOneByID(*string) error
}

type cartItemRepository struct {
	db *gorm.DB
}

func NewCartItemRepository(database *gorm.DB) CartItemRepository {
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

func (r *cartItemRepository) FindAll() ([]model.CartItem, error) {
	var cartItems []model.CartItem
	err := r.db.Preload("User").Preload("Product").Find(&cartItems).Error
	return cartItems, err
}

func (r *cartItemRepository) UpdateOneByID(id *string, newCartItem *model.CartItem) error {
	cartItem, err := r.FindOneByID(id)
	if err != nil {
		return err
	}

	cartItem.Quantity = newCartItem.Quantity
	cartItem.CheckoutID = newCartItem.CheckoutID

	err = r.db.Save(&cartItem).Error
	return err
}

func (r *cartItemRepository) DeleteOneByID(id *string) error {
	cartItem, err := r.FindOneByID(id)
	if err != nil {
		return err
	}

	err = r.db.Delete(&cartItem).Error
	return err
}
