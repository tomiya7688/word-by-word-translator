#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"
LOCK_FILE="$ROOT/tools/tkg-ja-en-import/source.lock"
OUTPUT="${1:-$ROOT/dist/dictionaries/mit/tkg-ja-en}"
PIN="$(tr -d '[:space:]' < "$LOCK_FILE")"

if [[ -z "$PIN" ]]; then
  echo "FAIL tkg-ja-en-import: empty source.lock" >&2
  exit 1
fi

CACHE_DIR="$ROOT/.tools-cache/tkg-ja-en/$PIN"
SOURCE="$CACHE_DIR/entries_index.json"
URL="https://raw.githubusercontent.com/tkgally/je-dict-1/$PIN/entries_index.json"
mkdir -p "$CACHE_DIR"

if [[ ! -s "$SOURCE" ]]; then
  TMP="$SOURCE.tmp"
  rm -f "$TMP"
  curl --fail --location --retry 3 --output "$TMP" "$URL"
  mv "$TMP" "$SOURCE"
fi

rm -rf "$OUTPUT"
(
  cd "$ROOT"
  go run ./tools/tkg-ja-en-import/script \
    -source "$SOURCE" \
    -output "$OUTPUT" \
    -revision "$PIN"
)

echo "OK tkg-ja-en-import: $OUTPUT"
