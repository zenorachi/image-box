package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
)

type User interface {
	SignUp(ctx *gin.Context, input models.SignUpInput) error
	SignIn(ctx *gin.Context, input models.SignInInput) (string, error)
}

type AuthHandler interface {
	signUp(ctx *gin.Context)
	signIn(ctx *gin.Context)
}

type Handler interface {
	AuthHandler
	//TODO: FileHandler
}

type handler struct {
	userService User
	//todo fileService File
}

func NewHandler(users User) Handler {
	return &handler{
		userService: users,
	}
}
