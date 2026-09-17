param(
    [string]$Output = ""
)

$ErrorActionPreference = "Stop"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$Root = (Resolve-Path (Join-Path $ScriptDir "../../..")).Path
$LockFile = Join-Path $Root "tools/ejdict-import/source.lock"
$Cache = Join-Path $Root ".tools-cache/ejdict-hand"
$Repo = "https://github.com/kujirahand/EJDict.git"
$Pin = (Get-Content -Raw $LockFile).Trim()

if ([string]::IsNullOrWhiteSpace($Pin)) {
    throw "ejdict-import: empty source.lock"
}

if ([string]::IsNullOrWhiteSpace($Output)) {
    $Output = Join-Path $Root "dist/dictionaries/mit/ejdict-hand"
}

function Invoke-Checked {
    param([scriptblock]$Command)
    & $Command
    if ($LASTEXITCODE -ne 0) {
        throw "command failed with exit code $LASTEXITCODE"
    }
}

if (-not (Test-Path (Join-Path $Cache ".git"))) {
    if (Test-Path $Cache) {
        Remove-Item -Recurse -Force $Cache
    }
    Invoke-Checked { git clone --filter=blob:none --no-checkout $Repo $Cache }
}

Invoke-Checked { git -C $Cache fetch --depth 1 origin $Pin }
Invoke-Checked { git -C $Cache sparse-checkout init --cone }
Invoke-Checked { git -C $Cache sparse-checkout set src }
Invoke-Checked { git -C $Cache checkout --force --detach $Pin }

if (Test-Path $Output) {
    Remove-Item -Recurse -Force $Output
}

Push-Location $Root
try {
    Invoke-Checked {
        go run ./tools/ejdict-import/script `
            -source (Join-Path $Cache "src") `
            -output $Output `
            -revision $Pin
    }
}
finally {
    Pop-Location
}

Write-Host "OK ejdict-import: $Output"
