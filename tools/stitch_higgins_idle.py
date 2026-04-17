"""Stitch 7 generated Higgins pose PNGs into a clean 7x1 idle sheet.

Each source image is an RGB PNG (1376x768) with a baked-in "transparent"
checkerboard (light grey + white) and a full-body Higgins sprite on top.

This tool:
  1. Color-keys checkerboard-grey AND pure-ish white to alpha 0.
  2. Finds the alpha bounding box of the character.
  3. Resizes the bbox to fit in a 256x256 cell, preserving aspect,
     anchored to the bottom (feet align across frames).
  4. Composites all 7 cells side-by-side into a 1792x256 transparent strip.
  5. Writes the strip to the target asset path.

Run from the project root:
    python tools/stitch_higgins_idle.py
"""
from __future__ import annotations

import os
import sys
from pathlib import Path

import numpy as np
from PIL import Image


ROOT = Path(__file__).resolve().parent.parent
POSE_DIR = Path(os.path.expandvars(r"%USERPROFILE%\.cursor\projects\c-go-workspace-src-bitbucket-org-Local-games-ClonedPP\assets"))
TARGET = ROOT / "assets/images/locations/camp/npc/higgins/npc_director_higgins_idle.png"

# Cell geometry matches the Higgins TALK sheet (172x384 per cell). Keeping
# the same per-cell aspect ratio ensures that aspect-preserving in-game
# rendering (see npc.drawScaled) produces the same character height and
# width for idle and talk poses. If cells were square (256x256) here the
# loader's min(scaleW, scaleH) would lock to the width and render Higgins
# at ~200px instead of the ~260px targeted by CHARACTERS.md.
CELL_W = 172
CELL_H = 384
N_CELLS = 7
SHEET_W = CELL_W * N_CELLS  # 1204
SHEET_H = CELL_H
BOTTOM_PAD = 2  # 2px under the feet so the loader inset doesn't crop boots


# Pose 6 is a visual duplicate of pose 1 (neutral), so cell 6 re-uses pose 1.
POSES = [
    ("higgins_pose_1_neutral.png", 1),
    ("higgins_pose_2_inhale.png", 2),
    ("higgins_pose_3_glance.png", 3),
    ("higgins_pose_4_glasses.png", 4),
    ("higgins_pose_5_tap.png", 5),
    ("higgins_pose_1_neutral.png", 6),  # reuse neutral for beat 6
    ("higgins_pose_7_smile.png", 7),
]


def color_key_to_alpha(rgb: np.ndarray) -> np.ndarray:
    """Return an HxWx4 RGBA array where the baked-in checkerboard (grey+white
    squares) is alpha=0 and the character silhouette is alpha=255.

    The image generator produces a fake "transparent" preview by baking a
    checkerboard of light-grey (~195) and pure-white (255) tiles under the
    sprite. Pixel-value-alone keying fails because Higgins' paper, hair, and
    lanyard badge have the SAME neutral value as the checkerboard tiles.

    Strategy: flood-fill from all four image edges across any pixel that
    looks like a checkerboard tile (neutral + light), stopping at any
    saturated pixel. Anything reached = background; anything not reached =
    character (even if it's white/neutral like the paper or hair).
    """
    from collections import deque

    r = rgb[:, :, 0].astype(np.int16)
    g = rgb[:, :, 1].astype(np.int16)
    b = rgb[:, :, 2].astype(np.int16)
    spread = np.maximum(np.maximum(r, g), b) - np.minimum(np.minimum(r, g), b)
    minc = np.minimum(np.minimum(r, g), b)

    # "Looks like a checkerboard tile": neutral (low spread) and light-ish
    # (minc >= 160). Generous to also catch interpolation pixels along tile
    # boundaries, while still rejecting saturated character colors (shirt
    # green spread ~50+, skin spread ~80+, khaki spread ~65+).
    looks_bg = (spread <= 16) & (minc >= 160)

    h, w = looks_bg.shape
    visited = np.zeros_like(looks_bg, dtype=bool)
    q: deque[tuple[int, int]] = deque()

    # Seed from every edge pixel that matches.
    for x in range(w):
        if looks_bg[0, x]:
            q.append((0, x))
            visited[0, x] = True
        if looks_bg[h - 1, x]:
            q.append((h - 1, x))
            visited[h - 1, x] = True
    for y in range(h):
        if looks_bg[y, 0]:
            q.append((y, 0))
            visited[y, 0] = True
        if looks_bg[y, w - 1]:
            q.append((y, w - 1))
            visited[y, w - 1] = True

    # 4-neighbour BFS - vectorised via scipy would be faster but BFS with
    # an explicit mask is ~1s per frame and clearer to reason about.
    while q:
        y, x = q.popleft()
        for ny, nx in ((y - 1, x), (y + 1, x), (y, x - 1), (y, x + 1)):
            if 0 <= ny < h and 0 <= nx < w and not visited[ny, nx] and looks_bg[ny, nx]:
                visited[ny, nx] = True
                q.append((ny, nx))

    bg = visited  # reached from edges over bg-looking pixels
    rgba = np.zeros((h, w, 4), dtype=np.uint8)
    rgba[:, :, :3] = rgb
    rgba[:, :, 3] = np.where(bg, 0, 255).astype(np.uint8)
    return rgba


def alpha_bbox(rgba: np.ndarray) -> tuple[int, int, int, int] | None:
    alpha = rgba[:, :, 3]
    mask = alpha > 16
    if not mask.any():
        return None
    ys, xs = np.where(mask)
    return int(xs.min()), int(ys.min()), int(xs.max()) + 1, int(ys.max()) + 1


def fit_into_cell(rgba_crop: np.ndarray) -> Image.Image:
    """Resize the character crop to fit in a CELL_W x CELL_H transparent cell.

    Preserves aspect ratio. Character is horizontally centered and anchored to
    the bottom with BOTTOM_PAD pixels of space beneath the feet.
    """
    h, w = rgba_crop.shape[:2]
    target_h = CELL_H - BOTTOM_PAD
    # We want the character tall but not overflowing horizontally.
    scale = target_h / h
    new_w = int(round(w * scale))
    if new_w > CELL_W - 4:
        # character is unusually wide (e.g. clipboard up arm extended) - fit width
        scale = (CELL_W - 4) / w
        new_w = CELL_W - 4
    new_h = int(round(h * scale))

    img = Image.fromarray(rgba_crop, mode="RGBA").resize(
        (new_w, new_h), Image.Resampling.LANCZOS
    )

    cell = Image.new("RGBA", (CELL_W, CELL_H), (0, 0, 0, 0))
    dst_x = (CELL_W - new_w) // 2
    dst_y = CELL_H - BOTTOM_PAD - new_h
    if dst_y < 0:
        dst_y = 0
    cell.alpha_composite(img, (dst_x, dst_y))
    return cell


def process_pose(path: Path, cell_idx: int) -> Image.Image:
    if not path.exists():
        raise SystemExit(f"[stitch] missing pose image: {path}")
    src = Image.open(path).convert("RGB")
    rgb = np.array(src)
    rgba = color_key_to_alpha(rgb)
    bbox = alpha_bbox(rgba)
    if bbox is None:
        raise SystemExit(f"[stitch] cell {cell_idx}: no non-transparent pixels in {path.name}")
    x0, y0, x1, y1 = bbox
    crop = rgba[y0:y1, x0:x1]
    cell = fit_into_cell(crop)
    print(
        f"[stitch] cell {cell_idx}: {path.name} bbox=({x0},{y0},{x1},{y1}) "
        f"size={x1 - x0}x{y1 - y0} -> placed"
    )
    return cell


def main() -> None:
    if not POSE_DIR.exists():
        print(f"[stitch] pose directory not found: {POSE_DIR}", file=sys.stderr)
        sys.exit(1)

    TARGET.parent.mkdir(parents=True, exist_ok=True)
    sheet = Image.new("RGBA", (SHEET_W, SHEET_H), (0, 0, 0, 0))

    for cell_idx, (fname, pose_num) in enumerate(POSES):
        src_path = POSE_DIR / fname
        cell = process_pose(src_path, cell_idx)
        sheet.alpha_composite(cell, (cell_idx * CELL_W, 0))

    sheet.save(TARGET, format="PNG")
    print(f"[stitch] wrote {TARGET} ({SHEET_W}x{SHEET_H}, 7x1, RGBA)")


if __name__ == "__main__":
    main()
