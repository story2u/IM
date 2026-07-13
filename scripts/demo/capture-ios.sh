#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
for tool in xcodebuild xcodegen xcrun; do command -v "$tool" >/dev/null || { echo "$tool is required; run this script on macOS with Xcode" >&2; exit 1; }; done
OUT="$ROOT/docs/assets/screenshots/ios"; BUILD="$ROOT/mobile/ios/.build/demo-derived"; DEVICE="iPhone 15 Pro"
mkdir -p "$OUT"; cd "$ROOT/mobile/ios"; xcodegen generate
xcrun simctl boot "$DEVICE" >/dev/null 2>&1 || true
xcodebuild -project OpportunityRadar.xcodeproj -scheme OpportunityRadar -configuration Debug -destination "platform=iOS Simulator,name=$DEVICE" -derivedDataPath "$BUILD" CODE_SIGNING_ALLOWED=NO build
APP="$BUILD/Build/Products/Debug-iphonesimulator/OpportunityRadar.app"
xcrun simctl install booted "$APP"
capture() {
  local screen="$1" file="$2"
  xcrun simctl terminate booted com.codeiy.im >/dev/null 2>&1 || true
  xcrun simctl launch booted com.codeiy.im -demo-screen "$screen"
  sleep 2
  xcrun simctl io booted screenshot "$OUT/$file"
  test -s "$OUT/$file"
}
capture login ios-login.png; capture dashboard ios-dashboard.png; capture opportunity-detail ios-opportunity-detail.png; capture settings ios-settings.png
