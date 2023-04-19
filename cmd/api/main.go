package main

import (
	"log"

	"github.com/AguilaMike/go-gin-hexagonal/cmd/api/bootstrap"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
