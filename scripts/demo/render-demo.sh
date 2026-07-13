#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
command -v ffmpeg >/dev/null || { echo "ffmpeg is required" >&2; exit 1; }
command -v ffprobe >/dev/null || { echo "ffprobe is required" >&2; exit 1; }
test -f "$ROOT/docs/assets/demo/raw/web-demo.webm" || { echo "missing Playwright recording; run make demo-record" >&2; exit 1; }
test -f "$ROOT/docs/assets/screenshots/web/web-home-desktop.png" || { echo "missing screenshots; run make demo-screenshots" >&2; exit 1; }

mkdir -p "$ROOT/demo-video/public/raw" "$ROOT/demo-video/public/screenshots" "$ROOT/demo-video/public/audio" "$ROOT/docs/assets/demo/render"
cp "$ROOT/docs/assets/demo/raw/web-demo.webm" "$ROOT/demo-video/public/raw/web-demo.webm"
cp "$ROOT/docs/assets/screenshots/web/"*.png "$ROOT/demo-video/public/screenshots/"

# Project-owned procedural ambient bed: three quiet sine waves, no third-party recording.
ffmpeg -hide_banner -loglevel error -y \
  -f lavfi -i "sine=frequency=110:sample_rate=48000:duration=150" \
  -f lavfi -i "sine=frequency=165:sample_rate=48000:duration=150" \
  -f lavfi -i "sine=frequency=220:sample_rate=48000:duration=150" \
  -filter_complex "[0:a]volume=0.035[a0];[1:a]volume=0.02[a1];[2:a]volume=0.012[a2];[a0][a1][a2]amix=inputs=3,afade=t=in:st=0:d=3,afade=t=out:st=146:d=4" \
  "$ROOT/demo-video/public/audio/background.wav"

(cd "$ROOT/demo-video" && corepack pnpm@10.25.0 install --frozen-lockfile && corepack pnpm@10.25.0 render)

INPUT="$ROOT/docs/assets/demo/render/product-demo.mp4"
OUTPUT="$ROOT/docs/assets/demo/im-opportunity-radar-demo.mp4"
ffmpeg -hide_banner -loglevel error -y -i "$INPUT" -c:v libx264 -preset medium -crf 24 -pix_fmt yuv420p -r 30 -movflags +faststart -c:a aac -b:a 128k "$OUTPUT"
ffprobe -v error -select_streams v:0 -show_entries stream=codec_name,width,height,r_frame_rate -of default=nw=1 "$OUTPUT"
