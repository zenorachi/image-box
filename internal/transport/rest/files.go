package rest

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/internal/transport/logger"
	"github.com/zenorachi/image-box/model"
	"github.com/zenorachi/image-box/pkg/storage"
	"github.com/zenorachi/url-shortener/pkg/base62"
	"net/http"
	"strconv"
	"strings"
)

const (
	uuidLen  = 7
	pageSize = 3
)

// @Summary Uploading a file to the storage
// @Tags storage
// @Description uploading a file to the minio storage
// @Security Auth
// @Accept  json
// @Produce  json
// @Param Authorization header string true "Bearer <token>" default("Bearer ")
// @Param file formData file true "File to upload"
// @Success 201 {object} string
// @Failure 400 {object} error
// @Failure 500 {object} string
// @Router /storage/upload [post]
func (h *handler) upload(ctx *gin.Context) {
	upload, _ := ctx.Get(uploadFile)
	uploadInput := upload.(storage.UploadInput)

	userIdCtx, _ := ctx.Get("userID")
	userID := userIdCtx.(uint)

	uploadInput.Name = base62.NewGenerator().GenerateCode(uuidLen) + uploadInput.Name
	if err := h.fileService.Upload(ctx.Request.Context(), userID, uploadInput); err != nil {
		logger.LogError(logger.UploadHandler, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "file uploaded successful"})
}

// @Summary Get list of files
// @Tags storage
// @Description get list of files for one user
// @Security Auth
// @Accept json
// @Produce  json
// @Param Authorization header string true "Bearer <token>" default("Bearer ")
// @Success 200 {object} []model.File
// @Success 204 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /storage/get [get]
func (h *handler) get(ctx *gin.Context) {
	userIdCtx, _ := ctx.Get("userID")
	userID := userIdCtx.(uint)

	files, err := h.fileService.Get(ctx.Request.Context(), userID)
	if files == nil {
		ctx.JSON(http.StatusNoContent, gin.H{"error": service.FilesNotFound})
		return
	}

	if err != nil {
		logger.LogError(logger.GetFilesHandler, err)
		if errors.Is(err, service.FilesNotFound) {
			ctx.JSON(http.StatusNoContent, gin.H{"error": err})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "error receiving files"})
		return
	}

	for i, file := range files {
		file.URL = strings.ReplaceAll(file.URL, "minio", "localhost")
		files[i] = file
	}

	pagedFiles, err := generatePagination(ctx, files)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unexpected error"})
		return
	}

	ctx.JSON(http.StatusOK, pagedFiles)
}

func generatePagination(ctx *gin.Context, files []model.File) ([]model.File, error) {
	page, err := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if err != nil {
		return nil, err
	}

	startIndex := (page - 1) * pageSize
	endIndex := startIndex + pageSize
	if endIndex > len(files) {
		endIndex = len(files)
	}

	return files[startIndex:endIndex], nil
}
