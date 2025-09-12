package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Name        string   `gorm:"not null;unique;size:100" json:"name"`
	Brand       string   `json:"brand"`
	Model2      string   `json:"model"`
	Description string   `gorm:"not null;size:255" json:"description"`
	Stock       int32    `json:"stock"`
	Price       float32  `json:"price"`
	StatusID    uint     `json:"status_id"`
	Status      Status   `gorm:"foreignKey:StatusID" json:"status"`
	CategoryID  uint     `gorm:"not null" json:"category_id"`
	Category    Category `gorm:"foreignKey:CategoryID" json:"category"`
}
