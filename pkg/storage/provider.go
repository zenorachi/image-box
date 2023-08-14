package storage

import (
	"io"

	"github.com/gin-gonic/gin"
)

type UploadInput struct {
	File        io.Reader
	Name        string
	Size        int64
	ContentType string
}

type Provider interface {
	Upload(ctx *gin.Context, input UploadInput) (string, error)
}

func NewUploadInput(file io.Reader, name string, size int64, contentType string) UploadInput {
	return UploadInput{
		File:        file,
		Name:        name,
		Size:        size,
		ContentType: contentType,
	}
}
