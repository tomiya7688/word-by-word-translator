#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
LOCK_FILE="$ROOT/tools/ejdict-import/source.lock"
CACHE="$ROOT/.tools-cache/ejdict-hand"
REPO="https://github.com/kujirahand/EJDict.git"
OUTPUT="${1:-$ROOT/dist/dictionaries/mit/ejdict-hand}"
PIN="$(tr -d '[:space:]' < "$LOCK_FILE")"

if [[ -z "$PIN" ]]; then
  echo "FAIL ejdict-import: empty source.lock" >&2
  exit 1
fi

if [[ ! -d "$CACHE/.git" ]]; then
  rm -rf "$CACHE"
  git clone --filter=blob:none --no-checkout "$REPO" "$CACHE"
fi

git -C "$CACHE" fetch --depth 1 origin "$PIN"
git -C "$CACHE" sparse-checkout init --cone
git -C "$CACHE" sparse-checkout set src
git -C "$CACHE" checkout --force --detach "$PIN"

rm -rf "$OUTPUT"
(
  cd "$ROOT"
  go run ./tools/ejdict-import/script \
    -source "$CACHE/src" \
    -output "$OUTPUT" \
    -revision "$PIN"
)

echo "OK ejdict-import: $OUTPUT"
