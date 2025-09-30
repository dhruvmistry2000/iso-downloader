#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "$0")/.." && pwd)
cd "$ROOT_DIR"

echo "Tidying modules..."
go mod tidy

echo "Building iso-downloader..."
mkdir -p dist
go build -o dist/iso-downloader ./cmd/iso-downloader

echo "Running iso-downloader..."
exec ./dist/iso-downloader "$@"


