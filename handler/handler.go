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

// @summary registers new user
// @id SignUp
// @accept json
// @param User body dto.User true "new user"
// @router /signup [post]
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

// @summary logs user in
// @id Login
// @accept json
// @param User body dto.User true "the user to be logged in"
// @router /login [get]
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

// @summary logs user out
// @id Logout
// @param username header string true "username of the user that wants to logout"
// @router /image [post]
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

// @summary uploads new image to server
// @id UploadImage
// @param username header string true "username of the user that wants to upload image"
// @param file formData file true "the file that the user wants to upload"
// @param filename formData string true "the filename of the corresponding file that the user wants to upload"
// @router /image [post]
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

// @summary renames existing image of a user
// @id RenameImage
// @accept json
// @param username header string true "username of the user that wants to rename his/her image"
// @param RenameFile body dto.RenameFile true "new file name"
// @router /image/:id [patch]
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

// @summary deletes existing image of a user
// @id DeleteImage
// @param username header string true "username of the user that wants to delete his/her image"
// @router /image/:id [delete]
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

// @summary lists images that belong to the user
// @id ListImages
// @param username header string true "username of the user that wants to list his/her image(s)"
// @router /images [get]
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
