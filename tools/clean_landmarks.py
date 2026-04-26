"""
Batch flood-fill color-key pass over every landmark icon under
`assets/images/ui/landmarks/`.

Landmark PNGs were authored as RGB with a solid background color, so even
though the runtime now loads them via SafeTextureFromPNGRaw (no runtime
color-key), they still render with a visible background around each icon.

This script runs the same flood-fill-from-every-edge approach used by
`tools/pack_atlas.py:apply_color_key` and overwrites each landmark PNG with
a clean RGBA version whose background has alpha=0. Interior colors (tower
steel, stone, copper) are preserved because flood-fill only reaches pixels
actually connected to an edge pixel of the background color.

Usage:
    python tools/clean_landmarks.py                 # processes every PNG in the folder
    python tools/clean_landmarks.py --dry-run       # report only, don't overwrite
    python tools/clean_landmarks.py some_file.png   # single file

Safe to re-run: already-transparent PNGs are detected and skipped.
"""

from __future__ import annotations

import argparse
import glob
import sys
from collections import Counter
from pathlib import Path

from PIL import Image, ImageDraw
import numpy as np


REPO = Path(__file__).resolve().parent.parent
LANDMARKS = REPO / "assets" / "images" / "ui" / "landmarks"

SENTINEL = (255, 0, 255)


def already_transparent(img: Image.Image) -> bool:
    """True if the image already has an alpha channel AND any pixel is
    transparent. Skips already-processed landmarks."""
    if img.mode != "RGBA":
        return False
    arr = np.array(img)
    return bool((arr[:, :, 3] == 0).any())


def apply_color_key(img: Image.Image, tol: int = 12) -> Image.Image:
    """Seed a flood-fill at every edge pixel that matches the most-common
    edge color. Returns an RGBA image with the flooded region made
    transparent. Identical logic to pack_atlas.py for consistency."""
    img = img.convert("RGBA")
    w, h = img.size
    rgb = img.convert("RGB").copy()

    edge: list[tuple[int, int, int]] = []
    for x in range(w):
        edge.append(rgb.getpixel((x, 0)))
        edge.append(rgb.getpixel((x, h - 1)))
    for y in range(h):
        edge.append(rgb.getpixel((0, y)))
        edge.append(rgb.getpixel((w - 1, y)))
    ref = Counter(edge).most_common(1)[0][0]

    def close(p):
        return (
            abs(p[0] - ref[0]) <= tol
            and abs(p[1] - ref[1]) <= tol
            and abs(p[2] - ref[2]) <= tol
        )

    for x in range(w):
        for y in (0, h - 1):
            if rgb.getpixel((x, y)) == SENTINEL:
                continue
            if close(rgb.getpixel((x, y))):
                ImageDraw.floodfill(rgb, (x, y), SENTINEL, thresh=tol)
    for y in range(h):
        for x in (0, w - 1):
            if rgb.getpixel((x, y)) == SENTINEL:
                continue
            if close(rgb.getpixel((x, y))):
                ImageDraw.floodfill(rgb, (x, y), SENTINEL, thresh=tol)

    arr = np.array(img)
    rgb_arr = np.array(rgb)
    mask = (
        (rgb_arr[:, :, 0] == SENTINEL[0])
        & (rgb_arr[:, :, 1] == SENTINEL[1])
        & (rgb_arr[:, :, 2] == SENTINEL[2])
    )
    arr[mask, 3] = 0
    return Image.fromarray(arr, "RGBA")


def process(path: Path, dry_run: bool) -> None:
    im = Image.open(path)
    if already_transparent(im):
        print(f"  skip (already transparent): {path.name}")
        return
    cleaned = apply_color_key(im, tol=12)
    arr = np.array(cleaned)
    stripped = int((arr[:, :, 3] == 0).sum())
    total = arr.shape[0] * arr.shape[1]
    pct = 100 * stripped / total
    if dry_run:
        print(f"  would strip {stripped:,}/{total:,} ({pct:.1f}%) from {path.name}")
        return
    cleaned.save(path, "PNG")
    print(f"  cleaned {path.name}: {stripped:,}/{total:,} ({pct:.1f}%) transparent")


def main() -> None:
    ap = argparse.ArgumentParser()
    ap.add_argument("files", nargs="*", help="optional specific PNG files")
    ap.add_argument("--dry-run", action="store_true")
    args = ap.parse_args()

    targets: list[Path]
    if args.files:
        targets = [Path(f) for f in args.files]
    else:
        targets = sorted(Path(p) for p in glob.glob(str(LANDMARKS / "*.png")))

    if not targets:
        print("no PNG files found", file=sys.stderr)
        sys.exit(1)

    print(f"processing {len(targets)} files{' (dry run)' if args.dry_run else ''}")
    for p in targets:
        process(p, args.dry_run)


if __name__ == "__main__":
    main()
