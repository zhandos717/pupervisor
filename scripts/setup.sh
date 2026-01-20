#!/bin/bash
# Setup development environment

set -e

echo "==> Setting up Pupervisor development environment..."

# Check Go version
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "==> Go version: $GO_VERSION"

# Download dependencies
echo "==> Downloading dependencies..."
go mod download

# Create config file if not exists
if [ ! -f "pupervisor.yaml" ]; then
    echo "==> Creating default config file..."
    cp configs/pupervisor.yaml.example pupervisor.yaml
fi

# Install development tools
echo "==> Installing development tools..."

if ! command -v golangci-lint &> /dev/null; then
    echo "    Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

if ! command -v air &> /dev/null; then
    echo "    Installing air (live reload)..."
    go install github.com/air-verse/air@latest
fi

echo ""
echo "==> Setup complete!"
echo ""
echo "Quick start:"
echo "  make run-dev    # Run in development mode"
echo "  make build      # Build binary"
echo "  make help       # Show all commands"
