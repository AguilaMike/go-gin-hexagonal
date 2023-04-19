package bootstrap

import "github.com/AguilaMike/go-gin-hexagonal/internal/platform/server"

const (
	host = "localhost"
	port = 8080
)

// Run starts the API server.
func Run() error {
	srv := server.New(host, port)
	return srv.Run()
}
