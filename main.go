package main

import (
	"github.com/bosspokin/image-storer/global"
	"github.com/bosspokin/image-storer/handler"
	"github.com/bosspokin/image-storer/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
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

	store := cookie.NewStore(global.Secret)
	r.Use(sessions.Sessions("mysession", store))
	handler := handler.NewHandler(db)

	r.POST("/signup", handler.SignUp)
	r.GET("/login", handler.Login)
	r.GET("/logout/:username", handler.Logout)

	// protected := r.Group("")
	// protected.Use()
	// protected.GET("/logout/:username", handler.Logout)

	r.Run()
}
