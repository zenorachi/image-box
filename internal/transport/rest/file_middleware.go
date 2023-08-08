package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/logger"
	"github.com/zenorachi/image-box/pkg/storage"
	"net/http"
)

func (h *handler) CheckUploadInput() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			logger.LogError(logger.FilesMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
