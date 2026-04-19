"""
pack_atlas.py — turn per-state sprite PNGs into a single texture atlas.

Reads a YAML manifest that names the source PNGs, their grid layouts, and
animation metadata (fps, optional take_row filter). Produces:

  assets/sprites/<name>.png    — packed atlas, one row per animation
  assets/sprites/<name>.json   — frame rectangles + fps per animation

Usage:
    python tools/pack_atlas.py tools/characters/marcus.yaml
    python tools/pack_atlas.py tools/characters/*.yaml
    python tools/pack_atlas.py  --all               # every manifest in tools/characters/

The tool never regenerates art — it only re-slices and re-packs the pixels
that already live under assets/images/. Visual style is preserved exactly.
"""

from __future__ import annotations

import argparse
import glob
import json
import os
import sys
from pathlib import Path

from PIL import Image, ImageDraw
import numpy as np
import yaml


SCRIPT_DIR = Path(__file__).resolve().parent
REPO_DIR = SCRIPT_DIR.parent
SPRITES_OUT = REPO_DIR / "assets" / "sprites"


def load_manifest(path: Path) -> dict:
    with path.open("r", encoding="utf-8") as f:
        return yaml.safe_load(f)


def apply_color_key(img: Image.Image, tol: int) -> Image.Image:
    """Flood-fill transparency starting from EVERY edge pixel that matches
    the most-common edge color.

    Seeding only from the four corners failed for two cases:
      1. Sheets where a corner is an off-color frame/border (e.g. Higgins'
         office_idle had a black-pixel corner so the flood never started).
      2. Sheets where a thin darker divider line is authored INSIDE a cell
         (Jake's idle has a faint gray column at x=21-24). The divider
         isolates the outer white strip from the inner white strip, and a
         corner-only flood stops at the divider.

    By seeding a flood at every edge pixel whose color is within tol of the
    most-common edge color, each isolated strip gets its own seed. Interior
    whites (eye whites, teeth, shirt highlights) are preserved because they
    aren't reachable from any edge pixel through sub-tol steps.
    """
    img = img.convert("RGBA")
    w, h = img.size
    rgb = img.convert("RGB").copy()
    sentinel = (255, 0, 255)

    # Find the most-common edge color — that's the background.
    edge_pixels: list[tuple[int, int, int]] = []
    for x in range(w):
        edge_pixels.append(rgb.getpixel((x, 0)))
        edge_pixels.append(rgb.getpixel((x, h - 1)))
    for y in range(h):
        edge_pixels.append(rgb.getpixel((0, y)))
        edge_pixels.append(rgb.getpixel((w - 1, y)))
    from collections import Counter
    ref = Counter(edge_pixels).most_common(1)[0][0]

    def close_to_ref(p):
        return (
            abs(p[0] - ref[0]) <= tol
            and abs(p[1] - ref[1]) <= tol
            and abs(p[2] - ref[2]) <= tol
        )

    # Seed floods from every edge pixel that matches the background.
    for x in range(w):
        for y in (0, h - 1):
            if rgb.getpixel((x, y)) == sentinel:
                continue
            if close_to_ref(rgb.getpixel((x, y))):
                ImageDraw.floodfill(rgb, (x, y), sentinel, thresh=tol)
    for y in range(h):
        for x in (0, w - 1):
            if rgb.getpixel((x, y)) == sentinel:
                continue
            if close_to_ref(rgb.getpixel((x, y))):
                ImageDraw.floodfill(rgb, (x, y), sentinel, thresh=tol)

    arr = np.array(img)
    rgb_arr = np.array(rgb)
    mask = (
        (rgb_arr[:, :, 0] == sentinel[0])
        & (rgb_arr[:, :, 1] == sentinel[1])
        & (rgb_arr[:, :, 2] == sentinel[2])
    )
    arr[mask, 3] = 0
    return Image.fromarray(arr, "RGBA")


def opaque_bbox(frame: Image.Image) -> tuple[int, int, int, int] | None:
    """Return (left, top, right, bottom) of the frame's non-transparent region,
    or None for a fully-transparent frame."""
    arr = np.array(frame)
    alpha = arr[:, :, 3]
    rows_with_content = np.where(np.any(alpha > 0, axis=1))[0]
    cols_with_content = np.where(np.any(alpha > 0, axis=0))[0]
    if len(rows_with_content) == 0 or len(cols_with_content) == 0:
        return None
    top = int(rows_with_content[0])
    bottom = int(rows_with_content[-1])
    left = int(cols_with_content[0])
    right = int(cols_with_content[-1])
    return (left, top, right, bottom)


def baseline_align(frames: list[Image.Image], baseline_margin: int = 8) -> list[Image.Image]:
    """Shift each frame vertically so its character's bottom (lowest opaque
    pixel) sits at frame_height - baseline_margin. Frame size is preserved;
    only content moves. Fully-transparent frames pass through unchanged.

    This fixes the 'rows don't stand on the same y line' problem: source art
    for idle and strange_idle is drawn with different top padding, so when the
    packer pastes rows directly into an atlas, the characters' feet sit at
    different atlas-y positions. Post-alignment every frame places feet at the
    same relative position, so the runtime can blit whole cells without
    per-animation anchor metadata.
    """
    aligned: list[Image.Image] = []
    for f in frames:
        bbox = opaque_bbox(f)
        if bbox is None:
            aligned.append(f)
            continue
        _, _, _, bbox_bottom = bbox
        target_bottom = f.height - 1 - baseline_margin
        shift = target_bottom - bbox_bottom
        if shift == 0:
            aligned.append(f)
            continue
        canvas = Image.new("RGBA", f.size, (0, 0, 0, 0))
        if shift > 0:
            canvas.paste(f, (0, shift))
        else:
            # Character drawn lower than target — lift by cropping top and
            # pasting at y=0. Never happens if baseline_margin is small.
            cropped = f.crop((0, -shift, f.width, f.height))
            canvas.paste(cropped, (0, 0))
        aligned.append(canvas)
    return aligned


def slice_sheet(img: Image.Image, cols: int, rows: int, take_row: int | None) -> list[Image.Image]:
    """Slice an image into cols*rows cells, row-major. If take_row is set,
    only that row's cells are returned."""
    w, h = img.size
    cw, ch = w // cols, h // rows
    frames: list[Image.Image] = []
    row_range = range(rows) if take_row is None else [take_row]
    for r in row_range:
        for c in range(cols):
            box = (c * cw, r * ch, c * cw + cw, r * ch + ch)
            frames.append(img.crop(box))
    return frames


def pack(manifest: dict, manifest_path: Path) -> tuple[Image.Image, dict]:
    """Build the atlas image and its metadata dict from a manifest.

    Alignment pipeline (default on):
      1. Slice source sheets into frames, color-key the background.
      2. Baseline-align within each animation: every frame shifted so its
         character's bottom lands at cell_h - 1 - baseline_margin.
      3. Normalize row heights: pick max frame height across all animations
         and pad shorter frames at the top so feet line up across rows.
      4. Pack as one row per animation, stacked vertically.

    Disable for a given character with `baseline_align: false` at the manifest
    root or per animation.
    """
    src_dir = (manifest_path.parent / manifest["source_dir"]).resolve()

    rows_meta: list[dict] = []
    max_row_w = 0

    # Default tolerance 8 matches engine.applyColorKey (non-kid sheets). Kid
    # sheets override this to 16 in their manifest because the pastel pinks
    # and beiges need a wider flood threshold to fully clear the backdrop.
    default_tol = int(manifest.get("color_key_tolerance", 8))
    baseline_margin = int(manifest.get("baseline_margin", 8))
    align_enabled = bool(manifest.get("baseline_align", True))

    for anim_name, anim_spec in manifest["animations"].items():
        src_path = src_dir / anim_spec["file"]
        if not src_path.exists():
            raise FileNotFoundError(f"{anim_name}: missing {src_path}")

        tol = int(anim_spec.get("color_key_tolerance", default_tol))
        cols, rows = anim_spec["grid"]
        take_row = anim_spec.get("take_row")
        anim_margin = int(anim_spec.get("baseline_margin", baseline_margin))
        do_align = anim_spec.get("baseline_align", align_enabled)

        # Slice FIRST, then color-key each cell independently. Whole-sheet
        # flood-fill was leaving 80%+ of some sheets opaque because a tall
        # character in every cell blocks the flood from reaching the
        # inter-character background strips. Per-cell keying isolates each
        # cell so the 4 corners reliably sit on background.
        raw_img = Image.open(src_path)
        frames = slice_sheet(raw_img, cols, rows, take_row)
        frames = [apply_color_key(f, tol) for f in frames]
        if do_align:
            frames = baseline_align(frames, baseline_margin=anim_margin)

        fw, fh = frames[0].size
        row_w = fw * len(frames)
        rows_meta.append({
            "name": anim_name,
            "frames": frames,
            "fw": fw,
            "fh": fh,
            "fps": anim_spec.get("fps", 8),
            "row_w": row_w,
            "aligned": do_align,
            "margin": anim_margin,
        })
        max_row_w = max(max_row_w, row_w)

    # Cross-row height normalization used to run here. It padded every
    # animation row to the tallest one in the character so the atlas PNG
    # looked aligned, but that ballooned frame dimensions — e.g. adding a
    # 768-tall static prop (Higgins "desk") forced every 384-tall walk row
    # to 768, and the runtime scale-to-bounds then rendered the character at
    # half size. Runtime doesn't need cross-row alignment (drawScaled
    # bottom-anchors each frame to the NPC's bounds), so the padding was
    # pure cost. Keeping rows at their natural heights.
    total_h = sum(rm["fh"] for rm in rows_meta)

    atlas = Image.new("RGBA", (max_row_w, total_h), (0, 0, 0, 0))
    meta = {"animations": {}}
    y = 0
    for rm in rows_meta:
        x = 0
        frame_rects = []
        for frame in rm["frames"]:
            atlas.paste(frame, (x, y))
            frame_rects.append({"x": x, "y": y, "w": rm["fw"], "h": rm["fh"]})
            x += rm["fw"]
        meta["animations"][rm["name"]] = {
            "fps": rm["fps"],
            "frame_w": rm["fw"],
            "frame_h": rm["fh"],
            "frames": frame_rects,
        }
        y += rm["fh"]

    return atlas, meta


def emit(atlas: Image.Image, meta: dict, name: str, subfolder: str = "") -> tuple[Path, Path]:
    out_dir = SPRITES_OUT if not subfolder else SPRITES_OUT / subfolder
    out_dir.mkdir(parents=True, exist_ok=True)
    png_path = out_dir / f"{name}.png"
    json_path = out_dir / f"{name}.json"

    atlas.save(png_path, "PNG", optimize=True)

    meta_out = {
        "image": f"{name}.png",
        "animations": meta["animations"],
    }
    with json_path.open("w", encoding="utf-8") as f:
        json.dump(meta_out, f, indent=2)

    return png_path, json_path


def process(manifest_path: Path) -> None:
    manifest = load_manifest(manifest_path)
    name = manifest["name"]
    subfolder = manifest.get("subfolder", "")

    atlas, meta = pack(manifest, manifest_path)
    png_path, json_path = emit(atlas, meta, name, subfolder)

    frame_count = sum(len(a["frames"]) for a in meta["animations"].values())
    aw, ah = atlas.size
    size_kb = png_path.stat().st_size // 1024
    print(
        f"{name}: {frame_count} frames across {len(meta['animations'])} anims, "
        f"{aw}x{ah} atlas, {size_kb} KB"
    )
    print(f"  -> {png_path.relative_to(REPO_DIR)}")
    print(f"  -> {json_path.relative_to(REPO_DIR)}")


def main():
    ap = argparse.ArgumentParser(description="Pack sprite sheets into atlases.")
    ap.add_argument("manifests", nargs="*", help="Manifest YAML files (or glob patterns)")
    ap.add_argument("--all", action="store_true", help="Pack every manifest in tools/characters/")
    args = ap.parse_args()

    paths: list[Path] = []
    if args.all:
        paths = sorted(Path(p) for p in glob.glob(str(SCRIPT_DIR / "characters" / "**" / "*.yaml"), recursive=True))
    else:
        for m in args.manifests:
            matched = sorted(glob.glob(m))
            if matched:
                paths.extend(Path(p) for p in matched)
            else:
                paths.append(Path(m))

    if not paths:
        ap.error("no manifests given (use --all or pass paths)")

    for p in paths:
        try:
            process(p)
        except Exception as e:
            print(f"{p}: FAILED — {e}", file=sys.stderr)
            sys.exit(2)


if __name__ == "__main__":
    main()
