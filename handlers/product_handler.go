package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lumiere11/pc-inventory-go/models"
	"github.com/lumiere11/pc-inventory-go/requests"
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
	
	fmt.Printf("Search query: %s, Status: %s\n", q, status)

	// Build the query step by step
	query := h.DB.WithContext(ctx).Model(&models.Product{}).
		Preload("Category").
		Preload("Status")

	// Add search filters if query parameter is provided
	if q != "" {
		searchPattern := "%" + q + "%"
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(brand) LIKE LOWER(?) OR LOWER(model2) LIKE LOWER(?)", 
			searchPattern, searchPattern, searchPattern)
		fmt.Printf("Applying search filter with pattern: %s\n", searchPattern)
	}

	// Filter by status using a subquery to find status ID
	if status != "" {
		var statusModel models.Status
		if err := h.DB.Where("name = ?", status).First(&statusModel).Error; err != nil {
			fmt.Printf("Status not found: %s, error: %v\n", status, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid status parameter",
				"details": fmt.Sprintf("Status '%s' not found", status),
			})
			return
		}
		query = query.Where("status_id = ?", statusModel.ID)
	}

	// Execute the query
	if err := query.Find(&products).Error; err != nil {
		fmt.Printf("Database error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Found %d products\n", len(products))
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": products,
		"count": len(products),
	})
}
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var productReq requests.ProductRequest
	ctx := context.Background()

	if err := c.ShouldBindJSON(&productReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Convert ProductRequest to Product model with proper type conversion
	product := models.Product{
		Name:        productReq.Name,
		Brand:       productReq.Brand,
		Model2:      productReq.Model2,
		Description: productReq.Description,
	}

	// Convert string fields to appropriate types
	if stock, err := strconv.Atoi(productReq.Stock); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid stock value",
		})
		return
	} else {
		product.Stock = int32(stock)
	}

	if price, err := strconv.ParseFloat(productReq.Price, 32); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid price value",
		})
		return
	} else {
		product.Price = float32(price)
	}

	if statusID, err := strconv.ParseUint(productReq.StatusID, 10, 32); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid status_id value",
		})
		return
	} else {
		product.StatusID = uint(statusID)
	}

	if categoryID, err := strconv.ParseUint(productReq.CategoryID, 10, 32); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid category_id value",
		})
		return
	} else {
		product.CategoryID = uint(categoryID)
	}

	// Create the product in database
	result := h.DB.WithContext(ctx).Create(&product)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"data":    gin.H{},
			"message": result.Error.Error(),
		})
		return
	}

	// Load the created product with its relationships
	if err := h.DB.WithContext(ctx).Preload("Category").Preload("Status").First(&product, product.ID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found after creation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
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
		"message": "Product created successfully",
	})
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	var product models.Product
	ctx := context.Background()
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	result := h.DB.WithContext(ctx).Updates(&product)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"data":    gin.H{},
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    product,
		"message": "Product Updated",
	})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	var product models.Product
	ctx := context.Background()
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	result := h.DB.WithContext(ctx).Delete(&product)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"data":    gin.H{},
			"message": result.Error.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    product,
		"message": "Product Deleted",
	})
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	productId := c.Param("id")
	var stockRequest requests.UpdateProductRequest
	ctx := context.Background()

	// Bind del JSON request
	if err := c.ShouldBindJSON(&stockRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validar que el stock sea válido
	if stockRequest.Stock < 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Stock cannot be negative",
		})
		return
	}

	// Actualizar el stock en la base de datos
	result := h.DB.WithContext(ctx).Model(&models.Product{}).Where("id = ?", productId).Update("stock", stockRequest.Stock)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"data":    gin.H{},
			"message": result.Error.Error(),
		})
		return
	}

	// Verificar si se encontró el producto
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"data":    gin.H{},
			"message": "Product not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"data":    gin.H{},
		"message": "Product stock updated successfully",
	})
}
