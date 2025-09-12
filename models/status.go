package models

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	Name string `gorm:"not null;unique;size:100" json:"name"`
}
