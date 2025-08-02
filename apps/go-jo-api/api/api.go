package api

import (
	"log"
	"net/http"

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
	// Load environment variables (if there is a .env file)
	_ = godotenv.Load()

	config, err := domain.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

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
		Addr:         ":" + config.API.DefaultPort,
		Handler:      apiRouter.GetRouter(),
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}

	return &API{
		config: config,
		router: apiRouter,
		server: server,
	}
}

// Start starts the API server
func (a *API) Start() error {
	log.Printf("Starting %s v%s on port %s", a.config.App.Name, domain.Version, a.config.API.DefaultPort)
	log.Printf("Build info: commit=%s, date=%s", domain.GitCommit, domain.BuildDate)
	log.Printf("Configuration loaded from: %s", "config.yaml")
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
