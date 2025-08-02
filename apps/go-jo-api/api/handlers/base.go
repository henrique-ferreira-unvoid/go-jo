package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/henrique-ferreira-unvoid/go-jo/apps/go-jo-api/api/domain"
)

// BaseHandler contains common functionality for all handlers
type BaseHandler struct {
	Config *domain.Config
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(config *domain.Config) *BaseHandler {
	return &BaseHandler{
		Config: config,
	}
}

// AuthMiddleware validates the authorization token
func (h *BaseHandler) AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			h.SendErrorResponse(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		if authHeader != h.Config.LicenseToken {
			h.SendErrorResponse(w, http.StatusUnauthorized, "Invalid authorization token")
			return
		}

		next(w, r)
	}
}

// SendJSONResponse sends a JSON response with the specified status and data
func (h *BaseHandler) SendJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// SendErrorResponse sends an error response
func (h *BaseHandler) SendErrorResponse(w http.ResponseWriter, status int, message string) {
	log.Printf("Error: %s", message)
	h.SendJSONResponse(w, status, domain.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}

// SendFileResponse sends a file response
func (h *BaseHandler) SendFileResponse(w http.ResponseWriter, filePath, filename string) {
	// Open and read the file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		h.SendErrorResponse(w, http.StatusInternalServerError, "Failed to read file: "+err.Error())
		return
	}

	// Get file info for size and mod time
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		h.SendErrorResponse(w, http.StatusInternalServerError, "Failed to get file info: "+err.Error())
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(fileData)))
	w.Header().Set("Last-Modified", fileInfo.ModTime().UTC().Format(http.TimeFormat))

	// Write the file data directly
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(fileData)
	if err != nil {
		log.Printf("Error writing file response: %v", err)
	}
}

// FetchFromGitHub makes authenticated requests to GitHub API
func (h *BaseHandler) FetchFromGitHub(url string, result interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), domain.REQUEST_TIMEOUT)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+h.Config.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}

	return nil
}
