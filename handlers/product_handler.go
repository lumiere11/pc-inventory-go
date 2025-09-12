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
func (h *ProductHandler) GetByProperty(c *gin.Context) {
	var products []models.Product
	ctx := context.Background()
	q := c.Query("q")
	status := c.Query("status")
	if status == "" {
		status = "stock"
	}
	if err := h.DB.WithContext(ctx).
		Joins("JOIN categories ON categories.id = products.category_id").
		Joins("JOIN statuses ON statuses.id = products.status_id").
		Where("(products.name LIKE ? OR products.brand LIKE ? OR categories.name LIKE ?)", "%"+q+"%", "%"+q+"%", "%"+q+"%").
		Where("statuses.name = ?", status).
		Preload("Category").
		Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
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
	if err := h.DB.Preload("Category").First(&product, product.ID).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product not found"})
		return

	}
	if err := h.DB.Preload("Status").First(&product, product.ID).Error; err != nil {
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
