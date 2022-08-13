package service

import (
	"fmt"

	"github.com/bosspokin/image-storer/entity"
	"github.com/bosspokin/image-storer/helper"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService interface {
	SignUp(user entity.User) error
	Login(ctx *gin.Context, user entity.User) error
	Logout(ctx *gin.Context, username string) error
}

type userService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) UserService {
	return &userService{
		db: db,
	}
}

func (service *userService) SignUp(user entity.User) error {
	hash, err := helper.HashPassword(user.Password)

	if err != nil {
		return err
	}

	user.Password = hash

	result := service.db.Create(&user)
	return result.Error
}

func (service *userService) Login(ctx *gin.Context, user entity.User) error {
	var userRecord entity.User

	if result := service.db.Where("username = ?", user.Username).First(&userRecord); result.Error != nil {
		return result.Error
	}

	if !helper.CheckPasswordHash(user.Password, userRecord.Password) {
		return fmt.Errorf("incorrect password")
	}

	session := sessions.Default(ctx)
	session.Set(userRecord.Username, userRecord.Username)

	err := session.Save()
	return err
}

func (service *userService) Logout(ctx *gin.Context, username string) error {
	session := sessions.Default(ctx)

	session.Delete(username)
	err := session.Save()
	return err
}
