# go-jo Monorepo

This monorepo contains three GoLang applications:

## Applications

### 1. go-jo

A simple GoLang application that prints its current version, git commit, and build date.

**Usage:**

```bash
go-jo
```

**Installation:**
Download the `.deb` file from the releases page and install it:

```bash
sudo dpkg -i go-jo_*.deb
```

### 2. go-jo-api (Coming Soon)

A REST API built with gorilla/mux that provides version and integration management.

### 3. go-jo-integration-installer (Coming Soon)

A CLI tool for downloading and installing go-jo integrations.

## Development

### Prerequisites

- Go 1.24+
- GoReleaser (for building releases)

### Project Structure

```
/
├── apps/
│   ├── go-jo/                 # Main go-jo application
│   ├── go-jo-api/            # API application (coming soon)
│   └── go-jo-integration-installer/ # CLI installer (coming soon)
├── .goreleaser.yml           # GoReleaser configuration
├── go.mod                    # Go module file
├── Makefile                  # Build shortcuts
└── .github/workflows/        # GitHub Actions (coming soon)
```

### Environment Variables

Copy `env.example` to `.env` and configure:

- `GITHUB_TOKEN`: GitHub token for API access
- `LICENSE_TOKEN`: License token for API authorization
- `API_URL`: API URL for the installer

## Building

### Available Make Targets

```bash
make help              # Show all available targets
make gojo-build        # Build go-jo binary for local development
make gojo-fast-release # Build go-jo binary using GoReleaser (snapshot)
make gojo-package      # Create go-jo .deb package for testing
make gojo-dev          # Run go-jo in development mode
make clean             # Clean build artifacts
```

### Local Development

```bash
# Run in development mode
make gojo-dev

# Build binary for local testing
make gojo-build
./dist/go-jo/go-jo

# Build with GoReleaser
make gojo-fast-release
./dist/go-jo_linux_amd64_v1/go-jo
```

### Creating Packages

```bash
# Create .deb package for testing
make gojo-package

# Install the package locally
sudo dpkg -i dist/go-jo_*.deb
go-jo
```

### Release

Releases are managed through GitHub Actions (coming soon).

## License

MIT
