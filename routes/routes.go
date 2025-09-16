package routes

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/lumiere11/pc-inventory-go/handlers"
	"github.com/lumiere11/pc-inventory-go/middlewares"
	"gorm.io/gorm"
)

// SetupRoutes configures all the application routes
func SetupRoutes(db *gorm.DB) *gin.Engine {
	// Initialize handlers
	productHandler := handlers.NewProductHandler(db)
	authHandler := handlers.NewAuthHandler(db)
	
	// Initialize Casbin enforcer
	enforcer, err := casbin.NewEnforcer("model.conf", "policy.csv")
	if err != nil {
		panic("Failed to initialize Casbin enforcer: " + err.Error())
	}

	// Initialize Gin router
	router := gin.Default()

	// Protected routes (require authentication and authorization)
	api := router.Group("/api/v1", middlewares.AuthMiddleware(), middlewares.CasbinMiddleware(enforcer))
	{
		api.POST("/products", productHandler.CreateProduct)
		api.GET("/products/search", productHandler.GetByProperty)
		api.PUT("/products/:id", productHandler.UpdateProduct)
		api.PUT("/products/:id/stock", productHandler.UpdateStock)
		api.DELETE("/products/:id", productHandler.Delete)
	}

	// Public routes (no authentication required)
	publicAPI := router.Group("/api/v1")
	{
		publicAPI.POST("/register", authHandler.Register)
		publicAPI.POST("/login", authHandler.Login)
	}

	return router
}
