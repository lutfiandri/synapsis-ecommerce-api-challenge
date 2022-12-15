package controller

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/repository"
	"gorm.io/gorm"
)

type CheckoutController interface {
	Route(*gin.Engine)
	Create(*gin.Context)
	FindOneByID(*gin.Context)
	UpdateOneByID(*gin.Context)
	DeleteOneByID(*gin.Context)
}

type checkoutController struct {
	repository         repository.CheckoutRepository
	cartItemRepository repository.CartItemRepository
}

func NewCheckoutController(repository *repository.CheckoutRepository, cartItemRepository *repository.CartItemRepository) CheckoutController {
	return &checkoutController{
		repository:         *repository,
		cartItemRepository: *cartItemRepository,
	}
}

func (c *checkoutController) Route(router *gin.Engine) {
	router.POST("/checkouts", c.Create)
	router.GET("/checkouts", c.FindAll)
	router.GET("/checkouts/:id", c.FindOneByID)
	router.PUT("/checkouts/:id", c.UpdateOneByID)
	router.DELETE("/checkouts/:id", c.DeleteOneByID)
}

type createRequest struct {
	UserID  uint   `json:"userID" binding:"required"`
	ItemsID []uint `json:"itemsID" binding:"required"`
}

func (c *checkoutController) Create(ctx *gin.Context) {
	var checkoutRequest createRequest

	err := ctx.BindJSON(&checkoutRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}

	// get cart items
	var cartItems []model.CartItem
	totalPrice := 0

	for _, itemID := range checkoutRequest.ItemsID {
		id := fmt.Sprintf("%d", itemID)
		cartItem, err := c.cartItemRepository.FindOneByID(&id)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"error": "Item with id " + id + " is not found",
				})
				return
			}
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err,
			})
			return
		}

		totalPrice += cartItem.Quantity * cartItem.Product.Price
		cartItems = append(cartItems, cartItem)
	}

	// create checkout -> get ID
	checkout := model.Checkout{
		UserID:     checkoutRequest.UserID,
		Paid:       false,
		TotalPrice: totalPrice,
	}

	err = c.repository.Create(&checkout)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	fmt.Println("checkout id: ", checkout.ID)

	// set checkoutID to each cartItem
	for _, cartItem := range cartItems {
		var p uint = 5
		cartItem.CheckoutID = &p
		cartItemID := fmt.Sprintf("%d", cartItem.ID)
		err := c.cartItemRepository.UpdateOneByID(&cartItemID, &cartItem)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"checkout": checkout,
	})
}

func (c *checkoutController) FindOneByID(ctx *gin.Context) {
	id := ctx.Param("id")

	checkout, err := c.repository.FindOneByID(&id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Checkout not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"checkout": checkout,
	})
}

func (c *checkoutController) FindAll(ctx *gin.Context) {
	checkouts, err := c.repository.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"checkouts": checkouts,
	})
}

type checkoutUpdateRequest struct {
	Paid bool
}

func (c *checkoutController) UpdateOneByID(ctx *gin.Context) {
	var updateRequest checkoutUpdateRequest

	err := ctx.BindJSON(&updateRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	checkout := model.Checkout{
		Paid: updateRequest.Paid,
	}

	id := ctx.Param("id")
	err = c.repository.UpdateOneByID(&id, &checkout)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Checkout not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"checkout": checkout,
	})
}

func (c *checkoutController) DeleteOneByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.repository.DeleteOneByID(&id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Checkout not found",
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