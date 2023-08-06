package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
)

type User interface {
	SignUp(ctx *gin.Context, input models.SignUpInput) error
	SignIn(ctx *gin.Context, input models.SignInInput) (string, string, error)
	ParseToken(ctx *gin.Context, token string) (uint, error)
	RefreshTokens(ctx *gin.Context, refreshToken string) (string, string, error)
}

type (
	AuthHandler interface {
		signUp(ctx *gin.Context)
		signIn(ctx *gin.Context)
		refresh(ctx *gin.Context)
	}

	FileHandler interface {
		files(ctx *gin.Context)
	}
)

type Handler interface {
	CheckToken() gin.HandlerFunc
	AuthHandler
	FileHandler
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
