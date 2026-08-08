#!/usr/bin/env bash
#
# Download the Windows (mingw) SDL2 development libraries that
# scripts/build_windows.sh links against, into .toolchain/ (gitignored).
set -euo pipefail

SDL2_VERSION="${SDL2_VERSION:-2.32.10}"

cd "$(dirname "$0")/.."

DEST="$PWD/.toolchain/sdl2-mingw"
# Probe for a real header rather than just the directory: .toolchain is
# gitignored, so `git clean -fdx` guts it while leaving empty dirs behind.
if [ -f "$DEST/x86_64-w64-mingw32/include/SDL2/SDL.h" ]; then
	echo "already present: $DEST (delete it to re-download)"
	exit 0
fi
rm -rf "$DEST"

TARBALL="$(mktemp -t sdl2-mingw).tar.gz"
trap 'rm -f "$TARBALL"' EXIT

URL="https://github.com/libsdl-org/SDL/releases/download/release-$SDL2_VERSION/SDL2-devel-$SDL2_VERSION-mingw.tar.gz"
echo "fetching $URL"
curl -fsSL -o "$TARBALL" "$URL"

mkdir -p "$PWD/.toolchain"
tar -xzf "$TARBALL" -C "$PWD/.toolchain"
mv "$PWD/.toolchain/SDL2-$SDL2_VERSION" "$DEST"

echo "installed SDL2 $SDL2_VERSION mingw libs into $DEST"
