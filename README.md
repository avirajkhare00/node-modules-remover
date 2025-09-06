# Node Modules Remover

A cross-platform CLI tool written in Go that automatically removes old `node_modules` directories to free up disk space.

## Features

- **Cross-platform**: Works on Windows, macOS, and Linux
- **Age-based filtering**: Only removes `node_modules` directories older than a specified duration (default: 3 months)
- **Safe operation**: Dry-run mode to preview what would be deleted
- **Flexible input**: Specify directories to scan via command-line arguments
- **Cron-friendly**: Quiet mode for automated scripts
- **Detailed reporting**: Shows size freed and number of directories processed

## Installation

### Download Pre-built Binaries

Download the latest release from the [Releases page](https://github.com/your-username/node-modules-remover/releases) and choose the appropriate binary for your operating system and architecture.

### Build from source

```bash
git clone <repository-url>
cd node-modules-remover

# Simple build
go build -o node-modules-remover .

# Or use the Makefile for cross-platform builds
make build-all
```

### Usage

```bash
# Basic usage - scan current directory for node_modules older than 3 months
./node-modules-remover

# Dry run to see what would be deleted
./node-modules-remover -dry-run -verbose

# Remove node_modules older than 7 days
./node-modules-remover -age 7d

# Scan specific directories
./node-modules-remover -dirs /path/to/projects,/another/path

# Quiet mode for cron jobs
./node-modules-remover -quiet

# Show help
./node-modules-remover -help
```

## Command Line Options

| Option | Description | Default |
|--------|-------------|---------|
| `-age duration` | Remove node_modules older than this duration (e.g., 3m, 90d, 24h) | 3m |
| `-dirs string` | Comma-separated list of directories to scan | Current directory |
| `-dry-run` | Show what would be deleted without actually deleting | false |
| `-verbose` | Show detailed output | false |
| `-quiet` | Minimal output (good for cron) | false |
| `-help` | Show help message | - |
| `-version` | Show version | - |

## Examples

### Safe testing with dry-run
```bash
./node-modules-remover -dry-run -verbose
```

### Clean up old node_modules in development folders
```bash
./node-modules-remover -dirs ~/projects,~/work -age 30d -verbose
```

### Automated cleanup via cron (runs weekly)
```bash
# Add to crontab
0 2 * * 0 /path/to/node-modules-remover -quiet -age 7d
```

### Remove very old node_modules (older than 6 months)
```bash
./node-modules-remover -age 6m -verbose
```

## Safety Features

- **Age filtering**: Only removes directories older than the specified threshold
- **Dry-run mode**: Preview operations before executing
- **Error handling**: Continues processing even if some directories can't be accessed
- **Size reporting**: Shows how much disk space would be freed

## Development

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run tests with coverage
make test-coverage

# Clean build artifacts
make clean
```

### CI/CD

This project uses GitHub Actions for continuous integration and deployment:

- **Tests**: Run on every push and pull request
- **Builds**: Cross-platform binaries are built for every push to main
- **Releases**: Automatic releases are created for tags starting with `v*`

Supported platforms:
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

## Requirements

- Go 1.21 or later (for building from source)
- No external dependencies

## License

MIT License
