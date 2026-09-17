$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$root = (Resolve-Path (Join-Path $scriptDir "../../..")).Path
$lockFile = Join-Path $root "tools/tkg-ja-en-import/source.lock"
$pin = (Get-Content -Raw $lockFile).Trim()

if ([string]::IsNullOrWhiteSpace($pin)) {
    throw "tkg-ja-en-import: empty source.lock"
}

if ($args.Count -gt 0) {
    $output = $args[0]
} else {
    $output = Join-Path $root "dist/dictionaries/mit/tkg-ja-en"
}

$cacheDir = Join-Path $root ".tools-cache/tkg-ja-en/$pin"
$source = Join-Path $cacheDir "entries_index.json"
$url = "https://raw.githubusercontent.com/tkgally/je-dict-1/$pin/entries_index.json"
New-Item -ItemType Directory -Force -Path $cacheDir | Out-Null

$needsDownload = !(Test-Path $source)
if (!$needsDownload) {
    $needsDownload = (Get-Item $source).Length -eq 0
}
if ($needsDownload) {
    $temp = "$source.tmp"
    Remove-Item -Force -ErrorAction SilentlyContinue $temp
    Invoke-WebRequest -Uri $url -OutFile $temp
    Move-Item -Force $temp $source
}

Remove-Item -Recurse -Force -ErrorAction SilentlyContinue $output
Push-Location $root
try {
    & go run ./tools/tkg-ja-en-import/script -source $source -output $output -revision $pin
    if ($LASTEXITCODE -ne 0) {
        throw "tkg-ja-en-import failed with exit code $LASTEXITCODE"
    }
}
finally {
    Pop-Location
}

Write-Host "OK tkg-ja-en-import: $output"
