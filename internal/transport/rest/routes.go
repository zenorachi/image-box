package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/rest/middleware"
)

func setupRoutes(router *gin.Engine, mw middleware.Middleware, handler Handler) {
	router.POST("/sign-up", mw.CheckBody(), mw.CheckJSONSignUp(), handler.signUp)
	router.POST("/sign-in", mw.CheckBody(), mw.CheckJSONSignIn(), handler.signIn)
	//todo: get files, upload file
}
