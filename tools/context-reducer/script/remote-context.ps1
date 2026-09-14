param(
    [string]$Remote = "origin",
    [string]$Branch = "main"
)

$ErrorActionPreference = "Stop"
$Root = (Resolve-Path (Join-Path $PSScriptRoot "../../..")).Path
$Target = "$Remote/$Branch"
$MaxCommits = 10
$MaxDiffLines = 120

Push-Location $Root
try {
    git fetch --prune $Remote | Out-Null
    if ($LASTEXITCODE -ne 0) { exit $LASTEXITCODE }

    git rev-parse --verify $Target 2>$null | Out-Null
    if ($LASTEXITCODE -ne 0) {
        Write-Error "remote ref not found: $Target"
        exit 1
    }

    $LocalHead = (git rev-parse --short HEAD).Trim()
    $RemoteHead = (git rev-parse --short $Target).Trim()
    $Counts = (git rev-list --left-right --count "HEAD...$Target").Trim() -split "\s+"
    $Ahead = $Counts[0]
    $Behind = $Counts[1]

    Write-Output "local=$LocalHead remote=$RemoteHead ahead=$Ahead behind=$Behind"
    Write-Output "-- remote commits --"
    git log --oneline --no-decorate "HEAD..$Target" -n $MaxCommits

    Write-Output "-- changed files --"
    git diff --name-status "HEAD...$Target"

    Write-Output "-- diff stat --"
    git diff --stat "HEAD...$Target"

    Write-Output "-- bounded diff excerpt (max $MaxDiffLines lines) --"
    git diff --no-ext-diff --unified=2 "HEAD...$Target" | Select-Object -First $MaxDiffLines
} finally {
    Pop-Location
}
