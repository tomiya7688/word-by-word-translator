#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
CACHE="$ROOT/.tools-cache/upd-commander-base-design"
REPO="https://github.com/tomiya7688/upd-commander-base-design.git"
PINNED_COMMIT="41143698dda8bf2bd1f985539dce15d82477ff9d"
TOOL="$CACHE/support_tools/go/upd_commander_checker"

mkdir -p "$ROOT/.tools-cache"

if [[ ! -d "$CACHE/.git" ]]; then
  git clone --filter=blob:none --no-checkout "$REPO" "$CACHE"
fi

CURRENT="$(git -C "$CACHE" rev-parse HEAD 2>/dev/null || true)"
if [[ "$CURRENT" != "$PINNED_COMMIT" || ! -d "$TOOL" ]]; then
  git -C "$CACHE" fetch --depth 1 origin "$PINNED_COMMIT"
  git -C "$CACHE" checkout --force --detach "$PINNED_COMMIT"
fi

(
  cd "$TOOL"
  go run ./cmd/upd-commander-check "$ROOT"
)
