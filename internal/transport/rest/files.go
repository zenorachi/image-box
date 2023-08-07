package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/pkg/storage"
	"log"
	"net/http"
)

func (h *handler) upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		log.Println("upload handler")
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "file is required"})
		return
	}

	uploadedFile, err := file.Open()
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
	}
	defer uploadedFile.Close()

	uploadInput := storage.NewUploadInput(uploadedFile, file.Filename,
		file.Size, file.Header.Get("Content-Type"))

	userIdCtx, ok := ctx.Get("userID")
	if !ok {
		log.Println("user id not found")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "user id not found"})
	}

	userID, ok := userIdCtx.(uint)
	if !ok {
		log.Println("user id invalid type")
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "user id not found"})
		return
	}

	if err = h.fileService.Upload(ctx, userID, uploadInput); err != nil {
		log.Println("upload handler upload failed", err)
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "user id not found"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "file uploaded successful"})
}
