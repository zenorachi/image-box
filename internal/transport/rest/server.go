package rest

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	handler Handler
}

func NewServer(h Handler) *Server {
	router := gin.Default()

	setupRoutes(router, h)

	return &Server{
		router:  router,
		handler: h,
	}
}

func (s *Server) Run(addr ...string) error {
	return s.router.Run(addr...)
}
