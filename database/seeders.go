package database

import (
	"context"

	"github.com/lumiere11/pc-inventory-go/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedData populates the database with initial data
func SeedData(db *gorm.DB) error {
	ctx := context.Background()

	// Seed categories
	categories := []models.Category{
		{Name: "Perifericos"},
		{Name: "Monitores"},
		{Name: "Gabinetes"},
		{Name: "Procesadores"},
		{Name: "Tarjetas Graficas"},
		{Name: "Memoria RAM"},
		{Name: "Placas Madre"},
		{Name: "Almacenamiento"},
		{Name: "Fuentes de Poder"},
		{Name: "Refrigeracion"},
		{Name: "Tarjetas de Sonido"},
		{Name: "Tarjetas de Red"},
		{Name: "Lectores Opticos"},
		{Name: "Cables y Conectores"},
		{Name: "Ventiladores"},
	}

	// Check if categories already exist
	var count int64
	db.Model(&models.Category{}).Count(&count)
	if count == 0 {
		if err := db.WithContext(ctx).Create(&categories).Error; err != nil {
			return err
		}
	}

	// Seed statuses
	statuses := []models.Status{
		{Name: "stock"},
		{Name: "sold out"},
	}

	// Check if statuses already exist
	db.Model(&models.Status{}).Count(&count)
	if count == 0 {
		if err := db.WithContext(ctx).Create(&statuses).Error; err != nil {
			return err
		}
	}

	// Seed admin user
	var userCount int64
	db.Model(&models.User{}).Where("email = ?", "admin@admin.com").Count(&userCount)
	if userCount == 0 {
		hash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := models.User{
			Email:    "admin@admin.com",
			Password: string(hash),
			Role:     "admin",
		}

		if err := db.WithContext(ctx).Create(&user).Error; err != nil {
			return err
		}
	}

	// Seed some sample products for testing
	var productCount int64
	db.Model(&models.Product{}).Count(&productCount)
	if productCount == 0 {
		// Get category and status IDs
		var peripheralsCategory models.Category
		var stockStatus models.Status
		
		db.Where("name = ?", "Perifericos").First(&peripheralsCategory)
		db.Where("name = ?", "stock").First(&stockStatus)
		
		if peripheralsCategory.ID > 0 && stockStatus.ID > 0 {
			sampleProducts := []models.Product{
				{
					Name:        "Razer DeathAdder V3 Mouse",
					Brand:       "Razer",
					Model2:      "RZ01-04910100-R3U1",
					Description: "Gaming mouse with ergonomic design",
					Stock:       15,
					Price:       89.99,
					StatusID:    stockStatus.ID,
					CategoryID:  peripheralsCategory.ID,
				},
				{
					Name:        "Razer BlackWidow V4",
					Brand:       "Razer",
					Model2:      "RZ03-04860100-R3U1",
					Description: "Mechanical gaming keyboard",
					Stock:       8,
					Price:       199.99,
					StatusID:    stockStatus.ID,
					CategoryID:  peripheralsCategory.ID,
				},
				{
					Name:        "Logitech G502 Mouse",
					Brand:       "Logitech",
					Model2:      "910-005550",
					Description: "High performance gaming mouse",
					Stock:       12,
					Price:       79.99,
					StatusID:    stockStatus.ID,
					CategoryID:  peripheralsCategory.ID,
				},
				{
					Name:        "Corsair M65 RGB Elite Mouse",
					Brand:       "Corsair",
					Model2:      "CH-9309011-NA",
					Description: "FPS gaming mouse with sniper button",
					Stock:       10,
					Price:       59.99,
					StatusID:    stockStatus.ID,
					CategoryID:  peripheralsCategory.ID,
				},
			}
			
			if err := db.WithContext(ctx).Create(&sampleProducts).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
