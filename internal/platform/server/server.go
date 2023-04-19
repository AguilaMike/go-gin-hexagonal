package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/AguilaMike/go-gin-hexagonal/internal/platform/server/handler/health"
)

// Server represents the API server.
type Server struct {
	httpAddr string
	engine   *gin.Engine
}

// New creates a new server.
func New(host string, port uint) Server {
	srv := Server{
		engine:   gin.New(),
		httpAddr: fmt.Sprintf("%s:%d", host, port),
	}

	srv.registerRoutes()
	return srv
}

func (s *Server) Run() error {
	log.Println("Server running on", s.httpAddr)
	return s.engine.Run(s.httpAddr)
}

func (s *Server) registerRoutes() {
	s.engine.GET("/health", health.CheckHandler())
}
