#!/bin/bash
# Build Luau static libraries for loveu Xcode targets.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
XCODE_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
LOVE_ROOT="$(cd "$XCODE_DIR/../.." && pwd)"
LUAU_SRC="$LOVE_ROOT/src/libraries/luau"

if [[ ! -f "$LUAU_SRC/CMakeLists.txt" ]]; then
  echo "error: Luau sources missing at $LUAU_SRC" >&2
  exit 1
fi

PLATFORM_NAME="${PLATFORM_NAME:-macosx}"
CONFIGURATION="${CONFIGURATION:-Release}"
SDKROOT="${SDKROOT:-}"
ARCHS="${ARCHS:-}"
DEPLOYMENT_TARGET="${MACOSX_DEPLOYMENT_TARGET:-${IPHONEOS_DEPLOYMENT_TARGET:-11.0}}"

OUT_DIR="${LUAU_OUT_DIR:-$XCODE_DIR/build/luau-$PLATFORM_NAME}"
BUILD_DIR="$OUT_DIR/cmake"

mkdir -p "$BUILD_DIR" "$OUT_DIR"

# Prefer Ninja when available (fast, predictable); fall back to Unix Makefiles.
GENERATOR=()
if command -v ninja >/dev/null 2>&1; then
  GENERATOR=(-G Ninja)
else
  GENERATOR=(-G "Unix Makefiles")
fi

CMAKE_ARGS=(
  "${GENERATOR[@]}"
  -S "$LUAU_SRC"
  -B "$BUILD_DIR"
  -DCMAKE_BUILD_TYPE=Release
  -DLUAU_BUILD_CLI=OFF
  -DLUAU_BUILD_TESTS=OFF
  -DLUAU_BUILD_WEB=OFF
  -DLUAU_EXTERN_C=ON
)

if [[ -n "$SDKROOT" ]]; then
  CMAKE_ARGS+=(-DCMAKE_OSX_SYSROOT="$SDKROOT")
fi

# Prefer single-arch for simulator/device clarity when Xcode provides CURRENT_ARCH.
if [[ -n "${CURRENT_ARCH:-}" && "$CURRENT_ARCH" != "undefined_arch" ]]; then
  CMAKE_ARGS+=(-DCMAKE_OSX_ARCHITECTURES="$CURRENT_ARCH")
elif [[ -n "$ARCHS" ]]; then
  # Xcode passes space-separated archs; CMake wants semicolon-separated.
  CMAKE_ARCHS="$(echo "$ARCHS" | tr ' ' ';')"
  CMAKE_ARGS+=(-DCMAKE_OSX_ARCHITECTURES="$CMAKE_ARCHS")
fi

case "$PLATFORM_NAME" in
  iphoneos|iphonesimulator)
    CMAKE_ARGS+=(-DCMAKE_OSX_DEPLOYMENT_TARGET="$DEPLOYMENT_TARGET")
    CMAKE_ARGS+=(-DCMAKE_SYSTEM_NAME=iOS)
    ;;
  macosx|*)
    CMAKE_ARGS+=(-DCMAKE_OSX_DEPLOYMENT_TARGET="$DEPLOYMENT_TARGET")
    ;;
esac

echo "Configuring Luau for $PLATFORM_NAME..."
cmake "${CMAKE_ARGS[@]}"
cmake --build "$BUILD_DIR" --config Release --target Luau.VM Luau.Compiler Luau.Ast Luau.Bytecode Luau.Common -j"$(sysctl -n hw.ncpu 2>/dev/null || echo 4)"

# Copy/link archives into a flat directory for OTHER_LDFLAGS -l search.
shopt -s nullglob
copied=0
while IFS= read -r lib; do
  base="$(basename "$lib")"
  cp -f "$lib" "$OUT_DIR/$base"
  copied=$((copied + 1))
done < <(find "$BUILD_DIR" -name 'libLuau.*.a' 2>/dev/null)

if [[ "$copied" -eq 0 ]]; then
  echo "error: no libLuau.*.a produced under $BUILD_DIR" >&2
  find "$BUILD_DIR" -name '*.a' 2>/dev/null | head -50 >&2 || true
  exit 1
fi

echo "Luau libraries installed to $OUT_DIR ($copied archives)"
ls -la "$OUT_DIR"/libLuau.*.a
