package router

import (
	"log"

	"github.com/gorilla/mux"
	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api/domain"
)

// Router manages the main application router
type Router struct {
	router           *mux.Router
	subrouterBuilder *SubrouterBuilder
}

// New creates a new router instance
func New(config *domain.Config) *Router {
	router := mux.NewRouter()
	subrouterBuilder := NewSubrouterBuilder(config)

	r := &Router{
		router:           router,
		subrouterBuilder: subrouterBuilder,
	}

	r.setupRoutes()
	return r
}

// setupRoutes configures all application routes
func (r *Router) setupRoutes() {
	log.Println("Setting up API routes...")

	// Build all subrouters
	r.subrouterBuilder.BuildVersionsSubrouter(r.router)
	r.subrouterBuilder.BuildIntegrationsSubrouter(r.router)
	r.subrouterBuilder.BuildDownloadSubrouter(r.router)
	r.subrouterBuilder.BuildHealthSubrouter(r.router)

	log.Println("API routes configured successfully")
}

// GetRouter returns the configured mux router
func (r *Router) GetRouter() *mux.Router {
	return r.router
}

// GetSubrouterBuilder returns the subrouter builder for accessing handlers
func (r *Router) GetSubrouterBuilder() *SubrouterBuilder {
	return r.subrouterBuilder
}

// LogRoutes logs all registered routes (useful for debugging)
func (r *Router) LogRoutes() {
	log.Println("Registered routes:")
	r.router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			methods, err := route.GetMethods()
			if err == nil {
				log.Printf("  %s %s", methods, pathTemplate)
			} else {
				log.Printf("  %s", pathTemplate)
			}
		}
		return nil
	})
}
