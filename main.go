package main

import (
	"github.com/bosspokin/image-storer/entity"
	"github.com/bosspokin/image-storer/global"
	"github.com/bosspokin/image-storer/handler"
	"github.com/bosspokin/image-storer/middleware"
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

	db.AutoMigrate(&entity.File{}, &entity.User{})

	r := gin.Default()

	store := cookie.NewStore(global.Secret)
	r.Use(sessions.Sessions("mysession", store))
	handler := handler.NewHandler(db)

	r.POST("/signup", handler.SignUp)
	r.GET("/login", handler.Login)

	protected := r.Group("")
	protected.Use(middleware.Auth)
	protected.GET("/logout", handler.Logout)
	protected.GET("/images", handler.ListImages)
	protected.POST("/upload", handler.UploadImage)
	protected.PATCH("/rename", handler.RenameImage)
	protected.DELETE("/image/:id", handler.DeleteImage)

	// protected := r.Group("")
	// protected.Use()
	// protected.GET("/logout/:username", handler.Logout)

	r.Run()
}
