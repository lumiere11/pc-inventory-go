package models

import "gorm.io/gorm"

type Category struct {
	Name     string
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
	gorm.Model
}
