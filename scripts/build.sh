#!/bin/zsh
set -e

echo "Building project..."
GOOS=linux GOARCH=amd64 go build -o ../bin/NetManager ../cmd/app/

echo "Build complete!"