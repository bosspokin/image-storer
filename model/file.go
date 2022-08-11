package model

import "gorm.io/gorm"

type File struct {
	gorm.Model
	Filename string `json:"filename"`
	Username string `json:"username"`
}
