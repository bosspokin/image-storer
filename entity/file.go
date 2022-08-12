package entity

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Filename string
	URL      string `gorm:"unique"`
	Username string
}
