package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/controller"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/database"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/repository"
)

func main() {
	godotenv.Load()

	dbHost, dbPort, dbUsername, dbPassword, dbDatabase := os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_USERNAME"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DATABASE")
	db := database.NewPostgresDatabase(dbHost, dbPort, dbUsername, dbPassword, dbDatabase)

	userRepository := repository.NewUserRepository(db)
	authController := controller.NewAuthController(&userRepository)

	productRepository := repository.NewProductRepository(db)
	productController := controller.NewProductController(&productRepository)

	cartItemRepository := repository.NewCartItemRepository(db)
	cartItemController := controller.NewCartItemController(&cartItemRepository)

	checkoutRepository := repository.NewCheckoutRepository(db)
	checkoutController := controller.NewCheckoutController(&checkoutRepository, &cartItemRepository)

	r := gin.Default()
	authController.Route(r)
	productController.Route(r)
	cartItemController.Route(r)
	checkoutController.Route(r)

	r.Run(":" + os.Getenv("APP_PORT"))
}
