#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
cd "$ROOT"

REMOTE="${1:-origin}"
BRANCH="${2:-main}"
TARGET="$REMOTE/$BRANCH"
MAX_COMMITS=10
MAX_DIFF_LINES=120

git fetch --prune "$REMOTE" >/dev/null

if ! git rev-parse --verify "$TARGET" >/dev/null 2>&1; then
  echo "ERROR remote ref not found: $TARGET" >&2
  exit 1
fi

LOCAL_HEAD="$(git rev-parse --short HEAD)"
REMOTE_HEAD="$(git rev-parse --short "$TARGET")"
COUNTS="$(git rev-list --left-right --count HEAD..."$TARGET")"
AHEAD="${COUNTS%%$'\t'*}"
BEHIND="${COUNTS##*$'\t'}"

printf 'local=%s remote=%s ahead=%s behind=%s\n' "$LOCAL_HEAD" "$REMOTE_HEAD" "$AHEAD" "$BEHIND"

echo "-- remote commits --"
git log --oneline --no-decorate HEAD.."$TARGET" -n "$MAX_COMMITS" || true

echo "-- changed files --"
git diff --name-status HEAD..."$TARGET" || true

echo "-- diff stat --"
git diff --stat HEAD..."$TARGET" || true

echo "-- bounded diff excerpt (max ${MAX_DIFF_LINES} lines) --"
git diff --no-ext-diff --unified=2 HEAD..."$TARGET" | head -n "$MAX_DIFF_LINES" || true
