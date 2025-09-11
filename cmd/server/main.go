package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lumiere11/pc-inventory-go/handlers"
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
	db.AutoMigrate(&models.Product{}, &models.Category{})
	categories := []models.Category{
		{Name: "Perifericos"},
		{Name: "Monitores"},
		{Name: "Gabinetes"},
	}
	ctx := context.Background()

	err = gorm.G[[]models.Category](db).Create(ctx, &categories) // pass pointer of data to Create

	productHandler := handlers.NewProductHandler(db)

	router := gin.Default()

	api := router.Group("/api/v1")
	{
		api.POST("/products", productHandler.CreateProduct)
	}
	router.Run() // listen and serve on 0.0.0.0:8080

}
