package handlers

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api/domain"
)

// DownloadHandler handles download-related requests
type DownloadHandler struct {
	*BaseHandler
	versionsHandler *VersionsHandler
}

// NewDownloadHandler creates a new download handler
func NewDownloadHandler(config *domain.Config, versionsHandler *VersionsHandler) *DownloadHandler {
	return &DownloadHandler{
		BaseHandler:     NewBaseHandler(config),
		versionsHandler: versionsHandler,
	}
}

// DownloadPackage handles GET /download/{app_version}/{integration} - Download combined package
func (h *DownloadHandler) DownloadPackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	appVersion := vars["app_version"]
	integration := vars["integration"]

	log.Printf("Download request: version=%s, integration=%s", appVersion, integration)

	// Handle "latest" version
	if appVersion == "latest" {
		latest, err := h.versionsHandler.GetLatestVersion()
		if err != nil {
			h.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get latest version: "+err.Error())
			return
		}
		appVersion = latest
		log.Printf("Resolved 'latest' to version: %s", appVersion)
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", domain.TEMP_DIR_PREFIX)
	if err != nil {
		h.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create temp directory: "+err.Error())
		return
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Download app .deb file
	debPath, err := h.downloadAppDeb(appVersion, tempDir)
	if err != nil {
		h.SendErrorResponse(w, http.StatusInternalServerError, "Failed to download app: "+err.Error())
		return
	}

	// Download integration branch as zip
	integrationZipPath, err := h.downloadIntegrationZip(integration, tempDir)
	if err != nil {
		h.SendErrorResponse(w, http.StatusInternalServerError, "Failed to download integration: "+err.Error())
		return
	}

	// Create combined zip
	combinedZipPath, err := h.createCombinedZip(debPath, integrationZipPath, integration, tempDir)
	if err != nil {
		h.SendErrorResponse(w, http.StatusInternalServerError, "Failed to create combined package: "+err.Error())
		return
	}

	// Send file
	h.SendFileResponse(w, combinedZipPath, fmt.Sprintf("go-jo-%s.zip", integration))
}

// downloadAppDeb downloads the .deb file for a specific version
func (h *DownloadHandler) downloadAppDeb(version, tempDir string) (string, error) {
	// Fetch release with assets
	url := fmt.Sprintf("%s/repos/%s/releases/tags/%s", domain.GITHUB_API_BASE_URL, domain.GO_JO_REPO, version)

	var release domain.GitHubReleaseWithAssets
	err := h.FetchFromGitHub(url, &release)
	if err != nil {
		return "", err
	}

	// Find .deb asset
	var debURL string
	for _, asset := range release.Assets {
		if strings.HasSuffix(asset.Name, ".deb") && strings.Contains(asset.Name, "go-jo") {
			debURL = asset.BrowserDownloadURL
			break
		}
	}

	if debURL == "" {
		return "", fmt.Errorf("no .deb file found in release %s", version)
	}

	// Download .deb file
	debPath := filepath.Join(tempDir, "go-jo-selected.deb")
	return debPath, h.downloadFile(debURL, debPath)
}

// downloadIntegrationZip downloads the integration branch as a zip file
func (h *DownloadHandler) downloadIntegrationZip(integration, tempDir string) (string, error) {
	// GitHub archive URL for branch
	url := fmt.Sprintf("https://github.com/%s/archive/refs/heads/%s.zip", domain.DOCKER_ENV_REPO, integration)

	zipPath := filepath.Join(tempDir, "integration.zip")
	return zipPath, h.downloadFile(url, zipPath)
}

// downloadFile downloads a file from a URL
func (h *DownloadHandler) downloadFile(url, filepath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), domain.REQUEST_TIMEOUT)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+h.Config.GitHubToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

// createCombinedZip creates a combined zip file with the .deb and integration files
func (h *DownloadHandler) createCombinedZip(debPath, integrationZipPath, integration, tempDir string) (string, error) {
	combinedZipPath := filepath.Join(tempDir, fmt.Sprintf("go-jo-%s.zip", integration))

	// Create new zip file
	combinedZip, err := os.Create(combinedZipPath)
	if err != nil {
		return "", err
	}
	defer combinedZip.Close()

	zipWriter := zip.NewWriter(combinedZip)
	defer zipWriter.Close()

	// Add the .deb file as "app-selected.deb"
	if err := h.addFileToZip(zipWriter, debPath, "app-selected.deb"); err != nil {
		return "", err
	}

	// Extract and add integration zip contents
	if err := h.addZipContentsToZip(zipWriter, integrationZipPath); err != nil {
		return "", err
	}

	return combinedZipPath, nil
}

// addFileToZip adds a file to a zip archive
func (h *DownloadHandler) addFileToZip(zipWriter *zip.Writer, filePath, nameInZip string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	header.Name = nameInZip

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

// addZipContentsToZip extracts content from one zip and adds it to another
func (h *DownloadHandler) addZipContentsToZip(zipWriter *zip.Writer, zipPath string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Skip root directory (usually repository-name-branch/)
		parts := strings.Split(file.Name, "/")
		if len(parts) > 1 {
			// Remove first part (root dir) and rejoin
			newName := strings.Join(parts[1:], "/")
			if newName != "" {
				if err := h.copyZipFile(zipWriter, file, newName); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// copyZipFile copies a file from one zip to another
func (h *DownloadHandler) copyZipFile(zipWriter *zip.Writer, file *zip.File, newName string) error {
	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer fileReader.Close()

	header := file.FileHeader
	header.Name = newName

	writer, err := zipWriter.CreateHeader(&header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, fileReader)
	return err
}
