# PowerShell script: Clone a GitHub repo to a temp directory and run install.ps1 if present

# You can override the repo by setting $env:REPO to <owner>/<repo>
$REPO = $env:REPO
if (-not $REPO) {
    $REPO = "dhruvmistry2000/iso-downloader"
}
$NAME = "iso-downloader"

# Check for git
if (-not (Get-Command git -ErrorAction SilentlyContinue)) {
    Write-Error "git is required"
    exit 1
}

# Use $env:TEMP for temp directory
$TMPFS_DIR = $env:TEMP
if (-not $TMPFS_DIR) {
    $TMPFS_DIR = "C:\Windows\Temp"
}

# Create a unique temp directory
$CLONE_DIR = New-Item -ItemType Directory -Path (Join-Path $TMPFS_DIR "$NAME.$([System.Guid]::NewGuid().ToString('N').Substring(0,8))") -Force
$CLONE_DIR = $CLONE_DIR.FullName

# Ensure cleanup on exit
$cleanup = {
    if (Test-Path $using:CLONE_DIR) {
        Remove-Item -Recurse -Force $using:CLONE_DIR
    }
}
Register-EngineEvent PowerShell.Exiting -Action $cleanup | Out-Null

Write-Host "Cloning $REPO into $CLONE_DIR..." -ForegroundColor Yellow
$gitUrl = "https://github.com/$REPO.git"
$gitClone = git clone --depth 1 $gitUrl $CLONE_DIR 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Error "Failed to clone repository: $REPO"
    exit 1
}

# Run install.ps1 if it exists
$installScript = Join-Path $CLONE_DIR "install.ps1"
if (Test-Path $installScript) {
    Write-Host "Running install.ps1..." -ForegroundColor Yellow
    & $installScript @args
    exit $LASTEXITCODE
} else {
    Write-Error "install.ps1 not found in the cloned repo."
    exit 1
}
