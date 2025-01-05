#!/bin/zsh
set -e

# Sprawdzenie argumentu
if [ "$#" -ne 1 ]; then
  echo "Error: Invalid arguments. Usage: $0 [macos|linux]"
  exit 1
fi

# Przypisanie argumentu do zmiennej
TARGET="$1"

# Warunki dla macOS i Linux
if [ "$TARGET" = "macos" ]; then
  echo "Building project for macOS..."
  go build -o ./bin/NetManager ./cmd/app/
  echo "Build for macOS complete!"
elif [ "$TARGET" = "linux" ]; then
  echo "Building project for Linux (amd64)..."
  GOOS=linux GOARCH=amd64 go build -o ./bin/NetManager ./cmd/app/
  echo "Build for Linux complete!"
elif [ "$TARGET" = "color" ]; then
  go run ./scripts/color.go
else
  echo "Error: Invalid target. Use 'macos' or 'linux'."
  exit 1
fi