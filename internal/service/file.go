package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
	"github.com/zenorachi/image-box/pkg/storage"
	"time"
)

var FilesNotFound = errors.New("files not found")

type FileRepository interface {
	Create(ctx *gin.Context, file models.File) error
	Get(ctx *gin.Context, userID uint) ([]models.File, error)
}

type Files struct {
	repository FileRepository
	provider   storage.Provider
}

func NewFiles(repository FileRepository, provider storage.Provider) *Files {
	return &Files{
		repository: repository,
		provider:   provider,
	}
}

func (f *Files) Upload(ctx *gin.Context, userID uint, input storage.UploadInput) error {
	url, err := f.provider.Upload(ctx, input)
	if err != nil {
		return err
	}

	file := models.CreateFile(userID, input.Name, url, input.Size, time.Now())
	return f.repository.Create(ctx, file)
}

func (f *Files) Get(ctx *gin.Context, userID uint) ([]models.File, error) {
	return f.repository.Get(ctx, userID)
}
