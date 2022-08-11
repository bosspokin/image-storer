package main

import (
	"github.com/bosspokin/image-storer/handler"
	"github.com/bosspokin/image-storer/model"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&model.File{}, &model.User{})

	r := gin.Default()

	handler := handler.NewHandler(db)

	r.POST("/signup", handler.SignUp)

	r.Run()
}
