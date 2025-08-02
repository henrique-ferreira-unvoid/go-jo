package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api/domain"
)

// VersionsHandler handles version-related requests
type VersionsHandler struct {
	*BaseHandler
}

// NewVersionsHandler creates a new versions handler
func NewVersionsHandler(config *domain.Config) *VersionsHandler {
	return &VersionsHandler{
		BaseHandler: NewBaseHandler(config),
	}
}

// GetVersions handles GET /versions - Get all available tagged versions of go-jo
func (h *VersionsHandler) GetVersions(w http.ResponseWriter, r *http.Request) {
	log.Printf("Fetching versions for repository: %s", domain.GO_JO_REPO)

	releases, err := h.fetchGitHubReleases(domain.GO_JO_REPO)
	if err != nil {
		h.SendErrorResponse(w, http.StatusInternalServerError, "Failed to fetch releases: "+err.Error())
		return
	}

	var versions []string
	for _, release := range releases {
		if !release.Draft {
			versions = append(versions, release.TagName)
		}
	}

	// Sort versions (newest first)
	sort.Slice(versions, func(i, j int) bool {
		return h.compareVersions(versions[i], versions[j]) > 0
	})

	response := domain.VersionResponse{Versions: versions}
	h.SendJSONResponse(w, http.StatusOK, response)
}

// GetLatestVersion returns the latest non-draft version
func (h *VersionsHandler) GetLatestVersion() (string, error) {
	releases, err := h.fetchGitHubReleases(domain.GO_JO_REPO)
	if err != nil {
		return "", err
	}

	for _, release := range releases {
		if !release.Draft {
			return release.TagName, nil
		}
	}

	return "", fmt.Errorf("no releases found")
}

// fetchGitHubReleases fetches releases from GitHub API
func (h *VersionsHandler) fetchGitHubReleases(repo string) ([]domain.GitHubRelease, error) {
	url := fmt.Sprintf("%s/repos/%s/releases", domain.GITHUB_API_BASE_URL, repo)

	var releases []domain.GitHubRelease
	err := h.FetchFromGitHub(url, &releases)
	return releases, err
}

// compareVersions provides basic version comparison
func (h *VersionsHandler) compareVersions(v1, v2 string) int {
	// Simple version comparison (for proper semver, use a library)
	// Remove 'v' prefix and compare
	clean1 := strings.TrimPrefix(v1, "v")
	clean2 := strings.TrimPrefix(v2, "v")

	if clean1 == clean2 {
		return 0
	}
	if clean1 > clean2 {
		return 1
	}
	return -1
}
