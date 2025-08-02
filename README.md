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

- `GITHUB_TOKEN`: GitHub token for API access
- `LICENSE_TOKEN`: License token for API authorization
- `API_URL`: API URL for the installer

## Building

### Available Make Targets

```bash
make help              # Show all available targets
make gojo-build        # Build go-jo binary for local development
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
- Auto-increments version (patch/minor/major)
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
2. Select "Release Management" workflow
3. Click "Run workflow"
4. Choose release type (patch/minor/major)
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
