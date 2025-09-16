package main

import (
	"log"

	"github.com/lumiere11/pc-inventory-go/database"
	"github.com/lumiere11/pc-inventory-go/routes"
)

func main() {
	// Initialize database
	db, err := database.InitDB()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Seed initial data
	if err := database.SeedData(db); err != nil {
		log.Fatal("Failed to seed database:", err)
	}

	// Setup routes
	router := routes.SetupRoutes(db)

	// Start server
	log.Println("Starting server on :8081...")
	if err := router.Run(":8081"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
