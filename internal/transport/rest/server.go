package rest

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	handler Handler
}

func NewServer() *Server {
	router := gin.Default()
	handler := New()

	setupRoutes(router, handler)

	return &Server{
		router:  router,
		handler: handler,
	}
}

func (s *Server) Run(addr ...string) error {
	return s.router.Run(addr...)
}
