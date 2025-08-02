package main

import (
	"log"

	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api"
)

func main() {
	// Create and start the API
	apiInstance := api.New()

	// Start the server (this blocks until the server stops)
	if err := apiInstance.Start(); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}
