package main

import (
	"log"

	"github.com/Shambou/golang-challenge/internal/server"
)

// Run - sets up our application
func Run() error {
	log.Println("Setting the app")

	handler := server.New()
	if err := handler.Serve(); err != nil {
		log.Fatal("failed to gracefully serve our application")
	}

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal("Error starting up REST API")
	}
}
