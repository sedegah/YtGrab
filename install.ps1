param(
    [string]$InstallDir = "$env:USERPROFILE\go\bin",
    [switch]$InstallDeps
)

$ErrorActionPreference = 'Stop'

function Write-Step {
    param([string]$Message)
    Write-Host "[yt-grab] $Message" -ForegroundColor Cyan
}

function Write-Warn {
    param([string]$Message)
    Write-Host "[yt-grab] $Message" -ForegroundColor Yellow
}

function Ensure-WingetPackage {
    param(
        [string]$CommandName,
        [string]$PackageId,
        [string]$FriendlyName
    )

    if (Get-Command $CommandName -ErrorAction SilentlyContinue) {
        Write-Step "$FriendlyName detected"
        return
    }

    if (-not $InstallDeps) {
        Write-Warn "$FriendlyName not found in PATH"
        Write-Warn "Run: winget install $PackageId"
        return
    }

    if (-not (Get-Command winget -ErrorAction SilentlyContinue)) {
        Write-Warn "winget not found. Install $FriendlyName manually, then restart PowerShell."
        return
    }

    Write-Step "Installing $FriendlyName via winget"
    winget install --id $PackageId --accept-source-agreements --accept-package-agreements

    if (Get-Command $CommandName -ErrorAction SilentlyContinue) {
        Write-Step "$FriendlyName installed successfully"
    } else {
        Write-Warn "$FriendlyName installation finished but command still not found in current shell. Restart PowerShell."
    }
}

$repoRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$binaryName = 'yt-grab.exe'

if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
    throw "Go is not installed or not available in PATH. Install Go first: https://go.dev/dl/"
}

Write-Step "Building $binaryName"
Push-Location $repoRoot
try {
    go build -o $binaryName ./cmd/yt-grab
} finally {
    Pop-Location
}

if (-not (Test-Path $InstallDir)) {
    Write-Step "Creating install directory: $InstallDir"
    New-Item -Path $InstallDir -ItemType Directory -Force | Out-Null
}

$sourceBinary = Join-Path $repoRoot $binaryName
$targetBinary = Join-Path $InstallDir $binaryName

Write-Step "Installing binary to $targetBinary"
Copy-Item -Path $sourceBinary -Destination $targetBinary -Force

$currentUserPath = [Environment]::GetEnvironmentVariable('Path', 'User')
$pathItems = @()
if ($currentUserPath) {
    $pathItems = $currentUserPath.Split(';') | Where-Object { $_ -and $_.Trim() -ne '' }
}

if (-not ($pathItems | Where-Object { $_.TrimEnd('\\') -ieq $InstallDir.TrimEnd('\\') })) {
    $newPath = if ([string]::IsNullOrWhiteSpace($currentUserPath)) {
        $InstallDir
    } else {
        "$currentUserPath;$InstallDir"
    }
    [Environment]::SetEnvironmentVariable('Path', $newPath, 'User')
    Write-Step "Added $InstallDir to your USER PATH"
} else {
    Write-Step "Install directory is already present in USER PATH"
}

Write-Host ""
Write-Host "Installation complete." -ForegroundColor Green
Write-Host "If this is a new PATH entry, restart PowerShell before using 'yt-grab'."
Write-Host "Then run: yt-grab --help"
Write-Host ""

Ensure-WingetPackage -CommandName "yt-dlp" -PackageId "yt-dlp.yt-dlp" -FriendlyName "yt-dlp"
Ensure-WingetPackage -CommandName "ffmpeg" -PackageId "Gyan.FFmpeg" -FriendlyName "ffmpeg"

Write-Host ""
Write-Host "Verification:" -ForegroundColor Cyan
Write-Host "  yt-grab doctor"
Write-Host "  yt-grab https://youtu.be/8ekJMC8OtGU"
Write-Host "  yt-grab https://youtu.be/8ekJMC8OtGU --audio"
