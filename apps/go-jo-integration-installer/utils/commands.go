package utils

import (
	"fmt"
	"os"
	"strings"
)

// ParseLicenseFlag parses the --license flag from command line arguments
func ParseLicenseFlag() (string, error) {
	if len(os.Args) < 1 {
		return "", fmt.Errorf("usage: %s --license=<path-to-license-file>", os.Args[0])
	}

	for _, arg := range os.Args {
		// Check for --license=<path> format
		if strings.HasPrefix(arg, "--license=") {
			licensePath := strings.TrimPrefix(arg, "--license=")
			if licensePath == "" {
				return "", fmt.Errorf("--license flag requires a file path")
			}
			return licensePath, nil
		}
	}

	// Check for --license <path> format (separate arguments)
	for i, arg := range os.Args {
		if arg == "--license" {
			if i+1 >= len(os.Args) {
				return "", fmt.Errorf("--license flag requires a file path")
			}
			return os.Args[i+1], nil
		}
	}

	return "", fmt.Errorf("--license flag is required")
}

// ParseVersionFlag parses the optional --version flag from command line arguments.
// Returns an empty string when the flag is not provided.
func ParseVersionFlag() (string, error) {
	for _, arg := range os.Args {
		// Check for --version=<tag> format
		if strings.HasPrefix(arg, "--version=") {
			version := strings.TrimPrefix(arg, "--version=")
			if version == "" {
				return "", fmt.Errorf("--version flag requires a tag")
			}
			return version, nil
		}
	}

	// Check for --version <tag> format (separate arguments)
	for i, arg := range os.Args {
		if arg == "--version" {
			if i+1 >= len(os.Args) {
				return "", fmt.Errorf("--version flag requires a tag")
			}

			version := os.Args[i+1]
			if strings.TrimSpace(version) == "" || strings.HasPrefix(version, "--") {
				return "", fmt.Errorf("--version flag requires a tag")
			}

			return version, nil
		}
	}

	return "", nil
}
