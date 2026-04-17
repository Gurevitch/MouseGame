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
