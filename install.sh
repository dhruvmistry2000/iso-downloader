#!/usr/bin/env bash
set -euo pipefail

REPO="${REPO:-dhruvmistry2000/iso-downloader}"
BIN_DIR="${BIN_DIR:-/usr/local/bin}"
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

echo "Installing to $BIN_DIR (sudo may be required)..."
if [ ! -w "$BIN_DIR" ]; then
  sudo mkdir -p "$BIN_DIR"
  sudo cp "$TMP_DIR/$NAME" "$BIN_DIR/$NAME"
else
  mkdir -p "$BIN_DIR"
  cp "$TMP_DIR/$NAME" "$BIN_DIR/$NAME"
fi

echo "Installed: $(command -v $NAME)"
"$BIN_DIR/$NAME" --help || true
