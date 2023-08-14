package rest

import (
	"net/http"
	"slices"

	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/logger"
	"github.com/zenorachi/image-box/pkg/storage"
)

const (
	allowedFileSize = 10 * 1024 * 1024
	contentTypeJPEG = "image/jpeg"
	contentTypePNG  = "image/png"
	contentTypeJPG  = "image/jpg"
)

var allowedMimeTypes = []string{
	contentTypeJPEG,
	contentTypePNG,
	contentTypeJPG,
}

func (h *handler) CheckUploadInput() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			logger.LogError(logger.FilesMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if file.Size > allowedFileSize {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "file is too large (size exceeds 10MB limit)"})
			return
		}

		if !slices.Contains(allowedMimeTypes, file.Header.Get("Content-Type")) {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid file type (supported types: jpg, jpeg, png)"})
			return
		}

		uploadedFile, err := file.Open()
		if err != nil {
			logger.LogError(logger.FilesMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		defer uploadedFile.Close()

		uploadInput := storage.NewUploadInput(uploadedFile, file.Filename,
			file.Size, file.Header.Get("Content-Type"))

		ctx.Set(uploadFile, uploadInput)
		ctx.Next()
	}
}
