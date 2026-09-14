#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$ROOT"

mapfile -t GO_FILES < <(git ls-files '*.go')
if (( ${#GO_FILES[@]} > 0 )); then
  UNFORMATTED="$(gofmt -l "${GO_FILES[@]}")"
  if [[ -n "$UNFORMATTED" ]]; then
    echo "FAIL gofmt"
    printf '%s\n' "$UNFORMATTED"
    exit 1
  fi
fi

echo "go vet"
go vet ./...

echo "go test"
go test ./...

echo "UPD checker"
bash "$ROOT/tools/upd-commander-checker/script/run.sh"

echo "OK"
