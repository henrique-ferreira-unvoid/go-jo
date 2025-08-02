# go-jo Monorepo Makefile
.PHONY: help gojo-build gojo-dev gojo-package clean

# Default target
help:
	@echo "Available targets:"
	@echo "  gojo-build        - Build go-jo binary for local development"
	@echo "  gojo-package      - Create go-jo .deb package (requires git tag or uses fake tag)"
	@echo "  gojo-dev          - Run go-jo in development mode"
	@echo "  clean             - Clean build artifacts"
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

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf dist/
	@echo "Build artifacts cleaned"

# Build everything
build-all: gojo-build

# Release everything (for CI/CD)
release-all: gojo-package 