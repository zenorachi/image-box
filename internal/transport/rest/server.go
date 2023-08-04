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

	setupRoutes(router, mw, h)

	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf("%s:%d", host, port),
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
