package db

import (
	"github.com/bosspokin/image-storer/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func InitNewGormStore() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&entity.File{}, &entity.User{})

	return db, nil
}
