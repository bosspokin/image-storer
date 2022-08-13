package main

import (
	"github.com/bosspokin/image-storer/config"
	"github.com/bosspokin/image-storer/db"
	"github.com/bosspokin/image-storer/handler"
	"github.com/bosspokin/image-storer/middleware"
	"github.com/bosspokin/image-storer/route"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.InitNewGormStore()

	if err != nil {
		panic("failed to connect database")
	}

	r := gin.Default()

	store := cookie.NewStore([]byte(config.EnvSecretKey()))
	r.Use(sessions.Sessions("mysession", store))
	handler := handler.NewHandler(db)

	public := r.Group("")
	route.PublicRoutes(public)
	// public.POST("/signup", handler.SignUp)
	// public.GET("/login", handler.Login)

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
