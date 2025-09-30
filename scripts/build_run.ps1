# PowerShell script: Build and run iso-downloader for Windows

# Set root directory to the parent of the script's directory
$ROOT_DIR = Split-Path -Parent $PSScriptRoot
Set-Location $ROOT_DIR

Write-Host "Tidying modules..."
go mod tidy

Write-Host "Building iso-downloader.exe..."
New-Item -ItemType Directory -Force -Path "dist" | Out-Null
go build -o dist/iso-downloader.exe ./cmd/iso-downloader

if (Test-Path "dist/iso-downloader.exe") {
    Write-Host "Running iso-downloader.exe..."
    & "dist/iso-downloader.exe" @args
    exit $LASTEXITCODE
} else {
    Write-Error "Build failed: dist/iso-downloader.exe not found."
    exit 1
}
