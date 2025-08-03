# go-jo Monorepo

This repository contains three Go applications in a monorepo structure:

## Applications

### 1. go-jo
A simple CLI application that prints the current version when executed.

**Location:** `apps/go-jo/`

**Installation:**
```bash
sudo dpkg -i go-jo_*_linux_amd64.deb
```

**Usage:**
```bash
go-jo
```

### 2. go-jo-api
A REST API service that provides endpoints for version management and integration downloads.

**Location:** `apps/go-jo-api/`

**Installation:**
```bash
sudo dpkg -i go-jo-api_*_linux_amd64.deb
sudo cp /etc/go-jo-api/.env.example /etc/go-jo-api/.env
sudo nano /etc/go-jo-api/.env  # Configure tokens
sudo systemctl start go-jo-api
```

**API Endpoints:**
- `GET /health` - Health check (no auth required)
- `GET /versions` - Get available versions (auth required)
- `GET /integrations` - Get available integrations (auth required)
- `GET /download/{version}/{integration}` - Download combined package (auth required)

### 3. go-jo-integration-installer
A CLI tool for downloading and installing go-jo integrations.

**Location:** `apps/go-jo-integration-installer/`

**Usage:**
```bash
./go-jo-integration-installer --license <path-to-license-file>
```

**Features:**
- Interactive version selection
- Interactive integration selection
- Automatic package download
- Environment-based configuration

## Development Prerequisites

- Go 1.24+
- Make
- Git

## Environment Variables

### go-jo-api
Create a `.env` file in the API directory or set these environment variables:

- `GITHUB_TOKEN`: GitHub API token for accessing repositories
- `LICENSE_TOKEN`: Authorization token for API endpoints
- `PORT`: API server port (default: 1207)
- `API_URL`: API base URL

### go-jo-integration-installer
- `API_URL`: The URL of the go-jo-api service (default: http://localhost:1207)

## Makefile Targets

### Building
```bash
make gojo-build          # Build go-jo binary
make api-build           # Build go-jo-api binary
make installer-build     # Build go-jo-integration-installer binary
make build-all           # Build all applications
```

### Packaging
```bash
make gojo-package        # Create go-jo .deb package
make api-package         # Create go-jo-api .deb package
make package-all         # Create all packages
```

### Development
```bash
make gojo-dev           # Run go-jo in development mode
make api-dev            # Run go-jo-api in development mode
make clean              # Clean build artifacts
```

## GitHub Actions Release Management

The repository uses automated GitHub Actions workflows for releases:

### Manual Release Workflow
Located at `.github/workflows/release.yml`, this workflow can be manually triggered with the following inputs:

- **Version Source**: Choose between "auto-increment" or "manual"
- **Release Type** (for auto-increment): "patch", "minor", or "major"
- **Manual Version** (for manual): Specify version like "v1.2.3"
- **Pre-release**: Mark as pre-release
- **Skip Merge**: Skip branch merging (use current main)

### Branch Merging Workflow
Located at `.github/workflows/branch-merge.yml`, this workflow automates branch merging between `main` and `dev`.

### Validation Workflow
Located at `.github/workflows/validate.yml`, this workflow runs on PRs and pushes to validate code quality.

## Release Process

1. **Branch Strategy**: Uses `main` (production) and `dev` (development) branches
2. **Automated Merging**: GitHub Actions can merge branches automatically
3. **Version Management**: Supports both auto-increment and manual versioning
4. **Unified Releases**: Single release tag and GitHub Release for all applications
5. **Artifact Publishing**: Automatically publishes `.deb` packages and binaries

### Release Artifacts

Each release includes:
- **go-jo**: `.deb` package for CLI installation
- **go-jo-api**: `.deb` package with systemd service
- **go-jo-integration-installer**: Binary file for manual execution

### Installation Instructions

After downloading a release:

1. **Install go-jo** (CLI):
   ```bash
   sudo dpkg -i go-jo_*_linux_amd64.deb
   ```

2. **Install go-jo-api** (Service):
   ```bash
   sudo dpkg -i go-jo-api_*_linux_amd64.deb
   sudo cp /etc/go-jo-api/.env.example /etc/go-jo-api/.env
   sudo nano /etc/go-jo-api/.env  # Configure tokens
   sudo systemctl start go-jo-api
   ```

3. **Use go-jo-integration-installer** (Binary):
   ```bash
   chmod +x go-jo-integration-installer_*_linux_amd64
   ./go-jo-integration-installer_*_linux_amd64 --license <path-to-license-file>
   ```

## Project Structure

```
.
├── apps/
│   ├── go-jo/                    # CLI application
│   ├── go-jo-api/               # REST API service
│   └── go-jo-integration-installer/  # Integration installer CLI
├── .github/workflows/            # GitHub Actions workflows
├── Makefile                     # Build and development commands
├── go.mod                       # Go module definition
└── README.md                    # This file
```
