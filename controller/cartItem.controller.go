package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/middleware"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/repository"
	"gorm.io/gorm"
)

type CartItemController interface {
	Route(*gin.Engine)
	Create(*gin.Context)
	FindOneByID(*gin.Context)
	FindAll(*gin.Context)
	UpdateOneByID(*gin.Context)
	DeleteOneByID(*gin.Context)
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
	router.POST("/cart-items", middleware.AuthorizeJWT(), middleware.AuthorizeUserRole("BUYER"), c.Create)
	router.GET("/cart-items/", middleware.AuthorizeJWT(), middleware.AuthorizeUserRole("BUYER"), c.FindAll)
	router.GET("/cart-items/:id", middleware.AuthorizeJWT(), middleware.AuthorizeUserRole("BUYER"), c.FindOneByID)
	router.PUT("/cart-items/:id", middleware.AuthorizeJWT(), middleware.AuthorizeUserRole("BUYER"), c.UpdateOneByID)
	router.DELETE("/cart-items/:id", middleware.AuthorizeJWT(), middleware.AuthorizeUserRole("BUYER"), c.DeleteOneByID)
}

type createCheckoutRequest struct {
	ProductID uint `binding:"required"`
	Quantity  int  `binding:"required"`
}

func (c *cartItemController) Create(ctx *gin.Context) {
	var cartItemRequest createCheckoutRequest
	err := ctx.BindJSON(&cartItemRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if cartItemRequest.Quantity < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "quantity must grater than 0",
		})
		return
	}

	userID := ctx.GetUint("UserID")
	cartItem := model.CartItem{
		UserID:    userID,
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

	userID := ctx.GetUint("UserID")
	userIDstr := fmt.Sprintf("%d", userID)

	cartItem, err := c.repository.FindOneByIDAndUserID(&id, &userIDstr)
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
	userID := ctx.GetUint("UserID")
	userIDstr := fmt.Sprintf("%d", userID)
	cartItems, err := c.repository.FindManyByUserID(&userIDstr)
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

	if cartItemRequest.Quantity < 1 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "quantity must grater than 0",
		})
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

func (c *cartItemController) DeleteOneByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.repository.DeleteOneByID(&id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Cart Item not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.Status(http.StatusOK)
}
