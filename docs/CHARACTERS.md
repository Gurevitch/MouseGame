# Characters & Scale

Single shared reference for **artists** (sheet content, proportions,
palette) and **coders** (draw sizes, per-scene scale). Pairs with:

- `PROMPTS.md` - how each sheet should be generated.
- `game/scene.go` / `game/npc.go` / `game/player.go` - how sheets are
  placed and drawn.

When art or layout changes, update this file first, then the code and
prompts follow.

---

## Canonical draw sizes (1.0x scene)

All heights are what the character should render at in pixels on a
**1400x800** screen when the scene's `characterScale` is **1.0**. See
`PROMPTS.md` -> "Framing & Scale" for the ratios this table descends
from.

### Pink Panther

- Draw size: **170 x 235**
- Palette: body `#E88BB5`, outline `#1C1C1C`, eyes `#FFFFFF` sclera with
  `#1C1C1C` pupil, gloves `#F5D547`, shoes `#E6B429`.
- Anchor: foot-center (bottom of the 235 px box sits on the walk path).
- Silhouette: tall slim pink panther, long tail, small eyes, relaxed
  slouch. Never draw him bigger than **1.25x** the tallest NPC in the
  scene.

### Director Higgins

- Draw size: **160 x 225** on `camp_entrance` and `camp_night`.
- Draw size: **160 x 225** on `camp_office` (was 180x280 — too tall;
  scene scale 0.9 brings the final render down a touch further).
- Palette: safari hat `#8B5A2B`, shirt `#A47148`, shorts `#C79A5A`,
  belt `#5A3A1E`, mustache `#2B1B10`.
- Props: clipboard in talk/shout frames; sits behind desk in
  `_office_*` sheets.

### Marcus (know-it-all)

- Draw size: **170 x 215**
- Palette: messy brown hair `#4A2E1B`, red bandage strip `#C4412A` on
  nose, tan shirt `#D8B47A`, green shorts `#4F7A3A`, sketchbook clamped
  under one arm.
- Strange state: same outline + palette, but eyes are wide, extra lines
  on sketchbook pages, loose strands of hair.

### Jake (tough kid)

- Draw size: **170 x 210**
- Palette: red cap `#C4412A`, freckles `#B06A3A`, navy t-shirt
  `#2B3A5C`, jeans `#4A5A78`, arms crossed or hands in pockets.
- Strange state: gripping a coin or rubbing something metallic; eyes
  fixed on an off-camera horizon.

### Lily (shy girl)

- Draw size: **160 x 195**
- Palette: lavender dress `#B89AC4`, braids `#5A3A1E`, daisy in hair
  `#FFFFFF` with `#F5D547` center, white socks.
- Receives-flower pose: hands cupped at chest height to receive the
  daisy. Used in the give-item animation.
- Strange state: petals scattered at her feet, distant stare, same
  dress/hair.

### Tommy (storyteller)

- Draw size: **170 x 200**
- Palette: small backpack `#3A5A7A`, blue t-shirt `#4A7AA8`, brown
  shorts `#8B5A2B`, floppy hair `#D8B47A`.
- Strange state: hand cupped to ear (listening), mouth agape.

### Danny (prankster)

- Draw size: **170 x 210**
- Palette: green cap `#4F7A3A` (backwards), yellow t-shirt `#F5D547`,
  red shorts `#C4412A`, slingshot in back pocket silhouette, mischievous
  grin.
- Strange state: dirt smudge across cheek, trembling hands, eyes look
  offscreen.

### Paris locals (`paris_street`, `paris_louvre`)

- Pierre the Artist: **160 x 215**. Beret `#1C1C1C`, striped shirt
  `#FFFFFF` + `#1C1C1C` stripes, easel in talk frames.
- Gendarme Claude: **160 x 225**. Navy uniform `#2B3A5C`, red armband
  `#C4412A`, kepi hat, whistle.
- French Guide: **140 x 220**. Blazer `#2B3A5C`, map in hand, glasses.
- Museum Curator: **130 x 220**. White shirt `#FFFFFF`, tweed jacket
  `#8B5A2B`, reading glasses.

### City locals (Jerusalem / Tokyo / Rio + BA / Rome / Mexico)

Every other adult local renders in the **210-225** band. Specific
per-character palette & silhouette tags are already in `PROMPTS.md` so
they are not duplicated here; what matters is the draw size is
consistent.

- Jerusalem: Eli **215**, Miriam **215**, Dov **220**, Gary (tourist) **215**.
- Tokyo: Hiro **210**, Kenji **215**, Oba-chan **200**, Gary **215**.
- Rio: Tio Jorge **220**, Marisa **215**, Padre Antonio **225**,
  Bruno **215**.
- Buenos Aires: Don Rafa **220**, Lucia **210**, Paco **215**,
  Gary **215**.
- Rome: Nonna Rosa **205**, Luca **215**, Dottor Bianchi **225**,
  Garibaldi (cat) **120**.
- Mexico City: Mariachi **215**, Abuela **200**, Vendor **215**.

---

## Per-scene `characterScale`

`characterScale` is a float on the `scene` struct that multiplies every
character's draw W/H at render time. **Hitboxes stay at authored size**
so click targets do not shrink with the scale. Props like the campfire
and the airplane are *not* characters and ignore this value.

Starter values:

- `camp_entrance`: **1.0** — wide park shot.
- `camp_grounds`: **1.0**
- `camp_lake`: **1.0**
- `camp_night`: **1.0** — the fire is the focal point, characters stay
  full-size.
- `camp_office`: **0.9** — tighter room, Higgins sits closer to camera.
- `tommy_room`, `jake_room`, `lily_room`, `marcus_room`,
  `danny_room`: **0.85** — cabin interior, PTP pub-style tight shot.
- `airplane_flight`: **1.0** (airplane itself is a prop and does not
  use the scale; PP is drawn inside the plane texture, not as a
  character sprite, so this is a no-op in practice).
- Paris `paris_street`, `paris_louvre`: **1.0**
- Jerusalem / Tokyo / Rio / BA / Rome / Mexico `*_street`: **1.0**
- Jerusalem / Tokyo / Rio / BA / Rome `*_interior` /
  `*_temple` / `*_bar` / `*_tango_school` / `*_colosseum` /
  `*_tunnel`: **0.9**

---

## Hard rules

1. **No standing NPC > 230 px** in a 1.0 scene. Beyond that the
   character dominates the frame and the PTP look breaks.
2. **No character > 1.25x PP's current draw height** (so no one taller
   than ~294 px including Higgins). Higgins specifically was drawn at
   280 tall — that fails this rule; bring him to 225 and let scene
   scale do the rest.
3. **At least 25% of screen height remains above the tallest character's
   head** so there is sky / ceiling room for the background to breathe.
4. When a new city gets drawn, artists target **215 px** as the default
   adult height. Deviations need a specific in-world reason (child
   character, seated character, priest standing on a step, etc).
5. Pair every new sheet with an entry in `PROMPTS.md` and update this
   file if the size is not in the 195-225 band.

---

## Per-NPC info entries

Single source of truth for every NPC's sprite paths, atlas yaml, grid
dimensions, factory bounds, scenes they appear in, and dialog handles.
When an NPC is regenerated, update both the visual identity above AND
the entry below so coders + artists stay in sync. Cross-link to
`EXTRA_PROMPTS.md` for the regen prompt instead of duplicating it here.

### Director Higgins

- **Role:** Camp director; gives PP the appointment letter, the travel
  map, and the various Day-2+ briefings.
- **Instances** (3 separate factories, distinct bounds + behavior):
  - **Entrance Higgins** — `newDirectorHiggins` (game/npc.go:225).
    bounds `(660, 345) 200×265`. dialog: `higginsDefaultDialog` →
    `higginsPostDialog`.
  - **Office Higgins** — `newOfficeHiggins` (game/npc.go:245). bounds
    `(1062, 357) 220×280`. dialog: `higginsWorriedDialog` →
    `higginsPostWorriedDialog`. silent until Day 2 starts. Owns the
    `give_map` one-shot animation registered on `oneShotAnims["give_map"]`.
  - **Grounds Higgins (hidden)** — `newGroundsHiggins` (game/npc.go:270).
    bounds `(1060, 570) 180×200`. Hidden + silent at scene load; revealed
    by the `higgins_walk_in` sequence after Lily's shy beat.
  - **Night Higgins** — `newNightHiggins` (game/npc.go:327). bounds
    `(1120, 430) 200×260`. silent; driven by `night_bedtime` sequence.
- **Sprite paths:** `assets/images/locations/camp/npc/higgins/`
  - `npc_director_higgins_idle.png` — clean 7×1 strip.
  - `npc_director_higgins_talk.png` — 8×2 (entrance/night load row 0
    only via `loadNPCGridRow`).
  - `npc_director_higgins_office_idle.png` — 6×2 (office loads row 0).
  - `npc_director_higgins_office_talk.png` — 6×2.
  - `npc_director_higgins_give_map.png` — 8×1 one-shot.
- **Atlas yaml:** `tools/characters/higgins.yaml` → packed at
  `assets/sprites/higgins.(png|json)`.
- **Animation speeds:** `talkFrameSpeed: 0.25`.
- **Regen prompt:** see `EXTRA_PROMPTS.md` §1, §2, §10, §18, §19.

### Marcus (know-it-all)

- **Role:** First "afflicted" kid PP heals (Postcard from the Louvre).
- **Instances:**
  - **Grounds Marcus** — `newMarcus` factory (camp_grounds NPC).
    bounds `(890, 395) 150×180`.
  - **Room Marcus** — `newRoomMarcus` (game/npc.go:311). bounds
    `(600, 260) 200×300`. Scene `marcus_room` uses `characterScale: 0.85`.
- **Sprite paths:** `assets/images/locations/camp/npc/kids/marcus/`
  - `npc_marcus_idle.png`, `npc_marcus_talk.png`,
    `npc_marcus_strange_idle.png`, `npc_marcus_strange_talk.png`,
    `npc_marcus_strange_alt.png` (inactivity-trigger pose, asset only —
    code hook still pending).
- **Atlas yaml:** `tools/characters/marcus.yaml` → packed at
  `assets/sprites/marcus.(png|json)`. Grids all 8×2.
- **Animation speeds:** `talkFrameSpeed: 0.10`,
  `strangeTalkFrameSpeed: 0.16`.
- **Dialog handles:** `marcusDefaultDialog` / `marcusPostDialog` /
  `marcusStrangeDialog` / `marcusPostStrangeDialog`.
- **Regen prompt:** see `EXTRA_PROMPTS.md` §13.

### Tommy (storyteller)

- **Role:** Second-tier camp kid; intro chatter only on Day 1, strange
  dialog gated on later chapter.
- **Sprite paths:** `assets/images/locations/camp/npc/kids/tommy/`
  - `npc_tommy_idle.png` / `_talk.png` / `_strange_idle.png` /
    `_strange_talk.png`.
- **Atlas yaml:** `tools/characters/tommy.yaml` → `assets/sprites/tommy.(png|json)`.
  Grids 8×2.
- **Bounds (camp_grounds):** `150×180`. Room version `170×260` with
  scene scale 0.85.
- **Animation speeds:** `talkFrameSpeed: 0.10`.
- **Dialog handles:** `tommyDialog` / `tommyPostDialog` /
  `tommyStrangeDialog` / `tommyPostStrangeDialog`.
- **Regen prompt:** `EXTRA_PROMPTS.md` §11.

### Jake (tough kid)

- **Role:** Cabin-bound kid healed by Coin Rubbing from Jerusalem.
- **Sprite paths:** `assets/images/locations/camp/npc/kids/jake/` (idle,
  talk, strange_idle, strange_talk).
- **Atlas yaml:** `tools/characters/jake.yaml`. Grids 8×2.
- **Bounds:** camp_grounds `150×180`; jake_room `170×260` with scale
  0.85.
- **Dialog handles:** `jakeDialog` / `jakePostDialog` /
  `jakeStrangeDialog` / `jakePostStrangeDialog`.
- **Regen prompt:** `EXTRA_PROMPTS.md` §12.

### Lily (shy girl)

- **Role:** Day 1 flower-handoff puzzle; Tokyo chapter target later.
- **Sprite paths:** `assets/images/locations/camp/npc/kids/lily/`.
- **Atlas yaml:** `tools/characters/lily.yaml`.
- **Bounds:** camp_grounds `150×180`; lily_room `170×260` with scale
  0.85.
- **Special state:** `hintState` field — `0` = unspoken, `1` = shy beat
  done (alt-dialog armed), `2` = flower handed over (post-dialog locked).
- **Dialog handles:** `lilyShyDialog` (initial) / `lilyFlowerDialog`
  (alt) / `lilyDialog` / `lilyPostDialog` / `lilyStrangeDialog` /
  `lilyPostStrangeDialog`.

### Danny (prankster)

- **Role:** Camp prankster; later afflicted, healed by Mexico-chapter
  item.
- **Sprite paths:** `assets/images/locations/camp/npc/kids/danny/`.
- **Atlas yaml:** `tools/characters/danny.yaml`.
- **Bounds:** camp_grounds `150×180`; danny_room `170×260` with scale
  0.85.
- **Dialog handles:** `dannyDialog` / `dannyPostDialog` /
  `dannyStrangeDialog` / `dannyPostStrangeDialog`.

### Madame Colette (French Guide)

- **Role:** First Paris NPC; flavor-dump tour guide that hints at the
  press-pass route.
- **Sprite paths:** `assets/images/locations/paris/npc/`
  - `npc_french_guide_idle.png` 8×2 + `npc_french_guide_talk.png` 8×1.
- **Atlas yaml:** `tools/characters/paris/french_guide.yaml` → packed at
  `assets/sprites/paris/french_guide.(png|json)`.
- **Factory:** `newFrenchGuide` (game/npc.go:797). bounds `(300, 440)
  140×240`.
- **Dialog handles:** `frenchGuideDialog` / `frenchGuidePostDialog`.

### Madame Poulain (Bakery Woman)

- **Role:** Hands over the baguette after PP returns her lost rolling
  pin (Paris bakery interior — see Chunk 7 in the most recent plan).
- **Sprite paths:** `assets/images/locations/paris/npc/npc_bakery_woman.png`
  (8×2).
- **Atlas yaml:** `tools/characters/paris/bakery_woman.yaml`.
- **Factory:** `newBakeryWoman` (game/npc.go:732). bounds `(540, 440)
  140×240`. Scene: `paris_bakery` (NOT paris_street since the rework).
- **Dialog handles:** `bakeryWomanLostPinDialog` (initial) /
  `bakeryWomanPinTradeDialog` (alt — fires when PP holds Rolling Pin)
  / `bakeryWomanPostDialog`.
- **Regen prompt:** `EXTRA_PROMPTS.md` §8.

### Pierre (Street Artist)

- **Role:** Trades baguette for press pass.
- **Sprite paths:** `assets/images/locations/paris/npc/npc_art_vendor.png`
  (8×2 — row 0 idle, row 1 talk).
- **Atlas yaml:** `tools/characters/paris/pierre_artist.yaml`.
- **Factory:** `newPierreArtist` (game/npc.go:874). bounds `(880, 440)
  130×240`.
- **Dialog handles:** `pierreArtistDialog` / `pierreArtistPostDialog`
  + altDialogFunc gated on `Baguette` in inventory.

### Nicolas (Press Photographer)

- **Role:** Flavor NPC near the Louvre steps; breadcrumbs PP toward the
  bakery.
- **Sprite paths:** `assets/images/locations/paris/npc/npc_press_photographer.png`
  (8×2).
- **Atlas yaml:** `tools/characters/paris/press_photographer.yaml`.
- **Factory:** `newPressPhotographer` (game/npc.go:772). bounds
  `(1010, 440) 110×240`.
- **Dialog handles:** `pressPhotographerDialog` /
  `pressPhotographerPostDialog`.
- **Regen prompt:** `EXTRA_PROMPTS.md` §9.

### Claude (Gendarme)

- **Role:** Louvre door blocker; trades museum ticket for press pass.
- **Sprite paths:** `assets/images/locations/paris/npc/` (gendarme art).
- **Atlas yaml:** `tools/characters/paris/gendarme_claude.yaml`.
- **Factory:** `newGendarmeClaude` (game/npc.go:910). bounds
  `(1120, 430) 120×250`.
- **Dialog handles:** `gendarmeDialog` / `gendarmePostDialog` +
  altDialogFunc gated on `Press Pass` in inventory.

### Curator Beaumont (Museum Curator)

- **Role:** Hands over the postcard PP brings back to Marcus.
- **Sprite paths:** `npc_museum_curator_idle.png` 8×1,
  `npc_museum_curator_talk.png` 4×2.
- **Atlas yaml:** `tools/characters/paris/museum_curator.yaml`.
- **Factory:** `newMuseumCurator` (game/npc.go:838). bounds
  `(500, 320) 130×250`. Scene: `paris_louvre` (interior, scale 0.9).
- **Dialog handles:** `museumCuratorDialog` / `museumCuratorPostDialog`.
- **Regen prompt:** `EXTRA_PROMPTS.md` §15.

---

When you regen ANY of the sheets above, update its row's "Sprite paths"
or "Grid" line in the same commit. Stale entries here have caused the
"frames cut in the middle" bugs more than once.
