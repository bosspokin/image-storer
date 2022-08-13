package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/bosspokin/image-storer/config"
	"github.com/bosspokin/image-storer/dto"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"golang.org/x/crypto/bcrypt"
)

func UploadImage(file dto.File) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cld, err := cloudinary.NewFromParams(config.EnvCloudName(), config.EnvCloudAPIKey(), config.EnvCloudAPISecret())
	if err != nil {
		return "", err
	}

	uploadParam, err := cld.Upload.Upload(ctx, file.File, uploader.UploadParams{
		Folder:      config.EnvCloudUploadFolder(),
		UseFilename: true,
		PublicID:    file.Filename,
	})

	if err != nil {
		return "", err
	}

	return uploadParam.SecureURL, nil
}

func RenameImage(oldPublicID, newPublicID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cld, err := cloudinary.NewFromParams(config.EnvCloudName(), config.EnvCloudAPIKey(), config.EnvCloudAPISecret())
	if err != nil {
		return "", err
	}

	uploadParam, err := cld.Upload.Rename(ctx, uploader.RenameParams{
		FromPublicID: fmt.Sprintf("files/%s", oldPublicID),
		ToPublicID:   fmt.Sprintf("files/%s", newPublicID),
	})

	if err != nil {
		return "", err
	}

	return uploadParam.SecureURL, nil
}

func DeleteImage(publicID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cld, err := cloudinary.NewFromParams(config.EnvCloudName(), config.EnvCloudAPIKey(), config.EnvCloudAPISecret())
	if err != nil {
		return err
	}

	_, err = cld.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: fmt.Sprintf("files/%s", publicID),
	})

	return err
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}
