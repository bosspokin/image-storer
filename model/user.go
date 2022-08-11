package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	Username  string `gorm:"primaryKey"`
	Password  string
	Files     []File `gorm:"foreignKey:Username"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
