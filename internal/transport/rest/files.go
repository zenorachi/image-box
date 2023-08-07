package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/internal/transport/rest/middleware"
	"github.com/zenorachi/image-box/pkg/storage"
	"net/http"
	"strings"
)

func (h *handler) upload(ctx *gin.Context) {
	upload, _ := ctx.Get(middleware.UploadInput)
	uploadInput := upload.(storage.UploadInput)

	userIdCtx, _ := ctx.Get("userID")
	userID := userIdCtx.(uint)

	if err := h.fileService.Upload(ctx, userID, uploadInput); err != nil {
		LogError(UploadHandler, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "file uploaded successful"})
}

func (h *handler) get(ctx *gin.Context) {
	userIdCtx, _ := ctx.Get("userID")
	userID := userIdCtx.(uint)

	files, err := h.fileService.Get(ctx, userID)
	if err != nil {
		LogError(GetFilesHandler, err)
		if errors.Is(err, service.FilesNotFound) {
			ctx.JSON(http.StatusNoContent, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error receiving files"})
		return
	}

	//TODO: DELETE THIS CRINGE
	for i, file := range files {
		file.URL = strings.ReplaceAll(file.URL, "minio", "localhost")
		files[i] = file
	}

	ctx.JSON(http.StatusOK, files)
}
