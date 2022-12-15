package repository

import (
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(*model.Product) error
	FindOneByID(*string) (model.Product, error)
	FindAll() ([]model.Product, error)
	UpdateOneByID(*string, *model.Product) error
	DeleteOneByID(*string) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(database *gorm.DB) ProductRepository {
	return &productRepository{
		db: database,
	}
}

func (r *productRepository) Create(product *model.Product) error {
	err := r.db.Create(&product).Error
	return err
}

func (r *productRepository) FindOneByID(id *string) (model.Product, error) {
	var product model.Product
	err := r.db.First(&product, "id = ?", id).Error
	return product, err
}

func (r *productRepository) FindAll() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepository) UpdateOneByID(id *string, newProduct *model.Product) error {
	product, err := r.FindOneByID(id)
	if err != nil {
		return err
	}

	product.Title = newProduct.Title
	product.Description = newProduct.Description
	product.Price = newProduct.Price

	err = r.db.Save(&product).Error
	return err
}

func (r *productRepository) DeleteOneByID(id *string) error {
	product, err := r.FindOneByID(id)
	if err != nil {
		return err
	}

	err = r.db.Delete(&product).Error
	return err
}
