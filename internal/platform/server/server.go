package server

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"

	mooc "github.com/AguilaMike/go-gin-hexagonal/internal"
	"github.com/AguilaMike/go-gin-hexagonal/internal/platform/server/handler/courses"
	"github.com/AguilaMike/go-gin-hexagonal/internal/platform/server/handler/health"
)

// Server represents the API server.
type Server struct {
	httpAddr string
	engine   *gin.Engine

	// deps
	courseRepository mooc.CourseRepository
}

// New creates a new server.
func New(host string, port uint, courseRepository mooc.CourseRepository) Server {
	srv := Server{
		engine:           gin.New(),
		httpAddr:         fmt.Sprintf("%s:%d", host, port),
		courseRepository: courseRepository,
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

	// courses
	courses.Router(s.engine, s.courseRepository)
}
