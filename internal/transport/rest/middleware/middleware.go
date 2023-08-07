package middleware

import "github.com/gin-gonic/gin"

const (
	InputSignUp = "inputSignUp"
	InputSignIn = "inputSignIn"
	requestBody = "requestBody"
	UploadInput = "uploadInput"
)

type (
	AuthMiddleware interface {
		CheckBody() gin.HandlerFunc
		CheckJSONSignUp() gin.HandlerFunc
		CheckJSONSignIn() gin.HandlerFunc
	}

	FilesMiddleware interface {
		CheckUploadInput() gin.HandlerFunc
	}

	Middleware interface {
		AuthMiddleware
		FilesMiddleware
	}
)

type middleware struct{}

func NewMiddleware() Middleware {
	return &middleware{}
}
