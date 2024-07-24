#!/bin/zsh
set -e

echo "Building project..."
go build -o ../bin/NetManager ../cmd/app/

echo "Build complete!"