#!/usr/bin/env bash
set -euo pipefail

# Entrypoint similar to linutil's start.sh: fetch latest binary and run it
# Source repo can be overridden with REPO=<owner>/<repo>

REPO="${REPO:-dhruvmistry2000/iso-downloader}"
NAME="iso-downloader"

command -v curl >/dev/null 2>&1 || { echo "curl is required"; exit 1; }

LATEST_URL="https://github.com/dhruvmistry2000/iso-downloader/releases/download/latest/iso-downloader"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Downloading latest binary..." >&2
BIN_PATH="$TMP_DIR/$NAME"
if ! curl -fSL "$LATEST_URL" -o "$BIN_PATH"; then
  echo "Failed to download binary from: $LATEST_URL" >&2
  exit 1
fi
chmod +x "$BIN_PATH"

echo "Running $NAME..." >&2
exec "$BIN_PATH" "$@"


