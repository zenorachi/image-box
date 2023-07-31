package handlers

import "github.com/gin-gonic/gin"

type AuthHandler interface {
	SignUp(ctx *gin.Context)
	SignIn(ctx *gin.Context)
}

type Handler interface {
	AuthHandler
	//TODO: FileHandler
}

type handler struct{}

func New() Handler {
	return &handler{}
}
