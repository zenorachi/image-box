package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/pkg/storage"
	"net/http"
	"strings"
)

func (h *handler) upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		LogError(UploadHandler, err)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		LogError(UploadHandler, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	defer uploadedFile.Close()

	uploadInput := storage.NewUploadInput(uploadedFile, file.Filename,
		file.Size, file.Header.Get("Content-Type"))

	userIdCtx, ok := ctx.Get("userID")
	if !ok {
		LogError(UploadHandler, errors.New("user id not found"))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user id not found"})
		return
	}

	userID, ok := userIdCtx.(uint)
	if !ok {
		LogError(UploadHandler, errors.New("user id invalid type"))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "user id not found"})
		return
	}

	if err = h.fileService.Upload(ctx, userID, uploadInput); err != nil {
		LogError(UploadHandler, err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "user id not found"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "file uploaded successful"})
}

func (h *handler) get(ctx *gin.Context) {
	userIdCtx, ok := ctx.Get("userID")
	if !ok {
		LogError(GetFilesHandler, errors.New("user id not found"))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "user id not found"})
		return
	}
	userID, ok := userIdCtx.(uint)
	if !ok {
		LogError(GetFilesHandler, errors.New("user id invalid type"))
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "user id not found"})
		return
	}
	files, err := h.fileService.Get(ctx, userID)
	if err != nil {
		LogError(GetFilesHandler, err)
		if errors.Is(err, service.FilesNotFound) {
			ctx.JSON(http.StatusNoContent, gin.H{"message": err})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "error receiving files"})
		return
	}

	//TODO: DELETE THIS CRINGE
	for i, file := range files {
		file.URL = strings.ReplaceAll(file.URL, "minio", "localhost")
		files[i] = file
	}

	ctx.JSON(http.StatusOK, files)
}
