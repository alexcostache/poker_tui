param(
    [ValidateSet("amd64", "arm64")]
    [string]$Arch = "amd64",
    [string]$Output = "poker_tui.exe"
)

$ErrorActionPreference = "Stop"

$env:GOOS = "windows"
$env:GOARCH = $Arch
$env:GOTELEMETRY = "off"

$cacheDir = Join-Path $PSScriptRoot ".gocache"
New-Item -ItemType Directory -Path $cacheDir -Force | Out-Null
$env:GOCACHE = $cacheDir

$goPathDir = Join-Path $PSScriptRoot ".gopath"
$goModCacheDir = Join-Path $PSScriptRoot ".gomodcache"
New-Item -ItemType Directory -Path $goPathDir -Force | Out-Null
New-Item -ItemType Directory -Path $goModCacheDir -Force | Out-Null
$env:GOPATH = $goPathDir
$env:GOMODCACHE = $goModCacheDir

Write-Host "==> Building $Output (GOOS=$env:GOOS GOARCH=$env:GOARCH)..."
go build -o $Output ./cmd/poker_tui
if ($LASTEXITCODE -ne 0) {
    throw "go build failed with exit code $LASTEXITCODE"
}
Write-Host "==> Done: .\$Output"
