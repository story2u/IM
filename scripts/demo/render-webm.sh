#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
INPUT="$ROOT/docs/assets/demo/im-opportunity-radar-demo.mp4"
OUTPUT="$ROOT/docs/assets/demo/im-opportunity-radar-demo.webm"
command -v ffmpeg >/dev/null || { echo "ffmpeg is required" >&2; exit 1; }
test -f "$INPUT" || { echo "missing $INPUT" >&2; exit 1; }
ffmpeg -hide_banner -loglevel error -y -i "$INPUT" -c:v libvpx-vp9 -crf 36 -b:v 0 -row-mt 1 -c:a libopus -b:a 96k "$OUTPUT"
