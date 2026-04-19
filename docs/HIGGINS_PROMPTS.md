# Director Higgins — Refined Regeneration Prompts

Paste any of the prompts below into your image generator. After generating,
save the output PNG to `assets/images/locations/camp/npc/higgins/<filename>`
and run:

```
python tools/pack_atlas.py tools/characters/higgins.yaml
```

The atlas (`assets/sprites/higgins.png`) refreshes automatically; no Go
changes needed.

---

## Aspect-ratio rule (non-negotiable)

Every sheet uses **tall rectangular cells**, never square. Match the 1997
*Pink Panther: Hokus Pokus Pink* and *Passport to Peril* sprite proportions:
head ~10% of cell height, body ~50%, legs ~35%, a little ankle-to-floor air
at the bottom. A square 512×512 prompt produces stocky "portrait photo"
characters that break this style — do not use square cells.

Existing Higgins sheets in the repo confirm the target format:

| Sheet          | Canvas    | Grid | Cell        |
| -------------- | --------- | ---- | ----------- |
| `idle`         | 1204×384  | 7×1  | **172×384** |
| `talk`         | 1376×768  | 6×1  | **229×768** |
| `walk`         | 1376×768  | 8×1  | **172×768** |
| `walk_back`    | 1376×768  | 8×1  | **172×768** |
| `shout`        | 1376×768  | 8×1  | **172×768** |
| `office_idle`  | 1024×256  | 7×1  | **146×256** |
| `office_talk`  | 1024×256  | 4×2  | **256×128** |
| `give_map`     | 1024×128  | 4×1  | **256×128** |

---

## Shared style lock (paste as a prefix to every prompt)

> **Style:** hand-drawn 1990s Saturday-morning cartoon, Pink Panther
> *Hokus Pokus Pink* (1997) / *Passport to Peril* (1996) look. Confident
> black ink linework around every silhouette (~3 px weight), flat
> saturated fills, no cross-hatching, no gradients inside the character,
> no airbrush, no anime styling, no photorealism, no 3D shading. Simple
> cel shading: one flat mid-tone plus one flat darker tone per color
> region, never more than two tones.
>
> **Background:** pure #FFFFFF white, zero gradient, zero texture, zero
> scenery, zero props. Each cell is separated from its neighbors by a
> thin charcoal 2 px vertical line only — no background art at all.
> Never include speech bubbles, captions, borders, shadows on the white,
> or a second character.
>
> **Character identity (Higgins):** tall thin adult man, early 40s. Round
> wire-rim glasses, short dark-brown hair neatly parted on the side, small
> brown mustache. Fair skin. Wears a full khaki ranger uniform: long-sleeve
> khaki button-down tucked into khaki trousers, brown leather belt, brown
> lace-up ankle boots. No hat. Holds a wooden clipboard in the left hand
> (metal clip, white paper on top). Lanky frame — shoulders narrow,
> silhouette **tall and thin**, not boxy.
>
> **Proportion lock (Hokus Pokus Pink reference):** head is ~10% of cell
> height, torso ~40%, legs ~40%, feet + ground clearance ~10%. Character
> **fills the cell top-to-bottom** with ~4 px inset. Baseline (soles of
> boots) sits at the **same Y across every frame** of the sheet — the
> engine loops frames and a shifting baseline makes him bounce.
>
> **Frame format:** one horizontal row (or 2×4 where noted) of cells at
> the **exact pixel dimensions listed per sheet**. Every frame is
> **portrait-oriented** (taller than wide). Do not crop to a square.

---

## 1. `npc_director_higgins_idle.png` — Entrance idle

**Canvas:** 1204×384. **Grid:** 7 cells × 1 row. **Cell:** 172×384
(tall, portrait orientation).

> [paste style lock above]
>
> **Animation:** 7-frame idle loop. Mouth stays closed in every frame;
> the motion is tiny weight-shift and breath beats, not a walk cycle.
> Same character front-facing in every frame, same costume, clipboard
> in left hand, baseline locked.
>
> - **Frame 1:** neutral standing, weight even, small closed-mouth smile.
> - **Frame 2:** gentle inhale — chest slightly raised, shoulders up 3 px.
> - **Frame 3:** exhale — shoulders drop 3 px below neutral, head tips 2° left.
> - **Frame 4:** glances down at clipboard — eyes roll down, clipboard
>   tilts up 5°. Body otherwise unchanged.
> - **Frame 5:** pushes glasses up the bridge of his nose with the free
>   right hand.
> - **Frame 6:** taps the clipboard edge with the right index finger.
> - **Frame 7:** slow blink — eyes closed or half-lidded. Returns to the
>   neutral silhouette of frame 1 so the loop is seamless.
>
> Output aspect must be **1204 wide × 384 tall**. Each cell is tall and
> narrow (172×384) so the character is rendered full-body, feet-to-head,
> not from the chest up.

---

## 2. `npc_director_higgins_walk.png` — Walk forward (toward camera)

**Canvas:** 1376×768. **Grid:** 8 cells × 1 row. **Cell:** 172×768.

> [paste style lock above]
>
> **Animation:** 8-frame walk cycle, character facing camera (walking
> toward the viewer). Both arms swing naturally; clipboard in left hand
> swings with that arm. Feet contact the same baseline on their contact
> frames so the loop doesn't bob.
>
> - **Frame 1:** contact — right foot forward (heel touches ground),
>   left foot back (toe off). Left arm swung forward.
> - **Frame 2:** down — right foot flat, weight on right, left foot lifting.
> - **Frame 3:** pass — legs cross, right leg bearing weight, left leg
>   mid-swing.
> - **Frame 4:** high — left knee lifted high, right foot planted.
> - **Frame 5:** contact — mirror of frame 1 (left foot forward).
> - **Frame 6:** down — mirror of frame 2.
> - **Frame 7:** pass — mirror of frame 3.
> - **Frame 8:** high — mirror of frame 4.
>
> Mouth closed, eyes forward, uniform unchanged. Full body visible
> top-to-bottom in every tall 172×768 cell.

---

## 3. `npc_director_higgins_walk_back.png` — Walk away (back turned)

**Canvas:** 1376×768. **Grid:** 8 cells × 1 row. **Cell:** 172×768.

> [paste style lock above]
>
> **Animation:** Same 8-frame gait as `walk`, but viewed from **behind**.
> Character walks away from camera, back of head and back of uniform
> visible, arms swinging seen from behind. Clipboard peeks out at the
> left hip on some frames. Timing matches `walk.png` frame for frame
> (contact → down → pass → high, repeat) so the two cycles swap cleanly
> in-engine.

---

## 4. `npc_director_higgins_shout.png` — Attention / shout

**Canvas:** 1376×768. **Grid:** 8 cells × 1 row. **Cell:** 172×768.

> [paste style lock above]
>
> **Animation:** 8-frame "calling campers to attention" loop. Character
> stands in place, facing camera, full-body visible, one arm raised.
> Authoritative, not angry.
>
> - **Frame 1:** mouth closed, right hand halfway up (wind-up).
> - **Frame 2:** right hand fully raised above head, fingers spread,
>   mouth opening.
> - **Frame 3:** mouth wide open mid-shout ("HEY!"), right hand at peak.
> - **Frame 4:** mouth open, head tilted slightly back.
> - **Frame 5:** mouth open, right hand starts lowering, shoulders
>   dropping.
> - **Frame 6:** mouth half-closed, hand mid-way down.
> - **Frame 7:** mouth closed, hand nearly at side, eyes forward.
> - **Frame 8:** fully returned to neutral; loops cleanly into frame 1.
>
> Clipboard stays in the left hand throughout, lowered at the side.
> Full body visible, baseline locked, tall 172×768 cells.

---

## 5. `npc_director_higgins_office_idle.png` — Office idle (seated)

**Canvas:** 1024×256. **Grid:** 7 cells × 1 row. **Cell:** 146×256.

> [paste style lock above] — with these framing overrides:
>
> **Framing override:** seated behind a desk. Each cell is still tall
> relative to its width (146×256) but shows Higgins from the
> **chest up only**. The desk's front edge is a flat horizontal line at
> y=220 in every cell. Background above the desk is pure white.
>
> **Animation:** 7-frame seated idle. Mouth closed, tiny motions only.
>
> - **Frame 1:** neutral, glancing down at paperwork (below cell edge).
> - **Frame 2:** looks up toward camera.
> - **Frame 3:** pushes glasses up nose with right index finger.
> - **Frame 4:** looks back down; right hand re-emerges at desk edge
>   holding a pen.
> - **Frame 5:** small nod, head tilts 3° forward.
> - **Frame 6:** eyes closed (slow blink), returns to neutral.
> - **Frame 7:** thought pose — chin rests on right fist, elbow on desk.
>   Loops cleanly back to frame 1.
>
> No clipboard (office scene). He uses a pen + papers on the desk.

---

## 6. `npc_director_higgins_office_talk.png` — Office talk (seated, 2 rows)

**Canvas:** 1024×256. **Grid:** 4 cells × 2 rows. **Cell:** 256×128.

> [paste style lock above] — chest-up framing, desk edge visible at the
> bottom of every cell, pure white above. Note these cells are **wider
> than tall** (256×128) because the seated pose reads best in landscape
> when torso/desk is the focus.
>
> **Animation:** 8 frames in a 4×2 grid.
>
> **Row 0 — explaining at desk** (4 frames, mouth open):
> 1. Open palm gesture, right hand raised to shoulder height.
> 2. Index finger pointing forward toward camera.
> 3. Thumb-up gesture.
> 4. Flat palm down (settling gesture).
>
> **Row 1 — pointing at map on desk** (4 frames, mouth open, head tilted
> down toward a map below the cell edge):
> 1. Right index finger tapping map (hand visible on desk surface).
> 2. Right index finger tracing a line across the map.
> 3. Right index lifting off the map.
> 4. Right index pointing down firmly at map.
>
> All 8 frames share the same baseline (chin line, desk edge, shoulder
> height). Only mouth shape and hand positions change.

---

## 7. `npc_director_higgins_give_map.png` — Handing the world map to PP

**Canvas:** 1024×128. **Grid:** 4 cells × 1 row. **Cell:** 256×128.

> [paste style lock above] — note these 4 cells are wider than tall
> (256×128) because the handoff reads best as a wide pose-sequence where
> the map fills the width. Full-body is **not** shown here — just torso
> up plus the arms holding the map.
>
> **Animation:** 4-frame handoff, plays once (not a loop). Higgins faces
> camera, torso-up, both hands unrolling a rolled world map and finally
> presenting it forward.
>
> - **Frame 1:** holding a rolled-up scroll/map in both hands at chest
>   height, map still rolled closed.
> - **Frame 2:** hands pulling apart, map partially unrolled, about 1/3
>   of a world map visible (continents in soft blues/greens).
> - **Frame 3:** map fully unrolled across both hands at chest height;
>   full world map readable (continent outlines, no text labels needed).
>   Map edges touch the character's wrists.
> - **Frame 4:** map extended forward toward camera, arms reaching out,
>   offering it to the viewer. Mouth closed (solemn handoff), eyes meet
>   camera.
>
> The map itself reads as a flat parchment with world outlines — simple
> and iconic, not photorealistic.

---

## Re-enabling a sheet in the atlas

Once the new PNG is saved:

1. Open `tools/characters/higgins.yaml`.
2. Uncomment the `idle:` line (or the line for whichever animation was
   regenerated) — or add a new entry if it's a new animation.
3. Run `python tools/pack_atlas.py tools/characters/higgins.yaml`.
4. `assets/sprites/higgins.png` and `.json` refresh; the game picks them
   up on next launch.

If a cell looks chopped in half after regeneration, the grid is wrong —
adjust `grid: [X, Y]` in the yaml and re-pack. No Go code changes.
