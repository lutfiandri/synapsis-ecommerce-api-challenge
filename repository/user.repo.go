package repository

import (
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(*model.User) error
	FindOne(*string) (model.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(database *gorm.DB) UserRepository {
	return &userRepository{
		db: database,
	}
}

func (r *userRepository) Create(user *model.User) error {
	err := r.db.Create(&user).Error
	return err
}

func (r *userRepository) FindOne(email *string) (model.User, error) {
	var user model.User
	err := r.db.First(&user, "email = ?", email).Error
	return user, err
}
