"""preview_cut.py <sheet.png> <cols> <rows> [out.png] [inset]
Render each grid cell's opaque-box crop the way the engine slices it, so you
can see whether a declared cols x rows slices cleanly or chops two poses
into one cell."""
import sys
from PIL import Image, ImageDraw

path = sys.argv[1]; cols = int(sys.argv[2]); rows = int(sys.argv[3])
out = sys.argv[4] if len(sys.argv) > 4 else '/tmp/preview_cut.png'
inset = int(sys.argv[5]) if len(sys.argv) > 5 else 3
im = Image.open(path).convert('RGBA'); w, h = im.size; px = im.load(); bg = px[0, 0]
cw = w//cols; ch = h//rows
def opaque(c):
    if c[3] < 24: return False
    return not (abs(c[0]-bg[0]) <= 8 and abs(c[1]-bg[1]) <= 8 and abs(c[2]-bg[2]) <= 8)
TW, TH = 150, 200
canvas = Image.new('RGBA', (cols*TW, rows*TH), (40, 40, 40, 255))
d = ImageDraw.Draw(canvas)
i = 0
for r in range(rows):
    for c in range(cols):
        L = c*cw+inset; T = r*ch+inset; R = (c+1)*cw-inset; B = (r+1)*ch-inset
        minx = miny = 10**9; maxx = maxy = -1
        for y in range(T, B, 2):
            for x in range(L, R, 2):
                if opaque(px[x, y]):
                    minx = min(minx, x); maxx = max(maxx, x)
                    miny = min(miny, y); maxy = max(maxy, y)
        if maxx >= 0:
            crop = im.crop((minx, miny, maxx+1, maxy+1)); crop.thumbnail((TW-8, TH-18))
            canvas.paste(crop, (c*TW+4, r*TH+4))
            d.text((c*TW+4, r*TH+TH-14), f'{i} {maxx-minx}x{maxy-miny}', fill=(255, 255, 0, 255))
        i += 1
canvas.save(out); print('saved', out, 'cell', cw, 'x', ch)
