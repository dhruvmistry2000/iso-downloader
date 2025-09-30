# PowerShell script: Install Go (if needed) and run build_run.sh for Windows

# Function to check if Go is installed
function Test-GoInstalled {
    $go = Get-Command go -ErrorAction SilentlyContinue
    return $null -ne $go
}

# Function to download and install Go for Windows
function Install-Go {
    if (Test-GoInstalled) {
        Write-Host "Go is already installed."
        return
    }

    $goVersion = "1.22.4"
    $arch = if ([Environment]::Is64BitOperatingSystem) { "amd64" } else { "386" }
    $goMsiUrl = "https://go.dev/dl/go${goVersion}.windows-${arch}.msi"
    $msiPath = Join-Path $env:TEMP "go${goVersion}.windows-${arch}.msi"

    Write-Host "Downloading Go $goVersion for Windows ($arch)..."
    Invoke-WebRequest -Uri $goMsiUrl -OutFile $msiPath

    Write-Host "Installing Go..."
    Start-Process msiexec.exe -Wait -ArgumentList "/i `"$msiPath`" /qn /norestart"

    # Add Go to PATH for current session if not already present
    $goBin = "C:\Program Files\Go\bin"
    if (-not ($env:PATH -split ";" | Where-Object { $_ -eq $goBin })) {
        $env:PATH = "$goBin;$env:PATH"
    }
    Write-Host "Go installed successfully."
}

# Main logic
Write-Host "Checking for Go installation..."
Install-Go

# Ensure Go is in PATH for this session
if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    $env:PATH = "C:\Program Files\Go\bin;$env:PATH"
}

# Run build_run.sh using WSL or Git Bash if available, else try to build/run natively
$buildScript = "scripts/build_run.sh"
if (Test-Path $buildScript) {
    # Prefer WSL if available
    if (Get-Command wsl.exe -ErrorAction SilentlyContinue) {
        Write-Host "Running build_run.sh using WSL..."
        wsl bash "$buildScript"
        exit $LASTEXITCODE
    }
    # Try Git Bash
    $gitBash = "${env:ProgramFiles}\Git\bin\bash.exe"
    if (Test-Path $gitBash) {
        Write-Host "Running build_run.sh using Git Bash..."
        & "$gitBash" "$buildScript"
        exit $LASTEXITCODE
    }
    # Fallback: try to build/run natively with Go
    Write-Host "No WSL or Git Bash found. Attempting to build and run with Go natively..."
    if (Test-Path "cmd/iso-downloader/main.go") {
        go build -o dist/iso-downloader.exe cmd/iso-downloader/main.go
        if (Test-Path "dist/iso-downloader.exe") {
            Write-Host "Running iso-downloader.exe..."
            & "dist/iso-downloader.exe"
            exit $LASTEXITCODE
        } else {
            Write-Error "Build failed: dist/iso-downloader.exe not found."
            exit 1
        }
    } else {
        Write-Error "cmd/iso-downloader/main.go not found!"
        exit 1
    }
} else {
    Write-Error "scripts/build_run.sh not found!"
    exit 1
}
