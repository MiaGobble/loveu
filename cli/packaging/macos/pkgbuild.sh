#!/usr/bin/env bash
# Build an unsigned macOS .pkg that installs loveu to /usr/local/bin.
# Usage: VERSION=0.1.0 BINARY=./loveu ./pkgbuild.sh
set -euo pipefail

VERSION="${VERSION:?VERSION required}"
BINARY="${BINARY:?BINARY path required}"
OUT="${OUT:-loveu-${VERSION}-macos.pkg}"
ROOT="$(mktemp -d)"
trap 'rm -rf "$ROOT"' EXIT

mkdir -p "$ROOT/usr/local/bin"
cp "$BINARY" "$ROOT/usr/local/bin/loveu"
chmod 755 "$ROOT/usr/local/bin/loveu"

pkgbuild \
  --root "$ROOT" \
  --identifier "org.loveu.cli" \
  --version "$VERSION" \
  --install-location "/" \
  "$OUT"

echo "wrote $OUT"
