package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/repository"
)

type UserController interface {
	Route(*gin.Engine)
}

type userController struct {
	repo *repository.UserRepository
}

func NewUserController(repo *repository.UserRepository) UserController {
	return &userController{
		repo: repo,
	}
}

func (c *userController) Route(router *gin.Engine) {
	// router.GET("/", c.FindAll)
	// router.POST("/", c.Create)
}
