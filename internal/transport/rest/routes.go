package rest

import (
	"github.com/gin-gonic/gin"
)

func setupRoutes(router *gin.Engine, handler Handler) {
	router.POST("/sign-up", checkJSONBody(), handler.signUp)
	//router.POST("/sign-in", handler.signIn)
	//todo: get files, upload file
}
