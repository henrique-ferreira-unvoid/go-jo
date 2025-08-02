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
├── Makefile                  # Build shortcuts (coming soon)
└── .github/workflows/        # GitHub Actions (coming soon)
```

### Environment Variables

Copy `env.example` to `.env` and configure:

- `GITHUB_TOKEN`: GitHub token for API access
- `LICENSE_TOKEN`: License token for API authorization
- `API_URL`: API URL for the installer

## Building

### Local Development

```bash
go run ./apps/go-jo
```

### Build with GoReleaser

```bash
goreleaser build --snapshot --clean
```

### Release

Releases are managed through GitHub Actions (coming soon).

## License

MIT
