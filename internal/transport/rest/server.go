package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swagFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	server *http.Server
}

func NewServer(h Handler, host string, port int) *Server {
	router := gin.New()

	server := &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: router,
		},
	}

	server.setupRoutes(router, h)

	return server
}

func (s *Server) setupRoutes(router *gin.Engine, h Handler) {
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})
	router.POST("/auth/sign-up", h.CheckBody(), h.CheckJSONSignUp(), h.signUp)
	router.POST("/auth/sign-in", h.CheckBody(), h.CheckJSONSignIn(), h.signIn)
	router.GET("/auth/refresh", h.refresh)
	router.POST("/storage/upload", h.CheckToken(), h.CheckUploadInput(), h.upload)
	router.GET("/storage/files", h.CheckToken(), h.get)
	router.GET("/docs/*any", ginSwagger.WrapHandler(swagFiles.Handler))
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
