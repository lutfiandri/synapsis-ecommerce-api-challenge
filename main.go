package main

import (
	"github.com/gin-gonic/gin"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/controller"
	"github.com/lutfiandri/synapsis-ecommerce-api-challenge/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := "host=localhost user=postgres password=root dbname=ecommerce port=8081 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	userRepository := repository.NewUserRepository(db)
	userController := controller.NewUserController(&userRepository)

	r := gin.Default()
	userController.Route(r)

	r.Run(":8080")
}
