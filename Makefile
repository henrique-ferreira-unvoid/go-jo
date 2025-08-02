# go-jo Monorepo Makefile
.PHONY: help gojo-build gojo-dev gojo-package clean api-build api-dev api-package

# Default target
help:
	@echo "Available targets:"
	@echo ""
	@echo "Go-jo app:"
	@echo "  gojo-build        - Build go-jo binary for local development"
	@echo "  gojo-package      - Create go-jo .deb package (requires git tag or uses fake tag)"
	@echo "  gojo-dev          - Run go-jo in development mode"
	@echo ""
	@echo "Go-jo API:"
	@echo "  api-build         - Build go-jo-api binary for local development"
	@echo "  api-package       - Create go-jo-api .deb package for testing"
	@echo "  api-dev           - Run go-jo-api in development mode"
	@echo ""
	@echo "General:"
	@echo "  clean             - Clean build artifacts"
	@echo "  build-all         - Build all applications"
	@echo "  package-all       - Package all applications"
	@echo "  help              - Show this help message"

# Build go-jo binary for local development
gojo-build:
	@echo "Building go-jo binary..."
	@mkdir -p dist/go-jo
	@go build -o dist/go-jo/go-jo ./apps/go-jo
	@echo "Binary built at: dist/go-jo/go-jo"

# Create go-jo .deb package (for testing packaging)
gojo-package:
	@echo "Creating go-jo .deb package..."
	@if ! git describe --tags --exact-match 2>/dev/null; then \
		echo "No git tag found, creating temporary tag for packaging..."; \
		git tag -a v0.1.0-test -m "Temporary tag for testing packaging" 2>/dev/null || true; \
		goreleaser release --snapshot --clean; \
		git tag -d v0.1.0-test 2>/dev/null || true; \
	else \
		goreleaser release --snapshot --clean; \
	fi
	@echo "Package created in dist/ directory"
	@find dist/ -name "*.deb" -type f | head -5

# Run go-jo in development mode
gojo-dev:
	@echo "Running go-jo in development mode..."
	@go run ./apps/go-jo

# Build go-jo-api binary for local development
api-build:
	@echo "Building go-jo-api binary..."
	@mkdir -p dist/api-bin
	@go build -o dist/api-bin/go-jo-api ./apps/go-jo-api
	@echo "Binary built at: dist/api-bin/go-jo-api"

# Create go-jo-api .deb package using GoReleaser
api-package:
	@echo "Creating go-jo-api .deb package using GoReleaser..."
	@if ! git describe --tags --exact-match 2>/dev/null; then \
		echo "No git tag found, creating temporary tag for packaging..."; \
		git tag -a v0.1.0-test -m "Temporary tag for testing packaging" 2>/dev/null || true; \
		goreleaser release --config .goreleaser-api.yml --snapshot --clean --skip=publish; \
		git tag -d v0.1.0-test 2>/dev/null || true; \
	else \
		goreleaser release --config .goreleaser-api.yml --snapshot --clean --skip=publish; \
	fi
	@echo "API package created in dist/ directory"
	@find dist/ -name "*api*.deb" -type f | head -5

# Run go-jo-api in development mode
api-dev:
	@echo "Running go-jo-api in development mode..."
	@echo "Make sure you have a .env file with GITHUB_TOKEN and LICENSE_TOKEN"
	@go run ./apps/go-jo-api

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf dist/
	@echo "Build artifacts cleaned"

# Build all applications
build-all: gojo-build api-build

# Package all applications  
package-all: gojo-package api-package

# Release everything (for CI/CD)
release-all: gojo-package 