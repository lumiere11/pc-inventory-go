package database

import (
	"fmt"
	"os"

	"github.com/lumiere11/pc-inventory-go/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// InitDB initializes the database connection and runs migrations
func InitDB() (*gorm.DB, error) {
	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "127.0.0.1")
	dbPort := getEnv("DB_PORT", "3306")
	dbUser := getEnv("DB_USER", "root")
	dbPassword := getEnv("DB_PASSWORD", "")
	dbName := getEnv("DB_NAME", "pc_inventory")

	// Create DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&collation=utf8mb4_general_ci",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	// Open database connection
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Run migrations
	err = db.AutoMigrate(&models.Product{}, &models.Category{}, &models.Status{}, &models.User{})
	if err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	// Seed initial data
	err = SeedData(db)
	if err != nil {
		return nil, fmt.Errorf("failed to seed data: %w", err)
	}

	return db, nil
}

// getEnv gets environment variable with fallback to default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
