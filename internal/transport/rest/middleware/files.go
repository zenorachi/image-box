package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/rest"
	"github.com/zenorachi/image-box/pkg/storage"
	"net/http"
)

func (mw *middleware) CheckUploadInput() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			rest.LogError(rest.FilesMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		uploadedFile, err := file.Open()
		if err != nil {
			rest.LogError(rest.FilesMiddleware, err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		defer uploadedFile.Close()

		uploadInput := storage.NewUploadInput(uploadedFile, file.Filename,
			file.Size, file.Header.Get("Content-Type"))

		ctx.Set(UploadInput, uploadInput)
		ctx.Next()
	}
}
