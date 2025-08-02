package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"

	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api/domain"
)

// IntegrationsHandler handles integration-related requests
type IntegrationsHandler struct {
	*BaseHandler
}

// NewIntegrationsHandler creates a new integrations handler
func NewIntegrationsHandler(config *domain.Config) *IntegrationsHandler {
	return &IntegrationsHandler{
		BaseHandler: NewBaseHandler(config),
	}
}

// GetIntegrations handles GET /integrations - Get all branches from go-jo-docker-environments
func (h *IntegrationsHandler) GetIntegrations(w http.ResponseWriter, r *http.Request) {
	log.Printf("Fetching integrations (branches) for repository: %s", domain.DOCKER_ENV_REPO)

	branches, err := h.fetchGitHubBranches(domain.DOCKER_ENV_REPO)
	if err != nil {
		h.SendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch branches: "+err.Error())
		return
	}

	var integrations []string
	for _, branch := range branches {
		// Filter out main/master branches if you only want integration branches
		if branch.Name != "main" && branch.Name != "master" {
			integrations = append(integrations, branch.Name)
		}
	}

	sort.Strings(integrations)

	response := domain.IntegrationsResponse{Integrations: integrations}
	h.SendJSONResponse(w, http.StatusOK, response)
}

// fetchGitHubBranches fetches branches from GitHub API
func (h *IntegrationsHandler) fetchGitHubBranches(repo string) ([]domain.GitHubBranch, error) {
	url := fmt.Sprintf("%s/repos/%s/branches", domain.GITHUB_API_BASE_URL, repo)

	var branches []domain.GitHubBranch
	err := h.FetchFromGitHub(url, &branches)
	return branches, err
}
