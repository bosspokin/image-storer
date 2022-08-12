package handler

import (
	"fmt"
	"net/http"

	"github.com/bosspokin/image-storer/entity"
	"github.com/bosspokin/image-storer/helper"
	"github.com/bosspokin/image-storer/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

type Handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{db: db}
}

func (handler *Handler) SignUp(ctx *gin.Context) {
	var user model.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	userRec := entity.User{}
	copier.Copy(&userRec, &user)

	if result := handler.db.Create(&userRec); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (handler *Handler) Login(ctx *gin.Context) {
	var user model.User
	var userRecord entity.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})
		return
	}

	if result := handler.db.Where("username = ?", user.Username).First(&userRecord); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if user.Password != userRecord.Password {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "incorrect password",
		})
		return
	}

	session := sessions.Default(ctx)
	session.Set(userRecord.Username, userRecord.Username)

	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
}

func (handler *Handler) Logout(ctx *gin.Context) {
	username := ctx.Request.Header[http.CanonicalHeaderKey("username")][0]
	session := sessions.Default(ctx)

	// user := session.Get(username)

	// if user == nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"error": "User is not logged in",
	// 	})
	// }

	session.Delete(username)
	if err := session.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot logout",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (handler *Handler) UploadImage(ctx *gin.Context) {
	formfile, _, err := ctx.Request.FormFile("file")
	filename := ctx.PostForm("filename")
	file := model.File{
		Filename: filename,
		File:     formfile,
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	uploadUrl, err := helper.UploadImage(file)
	if err != nil {
		fmt.Println("error of helper")
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	username := ctx.Request.Header[http.CanonicalHeaderKey("username")][0]

	fileRecord := entity.File{
		Filename: filename,
		URL:      uploadUrl,
		Username: username,
	}

	if result := handler.db.Create(&fileRecord); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": uploadUrl,
	})
}

func (handler *Handler) ChangeImageName(ctx *gin.Context) {

}

func (handler *Handler) DeleteImage(ctx *gin.Context) {

}

func (handler *Handler) ListImages(ctx *gin.Context) {
	username := ctx.Request.Header[http.CanonicalHeaderKey("username")][0]
	var files []entity.File

	if result := handler.db.Where("username = ?", username).Find(&files); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})

		return
	}

	filesRes := make([]model.File, len(files))

	if err := copier.Copy(&filesRes, &files); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, filesRes)
}
