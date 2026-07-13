#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
PORT="${DEMO_PORT:-3100}"
cd "$ROOT/frontend"
DEMO_MODE=true NEXT_PUBLIC_DEMO_MODE=true corepack pnpm@10.25.0 dev --hostname 127.0.0.1 --port "$PORT" > /tmp/opportunity-radar-demo-web.log 2>&1 &
SERVER_PID=$!
trap 'kill "$SERVER_PID" >/dev/null 2>&1 || true' EXIT
for _ in $(seq 1 60); do
  if curl -fsS "http://127.0.0.1:$PORT/" >/dev/null; then break; fi
  sleep 1
done
curl -fsS "http://127.0.0.1:$PORT/demo" >/dev/null || { cat /tmp/opportunity-radar-demo-web.log >&2; exit 1; }
DEMO_BASE_URL="http://127.0.0.1:$PORT" corepack pnpm@10.25.0 tsx e2e/demo/record-demo.ts
