package service

import (
	"context"
	"errors"
	"time"

	"github.com/zenorachi/image-box/model"
	"github.com/zenorachi/image-box/pkg/storage"
)

var FilesNotFound = errors.New("files not found")

type FileRepository interface {
	Create(ctx context.Context, file model.File) error
	Get(ctx context.Context, userID uint) ([]model.File, error)
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

func (f *Files) Upload(ctx context.Context, userID uint, input storage.UploadInput) error {
	url, err := f.provider.Upload(ctx, input)
	if err != nil {
		return err
	}

	file := model.CreateFile(userID, input.Name, url, input.Size, time.Now())
	return f.repository.Create(ctx, file)
}

func (f *Files) Get(ctx context.Context, userID uint) ([]model.File, error) {
	return f.repository.Get(ctx, userID)
}
