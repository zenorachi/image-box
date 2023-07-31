package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/rest/handlers"
)

func setupRoutes(router *gin.Engine, handler handlers.Handler) {
	router.POST("/sign-up", handler.SignUp)
	router.POST("/sign-in", handler.SignIn)
	//todo: get files, upload file
}
