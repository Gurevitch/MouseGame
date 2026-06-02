from PIL import Image
import glob, os, sys

def figs(path):
    im = Image.open(path).convert('RGBA'); w, h = im.size; px = im.load(); bg = px[0, 0]
    def empty(c):
        if c[3] < 24: return True
        return abs(c[0]-bg[0]) <= 30 and abs(c[1]-bg[1]) <= 30 and abs(c[2]-bg[2]) <= 30
    rowp = [sum(1 for x in range(0, w, 4) if not empty(px[x, y])) for y in range(h)]
    rr = []; s = None
    for y in range(h+1):
        on = y < h and rowp[y] > 3
        if on and s is None: s = y
        if not on and s is not None: rr.append((s, y-1)); s = None
    rr = [(a, b) for a, b in rr if b-a > h*0.1]
    nrows = max(1, len(rr))
    hh = h//nrows
    col = [sum(1 for y in range(0, hh, 2) if not empty(px[x, y])) for x in range(w)]
    cr = []; s = None
    for x in range(w+1):
        on = x < w and col[x] > 2
        if on and s is None: s = x
        if not on and s is not None: cr.append((s, x-1)); s = None
    big = [(a, b) for a, b in cr if b-a > 20]
    centers = [(a+b)//2 for a, b in big]
    gaps = [centers[i+1]-centers[i] for i in range(len(centers)-1)]
    uniform = (max(gaps)-min(gaps) <= 25) if len(gaps) >= 2 else True
    return w, h, nrows, len(big), uniform, gaps

if __name__ == '__main__':
    paths = sys.argv[1:] or sorted(glob.glob('assets/images/player/PP*.png'))
    for p in paths:
        try:
            w, h, r, c, u, g = figs(p)
            print(f'{os.path.basename(p):30s} {w}x{h} rows={r} figs/row={c} uniform={u} gaps={g}')
        except Exception as e:
            print(p, 'ERR', e)
