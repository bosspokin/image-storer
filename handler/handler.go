package handler

import (
	"net/http"

	"github.com/bosspokin/image-storer/model"
	"github.com/gin-gonic/gin"
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

	if result := handler.db.Create(&user); result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (handler *Handler) Login(ctx *gin.Context) {

}

func (handler *Handler) Logout(ctx *gin.Context) {

}

func (handler *Handler) UploadImage(ctx *gin.Context) {

}

func (handler *Handler) EditImageName(ctx *gin.Context) {

}

func (handler *Handler) DeleteImage(ctx *gin.Context) {

}

func (handler *Handler) ListImages(ctx *gin.Context) {

}
