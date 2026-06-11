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

### Playtest pass — museum (2026-06-05, focused)

Most of this pass's sprites were regenerated and their loaders updated to
match. The previously "parked" prompts (Higgins shout, Marcus unify, Colette
talk, PP talk front, etc.) are now re-queued with measured drift numbers in
the **§JIT batch** at the bottom of this file — no need to dig through git
history. §AB (walk-side stride) landed 2026-06-10 — see Done log.

#### §AB — PP walk side: full REGULAR walk cycle · DONE

**Path:** `assets/images/player/PP walk left.png` · **1536×1024, 8×1**
(8 cells of 192×1024). Loader auto-detects 8 cols from sheet width (≤1600px).
Generated + cleaned 2026-06-10; jitter audit clean on foot/center anchor
(ghost-piece warnings only on tail follow-through frames — expected).

===PROMPT START===
> A side-view WALK CYCLE film strip of the slim Pink Panther facing LEFT
> (plain pink paws, NO gloves, belly off-white #F2EFE5), EXACTLY 8 frames in
> ONE horizontal row on pure #FFFFFF at exactly 1536×1024 — 8 cells of
> 192×1024, ONE complete panther per cell, nothing touching any cell edge,
> no ghost or duplicate limbs, no separators, no extra portrait, no text.
> This is ONE continuous, REGULAR stride split into 8 EVENLY SPACED moments,
> in order left to right, looping perfectly (frame 8 flows straight back
> into frame 1). Use the classic animator's 8-pose cycle — contact, down,
> passing, up — then the mirrored half:
>   1 CONTACT A — left foot planted forward heel-down, right leg stretched
>     back on its toe; arms swing OPPOSITE the legs (right arm forward).
>   2 DOWN A — weight sinks onto the left leg (body at its LOWEST point),
>     right foot peeling off the ground behind him.
>   3 PASSING A — right leg swings under the body, legs nearly together,
>     body rising, both arms passing the hips.
>   4 UP A — body at its HIGHEST point, right leg reaching forward, left
>     toe pushing off behind.
>   5 CONTACT B — the exact MIRROR of frame 1: right foot planted forward
>     heel-down, left leg back on its toe, left arm forward.
>   6 DOWN B — mirror of frame 2: weight sinks onto the right leg.
>   7 PASSING B — mirror of frame 3: left leg under the body, rising.
>   8 UP B — mirror of frame 4: left leg reaching forward → loops into 1.
> The stride is REGULAR: identical step length both halves, no skipped,
> repeated or out-of-order poses, no half-steps. He walks IN PLACE — the
> body stays on the SAME spot at the SAME size in every cell, the planted
> foot contacts the SAME ground row every time; only the limbs, the long
> tail (trailing behind with follow-through) and a subtle head bob move.
> True side profile in every single frame — never three-quarter view.
===PROMPT END===

### Playtest pass 3 (museum/Paris polish) — 2026-06-05

Engine sides are done (grid counts, anchoring, positions, dialog flow). These
need regenerated PNGs; **keep each character's current design**.

#### §OD — Higgins office idle + talk WITHOUT his own desk — #16/#17

**Paths:** `npc_director_higgins_office_idle.png`, `npc_director_higgins_office_talk.png` · **Grid:** `6×2` (match current)

===PROMPT START===
> Regenerate Director Higgins's office IDLE and TALK sheets, 6 columns × 2 rows
> on pure #FFFFFF, KEEP his current design. CRITICAL: draw ONLY Higgins (head,
> torso, arms) — do NOT draw a desk/counter in the sprite. The scene background
> already has the office desk; the desk baked into the current sprites
> double-draws and clashes. Just the seated/leaning man, lower body implied,
> evenly spaced poses, no bleed. The talk sheet = mouth moving + gestures.
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

### Jerusalem task-chain art — DESIGN PENDING BUILD (see `docs/JERUSALEM_TASKS.md`)

These are for the full Jerusalem daisy-chain, which is **designed but not yet
wired**. Don't generate them until we start building that chain — the exact items
may still change. Listed here so the asset list lives with the other prompts.

#### §JW1 — PP one-shot: write the note

**Path:** `assets/images/player/PP write note.png` · single row, ~6 frames. Ends with PP pocketing the note (standing PP rule).

===PROMPT START===
> The Pink Panther writing a short note, side/three-quarter view, hand-drawn 90s
> cartoon style on pure #FFFFFF. A ~6-frame single-row strip: he holds a small
> slip of paper, presses a pencil to it and writes a line or two (hand moving
> across), then nods, folds the note once, and tucks it into his invisible hip
> pocket — ending empty-handed in a relaxed standing pose with a small smile. No
> gloves, no pure white on the panther. One pose per cell, even spacing, no
> separators or gridlines.
===PROMPT END===

#### §JW2 — PP one-shot: put the note in the wall

**Path:** `assets/images/player/PP put note in wall.png` · single row, ~6 frames.

===PROMPT START===
> The Pink Panther tucking a folded paper note into a crack between large ancient
> stone blocks (the Western Wall), side/three-quarter view, hand-drawn 90s cartoon
> style on pure #FFFFFF. A ~6-frame single-row strip: he reaches up to the wall
> with the folded note, presses it gently into a gap between stones, pats it in,
> and lowers his hand, ending in a quiet respectful standing pose. Only the
> panther and the small bit of stone wall he touches — no full background. No
> gloves, no pure white on the panther. One pose per cell, even spacing, no
> separators.
===PROMPT END===

#### §JC1 — Alley cat (market prop)

**Path:** `assets/images/locations/jerusalem/npc/alley_cat.png` · TRANSPARENT background · single-row strip, ~6 frames: a sit/idle loop then a "trot away" so it can walk off when lured.

===PROMPT START===
> A small scruffy Middle-Eastern street cat, side profile facing RIGHT, hand-drawn
> 90s cartoon style, TRANSPARENT background (NOT white). A ~6-frame single-row
> strip: frames 1-3 = sitting and flicking its tail / licking a paw (idle loop),
> frames 4-6 = standing and trotting away. Tabby/sandy fur, green eyes, simple and
> small (it's a background prop). One pose per cell, even spacing, no separators.
===PROMPT END===

#### §JI1 — Jerusalem item icons (inventory)

**Path:** `assets/images/items/<name>.png` each — one square icon per item, plain/
transparent background, bold 90s cartoon style: `guidebook` (a little travel
guidebook), `sardine_tin` (an opened tin of sardines), `glow_bug_jar` (a jar with
a faint glow inside), `charcoal_stick` (a stick of black charcoal), `paper` (a
blank slip), `pencil` (a short pencil), `note` (a folded paper note). Generate as
needed when the chain is built; keep the style consistent with the existing
`assets/images/items/*.png`.

#### §PR1 — PP generic RECEIVE one-shot (rule SKILL.md §8b) · DONE

**Path:** `assets/images/player/PP receive.png` · **1536×1024, 8×1** (cell
192×1024) · A reusable "PP receives a small item from someone" one-shot.
Registered as player one-shot `"receive_item"`; all five §PR1 call-sites in
`game.go` now use it (Pierre portrait, Beaumont postcard ×2, Poulain coffee
refill, Camille sketch). Jitter audit flags ghost pieces in frames 2–5
(the generic card prop reads as a detached piece — expected for this sheet).

===PROMPT START===
> An 8-frame one-shot of the slim Pink Panther (plain pink paws, NO gloves)
> RECEIVING a small flat item handed to him from his RIGHT, single row of 8
> frames on pure #FFFFFF at exactly 1536×1024 (cells 192×1024). Sequence:
> frames 1-2 he extends his right paw to the side at chest height, frames 3-4
> a small neutral rectangle (card/paper-sized, keep it generic) rests in his
> paw and he looks at it appreciatively, frames 5-6 he brings it to his hip,
> frames 7-8 it vanishes into his invisible hip pocket and he ends relaxed,
> empty-handed, with a satisfied look. ANCHOR LOCK — in EVERY cell his feet
> sit on the SAME pixel row and his body's vertical centerline on the SAME
> pixel column; same size every frame; nothing touches a cell edge. No
> separators, no extra portrait, no text.
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

#### §JIT-PP1 — PP idle front (foot drift 63px, center 40px)

**Path:** `assets/images/player/PP idle front.png` · **1536×1024, 8×2** (cell 192×512)

**Regen #1 (2026-06-10): foot drift FIXED; 31px horizontal drift remains —
borderline. Re-roll only if the slide is still visible in-game.**

===PROMPT START===
> A 16-frame front-facing IDLE loop of the slim Pink Panther (plain pink paws,
> NO gloves, belly off-white #F2EFE5), 8 columns × 2 rows on pure #FFFFFF at
> exactly 1536×1024 (cells 192×512). This is an ANIMATION: every frame must be
> CLEARLY different from its neighbours — classic 90s Pink Panther cool-cat
> personality, not 16 copies of one pose. Play 1→16 in order, looping:
>   Frames 1-4 — relaxed stance, chest rises and falls with a slow breath,
>     tail swishes left → right behind him with follow-through.
>   Frames 5-6 — he cocks his head slightly and glances LEFT, one eyebrow up,
>     tail curls up at the tip.
>   Frames 7-8 — lazy slow BLINK, head returns to center.
>   Frames 9-12 — he shifts his weight onto one hip (shoulders tilt, knee
>     bends, hip pops sideways), inspects the back of one paw, bored.
>   Frames 13-14 — glances RIGHT, ear twitch, tail flicks fast once.
>   Frames 15-16 — settles back to the frame-1 pose so the loop closes.
> ANCHOR LOCK — through ALL of this he stays nailed to ONE spot: in EVERY cell
> both feet contact the SAME pixel row and his body's vertical centerline
> stays on the SAME pixel column (the hip/shoulder action bends AROUND that
> line, the feet never slide). Same size every frame; nothing touches a cell
> edge. No separator lines, no extra portrait, no text.
===PROMPT END===

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

#### §JIT-MARCUS — Marcus talk + all strange sheets (regen #1 FAILED: figures straddled cell borders + ghost limbs; reverted to the previous art)

**Paths:** `assets/images/locations/camp/npc/kids/marcus/npc_marcus_talk.png`,
`npc_marcus_strange_idle.png`, `npc_marcus_strange_talk.png`,
`npc_marcus_strange_alt.png` · **1536×1024, 8×2 each** ·
`npc_marcus_idle.png` came back CLEAN — match its framing, centering and
character size EXACTLY. Paste once per sheet, swapping the [STATE] line.

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
> - TALK: sketching pose, mouth cycling natural speech shapes, light gestures.
> - STRANGE IDLE: hollow-eyed, compulsively drawing, slow eerie sway.
> - STRANGE TALK: hollow-eyed, mouth moving, never looking up from the page.
> - STRANGE ALT: pad held close to his face, scribbling faster and faster;
>   mid-loop he turns the pad and stares at what he drew, then resumes.
> ANCHOR LOCK — the boy is nailed to ONE spot: in EVERY cell his feet sit on
> the SAME pixel row and his body's vertical centerline on the SAME pixel
> column; only arms, face and the notepad move. No sliding, no size changes.
> No separators, no extra portrait, no text.
===PROMPT END===

#### §JIT-POULAIN — Madame Poulain idle + talk (REGEN #1 FAILED: drift got WORSE — 142/167px)

**Paths:** `assets/images/locations/paris/npc/coffee/npc_madame_poulain_idle.png`,
`npc_madame_poulain_talk.png` · **1536×1024, 8×2 each**. The WORK sheet came
back CLEAN on 2026-06-10 — only idle and talk need a re-roll. When re-rolling,
emphasise: her waist cutoff is a HARD straight line at the same row in all 16
cells, as if the bottom of the figure were sliced by a ruler. Paste once per
sheet, swapping the [STATE] line.

**Visibility rule (user 2026-06-10):** in-game she stands BEHIND the counter —
everything below her waist cutoff is hidden by the desk. So ALL action (hands,
dough, gestures) must happen at CHEST height or higher; anything drawn at
waist/counter level will be invisible in the scene.

===PROMPT START===
> A 16-frame animation sheet of Madame Poulain, a warm middle-aged Parisian
> baker (KEEP her current design: apron, hair in a bun — see reference),
> UPPER BODY behind-the-counter pose, 8 columns × 2 rows on pure #FFFFFF at
> exactly 1536×1024 (cells 192×512). IMPORTANT: in the game a counter hides
> everything below her waist — keep all hands/props at CHEST height or
> higher in every frame.
> [STATE — pick one:]
> - IDLE: arms loosely folded at chest height, small breathing motion, a
>   warm smile, an occasional glance to the side.
> - TALK: mouth cycling natural speech shapes, light hand gestures held UP
>   at chest level (never dropping to the waist).
> ANCHOR LOCK — she is nailed to ONE spot: in EVERY cell her waist cutoff sits
> on the SAME pixel row (bottom of the figure, where the counter hides her)
> and her body's vertical centerline on the SAME pixel column; only arms,
> face and the dough move. No bobbing up and down between cells (this sheet's
> current bug), no size changes, nothing touching cell edges. No separators,
> no extra portrait, no text.
===PROMPT END===

#### §JIT-COLETTE — Colette talk (regen #1 improved: 100→43px center, 37px foot remain)

**Path:** `assets/images/locations/paris/npc/outside/npc_madame_colette_talk.png`
· **1536×1024, 8×2** (regen #1 landed 2026-06-10 as a clean 8×2 and the
loader was switched back to 8×2 ✓; moderate drift remains — re-roll if
visible in-game).

===PROMPT START===
> A 16-frame TALK loop of Madame Colette, an elegant Parisian guide (KEEP her
> current design; cream shirt #E5DDC8 — NEVER pure white fabric), 8 columns ×
> 2 rows on pure #FFFFFF at exactly 1536×1024 (cells 192×512). Mouth cycles
> natural speech shapes, graceful hand gestures. ANCHOR LOCK — she is nailed
> to ONE spot: in EVERY cell her feet sit on the SAME pixel row and her
> body's vertical centerline on the SAME pixel column; no sliding sideways
> (this sheet's current bug), no size changes, nothing touching cell edges.
> All 16 cells filled — no blank last frame. No separators, no extras, no text.
===PROMPT END===

#### §JIT-JAKE — Jake strange talk (regen #1 improved: 115→87px foot, still drifting)

**Path:** `assets/images/locations/camp/npc/kids/jake/npc_jake_strange_talk.png`
· **1536×1024, 8×2** · Also note: jake strange IDLE drifts 36px — borderline,
fix in the same pass if re-rolling anyway.

===PROMPT START===
> A 16-frame STRANGE-STATE talk loop of Jake, the tough camp kid (KEEP his
> current strange design: vacant stare, clutching his coin collection),
> 8 columns × 2 rows on pure #FFFFFF at exactly 1536×1024 (cells 192×512).
> Mouth moves in slow, distant speech; he never lets go of the coins.
> ANCHOR LOCK — the boy is nailed to ONE spot: in EVERY cell his feet sit on
> the SAME pixel row and his body's vertical centerline on the SAME pixel
> column; only mouth, eyes and hands move. No vertical jumping between cells
> (this sheet's current bug), no size changes, nothing touching cell edges.
> No separators, no extra portrait, no text.
===PROMPT END===

#### §JIT-LILY — Lily idle + talk (regen #1 barely moved: now 53/65px foot drift)

**Paths:** `assets/images/locations/camp/npc/kids/lily/npc_lily_idle.png`,
`npc_lily_talk.png` · **1536×1024, 8×2 each** (§LL rule: KEEP her design, fix
anchors only). Paste once per sheet, swapping the [STATE] line.

===PROMPT START===
> A 16-frame animation sheet of Lily, the shy flower-loving camp kid (KEEP
> her current design exactly — see reference), 8 columns × 2 rows on pure
> #FFFFFF at exactly 1536×1024 (cells 192×512).
> [STATE — pick one:]
> - IDLE: shy stance, hands together, small sways, an occasional glance down.
> - TALK: quiet speech, small mouth shapes, bashful gestures.
> ANCHOR LOCK — she is nailed to ONE spot: in EVERY cell her feet sit on the
> SAME pixel row and her body's vertical centerline on the SAME pixel column;
> only face, hands and hair move. No vertical jumping between cells (this
> sheet's current bug), no size changes, nothing touching cell edges. No
> separators, no extra portrait, no text.
===PROMPT END===

#### §JIT-FLOWER — PP grab flower (regen #1 improved but still bad: 70px center, 23% pump)

**Path:** `assets/images/player/PP grab flower.png` · **1536×1024, 6×1** ·
Standing PP rules: NO gloves; the pickup ENDS with PP pocketing the flower.

===PROMPT START===
> A 6-frame one-shot of the slim Pink Panther (plain pink paws, NO gloves)
> picking a small flower, single row of 6 frames on pure #FFFFFF at exactly
> 1536×1024 (cells 256×1024). Sequence: stands relaxed → bends slightly at
> the waist reaching down-forward → plucks the flower → straightens, admires
> it for one frame → brings it to his hip → it vanishes into his invisible
> hip pocket and he ends relaxed, empty-handed. ANCHOR LOCK — in EVERY cell
> his feet sit on the SAME pixel row and his body's vertical centerline on
> the SAME pixel column, and he is drawn at the SAME size in every frame
> (the current sheet grows/shrinks 23% and slides 70px — that is the bug;
> the bend happens AROUND the anchor, the feet never move). Nothing touches
> a cell edge; no separators, no extra portrait, no text.
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

- §AC — PP talk front: natural speech · RETIRED 2026-06-10, superseded by
  §JIT-PP (talk-front bullet carries the natural-speech rules + dims fix).
- §BC — Beaumont talk match idle · DONE 2026-06-10: new 8×1 strips for both
  idle and talk landed 2026-06-05 (loader at npc.go newMuseumCurator), and
  the jitter audit measured both clean — the 1912×823 spec no longer exists.
- §SH / §MM / §CO / §PO / §AA (older parked regens) · SUPERSEDED 2026-06-10
  by the §JIT batch above, which re-queues the same sheets with measured
  drift numbers from tools/jitter_audit.
