#!/usr/bin/env bash
set -euo pipefail

REPO="${REPO:-dhruvmistry2000/iso-downloader}"
NAME="iso-downloader"

command -v curl >/dev/null 2>&1 || { echo "curl is required"; exit 1; }

LATEST_URL="https://api.github.com/repos/${REPO}/releases/latest"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Fetching latest release..."
ASSET_URL=$(curl -fsSL "$LATEST_URL" | grep browser_download_url | grep -E "/${NAME}$" | head -n1 | cut -d '"' -f4)
if [ -z "$ASSET_URL" ]; then
  echo "Could not find asset URL in latest release." >&2
  exit 1
fi

echo "Downloading $NAME..."
curl -fsSL "$ASSET_URL" -o "$TMP_DIR/$NAME"
chmod +x "$TMP_DIR/$NAME"

# Download the latest distros.json from the GitHub repository
DISTROS_URL="https://raw.githubusercontent.com/${REPO}/main/data/distros.json"
DISTROS_PATH="$TMP_DIR/distros.json"
echo "Fetching distros.json..."
curl -fsSL "$DISTROS_URL" -o "$DISTROS_PATH"

# Run the binary with the config file from GitHub
echo "Running $NAME with config from GitHub..."
ISO_DOWNLOADER_CONFIG="$DISTROS_PATH" "$TMP_DIR/$NAME"
