$ErrorActionPreference = "Stop"

$Root = [System.IO.Path]::GetFullPath((Join-Path $PSScriptRoot "../../.."))

powershell -ExecutionPolicy Bypass -File (Join-Path $Root "tools/ejdict-import/script/run.ps1")
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

powershell -ExecutionPolicy Bypass -File (Join-Path $Root "tools/tkg-ja-en-import/script/run.ps1")
if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

Push-Location $Root
try {
    go run ./tools/release-builder/script `
        -dictionaries (Join-Path $Root "dist/dictionaries") `
        -output (Join-Path $Root "dist/releases") `
        -app-license (Join-Path $Root "LICENSE") `
        -analyzer-notices (Join-Path $Root "licenses/analyzers") `
        -tier all
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }
}
finally {
    Pop-Location
}

Write-Host "OK release-builder: $(Join-Path $Root 'dist/releases')"
