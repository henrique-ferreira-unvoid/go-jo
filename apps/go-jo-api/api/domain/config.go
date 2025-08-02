package domain

import "time"

// Configuration
type Config struct {
	Port         string
	GitHubToken  string
	LicenseToken string
}

// Constants
const (
	DEFAULT_PORT        = "1207"
	GITHUB_API_BASE_URL = "https://api.github.com"
	GO_JO_REPO          = "henrique-ferreira-unvoid/go-jo"
	DOCKER_ENV_REPO     = "henrique-ferreira-unvoid/go-jo-docker-environments"
	TEMP_DIR_PREFIX     = "go-jo-api-"
	REQUEST_TIMEOUT     = 30 * time.Second
)

// Response structures
type VersionResponse struct {
	Versions []string `json:"versions"`
}

type IntegrationsResponse struct {
	Integrations []string `json:"integrations"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// GitHub API structures
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Draft   bool   `json:"draft"`
}

type GitHubBranch struct {
	Name string `json:"name"`
}

type GitHubAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type GitHubReleaseWithAssets struct {
	TagName string        `json:"tag_name"`
	Assets  []GitHubAsset `json:"assets"`
}
