# go-jo-integration-installer

A CLI tool for downloading and installing go-jo integrations.

## Usage

```bash
./go-jo-integration-installer --license=<path-to-license-file>
```

## Features

- **Robust Interactive Selection**: Modern terminal UI using bubbletea
- **License-based Authentication**: Secure API access with license tokens
- **Environment Configuration**: Configurable API endpoints
- **Automatic Downloads**: Downloads combined packages with descriptive names

## Requirements

- A running `go-jo-api` service
- A license file containing the authorization token

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

## Example

```bash
# Create a license file
echo "your-token-here" > license.txt

# Run the installer
./go-jo-integration-installer --license=license.txt
```

The tool will:
1. Fetch available versions from the API
2. Display a modern interactive interface for version selection
3. Fetch available integrations from the API
4. Display a modern interactive interface for integration selection
5. Download the combined package as `go-jo-{integration}.zip` 