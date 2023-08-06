package storage

import (
	"github.com/gin-gonic/gin"
	"io"
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
