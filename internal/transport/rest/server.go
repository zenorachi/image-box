package rest

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/rest/middleware"
	"net/http"
)

type Server struct {
	server *http.Server
}

func NewServer(h Handler, host string, port int) *Server {
	//TODO: gin.New() to create custom logger (logrus)
	//router := gin.New()
	router := gin.Default()
	mw := middleware.NewMiddleware()

	server := &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: router,
		},
	}

	server.setupRoutes(router, mw, h)

	return server
}

func (s *Server) setupRoutes(router *gin.Engine, mw middleware.Middleware, handler Handler) {
	router.POST("/sign-up", mw.CheckBody(), mw.CheckJSONSignUp(), handler.signUp)
	router.POST("/sign-in", mw.CheckBody(), mw.CheckJSONSignIn(), handler.signIn)
	router.POST("/upload", handler.CheckToken(), handler.upload)
	router.GET("/refresh", handler.refresh)
	router.GET("/files", handler.CheckToken(), handler.get)
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
