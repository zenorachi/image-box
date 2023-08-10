package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/models"
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
		SignUp(ctx *gin.Context, input models.SignUpInput) error
		SignIn(ctx *gin.Context, input models.SignInInput) (string, string, error)
		ParseToken(ctx *gin.Context, token string) (uint, error)
		RefreshTokens(ctx *gin.Context, refreshToken string) (string, string, error)
	}

	File interface {
		Upload(ctx *gin.Context, userID uint, input storage.UploadInput) error
		Get(ctx *gin.Context, userID uint) ([]models.File, error)
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

func NewHandler(users User, files File) Handler {
	return &handler{
		userService: users,
		fileService: files,
	}
}
