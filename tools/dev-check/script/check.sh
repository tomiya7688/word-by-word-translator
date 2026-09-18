#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$ROOT"

echo "Kagome module sums"
go mod download -json github.com/ikawaha/kagome/v2@v2.9.9

mapfile -t GO_FILES < <(git ls-files '*.go')
if (( ${#GO_FILES[@]} > 0 )); then
  UNFORMATTED="$(gofmt -l "${GO_FILES[@]}")"
  if [[ -n "$UNFORMATTED" ]]; then
    echo "FAIL gofmt"
    printf '%s\n' "$UNFORMATTED"
    exit 1
  fi

  echo "go vet"
  go vet ./...

  echo "go test"
  go test ./...
else
  echo "go: no source yet; format/vet/test skipped"
fi

echo "UPD checker"
bash "$ROOT/tools/upd-commander-checker/script/run.sh"

echo "OK"
