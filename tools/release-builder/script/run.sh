#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

bash "$ROOT/tools/ejdict-import/script/run.sh"
bash "$ROOT/tools/tkg-ja-en-import/script/run.sh"

(
  cd "$ROOT"
  go run ./tools/release-builder/script \
    -dictionaries "$ROOT/dist/dictionaries" \
    -output "$ROOT/dist/releases" \
    -app-license "$ROOT/LICENSE" \
    -analyzer-notices "$ROOT/licenses/analyzers" \
    -tier all
)

echo "OK release-builder: $ROOT/dist/releases"
