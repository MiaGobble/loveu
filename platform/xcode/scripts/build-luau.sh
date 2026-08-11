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

CMAKE_ARGS=(
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

cmake "${CMAKE_ARGS[@]}"
cmake --build "$BUILD_DIR" --config Release --target Luau.VM Luau.Compiler Luau.Ast Luau.Bytecode Luau.Common -j"$(sysctl -n hw.ncpu 2>/dev/null || echo 4)"

# Copy/link archives into a flat directory for OTHER_LDFLAGS -l search.
shopt -s nullglob
for lib in \
  "$BUILD_DIR"/libLuau.*.a \
  "$BUILD_DIR"/Release/libLuau.*.a \
  "$BUILD_DIR"/*/libLuau.*.a \
  "$BUILD_DIR"/*/*/libLuau.*.a
do
  base="$(basename "$lib")"
  cp -f "$lib" "$OUT_DIR/$base"
done

# Fallback: find by name
while IFS= read -r lib; do
  base="$(basename "$lib")"
  cp -f "$lib" "$OUT_DIR/$base"
done < <(find "$BUILD_DIR" -name 'libLuau.*.a' 2>/dev/null)

echo "Luau libraries installed to $OUT_DIR"
ls -la "$OUT_DIR"/libLuau.*.a
