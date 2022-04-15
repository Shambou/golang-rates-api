package main

import (
	"log"
)

// Run - sets up our application
func Run() error {
	log.Println("Setting the app")

	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatal("Error starting up REST API")
	}
}
