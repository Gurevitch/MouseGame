#!/usr/bin/env python3
"""detect_grid.py — infer the frame layout of a sprite sheet.

"Two frames at once" / "frames swiping" bugs come from a declared grid
(cols x rows) that doesn't match how the art was laid out. This reads a PNG,
treats the background (sampled from a corner, or true transparency) as empty,
and reports the number of content "islands" across X (columns) and Y (rows)
by finding the empty gutters between them.

    python tools/detect_grid.py <sheet.png> [more.png ...]

Note: column counts merge when neighbouring poses touch (no gutter), so treat
the result as a hint — uneven gutter spacing means the art needs a clean regen.
"""
import sys
import os
from PIL import Image


def analyze(path):
    im = Image.open(path).convert("RGBA")
    w, h = im.size
    px = im.load()
    bg = px[0, 0]

    def empty(c):
        if c[3] < 24:
            return True
        return abs(c[0] - bg[0]) <= 30 and abs(c[1] - bg[1]) <= 30 and abs(c[2] - bg[2]) <= 30

    def runs(vals, thresh):
        out, s = [], None
        for i, v in enumerate(vals + [0]):
            on = i < len(vals) and v > thresh
            if on and s is None:
                s = i
            if (not on) and s is not None:
                out.append((s, i - 1))
                s = None
        return out

    col = [sum(1 for y in range(0, h, 3) if not empty(px[x, y])) for x in range(w)]
    row = [sum(1 for x in range(0, w, 3) if not empty(px[x, y])) for y in range(h)]
    cruns = runs(col, max(1, (h // 3) // 40))
    rruns = runs(row, max(1, (w // 3) // 40))
    return w, h, cruns, rruns


def main(argv):
    if not argv:
        print(__doc__)
        return 1
    for p in argv:
        if not os.path.exists(p):
            print(f"{p}  MISSING")
            continue
        w, h, c, r = analyze(p)
        print(f"{os.path.basename(p)}  {w}x{h}  cols~{len(c)} rows~{len(r)}  "
              f"(cell @8={w/8:.0f} @7={w/7:.0f} @6={w/6:.0f})  col-gutters={[a for a, _ in c]}")
    return 0


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
