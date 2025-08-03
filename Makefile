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

# Build targets
gojo-build:
	@echo "Building go-jo..."
	@mkdir -p dist/go-jo
	go build -o dist/go-jo/go-jo ./apps/go-jo

api-build:
	@echo "Building go-jo-api..."
	@mkdir -p dist/api-bin
	go build -o dist/api-bin/go-jo-api ./apps/go-jo-api

installer-build:
	@echo "Building go-jo-integration-installer..."
	@mkdir -p dist/installer-bin
	go build -o dist/installer-bin/go-jo-integration-installer ./apps/go-jo-integration-installer

build-all: gojo-build api-build installer-build

# Package targets
gojo-package:
	@echo "Creating go-jo package..."
	@if [ -z "$$(git tag --points-at HEAD)" ]; then \
		echo "No git tag found, creating temporary tag for packaging..."; \
		TEMP_TAG="temp-package-$$(date +%s)"; \
		git tag $$TEMP_TAG; \
		goreleaser release --config .goreleaser.yml --clean --skip=publish --snapshot; \
		git tag -d $$TEMP_TAG; \
	else \
		goreleaser release --config .goreleaser.yml --clean --skip=publish --snapshot; \
	fi

api-package:
	@echo "Creating go-jo-api package..."
	@if [ -z "$$(git tag --points-at HEAD)" ]; then \
		echo "No git tag found, creating temporary tag for packaging..."; \
		TEMP_TAG="temp-package-$$(date +%s)"; \
		git tag $$TEMP_TAG; \
		goreleaser release --config .goreleaser-api.yml --clean --skip=publish --snapshot; \
		git tag -d $$TEMP_TAG; \
	else \
		goreleaser release --config .goreleaser-api.yml --clean --skip=publish --snapshot; \
	fi

package-all: gojo-package api-package

# Development targets
gojo-dev:
	@echo "Running go-jo in development mode..."
	@go run ./apps/go-jo

api-dev:
	@echo "Running go-jo-api in development mode..."
	@cd apps/go-jo-api && go run .

clean:
	@echo "Cleaning build artifacts..."
	@rm -rf dist/
	@rm -rf .goreleaser/

help:
	@echo "Available targets:"
	@echo "  gojo-build          - Build go-jo binary"
	@echo "  api-build           - Build go-jo-api binary"
	@echo "  installer-build     - Build go-jo-integration-installer binary"
	@echo "  build-all           - Build all applications"
	@echo "  gojo-package        - Create go-jo .deb package"
	@echo "  api-package         - Create go-jo-api .deb package"
	@echo "  package-all         - Create all packages"
	@echo "  gojo-dev            - Run go-jo in development mode"
	@echo "  api-dev             - Run go-jo-api in development mode"
	@echo "  clean               - Clean build artifacts"
	@echo "  help                - Show this help message" 