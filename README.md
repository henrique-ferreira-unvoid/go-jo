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

### 2. go-jo-api

A REST API built with gorilla/mux that provides version and integration management.

**Endpoints:**

- `GET /versions` - Get all tagged versions of go-jo (requires auth)
- `GET /integrations` - Get all integration branches (requires auth)
- `GET /download/{version}/{integration}` - Download combined package (requires auth)
- `GET /health` - Health check (no auth required)

**Installation:**

```bash
# Install the service
sudo dpkg -i go-jo-api_*.deb

# Configure environment
sudo cp /etc/go-jo-api/.env.example /etc/go-jo-api/.env
sudo nano /etc/go-jo-api/.env  # Add your GITHUB_TOKEN and LICENSE_TOKEN

# Enable and start service
sudo systemctl enable go-jo-api
sudo systemctl start go-jo-api
```

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
│   ├── go-jo-api/            # API application
│   │   ├── main.go           # Simple entry point
│   │   ├── go-jo-api.service # Systemd service file
│   │   └── api/              # Clean architecture
│   │       ├── api.go        # Main API orchestrator
│   │       ├── domain/       # Domain models and constants
│   │       ├── handlers/     # HTTP handlers (separated by feature)
│   │       └── router/       # Router and subrouter logic
│   └── go-jo-integration-installer/ # CLI installer (coming soon)
├── .github/workflows/        # GitHub Actions workflows
│   ├── branch-merge.yml      # Automated branch merging
│   ├── release.yml           # Release management
│   └── validate.yml          # Build validation
├── .goreleaser.yml           # GoReleaser configuration
├── go.mod                    # Go module file
└── Makefile                  # Build shortcuts
```

### Environment Variables

Copy `env.example` to `.env` and configure:

- `GITHUB_TOKEN`: GitHub token for API access (required for go-jo-api)
- `LICENSE_TOKEN`: License token for API authorization (required for go-jo-api)
- `PORT`: API port (defaults to 1207)
- `API_URL`: API URL for the installer

## Building

### Available Make Targets

#### Go-jo app:

```bash
make gojo-build        # Build go-jo binary for local development
make gojo-package      # Create go-jo .deb package for testing
make gojo-dev          # Run go-jo in development mode
```

#### Go-jo API:

```bash
make api-build         # Build go-jo-api binary for local development
make api-package       # Create go-jo-api package structure for testing
make api-dev           # Run go-jo-api in development mode
```

#### General:

```bash
make help              # Show all available targets
make build-all         # Build all applications
make package-all       # Package all applications
make clean             # Clean build artifacts
```

### Local Development

#### go-jo app:

```bash
# Run in development mode
make gojo-dev

# Build binary for local testing
make gojo-build
./dist/go-jo/go-jo
```

#### go-jo-api:

```bash
# Set up environment (required)
cp env.example .env
# Edit .env and add your GITHUB_TOKEN and LICENSE_TOKEN

# Run in development mode
make api-dev

# Build binary for local testing
make api-build
./dist/api-bin/go-jo-api
```

### Creating Packages

```bash
# Create .deb package for go-jo
make gojo-package

# Create package structure for go-jo-api
make api-package

# Create packages for all apps
make package-all
```

### API Usage Examples

```bash
# Get versions (requires LICENSE_TOKEN)
curl -H "Authorization: your_license_token" http://localhost:1207/versions

# Get integrations
curl -H "Authorization: your_license_token" http://localhost:1207/integrations

# Download latest version with syslog integration
curl -H "Authorization: your_license_token" \
     -o go-jo-syslog.zip \
     http://localhost:1207/download/latest/syslog

# Health check (no auth required)
curl http://localhost:1207/health
```

## Release Management

### Branch Strategy

- `main` - Production branch (stable releases)
- `dev` - Development branch (latest features)

### GitHub Actions Workflows

#### 1. Branch Merging (`branch-merge.yml`)

Manually triggered workflow for merging branches:

- **main-to-dev**: Merges main into dev (for hotfixes)
- **dev-to-main**: Merges dev into main (for releases)
- Creates pull requests with auto-merge enabled
- Follows the required branching strategy

#### 2. Release Management (`release.yml`)

Manually triggered workflow for creating releases:

- Merges dev to main (unless skipped)
- Auto-increments version (patch/minor/major) or uses manual version
- Creates git tags
- Builds and publishes releases with GoReleaser
- Supports pre-release marking

#### 3. Validation (`validate.yml`)

Automatically runs on pull requests and pushes:

- Validates Go code (vet, tests, build)
- Tests GoReleaser configuration
- Validates Makefile targets
- Provides coverage reports

### Creating a Release

1. Go to Actions tab in GitHub
2. Select "Generate Release" workflow
3. Click "Run workflow"
4. Choose version source:
   - **Auto-increment**: Select patch/minor/major
   - **Manual**: Enter custom version
5. Optionally mark as pre-release
6. The workflow will:
   - Merge dev → main (if not skipped)
   - Create new version tag
   - Build and publish release assets

### Manual Branch Merging

1. Go to Actions tab in GitHub
2. Select "Branch Merge Automation" workflow
3. Click "Run workflow"
4. Choose merge direction
5. Add optional description
6. The workflow will create and auto-merge the PR

## License

MIT
