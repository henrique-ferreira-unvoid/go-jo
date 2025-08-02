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

#### 1. Branch Merge Automation (`branch-merge.yml`)

Automatically merges branches (main ↔ dev) with proper PR creation and auto-merge.

#### 2. Release Management

Unified release system with separate build workflows:

##### Main Release Orchestrator (`release.yml`)
The primary release workflow that coordinates everything:

- **Branch Management**: Merges dev to main (unless skipped)
- **Version Management**: Auto-increments or uses manual version
- **Tag Creation**: Creates single `v*` git tag for both apps
- **Build Orchestration**: Calls individual app build workflows
- **Release Publishing**: Creates unified GitHub release with all assets

##### Individual Build Workflows
- **`release-main.yml`**: Builds go-jo app artifacts only
- **`release-api.yml`**: Builds go-jo-api artifacts only

#### 3. Validation (`validate.yml`)

Automatically runs on pull requests and pushes:

- Validates Go code (vet, tests, build)
- Tests GoReleaser configuration
- Validates Makefile targets
- Provides coverage reports

### Creating Releases

#### Unified Release Process

1. Go to Actions tab in GitHub
2. Select "Generate Release" workflow
3. Click "Run workflow"
4. Choose version source:
   - **Auto-increment**: Select patch/minor/major
   - **Manual**: Enter custom version
5. Optionally mark as pre-release
6. The workflow will:
   - Merge dev → main (if not skipped)
   - Create single `v*` version tag
   - Build both go-jo and go-jo-api in parallel
   - Publish unified GitHub release with all assets

#### Release Features

- **Single Version**: Both apps share the same version (e.g., `v1.2.3`)
- **Unified Release**: One GitHub release contains all application packages
- **Parallel Builds**: Apps are built simultaneously for faster releases
- **Professional Release Notes**: Includes installation instructions for both apps

### Local Development Testing

#### Go-jo App
```bash
# Test packaging
make gojo-package
```

#### Go-jo API
```bash
# Test packaging
make api-package
```

### Manual Branch Merging

1. Go to Actions tab in GitHub
2. Select "Branch Merge Automation" workflow
3. Click "Run workflow"
4. Choose merge direction
5. Add optional description
6. The workflow will create and auto-merge the PR

## License

MIT
