package handlers

import (
	"net/http"
	"time"

	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api/domain"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	*BaseHandler
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(config *domain.Config) *HealthHandler {
	return &HealthHandler{
		BaseHandler: NewBaseHandler(config),
	}
}

// HealthCheck handles GET /health - Health check endpoint (no auth required)
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := domain.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "1.0.0",
	}
	h.SendJSONResponse(w, http.StatusOK, response)
}
