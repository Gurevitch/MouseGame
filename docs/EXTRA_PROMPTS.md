# Extra Sprite Prompts — everything still needed for the current FIXME sweep

This file is read by ChatGPT inside Cursor. Each paste-ready prompt is wrapped
between `===PROMPT START===` and `===PROMPT END===` markers. **Workflow:**

1. Highlight everything BETWEEN the markers (the blockquote block itself,
   not the marker lines).
2. Paste into ChatGPT. Include the header below (style lock + standing
   rules) as context if ChatGPT doesn't already have it — those rules
   apply to every prompt in the file.
3. Save the resulting PNG at the path shown in that prompt's section.
4. Run the atlas re-pack (or restart the game for legacy loaders):

```
python tools/pack_atlas.py tools/characters/<name>.yaml
```

5. Move the section header into the **Done log** at the bottom of this
   file and delete the prompt body so the working part stays scannable.

---

**Style lock + standing rules below — feed these to ChatGPT once per session
so it doesn't violate them on the next prompt:**

> Hand-drawn 1990s Saturday-morning cartoon, Pink Panther *Hokus Pokus Pink*
> (1997) / *Passport to Peril* (1996). Confident black ink linework ~3 px,
> flat saturated fills, no cross-hatching, no gradients, no airbrush. Two
> cel tones max per color region. Pure #FFFFFF background, zero scenery.
> Every cell is **tall rectangular**, never square.

Canvas dimensions are locked per sheet; do **not** scale down to square.

**Standing PP design rules (apply to EVERY PP prompt):**

1. **No gloves of any color.** Pink Panther in this game has plain
   pink paws/hands — never yellow gloves, never any gloves.
2. **Every pickup sprite ends with PP pocketing the item.** The final
   1-2 frames show the item vanishing into his invisible hip pocket
   (the classic Pink Panther "magic pocket"); PP ends empty-handed in
   a relaxed standing pose with a small secretive smile.
3. **No pure white anywhere on the panther.** Belly uses ivory
   off-white `#F2EFE5`, eye sclera uses pale grey. Pure white pixels
   on PP get chroma-keyed away by the engine.

**Standing rule for ALL characters who need "white" in their design:**
the engine chroma-keys pure `#FFFFFF` plus a tolerance band. Use these in
order of preference for fabric / large white areas:
- **Cream `#E5DDC8`** ← USE THIS for "white shirts" or any large fabric.
- **Bone `#EDE5D3`** — paper, small label panels.
- **Pale grey `#C4C4C4`** — steam wisps, eye sclera.
- **Vanilla `#F2EFE5`** — only safe for tiny accents (a tooth, a button).

The **scene background** in the sprite cell still uses pure
`#FFFFFF` — that IS the chroma key; it's only the character /
foreground objects that must avoid pure white.

---

## Critical separator rule (applies to EVERY multi-frame sheet)

User 2026-05-24: several recent regens (café patrons, Marcus, Higgins office)
came back with **visible thin lines BETWEEN frame cells** — faint grey or
near-white seams that survive the chroma-key and render in-game as dark
verticals between animation frames.

**Fix language to include in every prompt that uses a grid:**

> The sheet is a **flat grid of cells with NO visible separators**: no
> drawn borders, no thin lines, no grey/black strips, no shadow gradients
> between cells. Cell boundaries are conceptual only — neighbouring cells
> meet directly with pure `#FFFFFF` background pixels on both sides. The
> exported PNG must look like ONE continuous canvas where each Nth × Mth
> rectangle happens to hold one frame; if you cropped any cell out you'd
> see only that frame on pure white, with no edge artefacts.

If you see a faint vertical/horizontal line in the preview, the generator
drew a separator — regenerate with the rule above emphasised.

---

## Open Prompts

All prompts below still need a PNG generated. When one lands, move its row
into the **Done log** at the bottom and delete the body.

---

### A. PP idle — regen to match talk design

**Canvas (front):** 1376×768. **Grid:** 8×2. **Cell:** 172×384.
**Path:** `assets/images/player/PP idle front.png`.

**Canvas (side):** 1672×941. **Grid:** 8×2. **Cell:** 209×470.
**Path:** `assets/images/player/PP idle side.png`.

User 2026-05-24: previously we regenerated TALK to match IDLE; the user has
now reversed the direction — the **talk** sheets are the canonical PP design
(better silhouette, cleaner paws, tail framing) and the **idle** sheets need
to be regenerated to look like the talk sheets but with mouth closed and
breathing-only motion instead of gestures.

Run this prompt twice (once with front canvas, once with side canvas).
Substitute the canvas/grid/cell numbers above for the run you're doing.

===PROMPT START===
> [style lock]
>
> Pink Panther standing in a relaxed idle pose, mouth closed. **Match
> the silhouette, paw shape, tail framing, and color palette of
> `PP talk front.png` / `PP talk side.png` exactly** — same body
> proportions, same tail curve, same pink fill `#E88BB5` with
> `#C4548A` cel shadow, same ivory off-white belly patch, same
> pale-grey eye sclera with black pupils, same `#1C1C1C` ink
> outline ~3 px. **Plain pink paws — NO gloves of any color.**
> **No pure white anywhere on the panther.**
>
> The difference from the talk sheet: mouth stays CLOSED in every
> idle frame (small relaxed smile is fine), and motion is restricted
> to subtle breathing, tail flick, eye blink, weight shift — no
> wide gestures, no raised paws.
>
> **Animation:** 16-frame idle loop (8 columns × 2 rows).
>
> Row 0 (frames 0–7): standard breathing loop. 0 neutral standing
> with mouth closed and small smile, 1 chest expanding slightly
> (breath in), 2 chest at peak, 3 chest releasing (breath out), 4
> back to neutral, 5 tail flick to one side, 6 tail returning, 7
> eye-blink frame (eyelids halfway).
>
> Row 1 (frames 8–15): same eight-beat breathing loop but with the
> weight subtly shifted to the other foot — alternates which leg
> bears weight so the loop has visual variety. Mouth stays closed
> throughout. One additional eye-blink frame somewhere in this row.
>
> **Critical canvas rules:** Cell dimensions exactly as specified
> above. PP's foot baseline locked to the bottom of the cell across
> all 16 frames — he does not change vertical position. Body
> centerline locked on every cell. Tail stays inside the cell. No
> frame is taller or wider than another. **Match `PP talk
> front.png` (or `PP talk side.png` for the side run) exactly for
> body proportions and overall PP size in each cell** so idle ↔
> talk swap is visually seamless.
>
> [SEPARATOR RULE — see top of file. No drawn lines between cells.
> Pure `#FFFFFF` between every neighbouring frame.]
===PROMPT END===

---

### B. Marcus idle — recolor to match the strange-idle palette

**Canvas:** 1376×768. **Grid:** 7×2. **Cell:** 196×384.
**Path:** `assets/images/locations/camp/npc/kids/marcus/npc_marcus_idle.png`.

User 2026-05-24: `npc_marcus_strange_idle.png` has a moodier / darker palette
(the "something-is-wrong" Marcus). The regular `npc_marcus_idle.png` is
brighter and reads like a different kid when the engine swaps from regular →
strange after the inactivity timer fires. Regenerate the regular idle so the
palette MATCHES the strange-idle exactly — same shirt tone, same shorts
tone, same hair tone, same skin tone. Only the pose/expression should differ
between the two sheets.

===PROMPT START===
> [style lock]
>
> Marcus the 10-year-old camp kid in regular relaxed idle. **Match
> the color palette of `npc_marcus_strange_idle.png` exactly** —
> sample the hair, skin, shirt, shorts, socks, shoes, and
> sketchbook colors from that sheet and reuse them here. Round
> wire-rim glasses `#1C1C1C`, small nose-bandage strip, holds his
> spiral sketchbook in his left hand. Black `#1C1C1C` ink outline
> ~3 px. **No pure white on the character** — sketchbook page uses
> bone `#EDE5D3`, socks use cream `#E5DDC8`.
>
> The difference from `npc_marcus_strange_idle.png`: Marcus is
> CALM and present — looking at the camera with a small friendly
> smile, normal posture, sketchbook held casually at his hip.
> NOT staring blankly, NOT clutching the sketchbook to his chest,
> NOT slumped — that's the strange-idle's job. This sheet is the
> "before the wrongness" Marcus.
>
> **Animation:** 14-frame idle loop (7 columns × 2 rows).
>
> Row 0 (frames 0–6): 0 neutral standing facing camera, 1 small
> breath in (chest rises), 2 breath peak, 3 breath out, 4 looks
> down at sketchbook briefly, 5 looks back up, 6 small adjust of
> the glasses.
>
> Row 1 (frames 7–13): 7 weight shift to other foot, 8 breath in,
> 9 breath peak, 10 breath out, 11 eye blink, 12 small head-turn
> to one side, 13 returns to neutral matching frame 0 so the loop
> seams.
>
> **Critical canvas rules:** Every cell exactly 196×384. Marcus's
> foot baseline locked to row pixel 380 in every frame. Body
> centerline matches the strange-idle sheet's centerline so
> swapping between the two sheets is seamless. Pure `#FFFFFF`
> background.
>
> [SEPARATOR RULE — see top of file. No drawn lines between cells.]
===PROMPT END===

---

### C. Higgins office idle + talk — match the entrance Higgins design

**Canvas (each):** 1376×768. **Grid:** 6×2. **Cell:** 229×384.
**Paths:**
- `assets/images/locations/camp/npc/higgins/npc_director_higgins_office_idle.png`
- `assets/images/locations/camp/npc/higgins/npc_director_higgins_office_talk.png`

User 2026-05-24: the office Higgins looks like a different character from
the entrance Higgins (`npc_director_higgins_idle.png` / `npc_director_higgins_talk.png`).
Different face shape, different mustache, different shirt color. Regenerate
BOTH office sheets so the design and color palette match the entrance
Higgins exactly — same lanky build, same round wire-rim glasses, same brown
mustache, same khaki safari shirt with epaulets, same brown belt, same cream
neckerchief. The ONLY difference between entrance and office Higgins is the
framing (standing full-body outside vs seated behind a desk).

===PROMPT START===
> [style lock — paste at the top of both sheet generations]
>
> Director Higgins seated behind his desk in the camp office.
> **Match the face, hair, mustache, glasses, shirt, belt, and
> neckerchief design and colors of `npc_director_higgins_idle.png`
> (the entrance standing-outside Higgins) exactly.** Lanky ranger
> build, round wire-rim glasses `#1C1C1C`, brown mustache
> `#4A2E1B`, short side-parted brown hair, khaki safari shirt
> `#A47148` with epaulets, brown leather belt, cream `#E5DDC8`
> neckerchief. Visible from desk-edge up — chair behind, desk edge
> at the bottom of the cell partially obscures lower torso. Black
> `#1C1C1C` ink outline ~3 px. **No pure white anywhere.**
> Pure `#FFFFFF` background.
>
> The difference from the entrance sheet: Higgins is SEATED behind
> a wooden desk, with both hands resting near the desktop. The
> desk top edge runs horizontally across the bottom of every cell
> at a consistent Y, so the engine can render him behind the
> office desk BG without head-clearance jumps.
>
> **Animation:** 6 frames per row (6 columns × 2 rows). Row 0 =
> idle, Row 1 = talk.
>
> Idle (row 0): 0 neutral with both hands on desk, 1 small
> shoulder shrug, 2 hand to mustache, 3 looks down at desk, 4
> looks up at camera, 5 returns to neutral matching frame 0.
>
> Talk (row 1): 0 mouth open mid-word both hands on desk, 1 right
> hand raised palm-up, 2 finger-point forward, 3 head-tilt
> emphatic, 4 both hands gesturing wide, 5 returns to neutral
> matching frame 0.
>
> **Critical canvas rules:** Every cell exactly 229×384. Higgins's
> chair-line and the desk-top edge locked across ALL 12 frames
> (idle + talk). Body centerline same in both rows so idle ↔ talk
> swap is seamless and matches the entrance sheet's centerline.
>
> [SEPARATOR RULE — see top of file. No drawn lines between cells.
> No grey or near-white seams between frames; the engine reads
> those as dark verticals after chroma-key.]
===PROMPT END===

---

### D. Café patron sheets — clean separator lines

**Canvas (each):** 1376×768. **Grid:** 8×2. **Cell:** 172×384.
**Paths:**
- `assets/images/locations/paris/npc/coffee/cafe_patron_yvette.png`
- `assets/images/locations/paris/npc/coffee/cafe_patron_bernard.png`
- `assets/images/locations/paris/npc/coffee/cafe_patron_camille.png`
- `assets/images/locations/paris/npc/coffee/cafe_patron_henri.png`
- `assets/images/locations/paris/npc/coffee/cafe_patron_lucien.png`

User 2026-05-24: each café patron sheet has visible thin lines BETWEEN the
animation frames — faint near-white / pale-grey seams that survive chroma-key
and render in-game as ugly vertical strips between cells. Regenerate each
patron with the strict separator rule.

Keep the existing pose set (chest-up upper-body, seated at café table, idle
row 0 / talk row 1). The only change vs the on-disk PNGs: **delete the
between-frame seams**.

===PROMPT START===
> [style lock]
>
> Paris café patron seated at a bistro table, viewed from
> chest-up (the table cloth in the BG covers their lower half).
> Match the existing patron design on disk for hair / clothing /
> accessories — this is a clean-up regen, NOT a redesign.
> Black `#1C1C1C` ink outline ~3 px. **No pure white on the
> character** — any "white" cloth uses cream `#E5DDC8`.
>
> **Animation:** 16 frames (8 columns × 2 rows). Row 0 = idle
> breathing + sip + glance loop. Row 1 = talk with mouth open in
> assorted shapes for dialog.
>
> **Critical canvas rules:** Every cell exactly 172×384. Patron's
> shoulder line locked to a consistent Y across all 16 frames so
> idle ↔ talk swap is seamless. Body centerline locked per cell.
>
> [SEPARATOR RULE — CRITICAL FOR THIS REGEN]
> The sheet is a flat grid of frames with **absolutely NO visible
> separators between cells**. No drawn borders, no thin grey/black
> lines, no faint near-white seams, no shadow gradients between
> frames. Cell boundaries are conceptual only — neighbouring cells
> meet directly with pure `#FFFFFF` background pixels on both
> sides. The exported PNG must look like ONE continuous canvas;
> if you cropped any cell out you'd see only that frame on pure
> white, with zero edge artefacts. The previous version of this
> sheet had faint vertical lines between frames — DO NOT
> reproduce them.
===PROMPT END===

---

## Done log — landed sprites (FYI only, no action needed)

These prompts produced PNGs that are now on disk and wired up. Listed for
record so we don't re-generate them by accident. If you need a variant, the
original prompt is in git history at `docs/EXTRA_PROMPTS.md` pre-2026-05-24.

| § | Sprite | Path | Landed |
|---|--------|------|--------|
| §1 | Higgins entrance idle | `npc_director_higgins_idle.png` | 2026-04 |
| §2 | Higgins walk back | `npc_director_higgins_walk_back.png` | 2026-04 |
| §4 | Marcus strange_alt | `npc_marcus_strange_alt.png` | 2026-04 |
| §6 | Campfire small loop | campfire frames | 2026-04 |
| §8 | Bakery Woman | `npc_bakery_woman.png` | 2026-04 |
| §9 | Press Photographer | `npc_press_photographer.png` | 2026-04 |
| §10 | Higgins entrance talk | `npc_director_higgins_talk.png` | 2026-04 |
| §18 | Higgins office idle + talk (v1) | `npc_director_higgins_office_*.png` | 2026-04 — superseded by §C above |
| §19 | Higgins give_map handoff | `npc_director_higgins_give_map.png` | 2026-04 |
| §Y | Paris Bakery BG v2 (door right + tablecloths + framed counter) | `paris_bakery.png` | 2026-05-23 |
| §E | Tommy walk_left | `npc_tommy_walk_left.png` | 2026-05-21 |
| §F | Jake walk_back | `npc_jake_walk_back.png` | 2026-05-21 |
| §M | Action cursor (cursor_point) | `cursor_point.png` | 2026-05-21 |
| §H | PP airplane (modern Cessna-style + pilot) | `pp_airplane.png` | 2026-05-23 |
| §7 | Café patrons combined sheets (v1, fringe issues) | `cafe_patron_<name>.png` | 2026-05 — superseded by §D above |
| §NEW Paris Clouds | Paris Clouds airplane sky | `paris_clouds.png` | 2026-05-23 |
| §I | Higgins throw-map one-shot | `npc_director_higgins_throw_map.png` | 2026-05-23 |
| §J | PP catch-map one-shot | `pp_catch_map.png` | 2026-05-23 |
| §K | Thrown-map projectile sprite | `inv_travel_map_throw.png` | 2026-05-23 |
| §L | Travel-map inventory icon | `travel_map_icon.png` | 2026-05 |
| §N | Item sprites (8 items) | `assets/images/items/*.png` | 2026-05-23 |
| §R | Café au Lait inventory item | `cafe_au_lait.png` | 2026-05-23 |
| §S | Confiture inventory item | `confiture.png` | 2026-05-23 |
| §T | Camille quick-sketch one-shot | `npc_camille_sketching.png` | 2026-05-23 |
| §V | Henri give-jam one-shot | `npc_henri_give_jam.png` | 2026-05-23 |
| Madame Colette | **DO NOT REGENERATE** — user 2026-05-23 prefers the current design | `npc_french_guide_*.png` | — |

**Removed in 2026-05-24 cleanup (low-priority / deferred):** previous PP
talk-front + talk-side regen prompts (user reversed direction — see §A),
previous Marcus talk regen (user wants idle recolor first — see §B),
previous Higgins office regen prompt (replaced by §C with the "match
entrance design" instruction), PP grab-flower regen, PP grab rolling pin,
Marcus strange-idle fringe touch-up, Windows .exe icon prompt. The git
history of this file before 2026-05-24 has the bodies if any of these
come back.
