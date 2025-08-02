package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api/domain"
	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api/router"
	"github.com/joho/godotenv"
)

// API represents the main API application
type API struct {
	config *domain.Config
	router *router.Router
	server *http.Server
}

// New creates a new API instance
func New() *API {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	config := loadConfig()

	// Validate required configuration
	if config.LicenseToken == "" {
		log.Fatal("LICENSE_TOKEN environment variable is required")
	}
	if config.GitHubToken == "" {
		log.Fatal("GITHUB_TOKEN environment variable is required")
	}

	// Create router
	apiRouter := router.New(config)

	// Create server
	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      apiRouter.GetRouter(),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	return &API{
		config: config,
		router: apiRouter,
		server: server,
	}
}

// Start starts the API server
func (a *API) Start() error {
	log.Printf("Starting go-jo-api on port %s", a.config.Port)
	log.Printf("Endpoints available:")
	log.Printf("  GET /versions")
	log.Printf("  GET /integrations")
	log.Printf("  GET /download/{app_version}/{integration}")
	log.Printf("  GET /health")

	// Optionally log all routes for debugging
	a.router.LogRoutes()

	return a.server.ListenAndServe()
}

// Stop gracefully stops the API server
func (a *API) Stop() error {
	log.Println("Stopping go-jo-api...")
	return a.server.Close()
}

// GetConfig returns the API configuration
func (a *API) GetConfig() *domain.Config {
	return a.config
}

// GetRouter returns the API router
func (a *API) GetRouter() *router.Router {
	return a.router
}

// loadConfig loads configuration from environment variables
func loadConfig() *domain.Config {
	return &domain.Config{
		Port:         getEnvOrDefault("PORT", domain.DEFAULT_PORT),
		GitHubToken:  os.Getenv("GITHUB_TOKEN"),
		LicenseToken: os.Getenv("LICENSE_TOKEN"),
	}
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
