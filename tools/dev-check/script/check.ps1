$ErrorActionPreference = "Stop"

$Root = (Resolve-Path (Join-Path $PSScriptRoot "../../..")).Path
Push-Location $Root
try {
    $GoFiles = @(git ls-files "*.go")
    if ($GoFiles.Count -gt 0) {
        $Unformatted = @(gofmt -l $GoFiles)
        if ($Unformatted.Count -gt 0) {
            Write-Output "FAIL gofmt"
            $Unformatted | ForEach-Object { Write-Output $_ }
            exit 1
        }
    }

    Write-Output "go vet"
    go vet ./...
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

    Write-Output "go test"
    go test ./...
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

    Write-Output "UPD checker"
    & (Join-Path $Root "tools/upd-commander-checker/script/run.ps1")
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

    Write-Output "OK"
} finally {
    Pop-Location
}
