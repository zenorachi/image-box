package rest

import (
	"github.com/gin-gonic/gin"
)

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	r := gin.Default()

	return &Server{
		router: r,
	}
}

func (s *Server) Run(addr ...string) error {
	return s.router.Run(addr...)
}
