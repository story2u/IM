#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
for tool in adb java; do command -v "$tool" >/dev/null || { echo "$tool is required; install Android SDK and JDK 17" >&2; exit 1; }; done
test -n "${ANDROID_HOME:-${ANDROID_SDK_ROOT:-}}" || { echo "ANDROID_HOME or ANDROID_SDK_ROOT is required" >&2; exit 1; }
adb get-state >/dev/null 2>&1 || { echo "start a Pixel emulator before running this script" >&2; exit 1; }
OUT="$ROOT/docs/assets/screenshots/android"; mkdir -p "$OUT"
(cd "$ROOT/mobile/android" && ./gradlew --no-daemon assembleDebug)
adb install -r "$ROOT/mobile/android/app/build/outputs/apk/debug/app-debug.apk" >/dev/null
capture() { local screen="$1" file="$2"; adb shell am force-stop com.codeiy.im; adb shell am start -n com.codeiy.im/.MainActivity --es demo-screen "$screen" >/dev/null; sleep 2; adb exec-out screencap -p > "$OUT/$file"; test -s "$OUT/$file"; }
capture login android-login.png; capture dashboard android-dashboard.png; capture opportunity-detail android-opportunity-detail.png; capture settings android-settings.png
