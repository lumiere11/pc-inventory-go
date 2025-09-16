package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lumiere11/pc-inventory-go/handlers"
	"github.com/lumiere11/pc-inventory-go/middlewares"
	"github.com/lumiere11/pc-inventory-go/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("file:test.db?cache=shared&_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&models.Product{}, &models.Category{}, &models.Status{}, &models.User{})
	categories := []models.Category{
		{Name: "Perifericos"},
		{Name: "Monitores"},
		{Name: "Gabinetes"},
	}
	ctx := context.Background()

	err = gorm.G[[]models.Category](db).Create(ctx, &categories) // pass pointer of data to Create
	statuses := []models.Status{
		{Name: "stock"},
		{Name: "sold out"},
	}

	err = gorm.G[[]models.Status](db).Create(ctx, &statuses) // pass pointer of data to Create

	productHandler := handlers.NewProductHandler(db)
	authHandler := handlers.NewAuthHandler(db)

	router := gin.Default()

	api := router.Group("/api/v1", middlewares.AuthMiddleware())
	{
		api.POST("/products", productHandler.CreateProduct)
		api.GET("/products/search", productHandler.GetByProperty)
	}
	apii := router.Group("/api/v1")
	{
		apii.POST("/register", authHandler.Register)
		apii.POST("/login", authHandler.Login)
	}
	router.Run() // listen and serve on 0.0.0.0:8080

}
