package handler

import (
	"net/http"
	"strconv"

	"github.com/bosspokin/image-storer/dto"
	"github.com/bosspokin/image-storer/entity"
	"github.com/bosspokin/image-storer/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

type Handler struct {
	imageService service.ImageService
	userService  service.UserService
}

func NewHandler(imageService service.ImageService, userService service.UserService) *Handler {
	return &Handler{
		imageService: imageService,
		userService:  userService,
	}
}

func (handler *Handler) SignUp(ctx *gin.Context) {
	var user dto.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	userRec := entity.User{}
	copier.Copy(&userRec, &user)

	err := handler.userService.SignUp(userRec)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{})
}

func (handler *Handler) Login(ctx *gin.Context) {
	var req dto.User

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	entityUser := entity.User{}
	copier.Copy(&entityUser, &req)

	err := handler.userService.Login(ctx, entityUser)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
}

func (handler *Handler) Logout(ctx *gin.Context) {
	username := ctx.Request.Header[http.CanonicalHeaderKey("username")][0]

	err := handler.userService.Logout(ctx, username)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (handler *Handler) UploadImage(ctx *gin.Context) {
	formfile, _, err := ctx.Request.FormFile("file")
	username := ctx.Request.Header[http.CanonicalHeaderKey("username")][0]

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	filename := ctx.PostForm("filename")
	file := entity.File{
		Filename: filename,
		File:     formfile,
		Username: username,
	}

	uploadUrl, err := handler.imageService.UploadImage(username, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": uploadUrl,
	})
}

func (handler *Handler) RenameImage(ctx *gin.Context) {
	renameReq := dto.RenameFile{}
	idParam := ctx.Param("id")
	username := ctx.Request.Header[http.CanonicalHeaderKey("username")][0]
	id, err := strconv.Atoi(idParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})

		return
	}

	if err := ctx.ShouldBindJSON(&renameReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "bad request",
		})

		return
	}

	newUrl, err := handler.imageService.RenameImage(username, uint(id), renameReq.New)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"url": newUrl,
	})
}

func (handler *Handler) DeleteImage(ctx *gin.Context) {
	idParam := ctx.Param("id")
	username := ctx.Request.Header[http.CanonicalHeaderKey("username")][0]
	id, err := strconv.Atoi(idParam)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	err = handler.imageService.DeleteImage(username, uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *Handler) ListImages(ctx *gin.Context) {
	username := ctx.Request.Header[http.CanonicalHeaderKey("username")][0]

	files, err := handler.imageService.ListImages(username)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	filesRes := make([]dto.File, len(files))

	if err := copier.Copy(&filesRes, &files); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, filesRes)
}
