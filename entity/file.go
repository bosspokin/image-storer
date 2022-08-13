package entity

import (
	"mime/multipart"

	"gorm.io/gorm"
)

type File struct {
	gorm.Model
	Filename string
	URL      string `gorm:"unique"`
	Username string
	File     multipart.File `gorm:"-"`
}
