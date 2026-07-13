#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
INPUT="$ROOT/docs/assets/demo/im-opportunity-radar-demo.mp4"
OUTPUT="$ROOT/docs/assets/demo/im-opportunity-radar-demo.gif"
PALETTE="$(mktemp --suffix=.png)"
trap 'rm -f "$PALETTE"' EXIT
command -v ffmpeg >/dev/null || { echo "ffmpeg is required" >&2; exit 1; }
test -f "$INPUT" || { echo "missing $INPUT" >&2; exit 1; }
ffmpeg -hide_banner -loglevel error -y -ss 2 -t 12 -i "$INPUT" -vf "fps=10,scale=960:-1:flags=lanczos,palettegen=max_colors=128" "$PALETTE"
ffmpeg -hide_banner -loglevel error -y -ss 2 -t 12 -i "$INPUT" -i "$PALETTE" -lavfi "fps=10,scale=960:-1:flags=lanczos[x];[x][1:v]paletteuse=dither=bayer:bayer_scale=4" "$OUTPUT"
