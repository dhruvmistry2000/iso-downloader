#!/usr/bin/env bash
set -euo pipefail

# Entrypoint: clone repo to tmpfs (so it's deleted on reboot) and run install.sh
# Source repo can be overridden with REPO=<owner>/<repo>

REPO="dhruvmistry2000/iso-downloader"
NAME="iso-downloader"

command -v git >/dev/null 2>&1 || { echo "git is required"; exit 1; }

# Use /dev/shm if available (tmpfs, deleted on reboot), else fallback to /tmp
if mountpoint -q /dev/shm; then
  TMPFS_DIR="/dev/shm"
else
  TMPFS_DIR="/tmp"
fi

CLONE_DIR="$(mktemp -d "$TMPFS_DIR/${NAME}.XXXXXX")"
trap 'rm -rf "$CLONE_DIR"' EXIT

echo "Cloning $REPO into $CLONE_DIR..." >&2
if ! git clone --depth 1 "https://github.com/$REPO.git" "$CLONE_DIR"; then
  echo "Failed to clone repository: $REPO" >&2
  exit 1
fi

# Run install.sh if it exists
if [[ -f "$CLONE_DIR/install.sh" ]]; then
  echo "Running install.sh..." >&2
  chmod +x "$CLONE_DIR/install.sh"
  exec "$CLONE_DIR/install.sh" "$@"
else
  echo "install.sh not found in the cloned repo." >&2
  exit 1
fi
