package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) signUp(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"mesSignUp": "ok"})
}

func (h *handler) signIn(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"mesSignIn": "ok"})
}
