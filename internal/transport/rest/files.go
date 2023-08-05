package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) files(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "files are here"})
}
