package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"not null;unique;size:255" json:"email"`
	Password string `gorm:"not null;size:40" json:"password"`
	Role     string `gorm:"not null;size:25" json:"role"`
}
