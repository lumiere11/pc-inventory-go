package database

import (
	"github.com/lumiere11/pc-inventory-go/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes the database connection and runs migrations
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("file:test.db?cache=shared&_foreign_keys=on"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Run migrations
	err = db.AutoMigrate(&models.Product{}, &models.Category{}, &models.Status{}, &models.User{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
