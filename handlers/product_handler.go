package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lumiere11/pc-inventory-go/models"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

func NewProductHandler(db *gorm.DB) *ProductHandler {
	return &ProductHandler{
		DB: db,
	}
}
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var product models.Product
	ctx := context.Background()

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := h.DB.WithContext(ctx).Create(&product) // pass pointer of data to Create
	if result.Error != nil {
		c.JSON(200, gin.H{
			"status":  "error",
			"data":    gin.H{},
			"message": result.Error.Error(),
		})
		return
	}
	err := h.DB.Preload("Category").First(&product, product.ID).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(200, gin.H{
		"status": "success",
		"data": gin.H{
			"id":          product.ID,
			"name":        product.Name,
			"category_id": product.CategoryID,
			"brand":       product.Brand,
			"model":       product.Model2,
			"description": product.Description,
			"stock":       product.Stock,
			"price":       product.Price,
			"status":      product.Status,
			"category":    product.Category,
		},
		"message": "Product retrived succesfully",
	})
}
