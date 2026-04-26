"""Clean up a freshly-generated sprite sheet before handing it to pack_atlas.

The image generator often bakes a thick black outer frame and 1-2px black
grid lines between cells. Those confuse pack_atlas.apply_color_key because:

  - The outer frame makes edge pixels solid black, so the "most common edge
    color" detection picks black as the background and the flood-fill
    would chew through the character's own outlines.
  - The interior grid lines sit between adjacent cells. After slice_sheet
    runs they become the first/last column (or row) of neighbouring cells,
    again biasing the edge-color estimate toward black.

Solution: run this cleaner FIRST. For each cell:
  1. Replace the outer 1-2px frame with pure white.
  2. Replace the first/last N columns/rows of each cell (the grid line
     pixels) with pure white.
  3. Leave the interior of each cell (character + any internal white) alone.

pack_atlas then sees a standard sheet with a clean white background and
its existing flood-fill logic works without modification.

Usage:
    python tools/clean_generated_sheet.py <src_png> <dst_png> <cols> <rows> [border_px]

Example:
    python tools/clean_generated_sheet.py \\
        "~/.cursor/.../assets/npc_director_higgins_talk.png" \\
        "assets/images/locations/camp/npc/higgins/npc_director_higgins_talk.png" \\
        8 2 3
"""
from __future__ import annotations

import argparse
import sys
from pathlib import Path

import numpy as np
from PIL import Image


def clean_sheet(src: Path, dst: Path, cols: int, rows: int, border: int = 3) -> None:
    im = Image.open(src).convert("RGB")
    arr = np.array(im)
    h, w = arr.shape[:2]
    if w % cols != 0 or h % rows != 0:
        raise SystemExit(
            f"{src.name}: canvas {w}x{h} does not divide evenly into {cols}x{rows}"
        )
    cw = w // cols
    ch = h // rows

    white = np.array([255, 255, 255], dtype=np.uint8)

    arr[:border, :, :] = white
    arr[h - border:, :, :] = white
    arr[:, :border, :] = white
    arr[:, w - border:, :] = white

    for c in range(1, cols):
        x = c * cw
        arr[:, max(0, x - border):min(w, x + border), :] = white
    for r in range(1, rows):
        y = r * ch
        arr[max(0, y - border):min(h, y + border), :, :] = white

    dst.parent.mkdir(parents=True, exist_ok=True)
    Image.fromarray(arr, "RGB").save(dst, "PNG")
    print(f"cleaned {src} -> {dst} ({cols}x{rows}, border={border}px)")


def main() -> None:
    ap = argparse.ArgumentParser(description="Clean generated sprite sheet.")
    ap.add_argument("src")
    ap.add_argument("dst")
    ap.add_argument("cols", type=int)
    ap.add_argument("rows", type=int)
    ap.add_argument("--border", type=int, default=3,
                    help="Pixels of black frame/grid line to overwrite with white (default 3)")
    args = ap.parse_args()

    src = Path(args.src).expanduser().resolve()
    dst = Path(args.dst).expanduser().resolve()
    if not src.exists():
        raise SystemExit(f"source not found: {src}")
    clean_sheet(src, dst, args.cols, args.rows, args.border)


if __name__ == "__main__":
    main()
