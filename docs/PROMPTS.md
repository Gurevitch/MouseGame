# Art Prompts

Standardized image-generation prompts for every character, background and
prop in **PP: Mystery on the Mountain**. The engine already strips white
backgrounds and grid lines on load (see `engine/texture.go` ->
`SpriteGridFromPNGClean`) so all prompts assume a white background with a
thin charcoal border between frames.

## Global art direction

- **Style**: hand-drawn 90s Saturday-morning cartoon — think *Pink Panther
  in Hokus Pokus Pink* (1997) and *Passport to Peril* (1996). Flat colors,
  confident black ink linework, slight exaggeration, no cross-hatching, no
  gradients inside the character.
- **Palette**: warm paper-white background (`#FFFFFF`), dark charcoal
  outlines (`#1C1C1C`), saturated but never neon fills. Drop shadows are a
  flat 40% black ellipse under the feet.
- **Aspect**: every character sheet is a **single horizontal row** of
  **8 frames** rendered at **512×512 per frame** (so the full PNG is
  **4096×512**). Each frame is centered inside its cell with a consistent
  ~32 px margin so the grid lines never touch the character.
- **Anatomy**: the character must occupy the **same Y baseline** across
  every frame. The engine trims 3 px from each cell edge on load — never
  rely on the border for visual information.
- **Talk sheets**: same grid, same baseline, but the mouth is **open in
  every frame** with a visible gesture variation (hands, eyes, head tilt).
  Idle sheets keep the mouth closed and use gentle shifts (blinks, weight
  shifts, collar flutter) so the engine can loop them at 0.12s per frame
  without the character looking twitchy.
- **Never** include speech bubbles, captions, background props, or second
  characters inside a cell.

## Framing & Scale (Passport-to-Peril reference)

Reference images live in the workspace assets folder:
`the-pink-panther-passport-to-peril_7-*.png` (London park, wide shot) and
`pub-*.png` (pub, tight shot). Both are the look we are matching.

Target ratios, measured off the PTP frames:

- Wide outdoor shot: characters occupy **25-30%** of the frame height. PP
  is a hair taller than the adult NPCs but never more than **1.2x** of
  an adult. In the park frame PP stands roughly 180 px tall against a
  480 px-tall frame, with the three bystanders at ~150 px.
- Tight indoor / pub shot: the camera crops in, so everyone renders
  **smaller** (~22-25% of frame height). See the pub image where the
  seated foreground NPCs are heads-and-shoulders only.
- **PP is never more than 1.2x** the height of an adult NPC, and **never
  more than 1.25x** the tallest character on screen.

Per-character target heights at our **1400x800** screen in the default
**1.0x** scene (see `CHARACTERS.md` for the canonical table):

- Pink Panther: **235 px** (current, do not grow past this).
- Adult NPCs (Director Higgins, all Paris locals, every other city
  adult): **210-225 px**.
- Kid NPCs (Marcus, Tommy, Jake, Lily, Danny): **195-215 px**.
- Ambient crowd / background figures: **160-190 px**.

Cabin interiors and "tight shot" scenes multiply every character by the
scene's `characterScale` (default **0.85** in rooms, **1.0** outdoors)
so kids in their rooms render at ~165-180 px.

Artist draw rule, non-negotiable:

1. The figure **fills the cell top-to-bottom** with a ~4 px inset and a
   consistent Y baseline across every frame of the sheet.
2. No big empty borders inside the cell. Empty space forces the engine
   to scale up, which brings back the "too big / pixelated" problem.
3. Background stays flat white (`#FFFFFF`). The charcoal 2 px border
   between cells is allowed but the engine erases it on load — do not
   rely on it for visual information.

## How to read each entry

Every entry lists:

```
path/to/sprite.png        8×1  (one row of 8 frames, total 4096×512)
```

followed by the full generation prompt. Any preferred city-specific path
(e.g. `assets/images/locations/jerusalem/npc/npc_eli_idle.png`) will be
picked up automatically by the engine — see `loadNPCGridPath` in
`game/npc.go`. If the preferred path is missing the game falls back to an
existing sprite (documented in the relevant `game/<city>.go` file), so new
art can be dropped in one character at a time without breaking the build.

---

## Core cast

### Pink Panther — walk / idle / talk (`assets/images/player/`)

- `pp_walk_front.png` 8×1 4096×512\
  Pink Panther walking *toward camera*, eight keyframes forming one gait
  cycle: 0) contact right; 1) down right; 2) pass right; 3) high right;
  4) contact left; 5) down left; 6) pass left; 7) high left. Long tail
  trails behind, yellow gloves, relaxed smile. Baseline Y = 440. No shadow
  baked into the frames.
- `pp_walk_back.png` 8×1 — same gait cycle from behind; tail visible.
- `pp_walk_left.png` / `pp_walk_right.png` 8×1 — side walk cycles. Left
  sheet must mirror the right sheet exactly (do not re-draw; flip on
  export). Face clear even in profile.
- `pp_idle.png` 8×1 — subtle breathing loop, arms at side, slow blink on
  frames 3 and 6, tail flicks on 5. Mouth closed.
- `pp_talk.png` 8×1 — half-body gestures with mouth open: 0) neutral;
  1) one paw raised; 2) shrug; 3) point forward; 4) sly grin; 5) both
  paws up; 6) finger to chin; 7) confident smile.
- `pp_sleeping.png` 8×2 (4096×1024) — first row: lying on back under
  blanket, arms at side, tail peeks out. Second row: gentle breathing
  (chest rises/falls), occasional Z symbol fading in/out on frame 5.
- `pp_waking.png` 8×2 — first row: stretch-and-blink wake cycle. Second
  row: sitting up rubbing eyes, then standing. Final frame matches the
  first frame of `pp_idle.png`.
- `pp_airplane.png` 4×3 (2048×1536) — Pink Panther piloting a tiny
  biplane over clouds. Row 0: biplane cruising (4 frames of wing tilt).
  Row 1: biplane banking left. Row 2: biplane banking right. PP always
  visible in the cockpit, scarf flapping.

### Director Higgins (`assets/images/locations/camp/npc/higgins/`)

- `npc_director_higgins_idle.png` 7×1 at 172×384 per cell (total
  1204×384) — tall thin man in khaki uniform, round glasses, clipboard
  in left hand. Canonical pose cycle: 0 neutral; 1 small inhale;
  2 glance at clipboard; 3 push glasses; 4 tap pencil on clipboard;
  5 neutral; 6 small closed-mouth smile. Mouth closed throughout.
  Cell geometry matches `_talk.png` so both animations render at the
  same on-screen size through `npc.drawScaled`.
- `npc_director_higgins_talk.png` 6×1 — same character, clipboard lowered,
  six gesture frames with mouth open.
- `npc_director_higgins_office_idle.png` 7×1 — seated behind a desk, only
  upper body visible. Idle shuffle of papers.
- `npc_director_higgins_office_talk.png` 4×2 — row 0: explaining gestures
  at desk; row 1: pointing at map on desk.
- `npc_director_higgins_give_map.png` 8×1 — standing, unrolling a world
  map across two hands. Final frame: map held up toward camera.
- `npc_director_higgins_walk.png` 8×1 and `walk_back.png` 8×1 — walking
  gait cycles matching PP's style.

---

## The five kids (Camp Chilly Wa Wa)

All five kids share a body template: ~1/3 PP's height, chunky shoes,
expressive eyebrows. Each kid gets **four** sprite sheets:

1. Normal idle (8×1) — mouth closed.
2. Normal talk (8×1) — mouth open, natural gestures.
3. Strange idle (8×1) — same pose but pupils slightly dilated, one hand
   subtly fidgeting, color slightly desaturated.
4. Strange talk (8×1) — mouth open, gestures more erratic, eyes wider.

| Kid     | Visual tag                                   | Palette                                  |
| ------- | -------------------------------------------- | ---------------------------------------- |
| Marcus  | curly black hair, sketchbook under arm       | olive shirt, brown shorts                |
| Jake    | buzz cut, belt with pouches and coins        | red tee, cargo shorts                    |
| Lily    | long blond braid, flower behind one ear      | lavender dress, white sneakers           |
| Tommy   | mullet, drumsticks in back pocket            | black band tee, ripped jeans             |
| Danny   | round glasses, sketchpad, Roman pendant      | blue hoodie, khakis                      |

Sheet paths live at `assets/images/locations/camp/npc/kids/<name>/`:

- `npc_<name>_idle.png` 8×1
- `npc_<name>_talk.png` 8×1 (Tommy / Jake are 7×2 — keep as-is)
- `npc_<name>_strange_idle.png` 4×2
- `npc_<name>_strange_talk.png` 7×2

---

## City chapters

Every city needs:

- 1 to 2 background PNGs at 1400×800 (engine's native resolution). No
  parallax layers are required; the engine draws the image fullscreen.
- 2 to 4 NPC sprite sheets, each 8×1 (idle) and 8×1 (talk) unless noted.
  The preferred path is listed in the relevant `game/<city>.go` file; if
  unwritten, the game falls back to a reused Paris sheet.

### Paris (`assets/images/locations/paris/`)

Already shipped. Existing art:

- `background/paris_street.png` + `paris_museum.png` + `paris_catacombs.png`
- `npc/npc_french_guide_idle.png` 8×2 — *Madame Colette* (keep)
- `npc/npc_french_guide_talk.png` 8×1 — *Madame Colette* talk (keep)
- `npc/npc_museum_curator_idle.png` 8×1 — *Curator Beaumont*
- `npc/npc_museum_curator_talk.png` 4×2 — *Curator Beaumont*
- `npc/npc_art_vendor.png` 8×2 — reused for *Pierre* (street artist)
- `npc/npc_security_guard.png` 6×2 — reused for *Claude* (gendarme)

No new Paris art required; the two "locals" share existing sheets.

### Jerusalem (`assets/images/locations/jerusalem/`)

Backgrounds:

- `background/jerusalem_street.png` 1400×800 — sunlit stone plaza in
  front of the Western Wall at midday. Limestone in foreground, the wall
  receding to the right. Kiosk with juice bottles and postcards on the
  left. No people.
- `background/jerusalem_tunnel.png` 1400×800 — dim tunnel beneath the
  plaza. Low stone ceiling with archeological scaffolding, a single
  flickering work lamp warming the center, shadows pooling at the edges.

NPCs (all 512×512 per frame, 8×1 sheets — add a second row for talk
where noted):

- `npc/npc_eli_idle.png` 8×2 — *Eli* the kiosk owner: short man, wide
  grin, white kippah, blue apron over striped shirt. Row 0: wiping
  counter, checking watch, waving hello, stroking beard, pointing at
  juice, gentle laugh, shrug, neutral (mouth closed). Row 1: same
  gestures with mouth open (talk cycle).
- `npc/npc_miriam_idle.png` 8×2 / `npc_miriam_talk.png` 8×1 — *Miriam*
  the archeologist: young woman, khaki work shirt, braid over one
  shoulder, dust on cheek, pencil behind ear. Idle: brushing dirt off
  wall, chiseling, stepping back to look, sketching in notebook.
- `npc/npc_dov_idle.png` 8×2 / `npc_dov_talk.png` 8×2 — *Dov* the kid
  with a flashlight: small boy, oversized green hoodie, blue ball cap
  backward, a dim flashlight held up. Idle: crouching, peering up,
  flashlight swept side-to-side.
- `npc/npc_gary_tourist.png` 6×2 — *Gary* the recurring tourist: older
  bald man, cargo shorts, sun hat, camera around neck, open guidebook.
  Same character reused in Tokyo, Buenos Aires and Mexico — use this as
  the master and mirror-palette for the others.

### Tokyo (`assets/images/locations/tokyo/`)

Backgrounds:

- `background/tokyo_street.png` — evening street under a red torii gate,
  rows of red paper lanterns, tiny ramen stall on the left, a calligraphy
  shop on the right, sakura petals in the air.
- `background/tokyo_temple.png` — quiet temple garden behind the torii:
  stone path, a koi pond, sakura tree in bloom, pressed-petal awning,
  wooden bench on the right.

NPCs:

- `npc/npc_hiro_idle.png` 8×2 — *Hiro* the ramen cook: middle-aged man,
  white bandana, blue yukata rolled to elbows, ladle in one hand. Idle:
  stirring pot, tasting, wiping brow.
- `npc/npc_kenji_idle.png` 6×2 — *Kenji* the calligraphy student: teen
  boy, gray hoodie under open haori, brush in hand, small table in front.
  Idle: dipping brush, painting a kanji, blowing on ink.
- `npc/npc_obachan_idle.png` 8×2 / `npc_obachan_talk.png` 8×1 —
  *Oba-chan* the elderly flower arranger: small gray-haired woman in a
  purple kimono with cherry-blossom embroidery, reading glasses on a
  chain, arranging a single sakura branch in a ceramic vase.
- `npc/npc_gary_idle.png` 6×2 — Gary again; same character as Jerusalem,
  identical palette.

### Rio de Janeiro (`assets/images/locations/rio/`)

Backgrounds:

- `background/rio_street.png` — sunset street in Copacabana with Christ
  the Redeemer visible on the Corcovado peak, palm trees, a yellow bar
  awning reading "Tio Jorge".
- `background/rio_bar.png` — interior of the bar: rattan chairs, a
  counter with an old radio, photos of Carnival kings on the wall, a
  wooden box of dance cards in the corner.

NPCs:

- `npc/npc_tio_jorge_idle.png` 8×2 — *Tio Jorge*: large cheerful man,
  white stained apron over a football shirt, a towel over one shoulder,
  beard. Idle: polishing a glass, laughing, wiping counter.
- `npc/npc_marisa_idle.png` 8×2 / `npc_marisa_talk.png` 8×1 — *Marisa*:
  mid-40s woman with long dark curls, red dress, tortoise-shell hair
  clip, dance card in her hand.
- `npc/npc_padre_idle.png` 6×2 — *Padre Antonio*: thin older priest in
  black cassock with white collar, reading a small prayer book, looking
  up between sentences.
- `npc/npc_bruno_kid_idle.png` 8×2 / `npc_bruno_kid_talk.png` 8×2 —
  *Bruno*: 10-year-old boy, yellow football tee, barefoot, juggling three
  oranges. Oranges must stay in frame across all 8 idle beats.

### Buenos Aires (`assets/images/locations/ba/`)

Backgrounds:

- `background/buenos_aires_street.png` — dusk on a balcony-lined street
  with jacaranda trees in purple bloom, cobbled pavement, a neon "TANGO"
  sign glowing warm above a doorway.
- `background/ba_tango_school.png` — warm wooden-floored studio: grand
  piano in the back, a mirror wall on the right, a red shawl draped over
  a chair on the left, a chalkboard with footwork diagrams.

NPCs:

- `npc/npc_don_rafa_idle.png` 8×2 — *Don Rafa*: silver-haired man in a
  tailored black suit with red pocket square, cane tapping the floor,
  small smile.
- `npc/npc_lucia_idle.png` 8×2 / `npc_lucia_talk.png` 8×1 — *Lucia*:
  mid-30s woman in an emerald tango dress with a rose at the waist, dark
  hair in a low bun, one hand on hip.
- `npc/npc_paco_idle.png` 6×2 — *Paco*: lanky young man in a white shirt
  and suspenders, shoes polished, sweat on his brow, practicing footwork.
- `npc/npc_gary_idle.png` 6×2 — Gary again (same master).

### Rome (`assets/images/locations/rome/`)

Backgrounds:

- `background/rome_street.png` — cobbled street with Nonna Rosa's pasta
  stall in the middle ground and the Colosseum silhouetted in warm haze
  on the right. Vespa parked against a wall.
- `background/rome_colosseum.png` — inside the arena: weathered stone
  arches, a rubbing tent on the central ground, soft afternoon light.

NPCs:

- `npc/npc_nonna_idle.png` 8×2 / `npc_nonna_talk.png` 8×1 — *Nonna Rosa*:
  sturdy older woman, flour-dusted blue apron, white hair in a bun,
  wooden spoon in hand, warm fierce expression. Idle: stirring pot,
  tasting, wagging spoon.
- `npc/npc_luca_idle.png` 8×2 — *Luca*: young man with long brown curls,
  round sunglasses on top of his head, red accordion. Idle: squeezing
  accordion in and out, tapping foot.
- `npc/npc_dottor_idle.png` 6×2 — *Dottor Bianchi*: elderly classicist,
  tweed jacket with leather elbow patches, white chalk on his fingers,
  small round glasses. Idle: making rubbing motion on paper against
  invisible stone.
- `npc/npc_garibaldi_cat_idle.png` 8×2 — *Garibaldi*: orange tabby cat,
  too dignified for his own good, sits on a crumbling step. Idle:
  blinking slow, flicking tail, one-ear twitch, pawing at a fly.

### Mexico City (`assets/images/locations/mexico/`)

Backgrounds:

- `background/mexico_plaza.png` — afternoon light on the Zocalo plaza:
  cathedral tower in the distance, papel picado bunting overhead, fruit
  carts, mariachis on the far side. The plaza feels celebratory.

NPCs:

- `npc/npc_mariachi_idle.png` 8×2 — *Mariachi*: man in full charro suit,
  wide sombrero, black guitar with silver trim. Idle: strumming, tipping
  hat, stepping to rhythm.
- `npc/npc_abuela_idle.png` 8×2 / `npc_abuela_talk.png` 8×1 — *Abuela*:
  elderly woman in a floral rebozo, silver hair braided with ribbons,
  wooden rosary around wrist. Idle: knitting, looking up, nodding.
- `npc/npc_vendor_idle.png` 6×2 — *Vendor*: young woman with a stack of
  elotes (corn on the cob) in a cart, yellow polo, red apron, cheerful
  wave.

---

## Props and items (`assets/images/items/`)

All item icons are **256×256 PNG**, centered, with a subtle drop shadow
inside the frame. White background (engine strips on load).

| ID                  | Prompt summary                                                                 |
| ------------------- | ------------------------------------------------------------------------------ |
| `travel_map`        | Old folded map with a red thumbtack on Paris. Paper creases, slight tear top.  |
| `flower`            | Single daisy with a water droplet on one petal. Cartoon, not photorealistic.   |
| `postcard`          | Vintage "Bonjour de Paris" postcard with a painting of a woman's face.         |
| `coin_rubbing`      | A yellowed paper rubbing of a Roman coin. Emperor's head is faint but visible. |
| `pressed_sakura`    | A single cherry-blossom petal pressed inside two sheets of rice paper.         |
| `dance_card`        | Two torn halves of a carnival dance card pinned side-by-side, "Marisa" ink.    |
| `inscription_rubbing` | A chalk-on-paper rubbing of Latin letters ("DANILLVS MARCVS") on stone.        |
| `magnifying_glass`  | Brass magnifier, glass dome catches a gleam.                                   |
| `comic_book`        | Retro-80s comic cover, title in Zapf Chancery: "SPACE GHOSTS!"                 |
| `letter`            | Folded letter with a wax seal and "PP" stamped into it.                        |
| `fishing_rod`       | Bamboo rod with a red-and-white bobber on a string.                            |

---

## Ambient Wildlife & Sky

Day 1 camp scenes need life moving in the background. The engine already
has a `particle` struct with `bird`, `insect` and `cloud` flags
(`game/scene.go`); these sheets upgrade the 3-pixel placeholder dots
into proper animated sprites. All three are **camp Day 1 only** — they
vanish on Day 2 and during night / interior scenes so the shift in tone
lands.

### `assets/images/ambient/bird_silhouette.png` — 1x8 384x32

- Single row of **8 frames**, each cell **48x32**, total image
  **384x32**. White background.
- A small bird silhouette mid-flight, rendered as a flat dark shape
  (`#1C1C1C` ink). Wing flap cycle across the 8 frames:
  0) wings up-up, 1) up-mid, 2) level, 3) down-mid, 4) down-down,
  5) down-mid, 6) level, 7) up-mid. No feet, no beak detail; the
  classic distant-bird "m" silhouette in 8 positions.
- Same Y baseline every frame. Figure occupies ~60% of cell height;
  surrounding ~20% top and ~20% bottom stays flat white so the trim
  does not clip the wingtips.

### `assets/images/ambient/butterfly_flutter.png` — 1x6 240x40

- Single row of **6 frames**, each cell **40x40**, total image
  **240x40**. White background.
- A monarch-style butterfly, orange body `#E8882B` with black outlines,
  wings opening and closing:
  0) wings closed, 1) ~30% open, 2) ~70% open, 3) wings fully spread,
  4) ~70% open, 5) ~30% open. Loop returns to 0.
- Tiny antennae and legs visible. Wings render symmetrically.
- Figure fills the cell centered; ~4 px inset top/bottom/left/right.

### `assets/images/ambient/cloud_puff.png` — single image 256x128

- Not a grid. Single soft cartoon cloud centered on a 256x128 canvas.
- Body color **off-white** `#FAFAFA` with subtle gray shadow
  underneath `#D8D8D8` — not pure `#FFFFFF`, because the engine's
  color-key erases pure white; the cloud must survive the trim.
- Shape: four or five bumps on top, flat bottom, classic Pink-Panther
  -era cartoon cloud. No outline.
- Transparent edges (the engine uploads as NRGBA with the color-key
  applied). Drifts are done by code via x-velocity, so no animation
  frames are needed here.

## Checklist for the artist

When a new sheet lands in `assets/images/locations/<city>/npc/` the
engine picks it up automatically — no code changes required, just drop
the file at the path listed in this document. If anything about the
grid size differs (for example your talk sheet ended up 7×1 instead of
8×1) please update the matching `loadNPCGridPath` call in
`game/<city>.go`.

Smoke test after adding each sheet:

```powershell
go build ./... ; go run ./cmd/pp
```

Walk to the character in-game and click — if the animation plays
smoothly and nothing floats or twitches, the sheet is good.

---

## Outstanding sprite regens (2026-04-16)

### `assets/images/locations/camp/npc/higgins/npc_director_higgins_idle.png`

Current sheet authors row 1 at a different frame count than row 0, so
interpreting it as a uniform grid produced a doubled-Higgins flicker.
We now only read row 0, but that leaves the sheet half wasted. Please
regen clean:

- 8 frames single horizontal row, 256x256 per cell, 2048x256 total.
- Transparent background, idle breathing with clipboard at chest (match
  `npc_director_higgins_talk.png` silhouette and palette exactly so the
  idle -> talk swap is seamless).
- Save in place at the same path, overwriting the current file.

After dropping the new file in, the `loadNPCGridRow(..., 8, 2, 0)` call
in `game/npc.go` still works — row 1 being empty is fine.

