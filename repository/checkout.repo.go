package repository

import (
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
	"gorm.io/gorm"
)

type CheckoutRepository interface {
	Create(*model.Checkout) error
	FindOneByID(*string) (model.Checkout, error)
	FindAll() ([]model.Checkout, error)
	UpdateOneByID(*string, *model.Checkout) error
	DeleteOneByID(*string) error
}

type checkoutRepository struct {
	db *gorm.DB
}

func NewCheckoutRepository(database *gorm.DB) CheckoutRepository {
	return &checkoutRepository{
		db: database,
	}
}

func (r *checkoutRepository) Create(checkout *model.Checkout) error {
	err := r.db.Create(&checkout).Error
	return err
}

func (r *checkoutRepository) FindOneByID(id *string) (model.Checkout, error) {
	var checkout model.Checkout
	err := r.db.Preload("CartItems.Product").Preload("User").First(&checkout, "id = ?", id).Error
	return checkout, err
}

func (r *checkoutRepository) FindAll() ([]model.Checkout, error) {
	var checkouts []model.Checkout
	err := r.db.Preload("CartItems.Product").Preload("User").Find(&checkouts).Error
	return checkouts, err
}

func (r *checkoutRepository) UpdateOneByID(id *string, newCheckout *model.Checkout) error {
	checkout, err := r.FindOneByID(id)
	if err != nil {
		return err
	}

	checkout.Paid = newCheckout.Paid

	err = r.db.Save(&checkout).Error
	return err
}

func (r *checkoutRepository) DeleteOneByID(id *string) error {
	checkout, err := r.FindOneByID(id)
	if err != nil {
		return err
	}

	err = r.db.Delete(&checkout).Error
	return err
}
