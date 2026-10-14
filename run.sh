#!/usr/bin/env bash
set -euo pipefail
root="$(cd "$(dirname "$0")" && pwd)"
cd "$root"

(cd apps/api && go run .) &
api_pid=$!
(pnpm --filter @learnhub/web exec next dev --hostname 127.0.0.1) &
web_pid=$!

cleanup() {
  kill "$api_pid" "$web_pid" 2>/dev/null || true
}
trap cleanup EXIT INT TERM
wait
