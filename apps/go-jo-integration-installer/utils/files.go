package utils

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ExtractZip extracts a zip file to the specified directory
func ExtractZip(zipPath, destDir string) error {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer reader.Close()

	for _, file := range reader.File {
		filePath := filepath.Join(destDir, file.Name)

		if file.FileInfo().IsDir() {
			os.MkdirAll(filePath, file.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// FindDeployDirectory finds the directory containing docker-compose.yml and Makefile
func FindDeployDirectory(baseDir string) (string, error) {
	// Look for docker-compose.yml and Makefile in the extracted directory
	var deployDir string
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			// Check if this directory contains docker-compose.yml and Makefile
			dockerComposePath := filepath.Join(path, "docker-compose.yml")
			makefilePath := filepath.Join(path, "Makefile")

			if _, err := os.Stat(dockerComposePath); err == nil {
				if _, err := os.Stat(makefilePath); err == nil {
					deployDir = path
					return filepath.SkipAll // Found it, stop walking
				}
			}
		}

		return nil
	})

	if deployDir == "" {
		return "", fmt.Errorf("could not find directory with docker-compose.yml and Makefile")
	}

	return deployDir, err
}

// ReadLicenseKey reads the license key from the specified file
func ReadLicenseKey(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read license file: %w", err)
	}

	licenseKey := string(content)
	licenseKey = strings.TrimSpace(licenseKey)
	if licenseKey == "" {
		return "", fmt.Errorf("license file is empty")
	}

	return licenseKey, nil
}
