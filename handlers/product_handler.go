package handlers

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

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

// levenshtein computes the Levenshtein distance between two UTF-8 strings.
func levenshtein(a, b string) int {
	if a == b {
		return 0
	}
	// Convert to runes to handle multi-byte characters
	ra, rb := []rune(a), []rune(b)
	if len(ra) == 0 {
		return len(rb)
	}
	if len(rb) == 0 {
		return len(ra)
	}

	// Ensure ra is the shorter slice
	if len(ra) > len(rb) {
		ra, rb = rb, ra
	}

	// Initialize the previous row
	prevRow := make([]int, len(ra)+1)
	for i := range prevRow {
		prevRow[i] = i
	}

	// Compute current row based on the previous row
	for i := 1; i <= len(rb); i++ {
		currRow := make([]int, len(ra)+1)
		currRow[0] = i
		for j := 1; j <= len(ra); j++ {
			cost := 1
			if rb[i-1] == ra[j-1] {
				cost = 0
			}
			currRow[j] = min(
				currRow[j-1]+1,    // Deletion
				prevRow[j]+1,      // Insertion
				prevRow[j-1]+cost, // Substitution
			)
		}
		prevRow = currRow
	}

	return prevRow[len(ra)]
}

// min is a helper function to find the minimum of three integers
func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
	} else {
		if b < c {
			return b
		}
	}
	return c
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
		query = query.Where("LOWER(name) LIKE LOWER(?) OR LOWER(brand) LIKE LOWER(?) OR LOWER(model2) LIKE LOWER(?) OR LOWER(description) LIKE LOWER(?)",
			searchPattern, searchPattern, searchPattern, searchPattern)
		fmt.Printf("Applying search filter with pattern: %s\n", searchPattern)
	}

	// Filter by status using a subquery to find status ID
	if status != "" {
		var statusModel models.Status
		if err := h.DB.Where("name = ?", status).First(&statusModel).Error; err != nil {
			fmt.Printf("Status not found: %s, error: %v\n", status, err)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid status parameter",
				"details": fmt.Sprintf("Status '%s' not found", status),
			})
			return
		}
		fmt.Printf("Found status: %s with ID: %d\n", statusModel.Name, statusModel.ID)
		query = query.Where("status_id = ?", statusModel.ID)
	}

	// Execute the query
	if err := query.Find(&products).Error; err != nil {
		fmt.Printf("Database error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("Query executed successfully. Found %d products before fuzzy search\n", len(products))

	// If no results and we have a query, attempt fuzzy matching fallback
	if len(products) == 0 && q != "" {
		fmt.Println("No exact matches found. Attempting fuzzy search fallback...")
		var candidates []models.Product

		// Reuse base query but without the text filter, keeping status and preloads
		baseQuery := h.DB.WithContext(ctx).Model(&models.Product{}).
			Preload("Category").
			Preload("Status")
		if status != "" {
			var statusModel models.Status
			if err := h.DB.Where("name = ?", status).First(&statusModel).Error; err == nil {
				baseQuery = baseQuery.Where("status_id = ?", statusModel.ID)
			}
		}

		if err := baseQuery.Find(&candidates).Error; err != nil {
			fmt.Printf("Fuzzy fallback DB error: %v\n", err)
		} else {
			lowerQ := strings.ToLower(q)

			type scoredProduct struct {
				Product models.Product
				Score   int
			}
			var results []scoredProduct

			// Dynamic threshold based on query length
			maxDist := 1
			if len(lowerQ) >= 5 {
				maxDist = 2
			}

			for _, p := range candidates {
				// Compute distance against several fields and take the minimum
				fields := []string{
					strings.ToLower(p.Name),
					strings.ToLower(p.Brand),
					strings.ToLower(p.Model2),
					strings.ToLower(p.Description),
				}
				bestScore := 1<<31 - 1
				for _, field := range fields {
					dist := levenshtein(lowerQ, field)
					if dist < bestScore {
						bestScore = dist
					}
					if bestScore == 0 { // Early exit on perfect match
						break
					}
				}

				if bestScore <= maxDist {
					results = append(results, scoredProduct{Product: p, Score: bestScore})
				}
			}

			// Sort results by score (best match first)
			sort.Slice(results, func(i, j int) bool {
				return results[i].Score < results[j].Score
			})

			// Limit to top 20 fuzzy matches
			limit := 20
			if len(results) < limit {
				limit = len(results)
			}

			// Clear the original products slice and populate with sorted, scored results
			products = products[:0]
			for i := 0; i < limit; i++ {
				products = append(products, results[i].Product)
			}

			if len(products) > 0 {
				fmt.Printf("Fuzzy search found %d matches\n", len(products))
			}
		}
	}

	fmt.Printf("Found %d products\n", len(products))
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   products,
		"count":  len(products),
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

	product.StatusID = 1

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
