# go-jo-integration-installer

A CLI tool for downloading and installing go-jo integrations with automatic Docker deployment.

## Usage

```bash
./go-jo-integration-installer --license=<path-to-license-file>
```

## Features

- **Docker Validation**: Automatically checks if Docker is running
- **Robust Interactive Selection**: Modern terminal UI using bubbletea
- **License-based Authentication**: Secure API access with license tokens
- **Environment Configuration**: Configurable API endpoints
- **Automatic Downloads**: Downloads combined packages with descriptive names
- **Automatic Deployment**: Extracts and deploys Docker environments automatically
- **Live Output**: Shows real-time output from make commands
- **Clean Architecture**: Well-organized code structure with separate packages

## Project Structure

```
apps/go-jo-integration-installer/
├── main.go                    # Application entry point
├── cli/                       # CLI interface and selection logic
│   └── cli.go
├── api/                       # API client for go-jo-api
│   └── client.go
├── config/                    # Configuration management
│   └── config.go
├── docker/                    # Docker-related functions
│   └── docker.go
├── utils/                     # Utility functions
│   ├── commands.go            # Command-line parsing
│   └── files.go               # File operations
└── README.md                  # This file
```

## Requirements

- A running `go-jo-api` service
- A license file containing the authorization token
- **Docker**: Must be installed and running
- **Make**: Required for building and starting containers

## Configuration

The application can be configured using environment variables:

- `API_URL`: The URL of the go-jo-api service (default: http://localhost:1207)

You can also create a `.env` file in the same directory as the binary:

```
API_URL=http://your-api-server:1207
```

## License File

The license file should contain only the authorization token that will be used as the `Authorization` header when making requests to the API.

Example license file content:
```
your-license-token-here
```

## Interactive Interface

The tool provides a modern and robust interactive selection interface:

- **Arrow Keys**: Navigate up and down through options
- **Vim Keys**: Use `j` and `k` for navigation
- **Enter/Space**: Select the highlighted option
- **Ctrl+C/Q**: Cancel the selection process

The interface features:
- Clean, modern design with proper styling
- Smooth navigation with multiple key options
- Professional color highlighting
- Robust error handling and graceful cancellation
- Built on the popular bubbletea framework

## Deployment Process

The installer automatically:

1. **Validates Docker**: Checks if Docker is running
2. **Downloads Package**: Fetches the selected version and integration
3. **Extracts Archive**: Creates a temporary directory and extracts the zip
4. **Finds Deployment Directory**: Locates the directory with `docker-compose.yml` and `Makefile`
5. **Runs Make Commands**: Executes `make build` and `make start` with live output
6. **Cleans Up**: Removes temporary files after deployment

## Example

```bash
# Create a license file
echo "your-token-here" > license.txt

# Run the installer
./go-jo-integration-installer --license=license.txt
```

The tool will:
1. Check if Docker is running
2. Fetch available versions from the API
3. Display a modern interactive interface for version selection
4. Fetch available integrations from the API
5. Display a modern interactive interface for integration selection
6. Download the combined package as `go-jo-{integration}.zip`
7. Extract and automatically deploy the Docker environment
8. Run `make build` and `make start` with live output
9. Clean up temporary files

## Error Handling

The installer provides clear error messages for common issues:

- **Docker not running**: Exits with clear error message
- **Invalid license file**: Shows helpful error details
- **Network issues**: Displays API connection errors
- **Deployment failures**: Shows make command output and errors 