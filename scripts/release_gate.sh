#!/usr/bin/env bash
# Release validation gate for the standalone Go + Next.js IM project.
# This is the neutral entrypoint for product-local readiness evidence.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
GO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

export ARTIFACT_DIR="${ARTIFACT_DIR:-$GO_ROOT/tmp/release-gate}"

"$SCRIPT_DIR/phase1_gate.sh" "$@"

if [[ -f "$ARTIFACT_DIR/phase1_gate_manifest.json" ]]; then
  cp "$ARTIFACT_DIR/phase1_gate_manifest.json" "$ARTIFACT_DIR/release_gate_manifest.json"
fi
