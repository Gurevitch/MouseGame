# Extra Sprite Prompts — everything still needed for the current FIXME sweep

Paste any prompt below into your image generator. After generation, save the
PNG at the path shown and run the usual atlas re-pack (or just restart the
game for legacy loaders):

```
python tools/pack_atlas.py tools/characters/<name>.yaml
```

All prompts inherit the style lock from `HIGGINS_PROMPTS.md`:

> Hand-drawn 1990s Saturday-morning cartoon, Pink Panther *Hokus Pokus Pink*
> (1997) / *Passport to Peril* (1996). Confident black ink linework ~3 px,
> flat saturated fills, no cross-hatching, no gradients, no airbrush. Two
> cel tones max per color region. Pure #FFFFFF background, zero scenery.
> Every cell is **tall rectangular**, never square.

Canvas dimensions are locked per sheet; do **not** scale down to square.

---

## Character sheet regen status — updated 2026-04-19 (campaign complete)

Full triage pass across every PNG referenced by the camp + Paris casts. Rows
marked **keep** are the canonical Hokus Pokus Pink cartoon references (same
look as Danny / Lily / PP idle_front). Rows marked **DONE** were regenerated
this session from the prompt in the § column. Rows marked **deferred** are
queued but out of scope for this sweep (currently only the two PP sheets).

### Camp cast

| Sheet | Status | § | Notes |
| --- | --- | --- | --- |
| `npc_director_higgins_idle.png` | keep | — | Canonical — 7×1, clean regen. |
| `npc_director_higgins_talk.png` | DONE | §10 | Regenerated to match new idle (silver hair, red lanyard, khaki). |
| `npc_director_higgins_walk.png` | keep | — | Matches idle style. |
| `npc_director_higgins_walk_back.png` | keep | — | Matches idle style. |
| `npc_director_higgins_shout.png` | keep | — | Matches idle style. |
| `npc_director_higgins_office_idle.png` | DONE | §18 | Pixel → cartoon. |
| `npc_director_higgins_office_talk.png` | DONE | §18 | Pixel → cartoon. |
| `npc_director_higgins_give_map.png` | DONE | §19 | Pixel-leaning → cartoon. |
| `npc_director_higgins_desk.png` | keep | — | Matches idle style. |
| Tommy `idle` / `talk` | DONE | §11 | Regenerated to match canonical strange_*. |
| Tommy `strange_idle` / `strange_talk` | keep | — | Already matches Danny/Lily canvas. |
| Jake `idle` / `talk` | DONE | §12 | Pixel → cartoon. |
| Jake `strange_idle` / `strange_talk` | DONE | §12 | Pixel → cartoon. |
| Marcus `idle` / `talk` | DONE | §13 | Regenerated against the canonical strange_idle identity (yellow polo, spiky brown hair, khaki cargo shorts). |
| Marcus `strange_idle` | keep | — | Cartoon. (NOTE: atlas shows a faint checkerboard artifact because the source PNG is RGB with a baked-in bg pattern — pre-existing, does not affect gameplay.) |
| Marcus `strange_talk` | DONE | §13 | Pixel → cartoon; matches canonical Marcus identity. |
| Marcus `strange_alt` | DONE | §4 | Landed; yaml has a 5th anim entry. Inactivity-trigger code hook still deferred. |
| Danny `idle` / `talk` / `strange_*` | keep | — | Canonical style. |
| Lily `idle` / `talk` / `strange_*` / `receive_flower` | keep | — | Canonical style. |
| `campfire_idle.png` | keep (size-only) | — | Existing art fine; §6 now authors the smaller companion loop. |
| `campfire_small.png` | DONE | §6 | New 1032×172 / 6×1 flame loop sized to drop into the (581,592)-(702,594) band. |

### Paris cast

| Sheet | Status | § | Notes |
| --- | --- | --- | --- |
| `npc_french_guide_idle.png` | DONE | §14 | Pixel → cartoon. |
| `npc_french_guide_talk.png` | DONE | §14 | Pixel → cartoon. |
| `npc_museum_curator_idle.png` | DONE | §15 | Pixel → cartoon; §15 canvas/cell dims updated to the actual 1376×768 / 8×1 / cell 172×768 that the generator produces (matches `loadNPCGrid(..., 8, 1)`). |
| `npc_museum_curator_talk.png` | DONE | §15 | Pixel → cartoon; §15 canvas updated to 1376×768 / 4×2 / cell 344×384 (matches `loadNPCGrid(..., 4, 2)`). |
| `npc_art_vendor.png` (Pierre) | keep | — | Canonical cartoon, 8×2 with Pierre talk row. |
| `npc_security_guard.png` (Claude) | keep | — | Canonical cartoon, 6×2. |
| `npc_mystery_figure.png` | keep | — | Canonical cartoon, mood-appropriate. |
| `npc_suspicious_dealer.png` | keep | — | Canonical cartoon, 8×2. |
| `npc_bakery_woman.png` | DONE | §8 | New sheet + `newBakeryWoman` switched from french_guide fallback to the real sheet (`loadNPCGridRow(..., 8, 2, 0/1)`). |
| `npc_press_photographer.png` | DONE | §9 | New sheet + `newPressPhotographer` factory + `press_photographer` NPC id added to `paris_street.json`. |
| `paris/ambient/cafe_patrons.png` | DONE | §7 | New folder + sheet (1376×768 / 4×2 / 8 distinct patrons). Paris ambient renderer hookup still deferred. |

### Player (Pink Panther)

| Sheet | Status | § | Notes |
| --- | --- | --- | --- |
| `PP idle front.png` / `PP talk front.png` | keep | — | Canonical reference. |
| `PP idle side.png` / `PP idle back.png` | keep | — | Canonical reference. |
| `PP walk front.png` | keep | — | Minor slicing artifact; art is good. |
| `PP walk back.png` | deferred | §3 | Needs clearer walk cycle (queued for next PP sweep). |
| `PP walk left.png` | keep | — | Slightly jaggy but cartoon-consistent. |
| `PP talk side.png` | keep | — | Canonical reference. |
| `PP grab.png` | keep | — | Canonical reference. |
| `PP grab flower.png` | deferred | §5 | Queued for next PP sweep. |
| `PP receive map.png` | keep | — | Canonical reference. |
| `PP celebrate.png` | keep | — | Canonical reference. |
| `PP sneak examine.png` / `PP sneak use.png` | keep | — | Canonical reference. |
| `pp_sleeping.png` / `pp_waking.png` / `pp_airplane.png` | keep | — | Canonical reference. |

**Takeaway:** campaign complete for all camp + Paris NPCs. Remaining open
items are the two deferred PP sheets (§3, §5) and two code-hook follow-ups
(Marcus inactivity swap to `strange_alt`, Paris cafe ambient renderer). See
`docs/FIXME.md` for the hookup tracking.

---

## 1. Higgins — entrance idle redesign (match `talk.png`)  *(DONE — canonical)*

**Canvas:** 1204×384. **Grid:** 7×1. **Cell:** 172×384.
**Path:** `assets/images/locations/camp/npc/higgins/npc_director_higgins_idle.png`.

Landed — the regenerated idle became the canonical Hokus Pokus Pink reference
sheet. Every Higgins regen below (§10, §18, §19) must match this file's
linework, face, mustache, and uniform palette.

> [style lock]
>
> Same Higgins as `npc_director_higgins_talk.png`: lanky ranger, round
> wire-rim glasses, small brown mustache, short side-parted brown hair, full
> khaki uniform tucked into khaki trousers, brown leather belt, brown
> lace-up ankle boots. Left hand holds a wooden clipboard. Facing camera.
>
> **Animation:** 7-frame idle loop, mouth closed throughout. Tiny motions
> only — no walking, no pose swap.
>
> - Frame 1: neutral, weight even, eyes forward.
> - Frame 2: inhale; shoulders up 3 px.
> - Frame 3: exhale; shoulders down 3 px, head tilts 2° left.
> - Frame 4: glance down at clipboard (eyes only; head stays forward).
> - Frame 5: pushes glasses up nose with right hand.
> - Frame 6: taps clipboard edge with right index finger.
> - Frame 7: slow blink (half-lidded eyes), silhouette matches frame 1 so
>   loop is seamless.
>
> Silhouette within 3 px of frame 1 in every frame. Baseline (soles of
> boots) locked across all 7 frames.

---

## 2. Higgins — walk back (used by the Lily walk-in sequence)

**Canvas:** 1376×768. **Grid:** 8×1. **Cell:** 172×768.
**Path:** `assets/images/locations/camp/npc/higgins/npc_director_higgins_walk_back.png`.

Currently missing. The `higgins_walk_in.json` sequence lerps his position
but keeps him on his `talk` animation — jarring. With a walk_back cycle we
can swap it in during the move step.

> [style lock]
>
> Same character identity as Higgins talk. **Viewed from behind** (he's
> walking away from camera toward the cabin path at left of frame).
> Clipboard visible at the left hip on some frames.
>
> **Animation:** Standard 8-frame gait cycle, each frame portrait-oriented:
> contact → down → pass → high, repeated left/right. Feet on the same
> baseline on their contact frames. Natural arm swing mirrored behind him.

---

## 3. PP — walk back (for "leaving camp" transition)

**Canvas:** 1376×768. **Grid:** 8×1. **Cell:** 172×768.
**Path:** `assets/images/player/PP walk back.png` (overwrites existing
placeholder).

User wants the "Enter Camp" hotspot to play PP walking away (getting
smaller) for a moment before the scene transitions. The current sheet exists
but reads static — regen with a clear walk cycle.

> [style lock]
>
> Pink Panther from behind, long cat-like silhouette, curved tail trailing
> just off the ground, yellow gloves, relaxed stance. **8-frame walk cycle
> away from camera**: contact right, down right, pass right, high right,
> contact left, down left, pass left, high left.
>
> Subtle hip sway — this is the casual saunter Pink Panther is known for,
> not a march. Baseline Y locked across contact frames (frames 1 and 5).
> Tail sways left-right with the opposite leg.

---

## 4. Marcus — secondary freakout (alt strange idle)

**Canvas:** 1200×180. **Grid:** 8×1. **Cell:** 150×180.
**Path:** `assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_alt.png`.

After the player ignores Marcus for a few seconds, his animation cycles
through this alt sheet once before returning to the normal strange idle.
Shows escalation without needing dialog to progress.

> [style lock]
>
> **Canonical Marcus identity** (match `npc_marcus_strange_idle.png`
> exactly): 10-year-old boy, spiky light-brown hair, fair skin, wide round
> wire-rim glasses, expressive brown eyes, freckled cheeks. Wears a
> **yellow short-sleeve polo shirt** with a small collar, **khaki cargo
> shorts**, white ankle socks, dark brown lace-up ankle shoes. Always
> holding a small spiral-bound **sketchbook** in one hand with a visible
> pencil-sketch on its top page (a pyramid, a face, a rooftop).
>
> **Animation:** 8-frame "he's losing it" cycle, mouth stays closed but
> contorted (grimace, clenched teeth).
>
> - Frame 1-2: rocks slightly, head tilted forward, looking at nothing,
>   sketchbook held loose in left hand at thigh.
> - Frame 3-4: drops sketchbook, brings both hands to temples.
> - Frame 5-6: hands claw through hair, shoulders hunched.
> - Frame 7: looks straight at camera, wide-eyed plea.
> - Frame 8: pulls hands away, loop-softens toward frame 1 silhouette.
>
> Baseline locked across all 8. Slightly muted / desaturated colors vs
> normal idle to suggest the psychic vision bleeding in.

---

## 5. PP — flower pickup

**Canvas:** 900×235. **Grid:** 6×1. **Cell:** 150×235.
**Path:** `assets/images/player/PP grab flower.png` (overwrites existing).

User complaint: picking up the flower currently has no animation — item
drops silently into inventory. Replace with a visible crouch + grab + rise.

> [style lock]
>
> Pink Panther at a cluster of daisies, three-quarter view. Six keyframes:
>
> - Frame 1: standing neutral, tail visible.
> - Frame 2: knees bending, upper body tilting forward.
> - Frame 3: crouch low, gloved right hand reaching toward a small yellow
>   daisy near feet.
> - Frame 4: hand has the flower, beginning to rise.
> - Frame 5: half-risen, flower held at chest with a small nod.
> - Frame 6: fully upright, flower held proudly at shoulder height,
>   satisfied half-lidded smile. Loops cleanly back to frame 1 silhouette
>   if played twice.
>
> Baseline Y of feet consistent across all 6 frames. Flower pixel art
> should read as "yellow daisy with green stem".

---

## 6. Campfire — smaller loop

**Canvas:** 1032×172. **Grid:** 6×1. **Cell:** 172×172.
**Path:** `assets/images/locations/camp/campfire_small.png`.

The current `campfire_idle.png` authoring renders too tall at 1.x scene
scale — user wants the visible flame to fit within roughly (581, 591) to
(702, 594) on screen. A smaller 172 px cell gives us that range when the
game draws it centered at (622, 573) at 1x scale.

> [style lock]
>
> Centered campfire flame over three charred log ends. Warm orange-yellow
> core, red tips, a faint ember glow at the base. **No ground, no grass**
> — just the flame + logs silhouette on pure white.
>
> **Animation:** 6-frame flicker loop. Flames wobble left-right, spark
> particles rise from the core on frames 3 and 5. Base (logs) stays still
> across every frame; only the flame moves. Logs about 80 px wide; flame
> tongue rises to ~160 px total height at peak.

---

## 7. Paris cafe ambient patrons

**Canvas:** 1376×768. **Grid:** 4×2 (8 patrons total).
**Cell:** 344×384.
**Path:** `assets/images/locations/paris/ambient/cafe_patrons.png` (new
folder — create it).

User asks for "background people that sit on the chairs and drink coffee in
loop". Authored small so they read as background detail, not interactive
NPCs.

> [style lock]
>
> Eight small seated Parisian patrons, each roughly 80 px tall, seen in
> front/three-quarter view. Each is an individual idle loop showing a
> different person sipping coffee, reading a newspaper, or chatting to an
> unseen companion. No two patrons look alike — vary hats, scarves, berets,
> coat colors.
>
> **Animation (per frame, each frame = a different patron, not a cycle):**
> each frame is its own 1-frame "pose" for that patron. The game will
> loop-render them in place with a tiny y-bob so they feel alive without
> dedicated per-patron cycles.
>
> Mouths closed, eyes half-lidded, each holds a small white coffee cup or
> pastry. No legs visible — they're seated at tables (the table edge is
> drawn as a thin line at the bottom of each cell).

---

## 8. Paris — Bakery Woman (new NPC for the pre-Louvre quest)

**Canvas:** 1376×768. **Grid:** 8×2 (take_row: 0 for idle, take_row: 1 for
talk). **Cell:** 172×384.
**Path:** `assets/images/locations/paris/npc/npc_bakery_woman.png`.

Stands under a red-striped boulangerie awning on the Paris street. PP buys
a baguette from her; she then hints about the press pass quest.

> [style lock]
>
> Madame Poulain: warm round face, flour dusted on apron, white baker's
> cap, dark hair in a bun, pleasant eyes, fifties. Wears white apron over a
> powder-blue dress. Holds a baguette across her forearm in every frame.
>
> **Row 0 (idle, 8 frames):** standing behind her counter, closed mouth.
> 1) neutral with baguette; 2) gentle inhale; 3) looks down at baguette;
> 4) taps baguette with free hand; 5) brushes flour off apron; 6) smiles
> slightly; 7) glances to the left (toward the street); 8) returns to
> neutral.
>
> **Row 1 (talk, 8 frames):** mouth open, same pose-to-pose beats but with
> speaking gestures — waves baguette once like a pointer on frame 4,
> nods on frame 6.
>
> Baseline locked across every frame of both rows. Use the same cell
> geometry as other kids so the atlas packer can apply `take_row: 0/1`
> cleanly.

---

## 9. Paris — Press Photographer (new NPC for the Louvre ticket quest)

**Canvas:** 1376×768. **Grid:** 8×2. **Cell:** 172×384.
**Path:** `assets/images/locations/paris/npc/npc_press_photographer.png`.

Lurks near the Louvre steps with a camera. PP trades a fresh baguette for
his press pass.

> [style lock]
>
> Nicolas, thin middle-aged man, rolled-sleeve shirt, suspenders, tweed
> flat cap. Holds a large vintage camera around his neck on a leather
> strap; a battered press-pass badge dangles on a ribbon at chest height.
>
> **Row 0 (idle):** 1) neutral, scanning the street; 2) lifts camera half
> way; 3) peers through viewfinder; 4) lowers camera; 5) wipes brow with
> handkerchief; 6) chuckles silently; 7) fidgets with press pass; 8) back
> to neutral.
>
> **Row 1 (talk):** mouth open, same poses but with gestures — points at
> camera on frame 3, taps press pass on frame 7.
>
> Baseline locked.

---

## 10. Higgins — entrance talk redesign (match new idle)

**Canvas:** 1376×768. **Grid:** 8×2 (`take_row: 0` is read).
**Cell:** 172×384.
**Path:** `assets/images/locations/camp/npc/higgins/npc_director_higgins_talk.png`.

Current file drifts from the regenerated idle — ruddier face, olive-green
pants instead of khaki, chunkier linework. Regenerate so idle + talk look
like the same day of shooting.

> [style lock]
>
> **Canonical identity** (match `npc_director_higgins_idle.png` exactly):
> older park ranger, ~60s, silver/gray combed-back hair, **silver-gray
> mustache**, round wire-rim glasses, ruddy weathered face. Wears a
> **dark-forest-green** short-sleeve camp shirt with two chest pockets, rolled
> cuffs, over an off-white undershirt collar. A **red lanyard** with a small
> white ID/badge hangs at mid-chest. Dark leather belt, light **khaki/tan**
> trousers with a small cargo pocket on the left thigh, dark brown lace-up
> ankle boots. Left hand holds a small **wooden clipboard** (with paper) at
> hip height.
>
> **Row 0 (talk, 8 frames) — mouth open, talking gestures:**
> 1) Neutral, eyebrows up, mouth open mid-word. 2) Right hand gestures wide
> toward camera. 3) Lifts clipboard to chest, points at it with free hand.
> 4) Palm up, head tilts right 4°. 5) Both hands forward in a "calm down"
> gesture. 6) Single firm finger-point toward camera, frown. 7) Shoulders
> relax, warm half-smile. 8) Returns to neutral silhouette of frame 1 so the
> loop is seamless.
>
> **Row 1 (optional alt poses / left empty)** — authored as padding so the
> sheet matches the existing 8×2 canvas; pack_atlas reads only row 0. You can
> leave row 1 transparent or repeat row 0 frames.
>
> Baseline (soles of boots) locked across all 8 frames. Silhouette within
> 4 px of the idle sheet's frame 1 so on-screen he does not "jump" when
> switching idle↔talk.

---

## 11. Tommy — full cast regen (idle + talk + strange_idle + strange_talk)

**Canvas (each):** 1376×768. **Grid:** 8×2 (packer reads both rows).
**Cell:** 172×384.
**Paths:**

- `assets/images/locations/camp/npc/kids/tommy/npc_tommy_idle.png`
- `assets/images/locations/camp/npc/kids/tommy/npc_tommy_talk.png`
- `assets/images/locations/camp/npc/kids/tommy/npc_tommy_strange_idle.png`
- `assets/images/locations/camp/npc/kids/tommy/npc_tommy_strange_talk.png`

Generate all four sheets in the same session with the same character
reference so face, hair, palette, and proportions stay identical across
animations.

> [style lock + character identity]
>
> Tommy: 10-year-old boy, tousled brown hair with a front tuft, warm-tan
> skin, green short-sleeve camp T-shirt with a small brown pine-tree motif
> centered on the chest, dark blue jeans rolled once at the ankle, **barefoot
> or socked only** (no shoes in any frame). Big earnest eyes, small closed
> mouth when idle. Always holds himself loose-limbed — this is the music-kid
> who hears songs no one else does.
>
> **Sheet A — `npc_tommy_idle.png`** (16 frames = 8×2, both rows used):
>
> Row 0 (8) — baseline fidget: neutral; weight shift left; weight shift
> right; scratches elbow; head tilts as if listening; blink; shoulders rise
> for a yawn; return to neutral. Mouth closed, baseline locked.
>
> Row 1 (8) — same anchor pose as row 0 frame 1 plus eight variant micro-
> motions (taps thigh, looks left, looks right, small sway, etc.) so the
> game can alternate rows without mesh-popping.
>
> **Sheet B — `npc_tommy_talk.png`** (16 frames = 8×2):
>
> Row 0 — mouth open, conversational gestures: neutral open mouth; points
> at camera; both palms up; counts on fingers; mild shrug; eager nod;
> gestures to side; returns to neutral.
>
> Row 1 — same beats shifted one notch (slight squash-stretch) so the
> talking read feels alive even at low FPS.
>
> **Sheet C — `npc_tommy_strange_idle.png`** (16 frames = 8×2):
>
> Like idle but with **faint colored music-note glyphs** drifting around his
> head (musical quaver, eighth-note, treble clef) as already done on the
> current sheet. Row 0 notes are small and subtle; row 1 notes are slightly
> more insistent.
>
> **Sheet D — `npc_tommy_strange_talk.png`** (16 frames = 8×2):
>
> Row 0 — excited sing-song gestures, hand waving, notes puff out of the
> mouth on frames 3, 5, 7. Row 1 — same but with a bigger bounce and notes
> trailing off the top of the cell.
>
> For all four sheets: baseline of bare feet locked across every frame;
> silhouette width within ±4 px of frame 1 in the same sheet. Keep the
> Danny/Lily linework weight.

---

## 12. Jake — full cast regen (idle + talk + strange_idle + strange_talk)

**Canvas (each):** 1376×768. **Grid:** 8×2.
**Cell:** 172×384.
**Paths:**

- `assets/images/locations/camp/npc/kids/jake/npc_jake_idle.png`
- `assets/images/locations/camp/npc/kids/jake/npc_jake_talk.png`
- `assets/images/locations/camp/npc/kids/jake/npc_jake_strange_idle.png`
- `assets/images/locations/camp/npc/kids/jake/npc_jake_strange_talk.png`

Current sheets are pure pixel art and must be converted to the cartoon look.
Generate all four sheets together using one identity reference.

> [style lock + character identity]
>
> Jake: sturdier 10-year-old boy, buzz-cut sandy-brown hair, round cheeks,
> warm skin, dark-green camp T-shirt with white "CAMP" chest decal, dark
> athletic shorts (knee-length), white socks, tan low-top sneakers. Stance
> is grounded — arms slightly away from the body. He is the skeptic kid, so
> keep his mouth line flat and eyes slightly narrowed in idle.
>
> **Sheet A — `npc_jake_idle.png`** (8×2): arms-folded fidget cycle. Row 0
> is the main 8-frame loop; row 1 is alternate takes so the engine can
> switch between them without desync.
>
> **Sheet B — `npc_jake_talk.png`** (8×2): arms unfold, hands gesture
> low-energy — "yeah, sure, whatever" reads. Mouth open for speech. Two
> rows.
>
> **Sheet C — `npc_jake_strange_idle.png`** (8×2): body stiffens, eyes go
> slightly unfocused, small **pale-green ripple** outlines his torso every
> few frames as if something is pulsing from inside his chest. No pixel
> dust, no sparkles — faint glowing outline only.
>
> **Sheet D — `npc_jake_strange_talk.png`** (8×2): wide alarmed eyes, mouth
> open as if speaking but hearing himself from outside. Hands flicker up to
> chest level on frames 3 and 6. Same green pulse outlines as sheet C but
> more frequent.
>
> Baseline (sneaker soles) locked across every frame and every sheet.

---

## 13. Marcus — full cast regen (idle + talk + strange_idle + strange_talk)

**Canvas (each):** 1376×768. **Grid:** 8×2. **Cell:** 172×384.
**Paths:**

- `assets/images/locations/camp/npc/kids/marcus/npc_marcus_idle.png`
- `assets/images/locations/camp/npc/kids/marcus/npc_marcus_talk.png`
- `assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_idle.png`  *(already on-style — regen only if a newer pass is needed)*
- `assets/images/locations/camp/npc/kids/marcus/npc_marcus_strange_talk.png`

Grids normalized to 8×2 across all four sheets during the 2026-04-19 regen
pass so `tools/characters/marcus.yaml` reads a single grid shape. Prior
layout was 6×2 for `idle`/`talk`; the regen pipeline's image generator
produced 1376×768 canvases consistently, so the yaml was updated to match.

> [style lock + character identity]
>
> Marcus: wiry boy with a pointy chin, messy upright brown hair with a wide
> front tuft, round wire glasses (kept throughout — he is the sketcher),
> yellow polo shirt with a small camp leaf on the chest, khaki cargo
> shorts, white socks, scuffed white low-top sneakers. Always carries a
> small spiral sketchbook in one hand and a pencil in the other.
>
> **Sheet A — `npc_marcus_idle.png`** (6×2 = 12 poses): alert sketching
> idle. Row 0 — glances at subject → draws → looks back → blinks → holds
> sketchbook up to examine → lowers it. Row 1 — six variant poses (tap
> pencil on chin, adjusts glasses, etc.) for engine row-switching.
>
> **Sheet B — `npc_marcus_talk.png`** (6×2): excited explain-it gestures.
> Row 0 — six key talk beats, Row 1 — six alt beats (different hand on
> sketchbook, etc.).
>
> **Sheet C — `npc_marcus_strange_idle.png`** (8×2): already on-style; only
> regen if you can produce a cleaner pass matching Danny/Lily linework.
> Keep the "shows the sketchbook with a face on it" motif.
>
> **Sheet D — `npc_marcus_strange_talk.png`** (8×2): replace the current
> pixel-leaning version. Glowing yellow eye-rings, wide-open mouth, both
> arms out in a "the painting speaks to me" gesture, camp-vest flapping.
> Row 0 — 8 escalating beats; row 1 — 8 come-down beats. **No floor tile,
> no shadow under feet** — cells must be on pure white.
>
> All four sheets: baseline of sneakers locked across every frame and every
> sheet.

---

## 14. French Guide (Madame Colette) — idle + talk redesign

**Canvas (each):** 1376×768. **Grid:** 8×2 for idle (cell 172×384), 8×1
for talk (cell 172×768). Matches the existing `loadNPCGrid(..., 8, 2)` /
`loadNPCGrid(..., 8, 1)` in `game/npc.go:newFrenchGuide`.
**Paths:**

- `assets/images/locations/paris/npc/npc_french_guide_idle.png`
- `assets/images/locations/paris/npc/npc_french_guide_talk.png`

Current sheets are pure pixel art and break the Paris street's cartoon look.
Convert to Hokus Pokus Pink style.

> [style lock + character identity]
>
> Madame Colette: French tour guide, mid-30s, tall, friendly-but-stern
> posture. Dark brown chin-length bob hair under a **navy beret**. Classic
> Parisienne outfit: white-and-red horizontal-striped three-quarter-sleeve
> top, scarlet neckerchief tied at the throat, dark navy wide-leg trousers,
> low-heel dark navy shoes. Warm smile, alert eyes, fair skin with a light
> blush. Stands like she has given this tour a thousand times.
>
> **Sheet A — `npc_french_guide_idle.png`** (8×2):
>
> Row 0 — 8-frame idle standing on the Paris street: breathing; glances
> right toward the Louvre; consults a small folded pamphlet; taps finger
> on pamphlet; peers over glasses; returns pamphlet to pocket; neutral;
> head tilt. Mouth closed.
>
> Row 1 — 8 alt takes (shifts weight, adjusts scarf, etc.).
>
> **Sheet B — `npc_french_guide_talk.png`** (8×1):
>
> One 8-frame row, mouth open, talking with her hands: welcomes visitor;
> gestures to the museum; counts off on fingers; palms-up "voilà"; leans
> slightly forward; nods; gestures to camera; returns to neutral.
>
> Baseline locked across every frame. Face her RIGHT on the page (toward
> the Louvre she's pointing at in-game; `newFrenchGuide` does not set
> `flipped=true`).

---

## 15. Museum Curator (Curator Beaumont) — idle + talk redesign

**Canvas (idle):** 1376×768. **Grid:** 8×1. **Cell:** 172×768.
**Canvas (talk):** 1376×768. **Grid:** 4×2. **Cell:** 344×384.
**Paths:**

- `assets/images/locations/paris/npc/npc_museum_curator_idle.png`
- `assets/images/locations/paris/npc/npc_museum_curator_talk.png`

Grid sizes match the existing `newMuseumCurator` constructor calls
(`loadNPCGrid(..., 8, 1)` for idle and `(..., 4, 2)` for talk). Do not
change dimensions unless we also update the constructor.

> [style lock + character identity]
>
> Curator Beaumont: distinguished Parisian art historian, late 50s. Short
> silver side-parted hair, neat silver mustache, black round glasses. Wears
> a dark charcoal three-piece suit, white dress shirt, black bowtie, a
> brass name-pin on the left lapel. Stands upright behind an ornate wooden
> podium (visible from the waist up — the podium is the bottom third of
> every cell). Holds a small brass magnifying glass in one hand on some
> frames.
>
> **Sheet A — `npc_museum_curator_idle.png`** (8×1): gentle "watchful
> caretaker" loop — adjusts glasses; taps podium edge; studies the
> magnifier; places magnifier down; folds hands; looks right; looks left;
> returns to neutral. Mouth closed. Baseline (podium edge) locked.
>
> **Sheet B — `npc_museum_curator_talk.png`** (4×2 = 8 poses): mouth open
> in conversational cadence; gestures with free hand; raises magnifier on
> frames 2 and 6; nods on frame 4. Same podium locked at the bottom of
> every cell.
>
> Face LEFT on the page (he points away from the podium toward the gallery
> hall in-game; `newMuseumCurator` does not flip).

---

## 16. Paris flavor NPCs — verification-driven regens *(currently skipped)*

The verification sweep marked these Paris sheets as **keep**:

- `npc_art_vendor.png` (Pierre, 8×2 — cartoon, fits style).
- `npc_security_guard.png` (Claude, 6×2 — cartoon, fits style).
- `npc_mystery_figure.png` (cartoon hooded figure, mood-appropriate).
- `npc_suspicious_dealer.png` (cartoon, fits style).

Do **not** regen these unless a later pass flags them. If a future change
demands one, copy the §14/§15 structure and anchor it to the existing grid
dimensions in `tools/characters/paris/*.yaml` (if we ever author them) or
the constructor's `loadNPCGrid` args in `game/npc.go`.

---

## 17. Pink Panther — full audit *(currently all sheets keep)*

The verification sweep confirms every PP sheet reads as canonical Hokus Pokus
Pink cartoon. The only queued regens are:

- **§3** `PP walk back.png` — clearer walk cycle for the "leaving camp"
  transition. Already in this document, unchanged.
- **§5** `PP grab flower.png` — visible crouch + grab + rise. Already in
  this document, unchanged.

No new PP prompts are added in this sweep. If future content needs an
"eating baguette" or "using press pass" take, author the prompt here and
update `CHARACTERS.md` with the cell dimensions.

---

## 18. Higgins — office idle + office talk redesign

**Canvas (each):** 1032×768. **Grid:** 6×2 (`take_row: 0` read for idle;
both rows read for talk).
**Cell:** 172×384.
**Paths:**

- `assets/images/locations/camp/npc/higgins/npc_director_higgins_office_idle.png`
- `assets/images/locations/camp/npc/higgins/npc_director_higgins_office_talk.png`

Current sheets are pixel art and clash with the cartoon entrance Higgins.
Regenerate in the cartoon style. Grid dimensions match
`tools/characters/higgins.yaml` — do not change them.

> [style lock]
>
> Same Higgins identity as §10 talk (silver/gray hair + mustache, wire-rim
> glasses, dark-forest-green shirt over cream undershirt, red lanyard with
> badge, khaki trousers, dark brown boots). In the office variant he is
> **seated behind a wooden desk** (desk edge visible at the bottom of every
> cell), green clipboard or open notebook on the desk, small brass lamp
> top-right corner of the desk. Chair arms visible at elbow height. No
> clipboard in his hand (it's on the desk).
>
> **Sheet A — `npc_director_higgins_office_idle.png`** (6×2):
>
> Row 0 — 6 calm reading / desk-work frames: reading notebook; turns page;
> adjusts glasses with index finger; leans back and stretches (arms up);
> sips from a brown mug; sets mug down.
>
> Row 1 — 6 alt micro-poses (tap pen on desk, look right, etc.).
>
> **Sheet B — `npc_director_higgins_office_talk.png`** (6×2):
>
> Row 0 — 6 conversational beats: open mouth explaining; points at
> notebook; palm up offering; palm down reassuring; leans forward; back to
> neutral.
>
> Row 1 — 6 alt takes.
>
> Baseline (desk edge on each cell's bottom) locked across every frame of
> both sheets. Desk silhouette identical in every frame — only Higgins
> animates.

---

## 19. Higgins — give_map handoff redesign

**Canvas:** 1376×384. **Grid:** 8×1. **Cell:** 172×384.
**Path:** `assets/images/locations/camp/npc/higgins/npc_director_higgins_give_map.png`.

Current sheet is pixel-leaning and inconsistent with the entrance Higgins.
Regenerate in the cartoon style so the Day-2 map-handoff beat flows into
PP's existing `receive map` animation without a style jump.

> [style lock]
>
> Same Higgins identity as §10 and §18. Seated behind the same office desk
> (desk edge at the bottom of every cell). He is handing a folded paper map
> across the desk toward the viewer.
>
> **Animation:** 8-frame handoff, mouth closed.
>
> - Frame 1: neutral, folded map sitting on desk in front of him.
> - Frame 2: picks up map with right hand.
> - Frame 3: begins unfolding map (partial fold visible).
> - Frame 4: map fully unfolded, held at chest height facing camera.
> - Frame 5: glances at map once, nods.
> - Frame 6: re-folds map with both hands.
> - Frame 7: extends folded map forward across the desk, eyes on camera.
> - Frame 8: holds the offer — map held steady, calm half-smile. Loops
>   cleanly back to frame 1 if played twice.
>
> Baseline (desk edge) locked. Silhouette within ±4 px of frame 1 on every
> frame.

---

## Re-enabling after regeneration

1. Save the PNG at the path in the section heading.
2. If the sheet has a matching entry in `tools/characters/*.yaml`, run:
   ```
   python tools/pack_atlas.py tools/characters/<name>.yaml
   ```
3. For legacy-path NPCs (Higgins entrance idle, Paris new NPCs), no
   manifest exists yet — just re-launch the game; the constructor loads
   directly from the PNG path.
4. For ambient patrons and campfire small: the consumer code is TBD as
   these land alongside Phase C and the camp-fire-size FIXME item.

If a cell looks chopped after regeneration, the grid is wrong — adjust
`grid: [X, Y]` in the manifest or the constructor's `loadNPCGrid(..., X,
Y)` call and re-run.
