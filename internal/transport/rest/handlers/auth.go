package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *handler) SignUp(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"mesSignUp": "ok"})
}

func (h *handler) SignIn(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"mesSignIn": "ok"})
}
