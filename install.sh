#!/usr/bin/env bash
set -euo pipefail

# Detect distro
detect_distro() {
    if [ -f /etc/os-release ]; then
        . /etc/os-release
        echo "$ID"
    elif command -v lsb_release >/dev/null 2>&1; then
        lsb_release -si | tr '[:upper:]' '[:lower:]'
    else
        uname -s | tr '[:upper:]' '[:lower:]'
    fi
}

install_go() {
    # Try to install Go using system package manager
    DISTRO="$1"
    if command -v go >/dev/null 2>&1; then
        echo "Go is already installed."
        return
    fi

    case "$DISTRO" in
        ubuntu|debian)
            sudo apt-get update
            sudo apt-get install -y golang
            ;;
        fedora)
            sudo dnf install -y golang
            ;;
        arch)
            sudo pacman -Sy --noconfirm go
            ;;
        opensuse*|suse)
            sudo zypper install -y go
            ;;
        alpine)
            sudo apk add go
            ;;
        *)
            echo "Unknown or unsupported distro: $DISTRO"
            echo "Trying to install Go via tarball..."
            GO_VERSION="1.22.4"
            ARCH=$(uname -m)
            case "$ARCH" in
                x86_64) GOARCH="amd64" ;;
                aarch64|arm64) GOARCH="arm64" ;;
                *) echo "Unsupported arch: $ARCH"; exit 1 ;;
            esac
            curl -LO "https://go.dev/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz"
            sudo tar -C /usr/local -xzf "go${GO_VERSION}.linux-${GOARCH}.tar.gz"
            export PATH="/usr/local/go/bin:$PATH"
            echo 'export PATH="/usr/local/go/bin:$PATH"' >> ~/.bashrc
            ;;
    esac
}

main() {
    DISTRO=$(detect_distro)
    echo "Detected distro: $DISTRO"
    install_go "$DISTRO"

    # Ensure go is in PATH
    if ! command -v go >/dev/null 2>&1; then
        export PATH="/usr/local/go/bin:$PATH"
    fi

    # Run build_run.sh
    if [ -f scripts/build_run.sh ]; then
        bash scripts/build_run.sh
    else
        echo "scripts/build_run.sh not found!"
        exit 1
    fi
}

main "$@"
