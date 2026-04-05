from PIL import Image
import numpy as np

src = r"c:\go-workspace\src\bitbucket.org\Local\games\ClonedPP\assets\images\player\chat\PP walk side.png"
img = Image.open(src)
arr = np.array(img)
h, w = arr.shape[:2]

alpha = arr[:, :, 3]

print(f"Image: {w}x{h}")

# Scan the right edge to find where the last character's content ends cleanly
# Look from the right for the first column with significant content
for c in range(w-1, 0, -1):
    density = (alpha[:, c] > 10).mean()
    if density > 0.01:
        print(f"Rightmost content at col {c}, density={density:.4f}")
        break

# The image may have been generated with characters spaced for 7 equal slots
# but the rightmost character bleeds past the canvas edge.
# Let's try equal slicing at different frame counts and see what's clean.

# If we treat it as 8 frames (192px each): 8*192 = 1536 -- exact fit
# If 7 frames: would need 1536/7 = 219.43 -- doesn't divide evenly

# Let me try 8-frame grid and extract each for visual check
frame_w = w // 8  # 192px
print(f"\n8-frame grid: {frame_w}px per frame")
for i in range(8):
    x0 = i * frame_w
    x1 = (i + 1) * frame_w
    frame_alpha = alpha[:, x0:x1]
    content_ratio = (frame_alpha > 10).mean()
    print(f"  Frame {i} ({x0}-{x1}): content={content_ratio:.3f}")

# Try 7-frame grid (truncated)
frame_w7 = w // 7  # 219px
print(f"\n7-frame grid: {frame_w7}px per frame (truncates {w - 7*frame_w7}px)")
for i in range(7):
    x0 = i * frame_w7
    x1 = (i + 1) * frame_w7
    frame_alpha = alpha[:, x0:x1]
    content_ratio = (frame_alpha > 10).mean()
    print(f"  Frame {i} ({x0}-{x1}): content={content_ratio:.3f}")

# Try cropping to first 6 poses: if spacing ~220px, that's 1320px for 6 frames
# 1320 / 6 = 220 each
# Or crop to 1200 for 6*200, 1050 for 7*150...
# Key insight: let's just cut off the rightmost ~200px and see if we get a clean 6-frame strip
crop_w = 1320
cropped = img.crop((0, 0, crop_w, h))
cropped.save(r"c:\go-workspace\src\bitbucket.org\Local\games\ClonedPP\assets\images\player\chat\test_6frame.png")
print(f"\nSaved test 6-frame crop: {crop_w}x{h} ({crop_w//6}px per frame)")
