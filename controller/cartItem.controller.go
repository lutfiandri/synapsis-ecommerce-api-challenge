package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/repository"
	"gorm.io/gorm"
)

type CartItemController interface {
	Route(*gin.Engine)
	Create(*gin.Context)
	FindOneByID(*gin.Context)
	FindAll(*gin.Context)
}

type cartItemController struct {
	repository repository.CartItemRepository
}

func NewCartItemController(repository *repository.CartItemRepository) CartItemController {
	return &cartItemController{
		repository: *repository,
	}
}

func (c *cartItemController) Route(router *gin.Engine) {
	router.POST("/cart-items", c.Create)
	router.GET("/cart-items/", c.FindAll)
	router.GET("/cart-items/:id", c.FindOneByID)
	router.PUT("/cart-items/:id", c.UpdateOneByID)
}

type createRequest struct {
	UserID    uint `binding:"required"`
	ProductID uint `binding:"required"`
	Quantity  int  `binding:"required"`
}

func (c *cartItemController) Create(ctx *gin.Context) {
	var cartItemRequest createRequest
	err := ctx.BindJSON(&cartItemRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	cartItem := model.CartItem{
		UserID:    cartItemRequest.UserID,
		ProductID: cartItemRequest.ProductID,
		Quantity:  cartItemRequest.Quantity,
	}

	err = c.repository.Create(&cartItem)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	id := fmt.Sprintf("%d", cartItem.ID)
	cartItemResponse, err := c.repository.FindOneByID(&id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "CartItem not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"cartItem": cartItemResponse,
	})
}

func (c *cartItemController) FindOneByID(ctx *gin.Context) {
	id := ctx.Param("id")

	cartItem, err := c.repository.FindOneByID(&id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "CartItem not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"cartItem": cartItem,
	})
}

func (c *cartItemController) FindAll(ctx *gin.Context) {
	cartItems, err := c.repository.FindAll()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "CartItem not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"cartItems": cartItems,
	})
}

type updateOneByIdRequest struct {
	Quantity int
}

func (c *cartItemController) UpdateOneByID(ctx *gin.Context) {
	var cartItemRequest updateOneByIdRequest

	err := ctx.BindJSON(&cartItemRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	cartItem := model.CartItem{
		Quantity: cartItemRequest.Quantity,
	}

	id := ctx.Param("id")
	err = c.repository.UpdateOneByID(&id, &cartItem)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "CartItem not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"cartItem": cartItem,
	})
}
