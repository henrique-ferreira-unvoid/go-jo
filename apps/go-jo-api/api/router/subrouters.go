package router

import (
	"github.com/gorilla/mux"
	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api/domain"
	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api/handlers"
)

// SubrouterBuilder contains handlers and configuration for building subrouters
type SubrouterBuilder struct {
	config              *domain.Config
	versionsHandler     *handlers.VersionsHandler
	integrationsHandler *handlers.IntegrationsHandler
	downloadHandler     *handlers.DownloadHandler
	healthHandler       *handlers.HealthHandler
}

// NewSubrouterBuilder creates a new subrouter builder
func NewSubrouterBuilder(config *domain.Config) *SubrouterBuilder {
	// Initialize handlers
	versionsHandler := handlers.NewVersionsHandler(config)
	integrationsHandler := handlers.NewIntegrationsHandler(config)
	downloadHandler := handlers.NewDownloadHandler(config, versionsHandler)
	healthHandler := handlers.NewHealthHandler(config)

	return &SubrouterBuilder{
		config:              config,
		versionsHandler:     versionsHandler,
		integrationsHandler: integrationsHandler,
		downloadHandler:     downloadHandler,
		healthHandler:       healthHandler,
	}
}

// BuildVersionsSubrouter builds the versions subrouter
func (sb *SubrouterBuilder) BuildVersionsSubrouter(router *mux.Router) {
	versionsRouter := router.PathPrefix("/versions").Subrouter()

	// GET /versions - Get all available tagged versions
	versionsRouter.HandleFunc("", sb.versionsHandler.AuthMiddleware(sb.versionsHandler.GetVersions)).Methods("GET")
}

// BuildIntegrationsSubrouter builds the integrations subrouter
func (sb *SubrouterBuilder) BuildIntegrationsSubrouter(router *mux.Router) {
	integrationsRouter := router.PathPrefix("/integrations").Subrouter()

	// GET /integrations - Get all available integrations
	integrationsRouter.HandleFunc("", sb.integrationsHandler.AuthMiddleware(sb.integrationsHandler.GetIntegrations)).Methods("GET")
}

// BuildDownloadSubrouter builds the download subrouter
func (sb *SubrouterBuilder) BuildDownloadSubrouter(router *mux.Router) {
	downloadRouter := router.PathPrefix("/download").Subrouter()

	// GET /download/{app_version}/{integration} - Download combined package
	downloadRouter.HandleFunc("/{app_version}/{integration}", sb.downloadHandler.AuthMiddleware(sb.downloadHandler.DownloadPackage)).Methods("GET")
}

// BuildHealthSubrouter builds the health check subrouter
func (sb *SubrouterBuilder) BuildHealthSubrouter(router *mux.Router) {
	healthRouter := router.PathPrefix("/health").Subrouter()

	// GET /health - Health check (no auth required)
	healthRouter.HandleFunc("", sb.healthHandler.HealthCheck).Methods("GET")
}

// GetHandlers returns the initialized handlers for external use if needed
func (sb *SubrouterBuilder) GetHandlers() (
	*handlers.VersionsHandler,
	*handlers.IntegrationsHandler,
	*handlers.DownloadHandler,
	*handlers.HealthHandler,
) {
	return sb.versionsHandler, sb.integrationsHandler, sb.downloadHandler, sb.healthHandler
}
