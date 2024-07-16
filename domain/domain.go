package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email     string `gorm:"unique;not null"`
	Password  string `gorm:"unique;not null"`
	IsAdmin   bool   `gorm:"default:false"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
