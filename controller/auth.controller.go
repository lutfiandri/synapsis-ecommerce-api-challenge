package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/helper"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthController interface {
	Route(*gin.Engine)
	SignUp(*gin.Context)
	SignIn(*gin.Context)
}

type authController struct {
	repository repository.UserRepository
}

func NewAuthController(repository *repository.UserRepository) AuthController {
	return &authController{
		repository: *repository,
	}
}

func (c *authController) Route(router *gin.Engine) {
	router.POST("/auth/signup", c.SignUp)
	router.POST("/auth/signin", c.SignIn)
}

func (c *authController) SignUp(ctx *gin.Context) {
	var user model.User
	err := ctx.BindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	user.Password = string(hashedPassword)

	err = c.repository.Create(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	fmt.Printf("%+v\n", user)
	ctx.JSON(http.StatusCreated, gin.H{
		"user": user,
	})
}

type signInRequest struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
}

func (c *authController) SignIn(ctx *gin.Context) {
	var userRequest signInRequest
	err := ctx.BindJSON(&userRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := c.repository.FindOne(&userRequest.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Email not found",
		})
		return
	}

	passwordTest := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userRequest.Password))
	if passwordTest != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Wrong password",
		})
		return
	}

	token, err := helper.GenerateJWT(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"user":        user,
		"accessToken": token,
	})
}
