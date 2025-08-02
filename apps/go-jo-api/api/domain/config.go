package domain

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Build-time variables (injected by GoReleaser)
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

// Configuration loaded from config.yaml and environment
type Config struct {
	GitHubToken  string
	LicenseToken string

	// Loaded from config.yaml
	API     APIConfig     `mapstructure:"api"`
	GitHub  GitHubConfig  `mapstructure:"github"`
	License LicenseConfig `mapstructure:"license"`
	Server  ServerConfig  `mapstructure:"server"`
	App     AppConfig     `mapstructure:"app"`
}

type APIConfig struct {
	DefaultPort    string `mapstructure:"port"`
	RequestTimeout string `mapstructure:"request_timeout"`
	TempDirPrefix  string `mapstructure:"temp_dir_prefix"`
}

type GitHubConfig struct {
	APIBaseURL   string             `mapstructure:"api_base_url"`
	Token        string             `mapstructure:"token"`
	Repositories RepositoriesConfig `mapstructure:"repositories"`
}

type LicenseConfig struct {
	Token string `mapstructure:"token"`
}

type RepositoriesConfig struct {
	GoJo               string `mapstructure:"go_jo"`
	DockerEnvironments string `mapstructure:"docker_environments"`
}

type ServerConfig struct {
	ReadTimeout  time.Duration `mapstructure:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout"`
}

type AppConfig struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

// Legacy constants for backward compatibility
// These are now available through config.API.*, config.GitHub.*, etc.
func (c *Config) GetDefaultPort() string {
	return c.API.DefaultPort
}

func (c *Config) GetGitHubAPIBaseURL() string {
	return c.GitHub.APIBaseURL
}

func (c *Config) GetGoJoRepo() string {
	return c.GitHub.Repositories.GoJo
}

func (c *Config) GetDockerEnvRepo() string {
	return c.GitHub.Repositories.DockerEnvironments
}

func (c *Config) GetRequestTimeout() time.Duration {
	timeout, err := time.ParseDuration(c.API.RequestTimeout)
	if err != nil {
		return 30 * time.Second
	}
	return timeout
}

func (c *Config) GetTempDirPrefix() string {
	return c.API.TempDirPrefix
}

// LoadConfig loads configuration from config.yaml and environment variables
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Add config paths for different environments
	viper.AddConfigPath("/etc/go-jo-api")                     // Production path
	viper.AddConfigPath("./apps/go-jo-api/deployment/config") // Development path
	viper.AddConfigPath(".")                                  // Current directory

	// Set environment variable mappings
	viper.SetEnvPrefix("GOJO")
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("api.port", "1207")
	viper.SetDefault("api.request_timeout", "30s")
	viper.SetDefault("api.temp_dir_prefix", "go-jo-api-")
	viper.SetDefault("github.api_base_url", "https://api.github.com")
	viper.SetDefault("github.token", "your-github-token-here")
	viper.SetDefault("github.repositories.go_jo", "henrique-ferreira-unvoid/go-jo")
	viper.SetDefault("github.repositories.docker_environments", "henrique-ferreira-unvoid/go-jo-docker-environments")
	viper.SetDefault("server.read_timeout", "15s")
	viper.SetDefault("server.write_timeout", "15s")
	viper.SetDefault("license.token", "your-license-token-here")

	// Read config file
	configFileUsed := ""
	if err := viper.ReadInConfig(); err != nil {
		// Config file not found is not an error, use defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		fmt.Println("No config file found, using defaults")
	} else {
		configFileUsed = viper.ConfigFileUsed()
		fmt.Printf("Using config file: %s\n", configFileUsed)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Load environment variables (these override config file values)
	config.GitHubToken = getEnvOrDefault("GITHUB_TOKEN", config.GitHub.Token)
	config.LicenseToken = getEnvOrDefault("LICENSE_TOKEN", config.License.Token)

	return &config, nil
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

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
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
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
