package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/repository"
)

type UserController interface {
	Route(*gin.Engine)
	Create(*gin.Context)
	FindAll(*gin.Context)
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
	router.GET("/", c.FindAll)
}

func (c *userController) Create(ctx *gin.Context) {
}

func (c *userController) FindAll(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "mantap bang",
	})
}
