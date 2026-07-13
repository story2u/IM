#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
required=(
  docs/assets/screenshots/web/web-home-desktop.png
  docs/assets/screenshots/web/dashboard-desktop.png
  docs/assets/screenshots/web/opportunity-detail.png
  docs/assets/demo/raw/web-demo.webm
  docs/assets/demo/im-opportunity-radar-demo.mp4
  docs/assets/demo/im-opportunity-radar-demo.webm
  docs/assets/demo/im-opportunity-radar-demo.gif
  docs/assets/demo/im-opportunity-radar-cover.png
)
for rel in "${required[@]}"; do test -s "$ROOT/$rel" || { echo "missing or empty: $rel" >&2; exit 1; }; done
echo "Demo asset manifest:"
for rel in "${required[@]}"; do size=$(du -h "$ROOT/$rel" | cut -f1); printf '%-68s %s\n' "$rel" "$size"; done
ffprobe -v error -show_entries format=duration,size -of default=nw=1 "$ROOT/docs/assets/demo/im-opportunity-radar-demo.mp4"
