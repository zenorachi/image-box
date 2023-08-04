package middleware

import "github.com/gin-gonic/gin"

const (
	InputSignUp = "inputSignUp"
	InputSignIn = "inputSignIn"
	requestBody = "requestBody"
)

type (
	AuthMiddleware interface {
		CheckJSONSignUp() gin.HandlerFunc
		CheckJSONSignIn() gin.HandlerFunc
	}

	FilesMiddleware interface {
		// TODO
	}

	Middleware interface {
		CheckBody() gin.HandlerFunc
		AuthMiddleware
		FilesMiddleware
	}
)

type middleware struct{}

func NewMiddleware() Middleware {
	return &middleware{}
}
