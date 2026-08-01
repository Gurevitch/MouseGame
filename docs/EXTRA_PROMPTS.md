# Extra Sprite Prompts — everything still needed for the current FIXME sweep

This file is read by ChatGPT inside Cursor. Each paste-ready prompt is wrapped
between `===PROMPT START===` and `===PROMPT END===` markers. **Workflow:**

1. Highlight everything BETWEEN the markers (the blockquote block itself,
   not the marker lines).
2. Paste into ChatGPT. Include the header below (style lock + standing
   rules) as context if ChatGPT doesn't already have it — those rules
   apply to every prompt in the file.
3. Save the resulting PNG at the path shown in that prompt's section —
   generated ON the flat blue `#B4D7EE` background, as prompted. The runtime
   loaders key the flat blue away; that is the default, no extra step.
4. **Only if the landed sheet shows bg residue in-game or in the audit**
   (blue/white spots around the figure or trapped between arm and body):
   add it to the allowlist in `tools/bg_bake/main.go` — mode per its runtime
   loader, strict `pocketTol` on near-white bgs — run
   `go run ./tools/bg_bake .`, and ALWAYS eyeball the result before moving
   on. Do NOT bake sheets that render fine (2026-08-01 lesson: a proactive
   sweep punched pixels out of PP's ivory and was rolled back).
5. Run the atlas re-pack (or restart the game for legacy loaders):

```
python tools/pack_atlas.py tools/characters/<name>.yaml
```

6. Move the section header into the **Done log** at the bottom of this
   file and delete the prompt body so the working part stays scannable.

---

**Style lock + standing rules below — feed these to ChatGPT once per session
so it doesn't violate them on the next prompt:**

> Hand-drawn 1990s Saturday-morning cartoon, Pink Panther _Hokus Pokus Pink_
> (1997) / _Passport to Peril_ (1996). Confident black ink linework ~3 px,
> flat saturated fills, no cross-hatching, no gradients, no airbrush. Two
> cel tones max per color region. Flat chroma-key background, zero scenery.
> Every cell is **tall rectangular**, never square.

Canvas dimensions are locked per sheet; do **not** scale down to square.

**Standing PP design rules (apply to EVERY PP prompt):**

1. **No gloves of any color.** Pink Panther in this game has plain
   pink paws/hands — never yellow gloves, never any gloves.
2. **Every pickup sprite ends with PP pocketing the item.** The final
   1-2 frames show the item vanishing into his invisible hip pocket
   (the classic Pink Panther "magic pocket"); PP ends empty-handed in
   a relaxed standing pose with a small secretive smile.
2b. **Every GIVE sprite STARTS with PP pulling the item OUT of that same
   hip pocket** (user 2026-07-24) — frame 1 is the reach to his hip,
   frame 2 the item appearing in his paw, THEN the offer. Never start a
   give with the item already floating in his hand.
3. **No pure white anywhere on the panther.** Belly uses ivory
   off-white `#F2EFE5`, eye sclera uses pale grey. Pure white pixels
   on PP get chroma-keyed away by the engine.

**★ HARD RULE (user 2026-06-30): ZERO pure `#FFFFFF` on ANY drawn element of
ANY sprite — characters, props, painter canvases, aprons/smocks, paper, teeth,
AND eyes (eye sclera = pale grey `#C4C4C4`, never white).** The engine
chroma-keys pure white and punches it into holes / white spots. This is why
Pierre is being recreated (his old canvas + smock + eyes were pure white).
Only a legacy sheet's scene-cell BACKGROUND may stay pure white. New one-shot
sheets use the blue chroma-key below. Verify any new sheet with
`SPRITE_SCAN=1 go test ./engine -run TestSpriteScan`.

**No white halo / fringe rule:** anti-aliasing on the sprite itself must NOT use
pure white or near-white. After the background is removed there must be no tiny
white rim around the character, eyes, props, flowers, canvas, smock, hands, or
outline. Use pink/cream/grey edge pixels that match the nearby fill instead
(for example pale-grey `#C4C4C4` for eye sclera and cream `#E5DDC8` /
bone `#EDE5D3` for light props). Pure `#FFFFFF` is allowed only in a legacy
sheet's empty cell background.

**Standing rule for ALL characters who need "white" in their design:**
the engine chroma-keys pure `#FFFFFF` plus a tolerance band. Use these in
order of preference for fabric / large white areas:

- **Cream `#E5DDC8`** ← USE THIS for "white shirts" or any large fabric.
- **Bone `#EDE5D3`** — paper, small label panels.
- **Pale grey `#C4C4C4`** — steam wisps, eye sclera.
- **Vanilla `#F2EFE5`** — only safe for tiny accents (a tooth, a button).

The **scene background** in a legacy sprite cell may still use pure
`#FFFFFF`; new one-shots use `#B4D7EE`. Both are sampled chroma keys. The
character / foreground objects must avoid pure white.

## One-shot sheet canvas rule (2026-07-17)

All new or re-rolled **one-shot** sheets use a perfectly flat light-blue
background `#B4D7EE` and **6 frames per row**. A single-row sheet is exactly
1536×1024 (six conceptual 256×1024 cells); multi-row one-shots keep six
columns (for example Higgins give-map is 6×2 at 1536×1024). Every background
pixel and every corner must be the same `#B4D7EE`: no vignette, gradient,
texture, or separator. The engine samples this color and globally keys it,
so cream paper/props survive while enclosed background gaps disappear.

Keep ≥15px full-height background gaps between figures, lock the anchor and
scale, and draw one figure per frame. The loader cuts at those gaps; the
conceptual 256px boundaries are not visible.

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
> meet directly with the sheet's flat chroma-key background on both sides. The
> exported PNG must look like ONE continuous canvas where each Nth × Mth
> rectangle happens to hold one frame; if you cropped any cell out you'd
> see only that frame on the same flat key color, with no edge artefacts.

If you see a faint vertical/horizontal line in the preview, the generator
drew a separator — regenerate with the rule above emphasised.

## Reference-anchoring rule (applies to EVERY PP sheet)

User 2026-06-12: a text-only prompt drifts off our established PP design.
Every PP prompt must be sent WITH reference images attached, and must open
with design-lock language:

> Use the attached images as the character reference: this is the SAME
> character — copy its exact design (head shape, eye/muzzle style, outline
> weight, colors, proportions). Do not restyle or modernize him.

What to attach:

1. **Always**: the canonical sheet for the same view —
   `PP idle front.png` for front-view sheets, `PP idle side.png` /
   `PP walk left.png` for side-view sheets.
2. **Optionally**: one or two stills from `assets/images/retro_frames/`
   (the original PTP captures) for pose attitude and era vibe — but the
   character design always comes from OUR sheet, not the still.

## No extras rule (applies to EVERY sheet)

User 2026-06-02: generators keep adding a large "hero" character **portrait**
beside the frame grid (and sometimes title text / labels). Include in every
prompt:

> Output ONLY the N×M grid of animation frames — nothing else. NO separate
> large character portrait or "hero" reference image beside or above the grid,
> NO title text, NO labels, NO watermark, NO colour swatches. Just the frames
> on the required flat chroma-key background.

## One character per cell / no ghosts rule (applies to EVERY sheet)

User 2026-06-10 (floating severed hand seen in-game on PP talk front): the
engine slices the sheet into exact N×M cells, so include in every prompt:

> Each cell contains EXACTLY ONE complete character drawing, fully inside its
> own cell, with a CLEAR EMPTY GAP of at least 15 pixels of pure background
> between neighbouring figures (and between the figures and the sheet edges).
> NEVER let any part of a drawing touch or cross a cell boundary, NEVER place
> a figure straddling two cells, and NEVER paint detached "ghost" duplicates
> of limbs (a second hand, a motion-trail arm from the previous pose)
> anywhere in a cell — one body, all parts connected to it.

**Why the gap matters (engine, 2026-06-10):** the loader now detects frames
by the EMPTY GAPS between figures and cuts at the gap midpoints — so figures
are never sliced, even when they're not perfectly centered in their cells.
When figures touch each other (no gap), the loader falls back to fixed grid
lines and anything crossing them gets cut. Clear gaps = uncuttable frames.

`tools/jitter_audit` flags violations (GHOST PIECES / CONTENT CROSSES);
`tools/sheet_clean` erases ghosts on PROP-FREE sheets as a stopgap (never run
it on sheets with legit separate objects: thrown map, handed items, pigeon).

---

## Open Prompts

## Done log

- **2026-08-01 — correction rerolls landed:** §HIGGINS-THROW-MAP-v3 now
  matches his green-shirt, glasses, and red-lanyard office design, keeps a
  warm expression, and shows the map through the camera-facing release.
  §PP-PUT-COFFEE-TABLE-v2 keeps PP upright and releases the cup at waist-high
  tabletop level. Both exact 1536×1024 sheets GAP-DETECT at their runtime
  grids and contain zero pure-white pixels.

- **2026-08-01 — final open prompt batch:** §PP-PUT-COFFEE-TABLE,
  §PP-GIVE-PEN, §PP-GIVE-COIN, §SHIMON-RECEIVE-PEN, and
  §HIGGINS-THROW-MAP-v2 landed at 1536×1024. All five sheets GAP-DETECT at
  their runtime grids and add no pure-white scan hits. The dedicated pen and
  coin gives replaced their temporary aliases, and all four new paths were
  added to both sprite audit manifests. **Same-day: §PP-PUT-COFFEE-TABLE v1
  REJECTED** (PP set the cup at ankle height — reads as the ground, not
  Henri's waist-high table); parked as `_rejected_v1`, v2 re-queued above
  with the waist-height rule. **§HIGGINS-THROW-MAP v2 also REJECTED** — it
  was generated from the pre-correction prompt (khaki + red tie, no
  glasses/lanyard, angry acting; the corrected prompt text was lost in the
  prompt-file cleanup before it was used); parked as `_rejected_v2`, the
  previous sheet restored, §HIGGINS-THROW-MAP-v3 re-queued above with the
  office-idle design lock + warm mood.

- **2026-07-31 — final prompt sweep (5 prompt groups / 8 sheets):**
  §PP-TEA-SPIN-PAIR, §TAKESHI-TALK-CHERRY-v2,
  §TAKESHI-TALK-GATES-v2, §JP-HIRO-COUNTER-IDLE-v2,
  §PP-SLEEP-WAKE-BLUE and §MARGAUX-RECEIVE-HEEL landed. All sheets are
  normalized to flat `#B4D7EE`, GAP-DETECTED at their runtime grids, and
  clear the pure-white scan. The sleep/wake and Margaux receive sheets were
  added to both audit manifests.

- **2026-07-31 — Japan/tea art batch landed (7 prompts cleared):**
  §JP-TAKESHI-DOORS (idle + open/close doors + all three tale cards),
  §PP-TEA-SEAT-SET (spin-to-sit, sit idle/talk, the new sit-drink),
  §TEA-MASTER-BLUE-v2, §JP-HIRO-COUNTER (idle/talk/work/give),
  §PP-SAKURA-SELFIE-v2, §PP-PICK-SAKURA and §WALL-PRAYER-BG are all on
  disk and wired (Takeshi's doors fire from onDialogStart/onDialogEnd with
  the per-tale talk swap; counter-Hiro's four sheets load in
  newRamenSeller). Three follow-ups from the playtest are re-queued above:
  the spin pair must actually change his clothes, Takeshi's cherry + gates
  talk sheets need the doors framing, and Hiro's counter idle must stop
  offering the bowl.


- **2026-07-28 — Japan story-v2 art batch:** completed
  §BAGEL-SELLER-REDESIGN, §JP-STREET-NOREN-CLEAR, §PP-SIT-TO-STAND,
  §PP-SAKURA-SELFIE, both §JP-LINE prompts, §KENJI-BRUSH-MENU,
  §JP-TAKESHI, §JP-ITEMS-V2, §PP-RAMEN-PAIR, and
  §OBACHAN-RECEIVE-RAMEN. All 15 animation sheets are exact
  1536×1024 strips on normalized `#B4D7EE`; the three item icons are
  exact 512×512. Every sheet GAP-DETECTS. Legitimate detached-object
  audit hits are the carts, tongs, food, menus, camera, stage/cards,
  briefcase, bundle, and bowls. The Tokyo background now has an open
  counter window; the ceremony stand-up and sakura selfie are wired.
  The bagel seller now uses separate idle/talk sheets with the old
  combined sheet retained only as a fallback.

- **2026-07-28 — §PP-TALK-SIDE-STUNG-v5:** generated the stronger
  hurt-surprise take, normalized its background to exact flat `#B4D7EE`,
  and replaced `PP talk side scared.png`. The 1536×1024 8×2 sheet
  GAP-DETECTS, has no pure-white or frame-leak hits, and passed visual
  inspection. The audit's detached-piece hit is PP's intentionally
  disconnected eyebrow/whisker ink; mathematical border crossings are
  harmless because the loader uses the detected full-height gaps.

- **2026-07-26 — §JAKE-GET-COIN-BLUE:** replaced the broken checkerboard
  sheet with an exact 1536×1024 flat-`#B4D7EE` 8×1 strip. It GAP-DETECTS,
  has no ghost, pure-white, or frame-leak hits, and the loader no longer
  needs the interim checkerboard tolerance.

- **2026-07-25 — five-sheet generated batch:** §BAGEL-RECEIVE-COFFEE-v2,
  all three §PP-KIMONO-SPLIT phases, and §PP-TALK-SIDE-STUNG (v3 — close
  but poses samey / brows neutral; v4 storyboard take queued above, v3
  stays wired as the interim) landed at
  1536×1024 on exact flat `#B4D7EE`. All five GAP-DETECT, have no
  pure-white or frame-leak hits, and were visually checked after the
  background normalization. The three new kimono sheets were added to
  both audit manifests; the bagel and talk entries were already present.

- **2026-07-24 — 17-sheet generated batch:** §PP-TALK-SIDE-SCARED (tone
  REJECTED twice: v1 too cowering, the BAFFLED v2 read as casual
  puzzlement — §PP-TALK-SIDE-STUNG v3 queued above with the dual
  confused + stung-by-the-tone acting; the v2 sheet stays wired as the
  interim),
  §PP-GIVE-BAGUETTE-BLUE, §PP-GIVE-CARD-v2, §PP-GET-PENCIL-POT,
  §PP-GET-BAGEL, §BAGEL-RECEIVE-COFFEE (design REJECTED same day — apron +
  missing cart, didn't match his idle; sheet deleted, v2 queued above),
  §BEAUMONT-IDLE-CALM,
  §BEAUMONT-GET-SKETCH, §KENJI-IDLE-SIT, all three §LOUVRE landscape
  props, and all five §JP landscape props landed. Beaumont's idle and
  receive-sketch keep his podium framing, and the received page shows
  Camille's canonical Mona Lisa replica; Kenji now remains seated. All are
  1536×1024 on exact `#B4D7EE` with zero pure white, and are wired into
  their runtime loaders and both audit manifests.

- **2026-07-23 — five-sheet PR prompt batch:** §POULAIN-GIVE-BAGUETTE-v2,
  §POULAIN-GIVE-COFFEE-v2, and §POULAIN-WORK-v2 landed as waist-cut 6×1
  strips with no baked counter/table. §PP-GIVE-COFFEE-v2 landed with a cream
  cup and moved to the 6×1 global blue-key loader. §PEOPLE-PRAY-1-v2 landed
  figures-only and restored ABAB entrance-crowd variety. All five are exact
  1536×1024 / `#B4D7EE`, GAP-DETECTED, with no pure-white or leak hits.

- **2026-07-23 — pre-PR sprite-check pass:** §PP-GET-WELL-WATER landed
  (lean → rope → lift the shining cup → steady → pocket → idle), gap-detects
  6×1 with no white/leak hits, and is wired via `receive_well_water`. All 13
  new working-tree sheets verified; the missing `give_cardamom` player-loader
  entry was wired (`PP give cardamom.png`).

- **2026-07-18 — PR-batch landing:** 21 of 25 prompts landed and were wired
  + verified (talk-front, pull-map(8), give pencil, the Pierre six, Camille
  sketch + receive-pencil, spice talk v3, Avi receive-bagel, people_pray2,
  Kenji talk + give-charm (8×1), open-door office day2/day3, the give/get
  card pair (get=8×1), cardamom/jerusalem-coffee/fire-striker/matcha-bowl
  gives, offering-bowl + voice-charm receives, note-in-paw wall beat). Zero
  gap-detect fallbacks across the manifest.

- **2026-07-21 — §AVI-RECEIVE-BAGEL design correction:** Rebuilt the receive
  handoff as the intended 6×1 strip, using Avi's idle rear view and his
  talk/give face and front costume as the design authorities. The six beats now
  turn/reach left, receive with both hands, bless, bite, chew, and lower/nod on
  exact flat `#B4D7EE`, with a locked baseline and zero pure white.
- **2026-07-20 — EXTRA_PROMPTS generation pass:** Landed and verified
  §PP-TALK-FRONT-BLUE (facial details restored), §PP-PULL-MAP-BLUE,
  §PP-GIVE-PENCIL-BLUE, §PP-GET-CARD-BLUE, §PP-GIVE-CARD,
  §PP-GIVE-CARDAMOM, §PP-GIVE-JER-COFFEE, §PP-NOTE-WALL-v3, §PIERRE-BLUE,
  §CAMILLE-SKETCH-BLUE, §CAMILLE-RECEIVE-PENCIL, §SPICE-TALK-v3,
  §AVI-RECEIVE-BAGEL, §PEOPLE-PRAY-2, §KENJI-TALK, §KENJI-GIVE-CHARM, and
  all five §JP-ITEMS sheets. Runtime-authoritative 8-frame variants were used
  where the older 6-frame prompt summary disagreed with the loaders. Every
  landed sheet gap-detects; newly accepted sheets contain zero pure-white
  pixels after fringe repair. The biker was rebuilt from the pristine sheet:
  all 124,829 non-background pixels match the reference exactly, with only
  the checkerboard background replaced by `#B4D7EE`. Camille sketching remains
  1536×1024 / 8×1; median content height 324px versus idle 335px (3.3%).
- **2026-07-18 — antique-girl follow-ups:** §ANTIQUE-GIRL-ALT-v2 and
  §ANTIQUE-GIRL-TALK-v2 landed as separate 6×1 strips matching her redesigned
  idle. The alternate dusts a held brass pot with no table; the talk strip
  keeps the same costume and scale. Both use exact `#B4D7EE` with zero pure
  white and are wired into the 6-frame loaders and audit manifests.
- **2026-07-18 — flat-blue conversion batch:** §PP-WALK-FRONT-BLUE,
  §PP-WALK-BACK-BLUE, §PP-GRAB-FLOWER-BLUE, §PP-GET-COFFEE-v2, and
  §FENCE-BLUE landed at 1536×1024 with exact `#B4D7EE` keys and zero pure
  white. The walk and flower frames were preserved; coffee now uses the
  small patterned finjan and ends pocketed. The fence now uses a global
  sampled key so blue enclosed between its bars is removed. All four
  registered player sheets GAP-DETECT and passed the sprite scan.
- **2026-07-17 — round-3 blue-background one-shots:** Nine existing give /
  receive sheets were re-rolled at 6×1; §PP-GET-CARDAMOM,
  §PP-GET-COFFEE, and §PP-GET-PAPER landed and were wired; and
  §SPICE-GIVE-v2, §ANTIQUE-GIRL-v2, and §HIGGINS-GIVE-MAP-v2 landed with
  their requested design fixes. All 15 sheets use exact flat `#B4D7EE`,
  were repacked with ≥15px full-height gaps, GAP-DETECT, and add no
  pure-white or frame-leak scan candidates. The 6×1/6×2 loader and checker
  manifests were updated.
- **2026-07-17 — round-2 art corrections:** §PP-GRAB-BASKET-v2 landed with
  no baked-in basket and its enclosed chroma-key whites repaired;
  §SPICE-TALK-RIGHT now faces and gestures toward viewer's right; and
  §ROOM-DOORS-LIGHT re-graded Jake/Tommy/Danny's open doors to the dusky
  day-2 camp palette. Both sprite sheets gap-detect; visual checks passed.
- **2026-06-30 — interim Pierre repair:** Current `npc_pierre_idle.png`,
  `npc_pierre_talk.png`, `npc_pierre_give_pass.png`, `npc_pierre_get_baguette.png`,
  and `npc_pierre_get_jam.png` were repaired/verified against the no-pure-white
  scan. This is not the final prompt-generated rebuild; §PIERRE-REBUILD remains
  open for the full five-sheet regen.

- **2026-07-23 — 2026-06-30 tail cleanup:** the old PR bug-sweep regen
  section was removed as completed/superseded: §PP-PUTNOTE (→ the landed
  §PP-NOTE-WALL-v3), §PP-GIVE-BG (sketch/postcard re-rolled blue 2026-07-17;
  give-baguette's pockets now cleared by the tol-24 global loader),
  §PP-GET-SKETCH-SIZE (→ the landed `PP receive sketch.png`), §FENCE-NOWHITE
  (→ the landed §FENCE-BLUE), §PP-RECEIVE-PRESSPASS (→ the landed
  `pp_get_card.png`). §WALL-PRAYER-BG moved up to the carried-over list.
