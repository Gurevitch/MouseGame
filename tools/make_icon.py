"""Make a square window/taskbar icon from PP's idle-front sprite.

SDL2 `window.SetIcon(surface)` on Windows uses the surface as the taskbar
icon. Source PNGs in the repo are full sprite sheets (e.g. 1376x768); pass
one to SetIcon directly and Windows displays a tiny slice of the sheet
because it can't auto-crop. This tool takes one clean pose from the sheet,
strips the background, centers it on a transparent 256x256 canvas, and
writes it to `assets/images/pp_icon.png` for main.go to load.
"""

from pathlib import Path
from PIL import Image, ImageDraw
import numpy as np


REPO = Path(__file__).resolve().parent.parent
SRC = REPO / "assets" / "images" / "player" / "PP idle front.png"
OUT = REPO / "assets" / "images" / "pp_icon.png"

# Sheet is 1376x768 laid out as 8x2 (172x384 per cell). Grab the first pose.
GRID_COLS = 8
GRID_ROWS = 2
TARGET = 256  # icon canvas size (square)


def flood_key(img: Image.Image, tol: int = 8) -> Image.Image:
    img = img.convert("RGBA")
    w, h = img.size
    rgb = img.convert("RGB").copy()
    sentinel = (255, 0, 255)
    from collections import Counter
    edge = []
    for x in range(w):
        edge.append(rgb.getpixel((x, 0)))
        edge.append(rgb.getpixel((x, h - 1)))
    for y in range(h):
        edge.append(rgb.getpixel((0, y)))
        edge.append(rgb.getpixel((w - 1, y)))
    ref = Counter(edge).most_common(1)[0][0]
    for x in range(w):
        for y in (0, h - 1):
            p = rgb.getpixel((x, y))
            if p == sentinel:
                continue
            if max(abs(p[0] - ref[0]), abs(p[1] - ref[1]), abs(p[2] - ref[2])) <= tol:
                ImageDraw.floodfill(rgb, (x, y), sentinel, thresh=tol)
    for y in range(h):
        for x in (0, w - 1):
            p = rgb.getpixel((x, y))
            if p == sentinel:
                continue
            if max(abs(p[0] - ref[0]), abs(p[1] - ref[1]), abs(p[2] - ref[2])) <= tol:
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


def main():
    sheet = Image.open(SRC)
    w, h = sheet.size
    cw, ch = w // GRID_COLS, h // GRID_ROWS
    # Pose 0 of row 0 — PP standing facing camera, neutral.
    pose = sheet.crop((0, 0, cw, ch))
    pose = flood_key(pose, tol=8)
    # Crop to content bbox so the character fills the icon.
    bbox = pose.getbbox()
    if bbox:
        pose = pose.crop(bbox)
    # Scale to fit TARGET square while preserving aspect.
    pw, ph = pose.size
    scale = min(TARGET / pw, TARGET / ph) * 0.9  # 10% margin
    new_w = max(1, int(pw * scale))
    new_h = max(1, int(ph * scale))
    pose = pose.resize((new_w, new_h), Image.LANCZOS)
    canvas = Image.new("RGBA", (TARGET, TARGET), (0, 0, 0, 0))
    canvas.paste(pose, ((TARGET - new_w) // 2, (TARGET - new_h) // 2), pose)
    canvas.save(OUT, "PNG")
    print(f"wrote {OUT} ({canvas.size}, from pose {cw}x{ch} -> {new_w}x{new_h})")


if __name__ == "__main__":
    main()
