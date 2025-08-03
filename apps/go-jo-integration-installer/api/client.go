package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// Client represents an API client for go-jo-api
type Client struct {
	baseURL    string
	licenseKey string
	httpClient *http.Client
}

// APIResponse represents the structure of API responses
type APIResponse struct {
	Versions     []string `json:"versions,omitempty"`
	Integrations []string `json:"integrations,omitempty"`
	Error        string   `json:"error,omitempty"`
	Message      string   `json:"message,omitempty"`
}

// NewClient creates a new API client
func NewClient(baseURL, licenseKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		licenseKey: licenseKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GetVersions fetches available versions from the API
func (c *Client) GetVersions() ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/versions", c.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.licenseKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to parse as JSON first
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err == nil {
		if len(apiResp.Versions) > 0 {
			return apiResp.Versions, nil
		}
	}

	// Fallback: try to parse as simple array
	var versions []string
	if err := json.Unmarshal(body, &versions); err == nil {
		return versions, nil
	}

	return nil, fmt.Errorf("failed to parse versions response: %s", string(body))
}

// GetIntegrations fetches available integrations from the API
func (c *Client) GetIntegrations() ([]string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/integrations", c.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.licenseKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Try to parse as JSON first
	var apiResp APIResponse
	if err := json.Unmarshal(body, &apiResp); err == nil {
		if len(apiResp.Integrations) > 0 {
			return apiResp.Integrations, nil
		}
	}

	// Fallback: try to parse as simple array
	var integrations []string
	if err := json.Unmarshal(body, &integrations); err == nil {
		return integrations, nil
	}

	return nil, fmt.Errorf("failed to parse integrations response: %s", string(body))
}

// DownloadPackage downloads a package for the specified version and integration
func (c *Client) DownloadPackage(version, integration, outputPath string) error {
	url := fmt.Sprintf("%s/download/%s/%s", c.baseURL, version, integration)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", c.licenseKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("download failed with status %d: %s", resp.StatusCode, string(body))
	}

	// Create output file
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Copy response body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
