package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	router.POST("/sign-up", h.CheckBody(), h.CheckJSONSignUp(), h.signUp)
	router.POST("/sign-in", h.CheckBody(), h.CheckJSONSignIn(), h.signIn)
	router.POST("/upload", h.CheckToken(), h.CheckUploadInput(), h.upload)
	router.GET("/refresh", h.refresh)
	router.GET("/files", h.CheckToken(), h.get)
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
