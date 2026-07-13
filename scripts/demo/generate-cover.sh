#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
INPUT="$ROOT/docs/assets/screenshots/web/demo-video-cover.png"
OUTPUT="$ROOT/docs/assets/demo/im-opportunity-radar-cover.png"
command -v ffmpeg >/dev/null || { echo "ffmpeg is required" >&2; exit 1; }
test -f "$INPUT" || { echo "missing $INPUT" >&2; exit 1; }
ffmpeg -hide_banner -loglevel error -y -i "$INPUT" -vf "scale=1920:1080:flags=lanczos" -frames:v 1 "$OUTPUT"
