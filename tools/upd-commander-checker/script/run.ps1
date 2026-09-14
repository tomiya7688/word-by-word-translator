$ErrorActionPreference = "Stop"

$Root = (Resolve-Path (Join-Path $PSScriptRoot "../../..")).Path
$CacheRoot = Join-Path $Root ".tools-cache"
$Cache = Join-Path $CacheRoot "upd-commander-base-design"
$Repo = "https://github.com/tomiya7688/upd-commander-base-design.git"
$PinnedCommit = "41143698dda8bf2bd1f985539dce15d82477ff9d"

function Assert-LastExitCode {
    if ($LASTEXITCODE -ne 0) {
        throw "command failed with exit code $LASTEXITCODE"
    }
}

New-Item -ItemType Directory -Force -Path $CacheRoot | Out-Null

if (-not (Test-Path (Join-Path $Cache ".git"))) {
    git clone --filter=blob:none --no-checkout $Repo $Cache
    Assert-LastExitCode
}

$Current = ""
try {
    $Current = (git -C $Cache rev-parse HEAD 2>$null).Trim()
} catch {
    $Current = ""
}

if ($Current -ne $PinnedCommit) {
    git -C $Cache fetch --depth 1 origin $PinnedCommit
    Assert-LastExitCode
    git -C $Cache checkout --detach $PinnedCommit
    Assert-LastExitCode
}

$Tool = Join-Path $Cache "support_tools/go/upd_commander_checker"
Push-Location $Tool
try {
    go run ./cmd/upd-commander-check $Root
    Assert-LastExitCode
} finally {
    Pop-Location
}
