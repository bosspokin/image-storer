package service

import (
	"fmt"

	"github.com/bosspokin/image-storer/entity"
	"github.com/bosspokin/image-storer/helper"
	"gorm.io/gorm"
)

type ImageService interface {
	ListImages(username string) ([]entity.File, error)
	UploadImage(username string, file entity.File) (string, error)
	RenameImage(username string, id uint, newName string) (string, error)
	DeleteImage(username string, id uint) error
}

type imageService struct {
	db *gorm.DB
}

func NewImageService(db *gorm.DB) ImageService {
	return &imageService{
		db: db,
	}
}

func (service *imageService) ListImages(username string) ([]entity.File, error) {
	var files []entity.File

	if result := service.db.Where("username = ?", username).Find(&files); result.Error != nil {
		return nil, result.Error
	}

	return files, nil
}

func (service *imageService) UploadImage(username string, file entity.File) (string, error) {
	uploadUrl, err := helper.UploadImage(file)
	if err != nil {
		return "", err
	}

	file.URL = uploadUrl

	if result := service.db.Create(&file); result.Error != nil {
		return "", err
	}

	return uploadUrl, nil
}

func (service *imageService) RenameImage(username string, id uint, newName string) (string, error) {
	var file entity.File

	if result := service.db.First(&file, id); result.Error != nil {
		return "", result.Error
	}

	if file.Username != username {
		return "", fmt.Errorf("the file %s does not belong to the user %s", file.Filename, username)
	}

	uploadUrl, err := helper.RenameImage(file.Filename, newName)
	if err != nil {
		return "", err
	}

	file.URL = uploadUrl
	file.Filename = newName

	if result := service.db.Save(&file); result.Error != nil {
		return "", err
	}

	return uploadUrl, nil
}

func (service *imageService) DeleteImage(username string, id uint) error {
	var file entity.File
	if result := service.db.First(&file, id); result.Error != nil {
		return result.Error
	}

	if file.Username != username {
		return fmt.Errorf("the file %s does not belong to the user %s", file.Filename, username)
	}

	err := helper.DeleteImage(file.Filename)
	if err != nil {
		return err
	}

	result := service.db.Unscoped().Delete(&file)
	return result.Error
}
