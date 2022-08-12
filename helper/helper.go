package helper

import (
	"context"
	"time"

	"github.com/bosspokin/image-storer/config"
	"github.com/bosspokin/image-storer/model"
	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploadImage(file model.File) (string, error) {
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

func GetImage(input interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cld, err := cloudinary.NewFromParams(config.EnvCloudName(), config.EnvCloudAPIKey(), config.EnvCloudAPISecret())
	if err != nil {
		return "", err
	}

	uploadParam, err := cld.Upload.Upload(ctx, input, uploader.UploadParams{Folder: config.EnvCloudUploadFolder()})
	if err != nil {
		return "", err
	}

	return uploadParam.SecureURL, nil
}
