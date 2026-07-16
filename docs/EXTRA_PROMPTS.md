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

> Hand-drawn 1990s Saturday-morning cartoon, Pink Panther _Hokus Pokus Pink_
> (1997) / _Passport to Peril_ (1996). Confident black ink linework ~3 px,
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

**★ HARD RULE (user 2026-06-30): ZERO pure `#FFFFFF` on ANY drawn element of
ANY sprite — characters, props, painter canvases, aprons/smocks, paper, teeth,
AND eyes (eye sclera = pale grey `#C4C4C4`, never white).** The engine
chroma-keys pure white and punches it into holes / white spots. This is why
Pierre is being recreated (his old canvas + smock + eyes were pure white).
Only the scene-cell BACKGROUND stays pure white. Verify any new sheet with
`SPRITE_SCAN=1 go test ./engine -run TestSpriteScan`.

**No white halo / fringe rule:** anti-aliasing on the sprite itself must NOT use
pure white or near-white. After the background is removed there must be no tiny
white rim around the character, eyes, props, flowers, canvas, smock, hands, or
outline. Use pink/cream/grey edge pixels that match the nearby fill instead
(for example pale-grey `#C4C4C4` for eye sclera and cream `#E5DDC8` /
bone `#EDE5D3` for light props). Pure `#FFFFFF` is allowed only in the empty
cell background.

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

### 2026-07-15 — PP full 12-frame walk cycles (front + back, 2×6)

User request: replace the near-static front/back walks with REAL 12-frame
cycles. Two sheets, each **1536×1024, 2 rows × 6 columns (cells 256×512)**,
play order row 0 left→right then row 1 left→right. Loader wiring exists
(current sheets load 8×2; when these land, update the grid args in
`game/player.go` to 6×2 — note in FIXME).

#### §PP-WALK-FRONT-12 — `PP walk front.png` · LANDED 2026-07-15 (verified, wired 6×2)

===PROMPT START===

> The classic 1964 cartoon Pink Panther (match the reference EXACTLY: slim
> pink panther, ivory `#F2EFE5` belly and muzzle, plain pink paws with NO
> gloves, long tail, half-lidded yellow eyes). A 12-frame WALK CYCLE seen
> from the FRONT (walking toward the camera), on pure white `#FFFFFF` at
> exactly 1536×1024 — 2 rows × 6 columns, cells 256×512, ≥15px clear white
> between figures and to every sheet edge. ANCHOR LOCK: feet contact the
> SAME pixel row and the body centerline stays on the SAME column in every
> cell; same size every frame. This is ONE continuous stride, every frame
> CLEARLY different:
> Row 0 — LEFT leg stride:
> 1. CONTACT — left foot planted forward (toward viewer), right toe back,
>    arms mid-swing (right arm forward), tail curls left.
> 2. DOWN — weight sinks onto the left leg, body drops a touch, both knees
>    soft, arms passing the hips.
> 3. PASS — right leg lifts and passes the left, body at its lowest,
>    tail swings behind the body.
> 4. UP — body rises, right knee high toward the camera, left heel
>    lifting, arms reversing.
> 5. REACH — right foot reaching forward-down, body tallest, arms at
>    full opposite swing (left arm forward).
> 6. CONTACT-PREP — right foot about to plant, left toe rolling off.
> Row 1 — RIGHT leg stride (mirror timing):
> 7. CONTACT — right foot planted forward, left toe back, left arm forward.
> 8. DOWN — weight sinks onto the right leg.
> 9. PASS — left leg lifts and passes, tail swings to the right.
> 10. UP — left knee high, right heel lifting.
> 11. REACH — left foot reaching forward-down, right arm at full swing.
> 12. CONTACT-PREP — left foot about to plant, loops cleanly back to 1.
> No separators, no text, no extra portrait, nothing pure white ON the
> character (belly/muzzle stay ivory `#F2EFE5`).

===PROMPT END===

#### §PP-WALK-BACK-12 — `PP walk back.png` · LANDED 2026-07-15 (verified, wired 6×2)

===PROMPT START===

> The classic 1964 cartoon Pink Panther seen from BEHIND (match the
> reference EXACTLY: slim pink panther, plain pink paws, NO gloves, long
> tail visible hanging behind the body, back of the head with both ears —
> NO facial features except whisker tips peeking past the cheeks). A
> 12-frame WALK CYCLE walking AWAY from the camera, on pure white `#FFFFFF`
> at exactly 1536×1024 — 2 rows × 6 columns, cells 256×512, ≥15px clear
> white between figures and to every sheet edge. ANCHOR LOCK: feet on the
> SAME pixel row, centerline on the SAME column, same size every frame. One
> continuous stride, every frame CLEARLY different:
> Row 0 — LEFT leg stride:
> 1. CONTACT — left foot stepping away (sole of the left foot showing),
>    right heel down, right arm swung forward out of view, tail sways left.
> 2. DOWN — body sinks onto the left leg, shoulders dip.
> 3. PASS — right leg passes under the body, sole hidden, body lowest,
>    tail hangs centered.
> 4. UP — body rises, right leg lifting behind (sole starting to show).
> 5. REACH — right foot reaching away-forward, body tallest, tail sways
>    right.
> 6. CONTACT-PREP — right foot about to plant.
> Row 1 — RIGHT leg stride (mirror timing):
> 7. CONTACT — right foot stepping away, left heel down, tail sways right.
> 8. DOWN — sink onto the right leg.
> 9. PASS — left leg passes under the body.
> 10. UP — left leg lifting behind.
> 11. REACH — left foot reaching away, tail sways left.
> 12. CONTACT-PREP — left foot about to plant, loops cleanly back to 7→1.
> No separators, no text, no extra portrait, nothing pure white ON the
> character.

===PROMPT END===

### 2026-06-24 bug-sweep — new give/receive beats + scene art

All code below is already wired with graceful fallbacks (missing sheet = the
generic give/grab reach, never a deadlock), so the game runs today; each PNG
that lands auto-upgrades the matching beat. Feed the global style/standing
rules above once per session. All grids are tall cells, no separators, one
figure per cell, PP design-locked to the reference sheets, output ONLY the grid.

**Match-the-item rule (user 2026-06-24).** When a hand-over sprite shows an
object being given/taken, that object MUST look like the inventory item the
player actually receives. Attach the matching item PNG from
`assets/images/items/` as a reference and copy its shape/colour exactly:

- coffee → `assets/images/items/cafe_au_lait.png`
- Camille's sketch → `assets/images/items/camille_sketch.png`
- press → `assets/images/items/press_pass.png`
  (and likewise for any other handed item).

**One-character-per-prompt rule (user 2026-06-24).** Each prompt below is for
ONE character only. A character's idle and talk sheets must be the SAME person
(same face/clothes/proportions) — never let a sheet drift into a different-looking
NPC. Where two NPCs share a stall they get fully separate prompt blocks.

**Separate idle + talk per NPC (user 2026-06-24 — now a standing rule).** EVERY
speaking NPC gets TWO sheets: a `*_idle.png` (resting loop) and a `*_talk.png`
(mouth/gesture loop), same character/size/anchor across both. Do NOT pack idle
on row 0 and talk on row 1 of one sheet — author them as two separate 8×1 files.
The loaders already look for the talk sheet and fall back to the idle until it
lands. (This rule is also captured in the `sprite-check` skill.)

#### §M16b — REGEN `npc_pierre_get_jam.png` — ghost second hand (user 2026-06-24) · SUPERSEDED

Superseded by the full five-sheet Pierre rebuild prompt at §PIERRE-REBUILD.
The currently landed sheet has been repaired enough to pass the no-white scan,
but the proper final fix is to regenerate the full Pierre set together.

#### §BG-OFFICE — Higgins office mid-dark + dark (FIXME #22) · NEEDED (2 BGs)

The camp darkens by mood level; the office needs matching variants (the loader
falls back to day1 until these exist). Save:

- `assets/images/locations/camp/background/day2/camp_office.png` — same office, dusk/uneasy lighting (mid-dark).
- `assets/images/locations/camp/background/day3/camp_office.png` — same office, night/wrong lighting (fully dark).

> Repaint the existing day1 camp_office background under [dusk | night] mood:
> cooler shadows, dimmer warm lamp, same composition and furniture. Full scene
> background (not a sprite), no characters, pure scene art.

#### §H5 — Higgins office talk blink · RE-ROLL (sprite-check diagnosed 2026-06-28)

`assets/images/locations/camp/npc/higgins/npc_director_higgins_office_talk.png` (6×2).

Sprite-check verdict: the sheet gap-detects cleanly and no cell is blank/half-cut,
so it's NOT a grid/cut bug — it's an **anchor-drift** problem the engine can't fully
cancel, so the body wobbles each frame (reads as a "blink"). Two concrete defects to
fix in the regen:

1. **Gesture slide (the main cause):** the body/head shifts horizontally between
   frames because the arm pose changes the figure's overall width/position. FIX:
   **anchor-lock** — Higgins's HEAD and TORSO must sit at the EXACT same column and
   the SAME size in all 12 frames; ONLY the arms/hands and mouth move. His body must
   not drift left/right or up/down at all between frames.
2. **Row offset:** the figures on the top row sit ~118px LOWER in their cells than
   the bottom row (top-row content starts ~y129, bottom ~y11). FIX: every cell across
   BOTH rows places the figure at the **same vertical position** (same waist-cutoff row).

Also: **no single closed-eye "blink" frame** — keep his eyes consistently open while
talking (if he blinks, spread it over 2–3 frames so it reads smooth, never 1 sharp frame).

Format: 6×2, upper-body/bust only (no desk in the sprite — the BG table covers his
lower half), ≥15px gaps, one figure per cell, cream `#E5DDC8` not pure white, same
character/size/anchor as `npc_director_higgins_office_idle.png`.

### 2026-06-24 Japan / Kyoto chapter (Batch 2)

The 3 scenes (torii → ramen → grove), dresser-geisha gag, Danny's call and the
lake heal are wired; the landed art is in the Done log. Remaining art below
upgrades it (all run on idle/placeholder until they land):

#### §JP-LILY-TALK — sad Lily TALK sheet · LANDED 2026-07-14 (auto-wired, verified)

Her idle landed; she needs a matching TALK sheet (until then she stays in the
sad idle while speaking). Save: `assets/images/locations/camp/npc/kids/lily/npc_lily_sad_talk.png` (8×1).

> Same sad Lily as `npc_lily_sad_idle.png` — seated at the dock end, seen from
> BEHIND hugging her knees — but now gently speaking: small shoulder/head moves
> as if murmuring, still facing away. Same size/anchor as the idle sheet.

#### §JP-HIGGINS-RUDE — rude Higgins TALK + WALK · NEEDED (2 sheets)

For the grounds intercept (he comes halfway, is curt, then leaves). Save under
`assets/images/locations/camp/npc/higgins/`, 8×1 (or 8×2 to match his other sheets):

- `npc_director_higgins_rude_talk.png` — Higgins talking but dismissive/impatient: arms crossed, waving PP off, glancing at his watch.
- `npc_director_higgins_walk_front.png` — Higgins walking toward the camera/PP (front-facing walk), to "come half the way" down the path.

#### §JP-TOURIST — Gary's upside-down-book gag · NEEDED (4 sheets)

User flow: Gary STARTS holding his guidebook UPSIDE-DOWN; when PP talks to him he
FLIPS it the right way and KEEPS it that way afterward. So he needs a before
state, the flip, and an after state. Save under `assets/images/locations/japan/npc/`,
all 8×1, same tourist character:

- `npc_gary_idle.png` — idle, guidebook held **UPSIDE-DOWN** (the start state; re-roll the current idle so the book is clearly inverted).
- `npc_gary_talk.png` — talking, book still upside-down (this plays during the chat, before the flip).
- `npc_gary_flip.png` — one-shot: he turns the book end-for-end the right way up, ending pleased (one continuous motion).
- `npc_gary_idle_flipped.png` + `npc_gary_talk_flipped.png` — idle + talk with the book now held CORRECTLY (the after state the code swaps to once he's flipped it).

#### §JP-LEAVES — falling leaves over the ramen-store tree · NEEDED

A live falling-leaf loop (user: the tree should drop leaves). Save:
`assets/images/locations/japan/props/leaf_fall.png` (**3 frames**, one row).

> THREE frames of a single small autumn/sakura leaf fluttering as it falls:
> frame 1 tilted left, frame 2 flat/spinning, frame 3 tilted right — a loop that
> reads as a leaf twirling down. Transparent (or pure-white) background; the leaf
> warm orange-pink, no pure white on the leaf itself.

#### §JP-DRESSER — Kiku the dresser TALK (idle landed) · NEEDED

`drawer.png` (idle) landed. Add her talk sheet. Save:
`assets/images/locations/japan/npc/drawer_talk.png` (8×1).

> Same geisha-styled kimono dresser as `drawer.png`, now talking with bright,
> bossy gestures (beckoning, measuring PP with her eyes, holding up silk). Same
> character/size/anchor as the idle.

#### §JP-KIMONO — PP spins into a kimono · NEEDED (the dresser gag)

Save: `assets/images/player/PP_kimono_spin.png` (8×1, PP design-locked).

> 8-frame one-shot of Pink Panther doing a fast spin-change: frames 1-2 he spins
> (motion-blur swirl), frames 3-6 he's mid-spin now wearing a **MEN'S kimono**
> (see the men's-kimono note below), frames 7-8 he spins again and ends back in
> his normal look, striking a tiny pose. One continuous spin; plain pink paws,
> no gloves; no pure white.

> **Men's-kimono rule (PP's outfit, user 2026-06-24):** PP is male, so every
> sprite where he wears a kimono (this gag + the tea-ceremony sheets) uses a
> MEN'S kimono — subdued colour (navy / charcoal / deep plum-brown), NARROW
> sleeves sewn to the body (NOT the long flowing furisole sleeves), a thin obi
> tied LOW on the hips, optionally a short _haori_ jacket. NOT a bright,
> long-sleeved, big-bow women's/geisha kimono.

#### §JP-NPC-TALK — talk sheet for the ramen cook · NEEDED

Per the separate idle/talk rule, Hiro still needs his talk sheet (idle landed).
Save `assets/images/locations/japan/npc/npc_hiro_talk.png` (8×1, same Hiro):
Hiro talking, ladle gestures over the pot. (Gary's talk sheets are in §JP-TOURIST.)

#### §JP-GAPS — Japan idle-sheet cut check · RESOLVED 2026-06-24

The earlier gap-fallback sheets are now fixed: `npc_hiro_idle`, `npc_obachan_idle`
and `npc_kenji_idle` were re-rolled and gap-detect cleanly at 8×1; `npc_gary_idle`
turned out to be drawn **6×2** (12 frames, two rows) and is now loaded at that
grid (gap-detects clean). No Japan NPC sheet currently mis-cuts. (Separate
pre-existing issue, not Japan: `PP put note in wall.png` still falls back at 6×1.)

#### §JP-HIRO-COUNTER — Hiro BEHIND the open counter, upper body · NEEDED (2 sheets)

When the stall opens, Hiro moves behind the counter (wired 2026-07-14 in
`openRamenStall`; he keeps his full-body street sprite until these land).
Behind-counter convention (same as office Higgins / Poulain): **upper body
only, waist cutoff, NO counter in the sprite**, all action at chest height.
Save under `assets/images/locations/japan/npc/`:

- `npc_hiro_counter_idle.png` — 1536×1024, **8×1**, one figure per cell,
  ≥15px empty gaps between figures and to the sheet edges. Hiro from the
  waist up facing the customer: wiping the counter, stirring the pot with his
  ladle, tasting the broth, a proud nod — every frame CLEARLY different.
  Same Hiro design as `npc_hiro_idle.png` (headband, apron). Anchor lock:
  waist cutoff row on the same pixel row every frame, centerline on the same
  column. Background pure white `#FFFFFF`; **NO pure white anywhere ON Hiro**
  (apron/towel = cream `#E5DDC8`, steam = pale grey `#C4C4C4`, eye whites
  off-white).
- `npc_hiro_counter_talk.png` — same framing/size/anchor, 8×1: talking to the
  customer — ladle gestures, offering a bowl forward, warm laugh.

#### §JP-RAMEN-STALL — dynamic stall: closed ↔ open prop · LANDED 2026-07-14

Art landed (`ramen_closed.png` + `ramen_open.png`, 1536×1024 each, baked
checkerboard keyed at load). Positioned over the painted BG stall at screen
(645, 530) scale 0.30. Kept for reference:

The ramen stall opens when PP returns Hiro's fire-striker. Wired as a prop over
the stall whose frame swaps on `jp_ramen_open`. Save under `assets/images/locations/japan/props/`:

- `ramen_closed.png` — the stall shuttered/dim: noren curtain furled, window
  boarded or dark, no steam.
- `ramen_open.png` — the SAME stall lit and open: noren hung, lantern glowing,
  steam rising, bowls out. Same size/position as the closed version so the swap
  is in-place. (Single frame each, or a short loop for the open steam.)

#### §JP-QUEUE — waiting line that sits when the stall opens · NEEDED (2 sheets)

A static line of customers outside that SITS at the counter when Hiro opens.
Save under `assets/images/locations/japan/npc/`:

- `customer_wait.png` — a Kyoto local standing patiently in line (1 figure; a
  4-frame idle sway is fine). Make 1-2 visual variants if easy so they're not clones.
- `customer_sit.png` — the same kind of local SEATED on a counter stool, slurping
  ramen. Same scale; this is what each queue figure swaps to on opening.

#### §JP-ITEMS — Kyoto quest item icons · NEEDED (4 icons)

Inventory icons (small, centered, transparent/pure-white bg). Save under `assets/images/items/`:

- `well_water.png` — a simple cup/ladle of water.
- `voice_charm.png` — a paper o-mamori charm brushed with the kanji for "voice".
- `fire_striker.png` — a flint-and-steel fire-striker (hiuchi).
- `offering_bowl.png` — a steaming ramen bowl (the blessed offering).

#### §JP-MATCHA — matcha ceremony (temple tea-house: BG + tea master + items + PP sit) · NEEDED

A required sub-quest gating the grove (`jp_tea_done`): grab matcha + a random
bowl in the flower store, whisk it at the street well, then climb UP to the
**temple tea-house** (new 5th scene) and share a bowl with the tea master.

- `assets/images/locations/japan/background/teahouse.png` — a quiet temple
  tea-house: tatami floor, low table, hanging scroll, shoji screens, soft light.
  Full scene BG, no characters. (Authentic: the tea ceremony grew out of Zen
  temple tea rooms.) Falls back to a flat tatami colour until it lands.
- `assets/images/locations/japan/npc/npc_tea_master_idle.png` + `..._talk.png` (8×1) — a serene
  old tea master, seiza/kneeling at the tatami, whisk in hand (separate idle + talk).
- **`assets/images/player/PP_spin_to_sit.png` (8×1) — THE TRANSFORM (the one
  you're missing).** PP goes from STANDING-NORMAL to SEATED-IN-KIMONO via a fast
  spin, in 8 frames, one continuous motion:
  - f1: standing normal, front.
  - f2-3: a fast spin — motion-blur swirl, body blurred mid-turn (like the
    Kiku kimono-spin gag).
  - f4-5: spin resolves and he's now wearing the **men's kimono** (subdued, narrow
    sleeves, low thin obi — see the men's-kimono rule under §JP-KIMONO), still upright.
  - f6-7: he folds his legs and lowers down.
  - f8: settled SEATED (seiza/kneeling) in the kimono — this final pose must
    match the seated idle below so the hand-off is seamless.
    (Code key `tea_sit`; also accepts `PP_tea_ceremony.png` / `PP_sit_down.png`.)
- `assets/images/player/PP_sit_idle.png` (8×1) — PP kneeling in the kimono, calm
  seated idle (slow breath, hands on knees). Same seated pose as `PP_spin_to_sit`
  frame 8.
- `assets/images/player/PP_sit_talk.png` (8×1) — PP kneeling, gently talking/
  sipping while seated — shown during the seated ceremony dialog.
- Item icons under `assets/images/items/`: `matcha.png` (tin of green powder),
  `tea_bowl.png` (an empty chawan), `matcha_bowl.png` (frothy whisked matcha).

#### §JP-SAKURA-BG — hidden sakura grove (4th scene, NEW) · DONE

The "follow me" payoff: Oba-chan opens a path from the flower store into a deep
pink cherry-blossom grove where PP picks the blossom himself at the old tree.
Wired as scene `tokyo_sakura` (right exit from `tokyo_temple`, gated on talking
to Oba-chan); runs on a flat-pink fallback until the BG lands. Save:
`assets/images/locations/japan/background/sakura_grove.png` (full scene, no characters).

> A secluded grove of old cherry trees in full bloom - a tunnel/clearing of deep
> pink sakura, petals drifting, soft dappled light, a mossy path. One especially
> large, ancient gnarled cherry tree as the centrepiece (right-of-centre, ~x560-880)
> where PP will pick a blossom. Warm pinks and greens, no pure white where PP stands.

#### §JP-BG-EDGES — connect the 3 Kyoto backgrounds · DONE (re-export)

User: each scene's left/right edge should "show what we can get" — a peek of the
adjacent scene so the exits read. Re-export the 3 BGs so:

- `start_of_tori` RIGHT edge hints the ramen street (a lantern/awning corner).
- `ramen-store` LEFT edge hints the torii gates; RIGHT edge hints the pink grove.
- `flower_store_near_forest` LEFT edge hints the ramen street.
  Keep them full-scene, same size; only the edges change.

### 2026-06-15 playtest — background / chroma-key re-rolls

Two sheets read with a leftover "background" in-game. Both are the
**white-on-white chroma-key problem**: the engine keys pure `#FFFFFF`, so any
prop the character holds that is ALSO pure white gets eaten where it touches
the background (or leaves enclosed white pockets / halos the edge-flood can't
reach). The biker (§BK2) and the pot pigeon already ship transparent and load
raw — those two are fine; these two are not.

#### §FLOWER-PICK — `PP grab flower.png` mis-cuts: frame 0 = the daisy alone (#1) · STILL NEEDED

**Path:** `assets/images/player/PP grab flower.png` · **6×1** single row · pure
`#FFFFFF` background (chroma key). The loader keys it at tol 36 already.

**Re-generated 2026-06-20 — STILL BROKEN, same two layout flaws as before:**

- `go test ./engine -run ContentGrid` gap-detects it **1×6 but the FIRST cell is
  only (0,0)-(82,1024)** — that 82px sliver is the daisy lying detached on the
  ground at the far LEFT of frame 1. So the animation's frame 0 is _just the
  daisy with no panther_, and every real PP pose is shoved one slot over → the
  pickup reads as broken/blinking in-game.
- `tools/jitter_audit` (fixed grid): **GHOST PIECES in 2 cells**, **CONTENT
  CROSSES 3 borders (13-16px)**, **"2 cells touch BOTH side edges"** — the 6 PP
  poses are packed so tightly they touch each other, so a fixed grid can't cut
  them either.

The engine cannot recover 6 clean frames from this layout (detached object +
figures with no gap between them). It needs a re-roll that fixes BOTH:
**(1) the daisy must be in PP's paw / within a few px of it in EVERY frame —
never a separate object on the ground at the cell edge; (2) the 6 PP poses need
≥15px of clear background between them so neither touches its neighbour.**

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
- **No white halo:** after background removal, PP and the daisy must have ZERO
  white or near-white fringe pixels clinging to their outlines. Anti-aliased
  edge pixels around PP should be pink/cream/grey matching the body, and daisy
  edge pixels should be bone/cream, never white.

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
> NO white or near-white anti-alias halo on PP or the flower after background
> removal: edge pixels around PP must be pink/cream/grey, and edge pixels around
> the daisy must be bone/cream, not white.
> Pure `#FFFFFF` background only, no ghost/duplicate limbs, no separators or
> gridlines, no portrait or labels. Tall rectangular cells, never square.
> ===PROMPT END===

#### §PIERRE-BOARD — Pierre's easel canvas vanishes in a frame (#6) · SUPERSEDED

Superseded by the full five-sheet Pierre rebuild prompt at §PIERRE-REBUILD.
The currently landed idle/talk sheets have been repaired enough to pass the
no-white scan, but the proper final fix is to regenerate the full Pierre set
together.

### 2026-06-20 bug-sweep art

#### §MARCUS-SLEEP — Marcus dozes off after the postcard heal (#19) · DONE 2026-06-20

LANDED: `npc_marcus_going_to_sleep.png` + `npc_marcus_sleeping.png` (sitting-doze),
both GAP-DETECTED 1×8 clean and wired in `newRoomMarcus`. Prompt kept below for re-rolls.

Two sheets. Once PP gives Marcus the Louvre postcard he's healed and finally calm,
and nods off (so Higgins's "sleeping soundly" line is true). Wired already - the
go-to-sleep one-shot ("sleep") plays at his spot, then his idle swaps to the
sleeping loop; both no-op until the art lands. Match `npc_marcus_idle.png` exactly
(brown hair, round glasses, golden-yellow polo `#EEB421`, brown shorts).

**Pose choice (important — read his room `marcus_room_day.png` first):** his cabin
has a bunk bed on the LEFT, a desk under the window, and his obsessive drawings
papered over every wall and the floor. He renders at his STANDING spot (centre of
the room, bottom-anchored), and the engine scales each frame by its tallest opaque
pose - so a flat, lying-in-bed figure would render stretched and wrong. Instead he
**sinks down and dozes off SITTING on the floor among his drawings**, head drooping
onto his chest, knees up, sketchbook sliding from his lap - a kid who finally
stopped and conked out mid-sketch. This stays vertical (anchors + scales cleanly at
his spot, no bed alignment needed) and reads far better than lying flat.

**Paths:** `assets/images/locations/camp/npc/kids/marcus/npc_marcus_sleep.png`
(doze-off one-shot) and `npc_marcus_sleeping_idle.png` (sleeping loop) ·
**8×1 each, 1536×1024** (cells 192×1024). **ATTACH** `npc_marcus_idle.png`.

===PROMPT START — SLEEP (one-shot)===

> Use the attached sheet as the reference: the SAME boy, Marcus - copy his exact
> design (brown hair, round glasses, golden-yellow polo `#EEB421`, brown shorts),
> size and framing. An 8-frame ONE-SHOT of him nodding off where he stands in his
> camp-cabin bedroom, single row on pure #FFFFFF at exactly 1536×1024 (cells
> 192×1024). KID-FRIENDLY, calm and content - he's just been soothed, not sad.
> Play 1→8: 1 standing relaxed, small smile; 2-3 a big sleepy yawn, rubbing one
> eye; 4-5 sinking down to sit on the floor, legs folding under/in front of him;
> 6 settled cross-legged, shoulders sagging, his little sketchbook sliding off his
> lap; 7 head drooping forward onto his chest, eyes closing; 8 fast asleep sitting
> up, head down, mouth softly closed, a tiny "z" starting above his head. He stays
> in ONE spot - his seated bottom lands on the SAME floor row as his standing feet.
> ANCHOR LOCK: body centerline on the SAME column every cell; one complete figure
> per cell, ≥15px clear background to both neighbours and the sheet edges, nothing
> crosses a cell boundary, no ghost/duplicate limbs, no separators, no portrait,
> no text. Pure #FFFFFF background only.
> ===PROMPT END===

===PROMPT START — SLEEPING IDLE (loop)===

> Use the attached sheet as the reference: the SAME boy, Marcus - copy his exact
> design and size. An 8-frame LOOPING idle of him ASLEEP, sitting cross-legged on
> his cabin floor with his head drooped forward onto his chest (the pose the
> doze-off one-shot ends on), single row on pure #FFFFFF at exactly 1536×1024
> (cells 192×1024). KID-FRIENDLY and peaceful. The only motion across the loop is
> gentle breathing (shoulders/chest rise and fall a few px) and a soft "z z z"
> that drifts up and fades in/out above his head. He does NOT wake or shift pose.
> ANCHOR LOCK: his seated bottom stays on the SAME floor row and his body
> centerline on the SAME column in every cell; one complete figure per cell, ≥15px
> clear background to both neighbours and the edges, nothing crosses a cell
> boundary, no ghost limbs, no separators, no portrait, no text. Pure #FFFFFF only.
> ===PROMPT END===

#### §GIVE-HEEL — `PP give heel.png` blinks (#15) · NEEDED (re-roll)

**Path:** `assets/images/player/PP give heel.png` · **8×1, 1536×1024** (cells
192×1024). `tools/jitter_audit`: PP's extended ARM + the bread heel reach into
the NEIGHBOURING cells (CONTENT CROSSES 18-30px, "cells touch both edges"), so
the gap slicer cuts unevenly → the hand-over blinks in-game. Re-roll with each
pose fully inside its own cell.

===PROMPT START===

> Use the attached `PP idle side.png` as the exact character reference (same
> head, muzzle, outline weight, pink fills, NO gloves, off-white `#F2EFE5` belly).
> An 8-frame one-shot of the Pink Panther, LEFT profile, handing a small bread
> "heel" (a stubby end of a baguette, golden-brown) to someone beside him. Play
> 1→8: 1 stands holding the heel low, 2-3 raises it, 4-5 extends his arm to offer
> it (arm and heel must stay WITHIN the cell - do not let them reach into the
> next frame), 6 the heel leaves his open paw, 7 pulls the arm back, 8 relaxed
> empty-handed stand with a small smile. CRITICAL: each of the 8 poses is fully
> inside its own cell with ≥15px clear background to both neighbours and the
> edges - no arm, paw, tail or heel may cross a cell boundary. ANCHOR LOCK - feet
> on the SAME pixel row, body centerline on the SAME column. Pure #FFFFFF
> background, no ghost limbs, no separators, no portrait, no text.
> ===PROMPT END===

### Backgrounds + ambient life (camp return, darkening, Jerusalem, bg life)

All backgrounds are **1376×768**, drawn in the game's hand-drawn 90s Pink Panther
cartoon style. Backgrounds contain **no characters** (PP/NPCs are drawn on top);
keep the lower third / foreground relatively clear so characters have a floor to
walk on. Ambient "objects" are separate transparent-background overlay sprites.

#### §AMB3 — Ambient: camp crow (lands on the airstrip sign) — #34 · LANDED 2026-07-14 (loader switched to KEYED for the baked checkerboard)

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
> ===PROMPT END===

spice_seller_give.png`(cardamom),`market/npc_coffee_seller_give.png`(the cup),`wall/npc_bagel_seller_give.png`(ka'ak),`wall/npc_praying_man_give.png` (the note slip).

**Fence prop:** `assets/images/locations/jerusalem/props/fence.png` — single-frame
transparent/white-bg overlay of a low separation fence Shimon stands by (plaza).

#### §JERUSALEM — Jerusalem chapter art (new coffee/spice/bagel/note/pen/coin flow, WIRED) · STILL NEEDED

The chain is wired in `game/jerusalem.go`; every NPC borrows a Paris/camp
placeholder sheet and every give-one-shot / icon no-ops or falls back until its
art lands. Art lives under `jerusalem/npc/wall/` (plaza + Wall NPCs) and
`jerusalem/npc/market/` (souk NPCs). Standard rules apply (pure `#FFFFFF`, anchor
lock, ≥15px gaps, no separators/extras). Backgrounds already exist.

**LANDED:** `wall/npc_shimon.png` (6×2 full body, idle row0 / talk row1) ·
`market/npc_spice_seller_idle.png`.

**NPC sheets still needed — FULL BODY, SEPARATE idle/talk for the sellers + kid:**

- `market/npc_spice_seller_talk.png` — **8×1**, matches the landed idle. APPEARANCE:
  earth-tone tunic/apron, small cap, short dark beard, behind a colourful spice stall
  (scooping/offering cardamom).
- `market/npc_coffee_seller_idle.png` + `_talk.png` — **8×1 each**. The coffee + spice
  sellers currently look IDENTICAL (shared placeholder) — make this man VISIBLY DIFFERENT
  from the spice seller: a **dark waistcoat over a cream shirt, a small red fez/cap, a big
  mustache (no beard)**, working a **brass finjan coffee pot** (brewing/pouring a tiny cup).
- `wall/npc_bagel_seller.png` — **6×2** (idle row0 / talk row1), ka'ak (sesame bagel) cart seller.
- `wall/npc_praying_man_idle.png` **8×2** (idle = facing the Wall, praying/swaying) + `wall/npc_praying_man_talk.png` **8×1** (turned toward PP).
- `wall/npc_wall_kid_idle.png` + `_talk.png` — **8×1 each**, a bar-mitzvah-age boy
  (used for both kids). **MUST NOT look like Marcus** — dark brown/black WAVY hair, a
  small KNITTED KIPPAH, a plain PALE-BLUE collared shirt (cream `#E5DDC8` if "white"),
  dark trousers. **NO round glasses, NO golden-yellow polo, NO notepad** (those read as
  Marcus). One kid may hold a small prayer book in the talk frames.

**GIVE one-shots (8×1, the §8b "both sides" rule):** `wall/npc_shimon_give.png`
(pen, then coin), `market/npc_
**Item icons (256×256, white bg):** `items/cardamom.png`, `items/jerusalem_coffee.png`,
`items/bagel.png`, `items/note_paper.png`, `items/pen.png`, `items/coin.png`.

**PP one-shots (wired, grab fallback):** `player/PP write note.png` (6×1, writes the
slip, pockets it) + `player/PP put note in wall.png` (6×1, tucks it into a Wall crack).

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
> ===PROMPT END===

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
> Step pattern per stride: foot plants (body sinks slightly, shoulders
> counter-rotate) → passing pose (legs together, body at its tallest) →
> other foot reaches and plants → repeat mirrored. Arms swing OPPOSITE the
> legs with loose cartoon overlap; the long tail snakes left-right behind
> him a half-beat behind the body; his head bobs subtly with each step.
> ANCHOR LOCK — he struts IN PLACE: in EVERY cell the planted foot contacts
> the SAME pixel row and his body's vertical centerline stays on the SAME
> pixel column (the body rises/sinks with the stride, but never slides
> sideways or drifts up the cell). Same size every frame; nothing touches a
> cell edge. No separators, no extra portrait, no text.
> ===PROMPT END===

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
> ===PROMPT END===

===PROMPT START — GIVE/SCATTER (optional)===

> The SAME pigeon lady (identical design + size as her idle sheet). A SINGLE
> ROW of 8 frames on pure #FFFFFF at exactly 1536×1024 (cells 192×1024) of
> her CALLING and scattering a big handful of crumbs to the side: 1-2 reaches
> into the bag, 3-5 a big underhand toss (crumbs visible mid-air), 6-8 claps
> the crumbs off her hands with a satisfied smile as pigeons flock in. Same
> palette/anchor rules as her idle. No separators, no extra portrait, no text.
> ===PROMPT END===

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
>
> - TALK (healed/normal): cheerful sketching pose, mouth cycling natural
>   speech shapes, light friendly gestures.
> - STRANGE ALT: the kid is a little OFF (eerie-sad, KID-FRIENDLY, NOT horror)
>   - absorbed in the notepad, drawing the same thing over and over, pausing
>     to gaze off, then back to the page; a small uneasy sway. Not smiling, but
>     not scared or distressed - just absent. NEVER bloodshot/sweaty/manic.
>     ANCHOR LOCK — the boy is nailed to ONE spot: in EVERY cell his feet sit on
>     the SAME pixel row and his body's vertical centerline on the SAME pixel
>     column; only arms, face and the notepad move. No sliding, no size changes.
>     No separators, no extra portrait, no text.
>     ===PROMPT END===

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
>
> - a faraway, distracted look - eyes a touch unfocused, not bloodshot, no
>   dark rings, no sweat, no wild stare.
> - absorbed in a notepad, drawing the same thing over and over; he pauses,
>   gazes off, then goes back to the page.
> - a small uneasy sway and the odd blink/shrug - quietly troubled, calm
>   hands (the game adds a faint quiver of its own, so the art stays gentle).
>   He is NOT smiling/cheery, but NOT scared or distressed either - just absent
>   and a bit melancholy. Every frame clearly different (it's an animation).
>   CRITICAL GRID RULE: each of the 16 cells holds EXACTLY ONE complete Marcus,
>   centered with clear white padding on every side and AT LEAST 15px of empty
>   background between neighbouring figures and to the sheet edges. Never two
>   figures touching or sharing a cell, never a figure cut by a cell boundary,
>   no detached ghost limbs or duplicate arms.
>   ANCHOR LOCK - feet on the SAME pixel row, body centerline on the SAME pixel
>   column in every cell; only the head, arms and notepad move. No sliding, no
>   size changes. No separators, no extra portrait, no text.
>   ===PROMPT END===

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
>
> - he speaks in a quiet, distant, distracted way - mouth cycling natural
>   speech shapes (closed - slightly open - wide "ah" - mid - narrow "oo" -
>   closed), but his eyes stay a touch unfocused, looking past you.
> - he keeps half his attention on the notepad he's drawing in, glancing
>   down at it between phrases as if he can't quite stop.
> - small uneasy head tilts and the odd slow blink; gentle, calm hands -
>   NO bloodshot eyes, no dark rings, no sweat, no wild or distressed stare.
>   Not cheery, not scared - just absent and a bit melancholy. Every frame
>   clearly different (it's an animation).
>   CRITICAL GRID RULE: each of the 16 cells holds EXACTLY ONE complete Marcus,
>   centered with clear white padding on every side and AT LEAST 15px of empty
>   background between neighbouring figures and to the sheet edges. Never two
>   figures touching or sharing a cell, never a figure cut by a cell boundary,
>   no detached ghost limbs or duplicate arms.
>   ANCHOR LOCK - feet on the SAME pixel row, body centerline on the SAME pixel
>   column in every cell; only the head, mouth, arms and notepad move. No
>   sliding, no size changes. No separators, no extra portrait, no text.
>   ===PROMPT END===

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
> 1 — stands relaxed, empty-handed.
> 2 — reaches across to his hip, eyes glancing down.
> 3 — paw "into" the invisible pocket at his hip (classic magic-pocket
> gag: the paw just disappears against his side).
> 4 — pulls out the folded map with a small flourish, eyebrows up.
> 5 — holds it in front of his chest with both paws.
> 6 — flicks it open one fold, leaning his head in with interest.
> 7 — map held up and open toward the camera, filling his paws.
> 8 — settles, holding the open map steady (the map screen takes over
> from this pose).
> The map paper is bone #EDE5D3 (never pure white).
> ANCHOR LOCK — in EVERY cell both feet contact the SAME pixel row and his
> body's vertical centerline stays on the SAME pixel column; same size
> every frame; ≥15px clear white between figures and to the sheet edges.
> No separators, no extra portrait, no text.
> ===PROMPT END===

#### §AMB5 — Paris street accordion player · RETIRED 2026-07-15 (user: the street is done, no new NPCs; wiring + PNG removed)

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
> ===PROMPT END===

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
> ===PROMPT END===

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
> Frames 1-8 (one full stride): left knee lifts toward camera (foot
> visibly raising, sole hinted) → body sinks slightly as the left foot
> plants → passing pose, legs crossing, body tallest → right knee lifts →
> right foot plants → passing pose. Shoulders sway opposite the stepping
> leg, arms swing loosely at his sides, head bobs subtly with each plant,
> and the long tail curls into view alternately left and right behind him.
> Frames 9-16: the mirrored second stride so the loop closes.
> ANCHOR LOCK — he walks IN PLACE: in EVERY cell the planted foot contacts
> the SAME pixel row and his body's vertical centerline stays on the SAME
> pixel column (the body rises/sinks ~10px max with the stride, never slides
> sideways). Same size every frame; nothing touches a cell edge; ≥15px white
> gaps between figures. No separators, no extra portrait, no text.
> ===PROMPT END===

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
> 1 — stands relaxed, daisy held low at his side.
> 2 — raises the daisy to chest height, looking at it fondly.
> 3 — turns slightly and begins extending his arm out to the side.
> 4 — arm fully extended, daisy offered, ears perked.
> 5 — the daisy starts leaving his paw (recipient's unseen pull), his
> fingers opening.
> 6 — paw empty and still extended, fingers spread, a happy blink.
> 7 — pulls the arm back toward his chest with cartoon follow-through.
> 8 — back to a relaxed stand, hands free, content smile.
> The daisy: yellow center, IVORY #F2EFE5 petals (never pure white).
> ANCHOR LOCK — in EVERY cell both feet contact the SAME pixel row and his
> body's vertical centerline stays on the SAME pixel column; same size every
> frame; ≥15px clear white between figures and to the sheet edges; no
> separators, no extra portrait, no text.
> ===PROMPT END===

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
>
> - Monsieur Bernard: bearded older Parisian in a flat cap and brown coat,
>   reading his folded Le Figaro newspaper, occasionally sipping coffee.
>   Loop: reads → page rustle → lifts cup and sips → lowers cup → reads.
> - Mademoiselle Camille: young art student, dark bob, red beret, green
>   blouse. Loop: holds her teacup in both hands → sips → lowers it →
>   glances dreamily aside → back to center.
>   CRITICAL GRID RULE: each of the 8 cells contains EXACTLY ONE complete
>   figure, centered in its own cell, with AT LEAST 15px of clear white
>   between neighbouring figures and to every sheet edge. Never two figures
>   touching or sharing a cell, no detached ghost limbs, no stray specks.
>   ANCHOR LOCK — the waist cutoff sits on the SAME pixel row and the body's
>   vertical centerline on the SAME pixel column in every cell; same size
>   every frame. No pure white on the character (cream #E5DDC8 for fabric).
>   No separators, no extra portrait, no text.
>   ===PROMPT END===

---

### 2026-06-27 playtest bug-sweep — new art prompts

**★ ART STYLE — READ FIRST (this is mandatory for EVERY sheet below):**
Flat **1990s Saturday-morning cartoon** style — the same look as the existing
sheets. **Bold, even black outlines; flat cel-shaded fills (1–2 shades per
colour, no gradients); simple cartoon shapes.** NOT realistic, NOT painterly,
NOT semi-realistic, NOT 3D/render, NOT anime, no soft airbrush, no detailed
rendering. When a reference sheet is named, **open it and copy its exact line
weight, colour palette, proportions and shading** — the new sheet must look
like it came off the SAME sheet, just a new pose. If in doubt, simpler and
flatter is correct.

**Rules reminder:** no pure white on characters; every PP pickup ends with him pocketing the item; PP has plain pink paws and plain pink feet (no gloves, no shoes); separate idle + talk sheets per NPC (8×1 each, 1536×1024, ≥15px gaps between frames).

**PP design reference:** Use `PP idle front.png` or `PP give sketch.png` as the exact visual reference for every PP sheet. Key traits: tall slim pink panther, off-white/cream chest, plain pink paws and feet, long pink tail, small beady dark eyes with pink eyelids, long snout, lanky cartoon proportions. Do NOT add gloves or shoes.

**Done (generated, listed for record):** §RAMEN-WAIT (ramen queue NPCs) · §PP-KIMONO-RUN (`PP_kimono_spin.png` already exists — spin-into-kimono and back).

---

#### §BUST-FRAMING — give/receive sheets render TINY · REGEN (diagnosed 2026-06-28)

Root cause of "the NPC shrinks during the give/receive" (#7 Poulain, #19 Camille,
#32 spice): the engine sizes every animation state so its tallest figure fills the
character's box. The IDLE/TALK sheets for these counter/seated NPCs are **bust /
waist-up** (figure ~44% of the cell), but their give/receive sheets were drawn
**FULL BODY** (figure fills the whole cell). Normalised to the same box, the
full-body give sheet's head/figure comes out tiny. There is no code fix — the
sheets must match the idle's framing.

Regenerate each of these to **EXACTLY match the matching idle sheet's framing** —
same waist-up bust crop, same head size, same scale, same anchor; ONLY the
arms/hands and the held item move. Do NOT draw the legs or full body, and **NO
table/desk/counter in the sprite** (the furniture is in the background art).

- `npc_bakery_woman_receive_rolling_pin.png` → match `npc_madame_poulain_idle.png`
- `cafe_patron_camille_give_sktach.png` (Camille give sketch) → match her seated
  café bust idle (`loadCafePatronGrids "camille"`); no table in the sprite
- `npc_spice_seller_give.png` (if drawn full-body) → match `npc_spice_seller_idle.png`
  All: 8×1, 1536×1024, transparent/keyed bg, no pure white, ≥15px gaps. The held item
  (rolling pin / sketch page / cardamom) may extend past the bust outline naturally.

---

#### §SPICE-SIDE — `npc_spice_seller_talk_side.png` · LANDED 2026-07-14 (wired as preferred talk sheet)

The spice seller is **seated behind his stall** (a low wooden table with small bowls of colourful spices). He currently only has a front-facing idle; generate a **side-facing talk** sheet where he turns his body and head to the RIGHT (toward the market centre), gesturing and speaking.

Exact design from `npc_spice_seller_idle.png`:

- Dark curly hair with a small dark cloth wrap/cap
- Short dark beard/stubble
- Cream/off-white loose shirt, dark earthy vest, apron
- Seated behind the stall table, spice bowls visible in front

8×1, full body including the stall, transparent background, no pure white.

#### §COFFEE-GIVE — `npc_coffee_seller_give_coffee.png` · LANDED 2026-07-14 (wired as preferred give sheet)

The coffee seller pours a small cup from his brass dallah pot and holds it out toward PP (off-screen left).

Exact design from `npc_coffee_seller_idle.png`:

- Red fez hat, dark mustache, medium-dark skin
- Dark vest over cream shirt, dark trousers
- Brass dallah coffee pot (tall, long-spouted) in one hand; small cup in the other

8×1, full body, transparent background, no pure white.

#### §JAKE-RECEIVE — `npc_jake_get_coin.png` · REGEN (wrong prop)

The current sheet shows Jake holding a **paper page / card with a coin printed on it** — WRONG. The anchor item is a **real, solid, round gold coin** (see `assets/images/items/coin.png`: a shiny gold/yellow disc with a sparkle, NOT paper, NOT a rubbing). Regenerate so Jake holds an actual coin.

**Open `npc_jake_idle.png` and `npc_jake_talk.png` and copy the design EXACTLY** — same flat cartoon style, same colours, same chubby build. Do not invent new colours.

Pose: Jake reaches out and takes a **small round gold coin** from PP (off-screen left), pinches it between thumb and finger and holds it up close, staring at it — eyes going wide as the face stamped on the coin matches the one from his nightmares. The coin is a clearly 3D metal disc (gold/yellow with a darker rim), held edge-on / face-on between his fingers — NEVER a flat sheet of paper or a printed card.

Exact design (from the live sheets):

- **Light/tan skin** (he is NOT dark-skinned), round chubby face
- **Very short buzz-cut**, light-brown stubble hair
- **Dark forest-green t-shirt** with **"CAMP" printed in black** across the chest
- **Dark-green / olive shorts** — NOT black
- **Green & white low-top sneakers**
- Stocky / chubby King-of-the-Hill build; arms uncrossing as he reaches out

8×1, full body, transparent background, no pure white.

#### §HIGGINS-GRUMPY — `npc_director_higgins_rude_idle.png` · REGEN (floating clipboard)

The current sheet has his **clipboard/notebook floating in mid-air** near his hip, wedged under the crossed arms but not actually held by a hand. **Remove the clipboard entirely from this pose** — for the rude crossed-arms stance both hands are tucked into the crossed arms, so there is NOTHING in his hands. (His clipboard only appears in the idle/talk/shout sheets where a hand clearly grips it; it must never float.)

**Open `npc_director_higgins_idle.png`, `npc_director_higgins_talk.png` AND `npc_director_higgins_rude_talk.png` and copy the design EXACTLY** — same flat cartoon style, same colours, same proportions.

Pose: Director Higgins standing with arms tightly and cleanly crossed over his chest (both hands visible, tucked at the elbows — empty hands, NO clipboard, NO object), frowning/scowling, weight shifted to one leg. Impatient resting pose — NOT gesturing (that's the talk sheet's job).

Exact design (do not omit any of these — they're all on his other sheets):

- Tall lean older man; **short neat grey hair** (slightly receding), **round glasses**, subtle grey mustache
- **Dark forest-green short-sleeve ranger/park button-up shirt** with a collar, button placket, and **TWO flap chest pockets**
- **Lanyard around the neck with a small green ID badge** on the chest
- **Brown leather crossbody satchel / messenger bag** worn across the chest down to the left hip (the strap is always visible)
- **Khaki/tan CARGO trousers** with **thigh side-pockets** (not plain slacks)
- **Brown hiking boots**
- NO hat, **NO clipboard in this pose**

8×1, 1536×1024, transparent background, no pure white, ≥15px gaps.

#### §LILY-SAD-TALK — `npc_lily_sad_talk.png` · NEEDED

**Open `npc_lily_sad_idle.png` and copy it EXACTLY** — same flat cartoon style, same colours, same from-behind sitting pose. This new sheet is literally that sheet plus a small head-turn; it must look like the next frames of the SAME animation. Also check `npc_lily_idle.png` for her colours.

Pose: Lily seen **from behind**, sitting hugging her knees on the dock. She turns her head ~45° to the side / over her shoulder to acknowledge PP talking to her, then turns back — still hugging her knees throughout. Subtle motion only.

Exact design (from the live sheets):

- Long **medium-brown** hair, middle-parted, flowing down her back
- **Pale mint-green sleeveless dress** (with a coral/peach V-front panel — barely visible from behind, but keep the dress that exact mint green)
- Fair skin; seen from behind, sitting, knees up

8×1, full body from behind, transparent background, no pure white.

#### §TEA-MASTER-IDLE — `npc_tea_master_idle.png` · NEEDED (regen — too realistic)

The previous idle came out **too realistic / over-rendered**. Regenerate it FLAT and CARTOONISH to match the talk sheet.

**Open `npc_tea_master_talk.png` and copy its exact flat-cartoon look** — bold even black outlines, flat cel-shaded fills, simple shapes. The idle must look like it came off the same sheet as the talk. Absolutely **no realistic skin/cloth rendering, no soft shading, no painterly detail**.

Design (an elderly Japanese woman — same character as the talk sheet):

- **Grey hair in a neat bun**
- **Sage / olive-green kimono** with a wide obi sash
- **Kneeling in seiza** (the same kneeling pose as the talk sheet)
- Hands resting calmly in her lap, or holding a tea whisk at rest — NOT gesturing

8×1, 1536×1024, transparent background, no pure white, ≥15px gaps.

#### §CAM-ROOM7 — `npc_camille_sketching_room7.png` · LANDED 2026-07-14 (auto-wired, verified)

Camille (red beret, art student) drawing the Room 7 Mona Lisa restoration sketch — **NOT a portrait of PP**. Show her waist-up, bent over a large sketchpad, charcoal in hand; the visible page shows a woman's face in a golden frame with a subtle hidden symbol (triangle or eye motif). **NO table/desk in the sprite** — she is upper-body only (the café table is in the background art; drawing one would double it up). 8×1 frames showing the drawing progress, ending with Camille holding the finished page toward the camera. 1536×1024, transparent background, no pure white.

#### §PP-RECEIVE-SKETCH — `PP receive sketch.png` · REGEN (extra/disembodied hand)

The current sheet has a **second, disembodied pink hand gripping the rolled sketch** (plus a leftover floating paw) — so PP looks like he has THREE hands. Regenerate with **only PP and his own two arms**.

Hard rules for this regen:

- **Do NOT draw any giving/offering hand.** The off-screen person who hands him the sketch is NEVER drawn — no hand, no arm, nothing entering from off-screen.
- PP uses **only his own two arms** in every frame. He reaches out with ONE open paw; the rolled sketch then simply **appears in that same paw** (he is now holding it) — there must never be a separate hand holding the sketch next to his reaching paw.
- No floating/detached paws anywhere on the sheet.

Sequence: PP reaches out with one paw → the rolled sketch is now in that paw → he glances at it with an approving nod → he slides it into his invisible side pocket (final frame: hand going into pocket, item gone).

**Match the design from `PP give sketch.png` exactly** — same proportions, same flat-cartoon style, mirrored direction (PP receives from the right, so he faces right), plain pink paws, off-white chest, long pink tail. 8×1, 1536×1024, transparent background, no pure white.

#### §PP-RUN-JAPAN — PP run sheets for the Japan kimono gag · NEEDED (4 sheets)

The gag plays in order: (1) normal PP sprints off right → (2) PP returns from right in kimono → (3) **PP stops centre-stage and shows off the kimono for a beat** → (4) PP sprints off right again in kimono (then walks back normal via existing walk sheet).

All four: **match `PP walk left.png` exactly** for proportions, palette and style (pink body, off-white/cream chest, plain pink paws and feet, long pink tail, pink eyelids). 8×1, 1536×1024, transparent background, no pure white.

**`PP_run_right.png`** — PP in normal clothes sprinting **to the right** (facing right). Exaggerated cartoon sprint: leaning forward, big pumping strides, tail streaming behind, arms swinging. Faster and more urgent-looking than the normal walk.

**`PP_run_kimono_left.png`** — PP in the **dusty mauve-purple kimono** (from `PP_kimono_spin.png` — same colour, same obi sash) running back **to the left** (facing left, entering from off-screen right). Because the kimono restricts his legs, the stride is short and frantic — legs shuffling rapidly, hem flapping, arms flailing for balance. Comedic urgency. This is a SHORT entrance — PP runs a few steps in and stops, so keep the run motion compact.

**`PP_kimono_pose.png`** — the MIDDLE "stay with the clothes" beat: PP standing still **facing the camera** in the same mauve-purple kimono, showing it off — a little posing/preening loop (smooths the obi, a small proud turn of the head, hands clasped or arms out presenting the outfit). NOT running, NOT walking — a stationary idle/pose the engine can hold for ~1.5 s between the run-in and the run-off. 8 frames of gentle posing, looping cleanly.

**`PP_run_kimono_right.png`** — same kimono, same frantic shuffle, but now facing and running **to the right** (exiting off-screen again).

---

### 2026-06-30 — §PIERRE-REBUILD — Pierre fully recreated (NO pure white) · NEEDED

Pierre's old sheets were scrapped: his easel CANVAS, his SMOCK, his EYES, and
some edge pixels around his body were pure/near-white, which the engine
chroma-keys into white spots / holes or leaves as a tiny white rim after the
background is removed. Regenerate ALL of his sheets to the design below with
**zero pure white anywhere on Pierre or his props** (per the ★ HARD RULE in this
file's header). Verify each with
`SPRITE_SCAN=1 go test ./engine -run TestSpriteScan` — the five rebuilt Pierre
sheets must not appear in the white-spot report.

**Pierre — canonical design (use on EVERY sheet, no pure white):**

- Cheerful Parisian street painter, slim, ~middle-aged, big curled brown
  **mustache**, warm tan skin.
- **Black beret `#1C1C1C`.**
- **Painter's smock in CREAM `#E5DDC8`** (NOT white) over a blue-and-cream
  striped shirt (`#2B6CB0` stripes on cream `#E5DDC8`, no white) with a small
  **blue neckerchief `#2B6CB0`**.
- Brown trousers `#5A4632`, brown shoes.
- **Eyes: dark pupils with PALE-GREY `#C4C4C4` sclera — never white, not even
  anti-aliased edge pixels.** Teeth (if he smiles) cream `#E5DDC8`.
- Stands at a **wooden easel** (frame `#8B5A2B`) holding a **CREAM `#E5DDC8`
  canvas** (NOT white) and a wooden palette. The canvas may show a few loose
  paint daubs but its base is cream, never white.

Style: the standing 1990s Saturday-morning cartoon look (bold ~3px black ink,
flat cel fills). 8×1 per sheet, 1536×1024, ≥15px gaps between frames, pure-white
BACKGROUND only (the chroma key) — nothing white on Pierre or the easel.

**Critical export / cleanup requirement:** after the white background is keyed
out, Pierre must not have any tiny white or near-white rim pixels left around
his silhouette, eyes, smock, canvas, palette, baguette, jam jar, or press.
No white anti-aliasing on the character or props. Light edge pixels must be
cream/bone/pale-grey that visibly differ from `#FFFFFF`.

**Sheets to generate (all 8×1, same path as the existing files):**

1. `npc_pierre_idle.png` — relaxed idle at the easel: dabs the cream canvas with
   the brush, leans back to consider it, small mustache-twitch / blink. Subtle
   loop, the figure stays planted.
2. `npc_pierre_talk.png` — same pose, talking: mouth moves, free hand gestures
   (points the brush, taps the palette). Same size/anchor as the idle.
3. `npc_pierre_give_pass.png` — PP (off-screen left) hands Pierre a baguette;
   Pierre reaches out, takes it, and in return holds out a **press / paper
   card** (bone `#EDE5D3`, not white) toward PP. Ends with the pass extended.
4. `npc_pierre_get_baguette.png` — Pierre receives a **baguette** from PP
   (off-screen left): reaches, takes the loaf, tucks it under one arm, nods.
5. `npc_pierre_get_jam.png` — Pierre receives a small **confiture / jam jar**
   from PP: reaches, takes the jar, holds it up approvingly. (Re-roll of the old
   ghost-hand sheet — only his two own hands, no detached/ghost limbs.)

Each give/receive sheet: keep Pierre's body/scale/anchor identical to his idle so
the swap mid-dialog doesn't jump; only the arms + the handed prop move.

---

### 2026-06-30 — §PP-TALK-BACK — PP back-facing talk sheet · LANDED 2026-07-14 (auto-wired 8×2; row-2 belly-patch nit, playtest)

Madame Poulain stands behind a back-of-scene bakery counter, ABOVE PP, so PP
talks UP toward her with his BACK to the camera. The code now selects a
back-talk sheet for this (`ppFaceBack`), falling back to PP's back IDLE
(`PP idle back.png`) until this lands — so PP already faces away, he just
doesn't have a talking cycle yet. Generate the talking version.

**File:** `assets/images/player/PP talk back.png`

**Prompt:** Pink Panther seen from BEHIND (back of his head, back, tail), in the
SAME size/stance/anchor as `PP idle back.png` so the swap mid-dialog doesn't
jump. A gentle TALK cycle: small head bob / weight shift / one arm making a
light conversational gesture at his side — readable as "speaking" from behind
(we never see his face). Plain pink paws, NO gloves. Same 1990s Saturday-morning
cartoon look (bold ~3px black ink, flat cel fills) as the other PP sheets.

- **8×2 grid, 1536×1024**, ≥15px gaps between frames so they gap-cut cleanly.
- **NO pure white `#FFFFFF` anywhere on PP** (per the ★ HARD RULE header) — the
  pure-white BACKGROUND is the only white (chroma key). Eyes aren't visible from
  behind, so no white-eye risk, but watch edge anti-aliasing.
- Match the back-idle's vertical baseline so his feet line up with the idle/side
  poses.

Verify: `go test ./engine -run ContentGrid` (clean cut) and
`SPRITE_SCAN=1 go test ./engine -run TestSpriteScan` (no white on PP).

---

## Done log

- **2026-06-30 — interim Pierre repair:** Current `npc_pierre_idle.png`,
  `npc_pierre_talk.png`, `npc_pierre_give_pass.png`, `npc_pierre_get_baguette.png`,
  and `npc_pierre_get_jam.png` were repaired/verified against the no-pure-white
  scan. This is not the final prompt-generated rebuild; §PIERRE-REBUILD remains
  open for the full five-sheet regen.

---

### 2026-06-30 — PR bug-sweep regens (NO pure white — see ★ header rule)

These came out of the Paris/Jerusalem/camp playtest PR. Code is fixed; these are
the art-only follow-ups. All obey the ★ hard rule: zero pure `#FFFFFF` on any
drawn element (eyes, paper, canvas, etc.); only the cell BACKGROUND is white.

#### §PP-PUTNOTE — `PP put note in wall.png` · REGEN (frames overlap)

6×1. The current sheet's 6 poses TOUCH, so the engine can't find a clean split
(it falls back to proportional cutting → mis-cut in-game). Re-draw with **≥15px
of clear background between every pose** and to the sheet edges. PP (side view)
writes on a slip of paper, then reaches up and tucks it into a crack in the
Western Wall; ends lowering his paw, empty-handed. Paper = bone `#EDE5D3`, never
white. 1536×1024.

#### §PP-GIVE-BG — give sheets show white BESIDE PP · REGEN (enclosed white bg)

`PP give flower.png`, `PP give baguette.png`, `PP give postcard.png`,
`PP give sketch.png` (8×1 each). The CUT is fine, but each frame has a patch of
pure-white BACKGROUND trapped beside PP (between his arm/body and the cell edge),
which the engine can't key away (it's enclosed) so it renders as a white block.
Re-draw so PP's silhouette stays OPEN to the surrounding background (no sealed
white pockets), and put the handed item in off-white (daisy/postcard/sketch =
bone `#EDE5D3`), never pure white. Keep ≥15px gaps. 1536×1024.

#### §PP-GET-SKETCH-SIZE — `pp_get_skatch.png` · REGEN (renders tiny)

8×1. PP receives Camille's sketch but renders much smaller than his idle because
PP is drawn larger/fuller in this sheet than in `PP idle side.png`. Re-draw PP at
**exactly his standard idle-side size and headroom** (use `PP idle side.png` as
the size+proportion reference) so he doesn't shrink when the receive plays. He
reaches, takes the rolled sketch, glances at it, pockets it. Sketch page = bone,
not white. 1536×1024.

#### §FENCE-NOWHITE — `jerusalem/props/fence.png` · REGEN (white between bars)

The metal crowd-barrier has pure-white BETWEEN the iron bars (enclosed gaps the
edge key can't reach), so white shows through in-game. Re-draw on a pure-white
background but with the gaps between the bars left as **transparent / true
background** (nothing painted there) and no pure white on the metal itself (use
pale grey `#C4C4C4` for highlights). Single static frame.

#### §WALL-PRAYER-BG — `jerusalem/npc/wall/praying_man.png` (+ `_man2`) · REGEN

The lone praying man at the Wall ships with a baked limestone BACKGROUND that
doesn't match the scene's wall, so a tan rectangle shows around him. Re-export
with a fully **transparent / pure-white background** (no baked stones) and no
pure white on the figure (his prayer shawl = cream `#E5DDC8`, not white). Keep
the 4-frame from-behind sway.

#### §PP-RECEIVE-PRESSPASS — `PP receive press.png` · NEEDED (optional)

8×1. When Pierre hands PP the press, PP currently just idles (the generic
grab fallback was removed). A small dedicated beat: PP (side) takes a flat
press-pass card (bone `#EDE5D3`), glances at it, pockets it. Matches
`PP give postcard.png` framing/size. Until it lands PP simply idles — no rush.
