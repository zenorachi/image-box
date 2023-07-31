package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/zenorachi/image-box/internal/transport/rest/handlers"
)

type Server struct {
	router  *gin.Engine
	handler handlers.Handler
}

func NewServer() *Server {
	router := gin.Default()
	handler := handlers.New()

	setupRoutes(router, handler)

	return &Server{
		router:  router,
		handler: handler,
	}
}

func (s *Server) Run(addr ...string) error {
	return s.router.Run(addr...)
}
