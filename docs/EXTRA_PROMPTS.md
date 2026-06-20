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
> on pure #FFFFFF.

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

All prompts below still need a PNG generated. When one lands, move its row
into the **Done log** at the bottom and delete the body.

### 2026-06-15 playtest — background / chroma-key re-rolls

Two sheets read with a leftover "background" in-game. Both are the
**white-on-white chroma-key problem**: the engine keys pure `#FFFFFF`, so any
prop the character holds that is ALSO pure white gets eaten where it touches
the background (or leaves enclosed white pockets / halos the edge-flood can't
reach). The biker (§BK2) and the pot pigeon already ship transparent and load
raw — those two are fine; these two are not.

#### §FLOWER-PICK — `PP grab flower.png` blinks + leftover white (#1) · STILL NEEDED

**Path:** `assets/images/player/PP grab flower.png` · **6×1** single row · pure
`#FFFFFF` background (chroma key). `tools/jitter_audit` on the current sheet:
**CONTENT CROSSES 3 cell borders (14-16px)** + the frame-1 daisy sits far from PP
as a detached 82px sliver, so gap-detection slices unevenly → frames show pieces
of their neighbours = the **blinking** in-game. Plus an enclosed white pocket
between his arm and body (the edge key can't reach it) reads as background.
Re-roll with the layout rules below; the loader keys it at tol 36 already.

Key fixes for the re-roll:
- **Even, self-contained frames:** 6 cells of equal width, each holding ONE
  cohesive PP-with-daisy silhouette. Keep the daisy WITHIN PP's reach in every
  frame (touching or within a few px of his paw) - never a separate object
  floating at the far edge of a cell. **≥15px empty gap between frames** and to
  the sheet edges; nothing crosses a cell boundary.
- **Anchor lock:** feet on the same pixel row and centerline on the same column
  in every frame; only the bend/reach changes. (The crouch may shorten him -
  that's fine, the engine anchors by feet - but the standing width/position must
  not drift.)
- **No enclosed white:** when his arm bends up holding the flower, leave the gap
  between forearm and chest OPEN to the background edge (don't seal it into a
  trapped white pocket), or close the arm flush against the body so there's no
  gap at all.
- Daisy petals in **bone `#EDE5D3`** (center golden yellow), never pure white.

===PROMPT START===
> 6-frame single-row pickup animation of the Pink Panther (use the attached
> `PP idle side.png` as the exact character reference - same head, muzzle,
> outline weight, pink fills, NO gloves, off-white `#F2EFE5` belly, pale-grey
> sclera). Side profile facing LEFT toward a small daisy. The 6 poses: 1)
> standing beside the daisy, 2) leaning/crouching toward it, 3) reaching down to
> it, 4) plucking it, 5) rising holding it at chest height, 6) tucking it into
> his invisible hip pocket, ending empty-handed with a small secretive smile.
> CRITICAL LAYOUT: 6 EQUAL-WIDTH cells in one row, each cell one complete
> PP-with-flower drawing fully inside it with a clear empty gap (>=15px) to both
> neighbours and to the sheet edges - no limb, tail, or daisy may cross a cell
> boundary, and the daisy must stay within PP's reach in every frame (never a
> separate object drifting to the cell edge). Feet on the same row, body on the
> same centerline, same standing size in all 6 frames. The arm/forearm gap when
> he holds the flower must stay OPEN to the background (no sealed white pocket).
> Daisy petals **bone `#EDE5D3`**, center golden yellow - never pure white.
> Pure `#FFFFFF` background only, no ghost/duplicate limbs, no separators or
> gridlines, no portrait or labels. Tall rectangular cells, never square.
===PROMPT END===

#### §PIERRE-BOARD — Pierre's easel canvas vanishes in a frame (#6) · STILL NEEDED

**Paths:** `assets/images/locations/paris/npc/outside/npc_pierre_idle.png` +
`npc_pierre_talk.png` · **8×1 each** · pure `#FFFFFF` background. Root cause of
"a single frame where he doesn't have his board": the easel CANVAS is pure
white, so where it abuts the white background the edge-connected key floods
through and erases it (he's left with just the wooden easel / no board). Re-roll
both sheets with the **canvas/board in cream `#E5DDC8`** (a primed-canvas tone),
clearly outlined so it never merges with the background. Keep the easel + canvas
present and identical in EVERY frame.

===PROMPT START===
> Two 8-frame single-row sprite sheets (IDLE and TALK) of Pierre, a Parisian
> street painter (black beret, khaki smock, blue-striped scarf, palette in
> hand), standing at his wooden easel in left profile - use the existing
> `npc_pierre_idle.png` as the exact design/size reference. CRITICAL: the easel
> CANVAS must be **cream `#E5DDC8`** with a clean dark outline, never pure white,
> and the easel + canvas must be fully present and unchanged in all 8 cells.
> IDLE row: small breathing/brush-dabbing gestures. TALK row: mouth cycling
> natural speech shapes, gesturing with the brush. Pure `#FFFFFF` background
> only, one complete figure (with his easel) per cell, clear gaps to every edge,
> no separators or gridlines, no portrait or labels. Tall rectangular cells.
===PROMPT END===

### Backgrounds + ambient life (camp return, darkening, Jerusalem, bg life)

All backgrounds are **1376×768**, drawn in the game's hand-drawn 90s Pink Panther
cartoon style. Backgrounds contain **no characters** (PP/NPCs are drawn on top);
keep the lower third / foreground relatively clear so characters have a floor to
walk on. Ambient "objects" are separate transparent-background overlay sprites.

#### §AMB3 — Ambient: camp crow (lands on the airstrip sign) — #34 · STILL NEEDED

**Path:** `assets/images/ambient/crow.png` · **TRANSPARENT background** · 8-frame single-row strip · **frames 0-5 = wing-flap loop, frames 6-7 = perched/standing pose** (the code flies frames 0-5 in and out, then holds the last frame while perched).

This is the bird the camp-landing scene already expects. Until it lands the crow
silently no-ops; drop the PNG in and it flaps in, sits on the CAMP sign, then
flies off, on a loop.

===PROMPT START===
> A small overlay sprite (TRANSPARENT background, NOT white) of a single black
> crow, side profile facing RIGHT, hand-drawn 90s cartoon style. Lay it out as an
> 8-frame single-row strip: frames 1-6 are a wing-flap FLYING loop (wings up
> through wings down, body level, legs tucked), and frames 7-8 are the bird
> PERCHED and standing still (wings folded, feet down, as if gripping a sign).
> Small and simple (background depth), glossy black with a hint of blue sheen, a
> small beak and eye. One pose per cell, even spacing, no separators or gridlines.
===PROMPT END===

#### §JN1 — Jerusalem: Eli the spice seller (souk vendor) — market scene

**Path:** `assets/images/locations/jerusalem/npc/npc_eli_idle.png` · **8×2 grid** (row 0 = idle, row 1 = talk), white background for color-key. Currently borrows the Paris art-vendor sheet as placeholder — this replaces it with a proper souk spice merchant.

===PROMPT START===
> A friendly Middle-Eastern spice merchant for a 90s point-and-click cartoon, full
> body, facing the viewer/slightly right, standing behind a market stall. Warm
> earth-tone tunic/apron, rolled sleeves, a small cap, short beard, weathered
> cheerful face. Lay out an 8-column × 2-row sprite sheet on a PLAIN WHITE
> background: ROW 1 = 8 idle poses (small gestures, scooping spice, wiping hands,
> a welcoming wave), ROW 2 = 8 talking poses (mouth open, hands presenting/
> gesturing as if describing his spices). Consistent character, size and baseline
> across all 16 cells, even spacing, no gridlines or separators. Hand-drawn 90s
> cartoon style, bold clean outlines.
===PROMPT END===

### Parked — Jerusalem task-chain art (don't generate yet; see `docs/JERUSALEM_TASKS.md`)

The chain is designed but NOT wired, and the exact items may still change.
One-line specs only — full prompt bodies get written when we build the chain
(older drafts recoverable from git history):

- **§JW1** `assets/images/player/PP write note.png` — ~6×1, PP writes a note,
  ends pocketing it.
- **§JW2** `assets/images/player/PP put note in wall.png` — ~6×1, PP tucks the
  note into a Western Wall crack.
- **§JC1** `assets/images/locations/jerusalem/npc/alley_cat.png` — ~6×1
  transparent, sit/idle loop + trot-away.
- **§JI1** item icons: guidebook, sardine_tin, glow_bug_jar, charcoal_stick,
  paper, pencil, note.
- **§JG1** four NPC give one-shots (Gary, Eli, Dov, Miriam) — 8×1 each,
  SKILL.md §8b rule.

#### §NIC1-v2b — Nicolas TALK sheet (idle landed + verified; talk still pending)

The split IDLE sheet landed 2026-06-12 and verified (gap-detected 1×8, full
camera routine, mouth closed). Until the talk sheet lands the loader falls
back to the OLD combined sheet's talk row. **Path:**
`assets/images/locations/paris/npc/outside/npc_press_photographer_talk.png`
· **8×1 at 1536×1024** (cells 192×1024).

===PROMPT START===
> TALK sheet: 8 frames in ONE row of Nicolas, a Parisian street photographer
> in his 30s (olive-green field vest, cream `#E5DDC8` shirt - never pure
> white, dark slacks, vintage camera hanging on its neck strap - identical
> design and size to his idle sheet, see reference). Pure #FFFFFF at exactly
> 1536×1024 (cells 192×1024), one complete figure per cell with clear gaps,
> nothing touching cell edges, no ghost limbs. The TALK loop, every cell
> clearly different: camera resting on his chest, MOUTH CYCLING natural
> speech shapes (closed - slightly open - wide "ah" - mid - narrow "oo" -
> closed), his free hand gesturing enthusiastically, an eyebrow raise and a
> small head tilt mid-loop. ANCHOR LOCK - feet on the SAME pixel row,
> centerline on the SAME column, same size in all 8 cells. No separators,
> no extra portrait, no text.
===PROMPT END===

#### §JG1 — Jerusalem NPC GIVE one-shots (rule SKILL.md §8b — needed BEFORE the chain is wired)

The Jerusalem daisy-chain (docs/JERUSALEM_TASKS.md) has FOUR person-to-PP
hand-overs, and per §8b each giving NPC needs a give one-shot. Queue these
alongside the §JI1 item icons when the chain gets built. All four: **1536×1024,
8×1** strips, pure #FFFFFF, matching each NPC's existing design, ANCHOR LOCK
(feet/waist on the same pixel row, centerline on the same column, only arms
and head move, nothing touches cell edges, no separators/extras/text):

- `assets/images/locations/jerusalem/npc/npc_gary_give.png` — Gary the
  tourist digs the Pencil + Sardine Tin out of his daypack and holds them out.
- `assets/images/locations/jerusalem/npc/npc_eli_give.png` — Eli the spice
  seller tears a paper slip off his wrapping roll and offers it.
- `assets/images/locations/jerusalem/npc/npc_dov_give.png` — Dov hands over
  the Charcoal Stick from his tool belt.
- `assets/images/locations/jerusalem/npc/npc_miriam_give.png` — Miriam the
  archeologist carefully presents the finished Coin Rubbing with both hands.

---

## §JIT — Jitter regen batch (2026-06-10 automated audit)

**YES — the whole point of every prompt in this batch is putting the
character in the SAME position in every cell.** `go run ./tools/jitter_audit`
measured these sheets drifting: the feet line and/or horizontal center moves
from cell to cell, which renders in-game as the sprite jumping/sliding while
it animates. Each prompt below embeds the anchor-lock language — paste a
block as-is, regenerate, drop the PNG in, and re-run the audit tool: a fixed
sheet comes back with no FOOT/CENTER-X warnings.

#### §JIT-PP2 — PP walk back (foot drift 60px, center 44px)

**Path:** `assets/images/player/PP walk back.png` · **1536×1024, 8×2**

**Regen #1 (2026-06-10) FAILED — foot drift got WORSE (60 → 97px). The
stride's rise/sink overshot. Re-roll; tell the generator the body may rise
and sink only ~10px between frames, and the planted foot must stay put.**

===PROMPT START===
> A 16-frame BACK-VIEW walk cycle of the slim Pink Panther (seen from behind,
> walking away from camera; plain pink paws, NO gloves), 8 columns × 2 rows on
> pure #FFFFFF at exactly 1536×1024 (cells 192×512). This is an ANIMATION:
> every frame is a clearly different moment of his famous strut — smooth,
> confident, a little smug. He walks IN PLACE; frames 1→16 are TWO full
> strides (8 steps each), looping cleanly back to frame 1:
>   Step pattern per stride: foot plants (body sinks slightly, shoulders
>   counter-rotate) → passing pose (legs together, body at its tallest) →
>   other foot reaches and plants → repeat mirrored. Arms swing OPPOSITE the
>   legs with loose cartoon overlap; the long tail snakes left-right behind
>   him a half-beat behind the body; his head bobs subtly with each step.
> ANCHOR LOCK — he struts IN PLACE: in EVERY cell the planted foot contacts
> the SAME pixel row and his body's vertical centerline stays on the SAME
> pixel column (the body rises/sinks with the stride, but never slides
> sideways or drifts up the cell). Same size every frame; nothing touches a
> cell edge. No separators, no extra portrait, no text.
===PROMPT END===

#### §PIGEON-LADY — Madame Margaux, the pigeon lady (2026-06-12, new quest NPC)

She stands on the LEFT side of paris_street (opposite Pierre), feeding the
pigeons, and lures the flower-pot guard pigeon off when PP brings her the
day-old Baguette Heel. Sized like Pierre (mid-distance). **Wired but
invisible until her art lands** (clickable via bounds meanwhile), so this is
the priority sheet.

**Paths:** `assets/images/locations/paris/npc/outside/npc_pigeon_lady_idle.png`
(required) and `npc_pigeon_lady_give.png` (optional scatter one-shot) ·
**8×1 each, 1536×1024** (cells 192×1024).

**ATTACH** `npc_pierre_idle.png` for SIZE/scale match (she should read the
same mid-distance size as Pierre).

===PROMPT START — IDLE===
> An original 1990s point-and-click cartoon character: a kindly, plump
> elderly Parisian "pigeon lady" - plum coat, knitted shawl, a little hat
> with a flower, a paper bag of crumbs in one hand. Match the SIZE and full-
> body framing of the attached Pierre reference (same mid-distance scale).
> A SINGLE ROW of 8 IDLE frames on pure #FFFFFF at exactly 1536×1024 (cells
> 192×1024): she sprinkles a few crumbs by her feet, looks down fondly, a
> gentle sway, a soft "coo-coo" mouth - one or two small grey pigeons peck
> near her hem (they bob on different frames). Every frame clearly different.
> Bold dark outlines, flat saturated colours. NEVER pure white on her (cream
> #E5DDC8 for the shawl, bone #EDE5D3 for the paper bag). ONE figure per
> cell, ≥15px clear white between figures and to the edges, ANCHOR LOCK (feet
> same pixel row, centre same column). No separators, no extra portrait, no
> text.
===PROMPT END===

===PROMPT START — GIVE/SCATTER (optional)===
> The SAME pigeon lady (identical design + size as her idle sheet). A SINGLE
> ROW of 8 frames on pure #FFFFFF at exactly 1536×1024 (cells 192×1024) of
> her CALLING and scattering a big handful of crumbs to the side: 1-2 reaches
> into the bag, 3-5 a big underhand toss (crumbs visible mid-air), 6-8 claps
> the crumbs off her hands with a satisfied smile as pigeons flock in. Same
> palette/anchor rules as her idle. No separators, no extra portrait, no text.
===PROMPT END===

#### §JIT-MARCUS — Marcus NORMAL talk + strange-alt (strange idle/talk have their own standalone prompts above)

**Paths:** `assets/images/locations/camp/npc/kids/marcus/npc_marcus_talk.png`
(his NORMAL, healed talk) and `npc_marcus_strange_alt.png` (the strange
fidget sheet, currently doubling as the strange idle in code) · **1536×1024,
8×2 each** · `npc_marcus_idle.png` came back CLEAN — match its framing,
centering and size EXACTLY. Paste once per sheet, swapping the [STATE] line.

===PROMPT START===
> A 16-frame animation sheet of Marcus, the know-it-all camp kid (KEEP his
> current canonical design — see reference), EXACTLY 8 columns × 2 rows on
> pure #FFFFFF at exactly 1536×1024 (cells 192×512).
> CRITICAL GRID RULE: each of the 16 cells contains EXACTLY ONE complete
> Marcus, centered in its own cell with clear white padding on every side.
> Never two figures sharing a cell, never a figure cut by a cell boundary,
> never detached ghost limbs or duplicate arms floating in a cell — one boy,
> all parts connected, sixteen times.
> [STATE — pick one:]
> - TALK (healed/normal): cheerful sketching pose, mouth cycling natural
>   speech shapes, light friendly gestures.
> - STRANGE ALT: the kid is a little OFF (eerie-sad, KID-FRIENDLY, NOT horror)
>   - absorbed in the notepad, drawing the same thing over and over, pausing
>   to gaze off, then back to the page; a small uneasy sway. Not smiling, but
>   not scared or distressed - just absent. NEVER bloodshot/sweaty/manic.
> ANCHOR LOCK — the boy is nailed to ONE spot: in EVERY cell his feet sit on
> the SAME pixel row and his body's vertical centerline on the SAME pixel
> column; only arms, face and the notepad move. No sliding, no size changes.
> No separators, no extra portrait, no text.
===PROMPT END===

#### §MARCUS-STRANGE-IDLE — Marcus room strange idle · DAY + NIGHT, LANDED 2026-06-12

The user split this into two lighting variants - `npc_marcus_strange_idle_day.png`
and `npc_marcus_strange_idle_night.png` (8×2 each) - and both landed clean
(gap-detected). The loader (newRoomMarcus) loads both; `setStrangeVariant`
picks night during the cutscene / day on Day 2, in step with the cabin bg
swap. Prompt kept below for future re-rolls (run it once per variant, adding
"warm daytime cabin light" / "dim night-time cabin light, cooler tones").

**Paths:** `npc_marcus_strange_idle_day.png`, `npc_marcus_strange_idle_night.png`
· **8×2, 1536×1024** (cells 192×512). **ATTACH** `npc_marcus_idle.png`.

===PROMPT START===
> Use the attached sheet as the reference: the SAME boy, Marcus the camp kid -
> copy his exact design, size and framing. Produce a 16-frame STRANGE-STATE
> IDLE loop, EXACTLY 8 columns × 2 rows on pure #FFFFFF at exactly 1536×1024
> (cells 192×512). KEEP IT KID-FRIENDLY - this is a gentle adventure game, so
> the mood is "this kid is a little OFF / not himself", eerie-sad, NOT horror.
> Through the loop:
>   - a faraway, distracted look - eyes a touch unfocused, not bloodshot, no
>     dark rings, no sweat, no wild stare.
>   - absorbed in a notepad, drawing the same thing over and over; he pauses,
>     gazes off, then goes back to the page.
>   - a small uneasy sway and the odd blink/shrug - quietly troubled, calm
>     hands (the game adds a faint quiver of its own, so the art stays gentle).
> He is NOT smiling/cheery, but NOT scared or distressed either - just absent
> and a bit melancholy. Every frame clearly different (it's an animation).
> CRITICAL GRID RULE: each of the 16 cells holds EXACTLY ONE complete Marcus,
> centered with clear white padding on every side and AT LEAST 15px of empty
> background between neighbouring figures and to the sheet edges. Never two
> figures touching or sharing a cell, never a figure cut by a cell boundary,
> no detached ghost limbs or duplicate arms.
> ANCHOR LOCK - feet on the SAME pixel row, body centerline on the SAME pixel
> column in every cell; only the head, arms and notepad move. No sliding, no
> size changes. No separators, no extra portrait, no text.
===PROMPT END===

#### §MARCUS-STRANGE-TALK — Marcus room strange talk (2026-06-12, standalone, matches the softened idle)

**Path:** `assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_talk.png`
· **8×2, 1536×1024** (cells 192×512). Plays while strange Marcus is SPEAKING
(before he's healed). Same eerie-sad, kid-friendly mood as §MARCUS-STRANGE-IDLE
— NOT horror. Match `npc_marcus_idle.png`'s framing, centering and size EXACTLY,
and keep him identical in design/size to his strange-idle sheet.

**ATTACH** `npc_marcus_idle.png` (design + framing lock).

===PROMPT START===
> Use the attached sheet as the reference: the SAME boy, Marcus the camp kid -
> copy his exact design, size and framing. Produce a 16-frame STRANGE-STATE
> TALK loop, EXACTLY 8 columns × 2 rows on pure #FFFFFF at exactly 1536×1024
> (cells 192×512). KEEP IT KID-FRIENDLY - the mood is "this kid is a little
> OFF / not himself", eerie-sad, NOT horror. Through the loop:
>   - he speaks in a quiet, distant, distracted way - mouth cycling natural
>     speech shapes (closed - slightly open - wide "ah" - mid - narrow "oo" -
>     closed), but his eyes stay a touch unfocused, looking past you.
>   - he keeps half his attention on the notepad he's drawing in, glancing
>     down at it between phrases as if he can't quite stop.
>   - small uneasy head tilts and the odd slow blink; gentle, calm hands -
>     NO bloodshot eyes, no dark rings, no sweat, no wild or distressed stare.
> Not cheery, not scared - just absent and a bit melancholy. Every frame
> clearly different (it's an animation).
> CRITICAL GRID RULE: each of the 16 cells holds EXACTLY ONE complete Marcus,
> centered with clear white padding on every side and AT LEAST 15px of empty
> background between neighbouring figures and to the sheet edges. Never two
> figures touching or sharing a cell, never a figure cut by a cell boundary,
> no detached ghost limbs or duplicate arms.
> ANCHOR LOCK - feet on the SAME pixel row, body centerline on the SAME pixel
> column in every cell; only the head, mouth, arms and notepad move. No
> sliding, no size changes. No separators, no extra portrait, no text.
===PROMPT END===

#### §PM1 — PP pulls the travel map from his pocket (2026-06-12: plays before the map screen opens)

**Path:** `assets/images/player/PP pull map.png` · **1536×1024, 8×1**
(cells 192×1024). Pre-wired: `openTravelMap` plays the `pull_map` one-shot
(~0.9s) and the globe opens when it ends — drop the PNG in and the beat
appears, no code change.

**ATTACH AS REFERENCES (standing rule):**
1. `assets/images/player/PP idle front.png` — the canonical design to copy
   exactly.
2. `assets/images/player/PP receive map.png` — the folded travel map prop
   (tan paper, red ribbon) so the map matches the hand-over art.

===PROMPT START===
> Use the attached images as the character reference: this is the SAME
> character as the first reference sheet — copy its exact design (head
> shape, eye/muzzle style, outline weight, the same pink, off-white belly
> #F2EFE5, plain paws with NO gloves, same proportions). Do not restyle
> him. The folded map prop must match the second reference.
> Produce an 8-frame one-shot of him pulling the folded travel map out of
> his invisible hip pocket, single row of 8 frames on pure #FFFFFF at
> exactly 1536×1024 (cells 192×1024). This is an ANIMATION: every frame
> must be CLEARLY different — one continuous motion, not held poses. Play
> 1→8 in order:
>   1 — stands relaxed, empty-handed.
>   2 — reaches across to his hip, eyes glancing down.
>   3 — paw "into" the invisible pocket at his hip (classic magic-pocket
>       gag: the paw just disappears against his side).
>   4 — pulls out the folded map with a small flourish, eyebrows up.
>   5 — holds it in front of his chest with both paws.
>   6 — flicks it open one fold, leaning his head in with interest.
>   7 — map held up and open toward the camera, filling his paws.
>   8 — settles, holding the open map steady (the map screen takes over
>       from this pose).
> The map paper is bone #EDE5D3 (never pure white).
> ANCHOR LOCK — in EVERY cell both feet contact the SAME pixel row and his
> body's vertical centerline stays on the SAME pixel column; same size
> every frame; ≥15px clear white between figures and to the sheet edges.
> No separators, no extra portrait, no text.
===PROMPT END===

#### §AMB5 — Paris street accordion player (retro plan #5: street density)

**Path:** `assets/images/locations/paris/npc/outside/ambient_accordion_player.png`
· **8×1 strip, 1536×1024** (cells 192×1024). Pre-wired in
`decorateParisStreetSprites` (left side of the street, x≈120, ground y≈470,
scale 0.85) — drop the PNG in and he appears, no code change.

===PROMPT START===
> An 8-frame in-place loop of an original cartoon Parisian street musician,
> single row of 8 frames on pure #FFFFFF at exactly 1536×1024 (cells
> 192×1024). 1990s point-and-click adventure style: bold dark outlines, flat
> saturated colors, exaggerated friendly proportions. A round, mustachioed
> man in a navy waistcoat, rolled sleeves and a flat cap, playing a small
> red accordion: the loop is the bellows stretching open and squeezing shut
> (clearly different hand spacing each frame), his shoulders rocking gently
> with the rhythm, one foot tapping, eyes closed blissfully on the squeeze
> frames. NEVER pure white on the character (cream #E5DDC8 for the shirt).
> ANCHOR LOCK — feet on the SAME pixel row and body centerline on the SAME
> pixel column in every cell; same size every frame; ≥15px clear white
> between figures and to sheet edges. No separators, no extra portrait, no
> text.
===PROMPT END===

#### §AMB6 — Paris street pigeon lady (retro plan #5: street density)

**Path:** `assets/images/locations/paris/npc/outside/ambient_crumb_lady.png`
· **8×1 strip, 1536×1024** (cells 192×1024). Pre-wired in
`decorateParisStreetSprites` (right side near the lamppost, x≈1080, ground
y≈480, scale 0.8) — drop the PNG in and she appears, no code change.

===PROMPT START===
> An 8-frame in-place loop of an original cartoon elderly Parisian lady
> feeding pigeons, single row of 8 frames on pure #FFFFFF at exactly
> 1536×1024 (cells 192×1024). 1990s point-and-click adventure style: bold
> dark outlines, flat saturated colors, kind face. She wears a plum coat,
> a knitted shawl and a tiny hat with a flower, holding a paper bag of
> crumbs: the loop is her reaching into the bag, scattering crumbs with a
> gentle underhand toss (crumbs visible mid-air on the toss frames), then
> smiling down at two small grey pigeons pecking by her hem (the pigeons
> bob on different frames). NEVER pure white on the character (cream
> #E5DDC8 for the shawl, bone #EDE5D3 for the paper bag).
> ANCHOR LOCK — feet on the SAME pixel row and body centerline on the SAME
> pixel column in every cell; same size every frame; ≥15px clear white
> between figures and to sheet edges. No separators, no extra portrait, no
> text.
===PROMPT END===

#### §JIT-WALKFRONT — PP walk front (2026-06-12: not a walk at all — 16 near-identical standing poses)

**Path:** `assets/images/player/PP walk front.png` · **1536×1024, 8×2** (cells 192×512).

The current sheet is 16 copies of PP standing facing camera with barely any
leg motion, so walking toward the camera reads as PP gliding. Needs a real
full walk cycle.

**ATTACH AS REFERENCES (user 2026-06-12: the prompt alone drifts off our
design — anchor it with images):**
1. `assets/images/player/PP idle front.png` — the CANONICAL front-view
   design to copy exactly (head, eyes, line weight, colors, proportions).
2. `assets/images/player/PP walk left.png` — how his stride reads in our set.
3. Optionally one still from `assets/images/retro_frames/` (e.g.
   `clip_t01m00s.png`) for the era's walk attitude only — design still
   comes from reference 1.

===PROMPT START===
> Use the attached images as the character reference: this is the SAME
> character as the first reference sheet — copy its exact design (head
> shape, eye/muzzle style, outline weight, the same pink, off-white belly
> #F2EFE5, plain paws with NO gloves, same body proportions). Do not
> restyle or modernize him.
> Produce a 16-frame FRONT-VIEW walk cycle of him walking TOWARD the
> camera, 8 columns × 2 rows on pure #FFFFFF at exactly 1536×1024 (cells
> 192×512). This is an ANIMATION: every frame must be CLEARLY different —
> a complete, readable walk cycle, looping cleanly 1→16→1. He walks IN
> PLACE:
>   Frames 1-8 (one full stride): left knee lifts toward camera (foot
>   visibly raising, sole hinted) → body sinks slightly as the left foot
>   plants → passing pose, legs crossing, body tallest → right knee lifts →
>   right foot plants → passing pose. Shoulders sway opposite the stepping
>   leg, arms swing loosely at his sides, head bobs subtly with each plant,
>   and the long tail curls into view alternately left and right behind him.
>   Frames 9-16: the mirrored second stride so the loop closes.
> ANCHOR LOCK — he walks IN PLACE: in EVERY cell the planted foot contacts
> the SAME pixel row and his body's vertical centerline stays on the SAME
> pixel column (the body rises/sinks ~10px max with the stride, never slides
> sideways). Same size every frame; nothing touches a cell edge; ≥15px white
> gaps between figures. No separators, no extra portrait, no text.
===PROMPT END===

#### §JIT-GIVEFLOWER — PP give flower (2026-06-12 PR#2: "not smooth" — half the frames are near-duplicates)

**Path:** `assets/images/player/PP give flower.png` · **1536×1024, 8×1** (cells 192×1024).

The current sheet reads as a 4-pose animation: frames 1-2 are the same
stand-with-flower and frames 4-6 are the same extended-arm hold, so the
hand-over pops between a few poses instead of flowing. (The engine-side
white-petal erasure was fixed separately — the daisy now survives the
color key — but the motion itself needs distinct in-between frames.)

**ATTACH AS REFERENCES (user 2026-06-12: anchor prompts with our own art):**
1. `assets/images/player/PP idle side.png` — the canonical design to copy
   exactly.
2. `assets/images/player/PP give flower.png` — the current sheet, for pose
   framing only (its motion is what we're fixing).

===PROMPT START===
> Use the attached images as the character reference: this is the SAME
> character as the first reference sheet — copy its exact design (head
> shape, eye/muzzle style, outline weight, the same pink, off-white belly
> #F2EFE5, plain paws with NO gloves, same proportions). Do not restyle him.
> Produce an 8-frame one-shot of him HANDING a small daisy to someone
> beside him, single row of 8 frames on pure #FFFFFF at exactly 1536×1024
> (cells 192×1024). This is an ANIMATION: every frame must be CLEARLY
> different from its neighbours — one continuous motion with in-betweens,
> not held poses. Play 1→8 in order:
>   1 — stands relaxed, daisy held low at his side.
>   2 — raises the daisy to chest height, looking at it fondly.
>   3 — turns slightly and begins extending his arm out to the side.
>   4 — arm fully extended, daisy offered, ears perked.
>   5 — the daisy starts leaving his paw (recipient's unseen pull), his
>       fingers opening.
>   6 — paw empty and still extended, fingers spread, a happy blink.
>   7 — pulls the arm back toward his chest with cartoon follow-through.
>   8 — back to a relaxed stand, hands free, content smile.
> The daisy: yellow center, IVORY #F2EFE5 petals (never pure white).
> ANCHOR LOCK — in EVERY cell both feet contact the SAME pixel row and his
> body's vertical centerline stays on the SAME pixel column; same size every
> frame; ≥15px clear white between figures and to the sheet edges; no
> separators, no extra portrait, no text.
===PROMPT END===

#### §JIT-PATRONS — Bernard + Camille idle (2026-06-12 sprite-check: two figures TOUCH, gap split broke)

**Paths:** `assets/images/locations/paris/npc/coffee/cafe_patron_bernard_idle.png`,
`cafe_patron_camille.png` · **1536×1024, 8×1 each** (cells 192×1024).

The gap-based slicer found a stray speck plus a MERGED double-figure run in
both sheets (Bernard cells 2-3 share one 377px run, frame 0 is a 3px sliver;
Camille mirrors it with the sliver at frame 7). In-game that's one blink
frame and one frame showing two copies of the patron. The other patron
sheets are fine. Match each character's current outfit/colors exactly.
Paste once per sheet, swapping the [CHARACTER] line.

===PROMPT START===
> An 8-frame seated IDLE loop, single row of 8 frames on pure #FFFFFF at
> exactly 1536×1024 (cells 192×1024). WAIST-UP BUST framing, same as the
> current art: head + torso + hands only, waist cutoff on the same flat
> bottom row in every frame, no chair, no table.
> [CHARACTER — pick one:]
> - Monsieur Bernard: bearded older Parisian in a flat cap and brown coat,
>   reading his folded Le Figaro newspaper, occasionally sipping coffee.
>   Loop: reads → page rustle → lifts cup and sips → lowers cup → reads.
> - Mademoiselle Camille: young art student, dark bob, red beret, green
>   blouse. Loop: holds her teacup in both hands → sips → lowers it →
>   glances dreamily aside → back to center.
> CRITICAL GRID RULE: each of the 8 cells contains EXACTLY ONE complete
> figure, centered in its own cell, with AT LEAST 15px of clear white
> between neighbouring figures and to every sheet edge. Never two figures
> touching or sharing a cell, no detached ghost limbs, no stray specks.
> ANCHOR LOCK — the waist cutoff sits on the SAME pixel row and the body's
> vertical centerline on the SAME pixel column in every cell; same size
> every frame. No pure white on the character (cream #E5DDC8 for fabric).
> No separators, no extra portrait, no text.
===PROMPT END===

---

## Done / Retired log

Headers moved here once the PNG landed or the prompt was superseded; bodies
deleted (recover from git history if ever needed).

**2026-06-10 generation batch (user) — landed, wired, audit-verified clean:**

- §PI1 — 4 Paris quest item icons (charcoal_pencil, camille_sketch,
  baguette_heel, mini_portrait) · DONE, wired in items.json.
- §PA1 — pigeon-lands one-shot · DONE, wired as Pierre's "pigeon" (sequenced
  before his "give" — playOneShotAnim replaces the active anim, so they must
  not fire back-to-back).
- §PA2 — flower-pot props (pigeon / pencil states) · DONE, wired as the
  pencil floorItem texture + swap on Pierre's favor.
- §AB — PP walk side (8×1 regular cycle) · DONE 2026-06-10, loader auto-detects 8 cols.
- §PR1 — PP generic receive · DONE 2026-06-10, wired as `receive_item`
  (5 call-sites in game.go).
- §PR2 — Pierre give · DONE, wired (npc.go).
- §PR3 — Beaumont give · DONE, wired at both postcard beats.
- §JIT PP talk front + PP talk side + PP grab + PP receive map · DONE, clean
  (talk front needed a tools/sheet_clean pass — ghost hands erased 2026-06-10).
- §JIT Marcus idle · DONE, clean. Marcus talk + strange sheets' regen #1 had
  figures STRADDLING cell borders + ghost limbs (masked the morning audit);
  reverted to the previous art — re-roll live at §JIT-MARCUS.
- §JIT Higgins shout + give-map (seated) · DONE, both clean (12-frame 6×2
  give-map, 16-frame shout, all cells filled).
- §JIT Poulain work · DONE, clean (idle/talk re-roll still live).

**2026-06-11 generation batch (user) — landed:**

- §PG1 v1 — PP generic give · landed as `PP give.png`, now the FALLBACK
  behind the per-item §PG1-v2 sheets ("give what?" - each trade should show
  the real item).
- §PA2-v2 — flower pot props · DONE 2026-06-11, regenerated + real alpha.
- §BK1 — interactive biker · DONE 2026-06-11 (ride 1-6 + braked 7-8; pause
  shows the braked pose, ride loop cycles only 1-6).
- §NIC1-v2 idle — Nicolas split IDLE sheet · DONE 2026-06-12, verified
  (gap-detected, camera routine, mouth closed). TALK pending → §NIC1-v2b.
- §PG1-v2 — PP per-item give sheets · DONE, nine `PP give *.png` sheets
  landed and are loaded automatically by `player.playGive`.
- §NIC1 v1 — combined 8×2 photographer sheet · landed but idle/talk rows
  came out ambiguous → superseded by §NIC1-v2 (two separate 8×1 sheets;
  loader prefers the split files).
- §PA2-v2 — flower-pot pigeon/pencil props · DONE, replaced both visible pot states.
- §BK1 — interactive crossing biker · DONE, replaced the 8×1 cyclist strip.

**2026-06-12 PR batch — landed, fitted + wired, bodies removed:**

- §PIERRE-IDLE — Pierre split into `npc_pierre_idle.png` + `npc_pierre_talk.png`
  (8×1 each); loader prefers the split files. Both gap-detect clean.
- §POULAIN-GIVE — `npc_madame_poulain_give.png` re-rolled to a baguette (was a
  wrapped present); code flips her to face PP.
- §JUMPBACK — `PP jump back.png` landed; wired into the biker bump (jump_back
  one-shot, flinch fallback).
- §CAM2 — `cafe_patron_camille_lostpencil.png` landed; wired as her
  `lost_pencil` one-shot, plays (with hold) on her lost-pencil ask.
- §BK2 — biker re-rolled with a TRANSPARENT background; loader switched to the
  raw (no-key) path. Interior bike-frame pockets are see-through now.

**2026-06-12 prune (user: "too much in the file that is not relevant any
more") — bodies deleted, recover from git if a re-open is ever needed:**

- §AB — PP walk side body · already DONE in the 06-10 batch; body removed.
- §OD — Higgins office no-desk regen · RETIRED: the current office sheets
  have been accepted through every playtest since 06-05 (and got a tol-4
  color-key fix 06-12); re-open only if the desk clash resurfaces.
- §PR1 — PP generic receive body · already DONE in the 06-10 batch.
- §JW1/§JW2/§JC1/§JI1 — Jerusalem chain bodies compressed to the one-line
  "Parked" list (chain not built yet; prompts get rewritten when it is).
- §JIT-PP1 — PP idle front · RETIRED: regen #1 fixed the foot drift; the
  31px horizontal remainder hasn't been visible in any playtest since.
- §JIT-POULAIN — idle/talk re-roll · RETIRED: her sheets render correctly
  after the 06-12 bust-scale fix; re-open if the idle-vs-talk outfit
  mismatch (#26) still reads in-game.
- §JIT-COLETTE / §JIT-JAKE / §JIT-LILY / §JIT-FLOWER (grab) · RETIRED to
  low-priority: their remaining numbers are foot/center DRIFT, which the
  per-frame feet-anchoring renderer (2026-06-10) cancels on screen — art
  polish only, nothing visibly wrong in-game.

- §AC — PP talk front: natural speech · RETIRED 2026-06-10, superseded by
  §JIT-PP (talk-front bullet carries the natural-speech rules + dims fix).
- §BC — Beaumont talk match idle · DONE 2026-06-10: new 8×1 strips for both
  idle and talk landed 2026-06-05 (loader at npc.go newMuseumCurator), and
  the jitter audit measured both clean — the 1912×823 spec no longer exists.
- §SH / §MM / §CO / §PO / §AA (older parked regens) · SUPERSEDED 2026-06-10
  by the §JIT batch above, which re-queues the same sheets with measured
  drift numbers from tools/jitter_audit.
