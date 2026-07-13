# Demo Asset Runbook

## Prerequisites

- Node 22+, Corepack, Chrome/Chromium, FFmpeg/ffprobe.
- About 4 GB free disk and several minutes for the 4500-frame Remotion render.
- No production credentials are required or accepted by the Demo routes.

## Generate

```bash
make demo-screenshots
make demo-record
make demo-video
# or all stages:
make demo-assets
```

The screenshot command starts a temporary Next dev server with `DEMO_MODE=true`. Recording starts its own
server and closes the browser context before moving Playwright's WebM. Remotion copies only generated demo
inputs, synthesizes a project-owned ambient WAV with FFmpeg, and renders a 150-second Chinese timeline.

## Mobile

```bash
make demo-ios-screenshots       # macOS, Xcode, xcodegen, iPhone 15 Pro simulator
make demo-android-screenshots   # JDK 17, Android SDK, running Pixel emulator
```

Both use Debug-only routes. Release builds ignore the launch argument/Intent extra and retain authentication.

## Safety Review

1. Run `rg -n '(gho_|sk-|session_string|api_hash|@gmail|1[3-9][0-9]{9})' docs/assets demo-video frontend/lib/demo`.
2. Inspect social preview, dashboard, detail, cover and sampled video frames.
3. Run `scripts/demo/check-assets.sh` and FFmpeg black-frame detection.
4. Upload MP4/WebM through a GitHub Release; do not commit repeatedly rendered full videos.

The ambient bed is generated from three sine waves by `scripts/demo/render-demo.sh`; it is original project
output and has no third-party source or license dependency.
