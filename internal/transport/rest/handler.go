package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
)

type User interface {
	SignUp(ctx *gin.Context, input models.SignUpInput) error
	//TODO: SignIn(ctx *gin.Context, input models.SignInInput)
}

type AuthHandler interface {
	signUp(ctx *gin.Context)
	signIn(ctx *gin.Context)
}

type Handler interface {
	AuthHandler
	//TODO: FileHandler
}

type handler struct{}

func New() Handler {
	return &handler{}
}
