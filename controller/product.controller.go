package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/middleware"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/model"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/repository"
	"gorm.io/gorm"
)

type ProductController interface {
	Route(*gin.Engine)
	Create(*gin.Context)
	FindOneByID(*gin.Context)
	FindAll(*gin.Context)
	UpdateOneByID(*gin.Context)
	DeleteOneByID(*gin.Context)
}

type productController struct {
	repository repository.ProductRepository
}

func NewProductController(repository *repository.ProductRepository) ProductController {
	return &productController{
		repository: *repository,
	}
}

func (c *productController) Route(router *gin.Engine) {
	router.POST("/products", middleware.AuthorizeJWT(), middleware.AuthorizeUserRole("SELLER"), c.Create)
	router.GET("/products", c.FindAll)
	router.GET("/products/:id", c.FindOneByID)
	router.PUT("/products/:id", middleware.AuthorizeJWT(), middleware.AuthorizeUserRole("SELLER"), c.UpdateOneByID)
	router.DELETE("/products/:id", middleware.AuthorizeJWT(), middleware.AuthorizeUserRole("SELLER"), c.DeleteOneByID)
}

func (c *productController) Create(ctx *gin.Context) {
	var product model.Product
	err := ctx.BindJSON(&product)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	err = c.repository.Create(&product)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"product": product,
	})
}

func (c *productController) FindOneByID(ctx *gin.Context) {
	id := ctx.Param("id")

	product, err := c.repository.FindOneByID(&id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (c *productController) FindAll(ctx *gin.Context) {
	products, err := c.repository.FindAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}

func (c *productController) UpdateOneByID(ctx *gin.Context) {
	var product model.Product

	err := ctx.BindJSON(&product)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, err)
		return
	}

	id := ctx.Param("id")
	err = c.repository.UpdateOneByID(&id, &product)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"product": product,
	})
}

func (c *productController) DeleteOneByID(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.repository.DeleteOneByID(&id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
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
