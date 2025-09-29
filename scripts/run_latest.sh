#!/usr/bin/env bash
set -euo pipefail

# One-shot runner: fetch latest release binary and run it (no install)

REPO="${REPO:-dhruvmistry2000/iso-downloader}"
NAME="iso-downloader"

command -v curl >/dev/null 2>&1 || { echo "curl is required"; exit 1; }

LATEST_URL="https://api.github.com/repos/${REPO}/releases/latest"
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "$TMP_DIR"' EXIT

echo "Fetching latest release..." >&2
ASSET_URL=$(curl -fsSL "$LATEST_URL" | grep browser_download_url | grep -E "/${NAME}$" | head -n1 | cut -d '"' -f4)
if [ -z "$ASSET_URL" ]; then
  echo "Could not find asset URL in latest release." >&2
  exit 1
fi

BIN_PATH="$TMP_DIR/$NAME"
echo "Downloading $NAME..." >&2
curl -fsSL "$ASSET_URL" -o "$BIN_PATH"
chmod +x "$BIN_PATH"

echo "Running $NAME..." >&2
exec "$BIN_PATH" "$@"


