package rest

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine, handler Handler) {
	router.POST("/sign-up", checkBody(), checkJSONSignUp(), handler.signUp)
	router.POST("/sign-in", checkBody(), checkJSONSignIn(), handler.signIn)
	//todo: get files, upload file
}
