#!/usr/bin/env bash
#
# Cross-compile PP.exe for Windows (x86_64) from macOS/Linux.
#
# go-sdl2 is a cgo package, so a Windows build needs a Windows C toolchain plus
# Windows SDL2 headers and import libraries. Without them cgo is skipped and
# every sdl.* type reads as undefined. Prerequisites:
#
#   brew install mingw-w64
#   scripts/fetch_sdl2_mingw.sh
#
# On Windows itself, build natively instead: GOTMPDIR=.gotmp go build -o PP.exe .
set -euo pipefail

cd "$(dirname "$0")/.."

SDL_PREFIX="$PWD/.toolchain/sdl2-mingw/x86_64-w64-mingw32"
if [ ! -d "$SDL_PREFIX" ]; then
	echo "error: SDL2 mingw libs missing. Run scripts/fetch_sdl2_mingw.sh first." >&2
	exit 1
fi
if ! command -v x86_64-w64-mingw32-gcc >/dev/null; then
	echo "error: x86_64-w64-mingw32-gcc not found. Run: brew install mingw-w64" >&2
	exit 1
fi

# CGO_CFLAGS/CGO_LDFLAGS are set explicitly rather than appended: `go env` here
# carries host -I/opt/homebrew/include and -L/opt/homebrew/lib defaults, and
# feeding those macOS paths to the mingw compiler breaks the build.
#
# Two include paths are needed because go-sdl2 includes <SDL.h> unqualified
# (matching what pkg-config hands it on Unix) while the tarball nests the
# headers under include/SDL2.
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++
export CGO_CFLAGS="-I$SDL_PREFIX/include -I$SDL_PREFIX/include/SDL2 -D_REENTRANT"
export CGO_LDFLAGS="-L$SDL_PREFIX/lib -lSDL2"
export GOTMPDIR="$PWD/.gotmp"
mkdir -p "$GOTMPDIR"

# -H windowsgui suppresses the console window that would otherwise open
# alongside the game.
go build -ldflags "-H windowsgui" -o PP.exe "$@" .

echo "built PP.exe"
echo "note: ship SDL2.dll next to it - copy from"
echo "  $SDL_PREFIX/bin/SDL2.dll"
