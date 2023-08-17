package rest

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/service"
	"github.com/zenorachi/image-box/model"
	"github.com/zenorachi/image-box/pkg/storage"
)

const (
	inputSignUp = "inputSignUp"
	inputSignIn = "inputSignIn"
	requestBody = "requestBody"
	uploadFile  = "uploadInput"
)

type (
	User interface {
		SignUp(ctx context.Context, input model.SignUpInput) error
		SignIn(ctx context.Context, input model.SignInInput) (string, string, error)
		ParseToken(ctx context.Context, token string) (uint, error)
		RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
	}

	File interface {
		Upload(ctx context.Context, userID uint, input storage.UploadInput) error
		Get(ctx context.Context, userID uint) ([]model.File, error)
	}
)

type (
	AuthHandler interface {
		signUp(ctx *gin.Context)
		signIn(ctx *gin.Context)
		refresh(ctx *gin.Context)
	}

	AuthMiddleware interface {
		CheckBody() gin.HandlerFunc
		CheckJSONSignUp() gin.HandlerFunc
		CheckJSONSignIn() gin.HandlerFunc
		CheckToken() gin.HandlerFunc
	}

	FileHandler interface {
		upload(ctx *gin.Context)
		get(ctx *gin.Context)
	}

	FileMiddleware interface {
		CheckUploadInput() gin.HandlerFunc
	}
)

type Handler interface {
	AuthHandler
	AuthMiddleware
	FileHandler
	FileMiddleware
}

type handler struct {
	userService User
	fileService File
}

func NewHandler(users *service.Users, files *service.Files) Handler {
	return &handler{
		userService: users,
		fileService: files,
	}
}
