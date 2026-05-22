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
| `paris/ambient/cafe_patrons.png` | superseded | §7 | First-pass mock-up (8 single-frame poses with baked-in white table edge) replaced by six per-patron sheets (`cafe_patron_<name>.png`) — see §7.1–§7.6. Renderer hookup still deferred. |
| `paris/npc/coffee/cafe_patron_<name>_<idle\|talk>.png` (×12 sheets) | TODO | §7.1–§7.6 | User 2026-05-20: replaced the old single-frame patrons with **twelve** sheets — six patrons × two animations each (idle + talk, 8×1 strip per sheet, 100×170 per cell). New folder: `paris/npc/coffee/` (existing single-pose PNGs there will be overwritten by the new `_idle.png` per patron). |
| `paris/background/paris_bakery.png` | TODO (regen) | §NEW Paris Bakery | Café-corner rework: counter shifts right, three bistro tables + chairs added on the left, rolling-pin floor patch moves to `(740, 720)`. JSON wiring follow-up tracked in FIXME. |
| `paris/background/paris_clouds.png` | TODO (regen) | §NEW Paris Clouds | Replaces the static transparent-bg cloud row with a full 1400×800 sky background for the airplane flight cutscene. |

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

## 7. Paris cafe ambient patrons (six individual sheets)

**Replaces** the original single 4×2 `cafe_patrons.png` mock-up. That sheet
had every patron as a single-frame pose with a baked-in white table edge,
which clashes once the bakery BG paints real oak bistro tables in (see
§NEW Paris Bakery below). These six sheets are **chest-up only, no table,
no saucer, no chair back** — the engine drops them on top of the painted
tables in `paris_bakery.png` and the BG's table reads behind their hands.

**Per-sheet canvas:** 800 × 340. **Grid:** 8 × 2 (row 0 idle, row 1 talk).
**Cell:** 100 × 170 — sized to the engine bounds proposed in
`paris_bakery.json`.

**Shared style lock (paste at the top of every prompt below):**

> Hand-drawn 1990s Saturday-morning cartoon, Pink Panther *Hokus Pokus
> Pink* (1997) / *Passport to Peril* (1996). Confident black ink linework
> ~2 px, flat saturated fills, two cel tones max per color region, no
> airbrush, no gradients, no pixel art. **Pure #FFFFFF background.**
>
> Each cell shows the patron **from the waist up only** in a seated pose.
> **Do NOT draw the table, the chair back, the saucer, the chair seat,
> any table props, or any floor surface** — the engine composites the
> patron over the bakery background which already paints those in. Only
> the patron's body, clothing, hands, and the item they hold (cup,
> newspaper, etc.) appear in the cell. Background behind the patron is
> pure white.
>
> Baseline (bottom of the patron's silhouette — typically forearms or
> jacket hem) locked across all 16 frames of the sheet. Silhouette width
> within ±3 px of frame 1. Hands stay at the same Y position across the
> sheet so the BG's table-edge line reads continuous behind them.

---

### 7.1 Madame Yvette — beret + pearls, sipping tea

**Canvas (each):** 800×170. **Grid:** 8×1. **Cell:** 100×170.
**Paths (TWO files):**
- `assets/images/locations/paris/npc/coffee/cafe_patron_yvette_idle.png`
- `assets/images/locations/paris/npc/coffee/cafe_patron_yvette_talk.png`

> [style lock — character identity, paste into BOTH idle and talk
> prompts]
>
> Elderly Parisienne, late 70s, fair skin with soft wrinkles. Short
> curled silver hair under a **black wool beret** tilted slightly
> right. Single strand of pearls at the throat. Wears a
> **mustard-yellow knit cardigan** over a cream blouse `#F2EFE5` (NOT
> white) with a brown brooch. Warm dignified expression. Holds a
> small cream teacup `#E5DDC8` (NOT white) in both hands at chest
> height. Seated, chest-up framing. Black `#1C1C1C` ink outline ~3 px.
> **No pure white on the character or teacup.** Pure white `#FFFFFF`
> background, no scenery. Every cell exactly 100×170, foot/chair
> baseline locked across all 8 cells.

> **`cafe_patron_yvette_idle.png` — 8-frame idle loop, mouth closed.**
>
> 1) cup at chest, eyes lowered to it
> 2) lifts cup toward mouth
> 3) sips (eyes closed contented)
> 4) lowers cup back to chest height
> 5) free finger touches lip
> 6) gentle nod, eyes down
> 7) glances right
> 8) returns to neutral matching frame 1 silhouette so the loop seams

> **`cafe_patron_yvette_talk.png` — 8-frame talk loop, mouth open.**
>
> Same pose anchor as idle, cup still at chest height in left hand.
>
> 1) mouth open mid-word, both hands on cup
> 2) free hand rises in a soft palm-up gesture
> 3) wags finger gently
> 4) nods while speaking
> 5) raises eyebrows
> 6) tilts head 4° left
> 7) chuckles silently (eyes crinkle)
> 8) returns to neutral

---

### 7.2 Monsieur Bernard — bearded man reading *Le Figaro*

**Canvas (each):** 800×170. **Grid:** 8×1. **Cell:** 100×170.
**Paths (TWO files):**
- `assets/images/locations/paris/npc/coffee/cafe_patron_bernard_idle.png`
- `assets/images/locations/paris/npc/coffee/cafe_patron_bernard_talk.png`

> [style lock — character identity, paste into BOTH idle and talk
> prompts]
>
> Mid-50s man, warm tan skin, **full salt-and-pepper beard**, brown
> eyes, **brown tweed flat cap**. Wears a **mustard-brown corduroy
> jacket** over a charcoal shirt `#2A2A2A`. Holds a folded broadsheet
> newspaper open at chest height in both hands — paper body cream
> `#E5DDC8` (NOT white) with black `#1C1C1C` "LE FIGARO" header
> readable on the top edge. Seated, chest-up framing. Black `#1C1C1C`
> ink outline ~3 px. **No pure white on the character or newspaper.**
> Pure white `#FFFFFF` background, no scenery. Every cell exactly
> 100×170, chair baseline locked across all 8 cells.

> **`cafe_patron_bernard_idle.png` — 8-frame idle loop, mouth closed.**
>
> 1) reading paper, eyes scanning
> 2) turns a page (paper crinkles wider)
> 3) eyebrows raise at headline
> 4) lowers paper to lap level (paper still in both hands)
> 5) raises paper back to chest
> 6) reads, slight frown
> 7) folds paper edge in
> 8) returns to neutral matching frame 1 silhouette

> **`cafe_patron_bernard_talk.png` — 8-frame talk loop, mouth open.**
>
> Paper lowered to lap for all 8 frames so the face is visible.
>
> 1) mouth open, right hand gestures to paper at lap
> 2) taps headline with index finger
> 3) shakes head slightly
> 4) palms up "can you believe it"
> 5) leans forward slightly
> 6) shrugs
> 7) huffs (mouth in O-shape)
> 8) returns to neutral, paper at lap

---

### 7.3 Mademoiselle Camille — red-beret art student

**Canvas (each):** 800×170. **Grid:** 8×1. **Cell:** 100×170.
**Paths (TWO files):**
- `assets/images/locations/paris/npc/coffee/cafe_patron_camille_idle.png`
- `assets/images/locations/paris/npc/coffee/cafe_patron_camille_talk.png`

> [style lock — character identity, paste into BOTH idle and talk
> prompts]
>
> Young woman, early 20s, fair skin with light freckles, dark brown
> chin-length bob. **Bright scarlet beret `#C4412A`** tilted left.
> **Emerald-green wrap dress `#2E7D5A`** with a small gold collar
> pin. Cradles a small cream cappuccino cup `#E5DDC8` (NOT white) in
> both hands at chin height — a tan-foam heart `#C9A878` visible from
> above. Seated, chest-up framing. Black `#1C1C1C` ink outline ~3 px.
> **No pure white on the character or cup.** Pure white `#FFFFFF`
> background, no scenery. Every cell exactly 100×170, chair baseline
> locked across all 8 cells.

> **`cafe_patron_camille_idle.png` — 8-frame idle loop, mouth closed.**
>
> 1) cup at chin, eyes lowered to foam
> 2) blows softly on cup (lips pursed, NOT mouth-open)
> 3) sips (eyes briefly closed)
> 4) lowers cup to chest
> 5) tucks hair behind ear with right hand (cup in left at chest)
> 6) raises cup again
> 7) glances down
> 8) returns to neutral matching frame 1 silhouette

> **`cafe_patron_camille_talk.png` — 8-frame talk loop, mouth open.**
>
> Cup in left hand at chest height for all 8 frames so the right hand
> is free to gesture.
>
> 1) mouth open, right hand gestures forward
> 2) excited finger-point
> 3) both hands raise palms up (cup briefly in left fingertips only)
> 4) laughs (head tips back 5°)
> 5) eyes wide
> 6) shrugs
> 7) leans on left elbow
> 8) returns to neutral with cup at chest

---

### 7.4 Monsieur Henri — silver-haired gentleman with pastry

**Canvas (each):** 800×170. **Grid:** 8×1. **Cell:** 100×170.
**Paths (TWO files):**
- `assets/images/locations/paris/npc/coffee/cafe_patron_henri_idle.png`
- `assets/images/locations/paris/npc/coffee/cafe_patron_henri_talk.png`

> [style lock — character identity, paste into BOTH idle and talk
> prompts]
>
> Distinguished older man, late 60s, fair skin, **thick silver
> handlebar mustache** (use silver `#C0C0C0`, NOT white), neatly
> combed silver hair parted on the side, small round wire-rim
> glasses. **Navy three-piece suit `#2B3A5C`** with a **burgundy
> bowtie `#7A2A2A`** and a small cream pocket square `#E5DDC8` (NOT
> white). Holds a **golden croissant `#B97E3F`** in his right hand at
> chest height; left hand rests at chest height as if resting near a
> cup. Seated, chest-up framing. Black `#1C1C1C` ink outline ~3 px.
> **No pure white on the character or pocket square.** Pure white
> `#FFFFFF` background, no scenery. Every cell exactly 100×170,
> chair baseline locked across all 8 cells.

> **`cafe_patron_henri_idle.png` — 8-frame idle loop, mouth closed.**
>
> 1) breaks off a corner of croissant
> 2) brings piece to mouth (lips closed around it, NOT mouth-open)
> 3) chews (mustache twitches)
> 4) lowers croissant to chest
> 5) raises left hand as if lifting a cup (no cup drawn)
> 6) lowers left hand
> 7) brushes mustache with right index finger (croissant briefly in
>    left hand)
> 8) returns to neutral with croissant in right hand at chest

> **`cafe_patron_henri_talk.png` — 8-frame talk loop, mouth open.**
>
> Croissant in left hand at chest height for all 8 frames so the
> right hand is free to gesture.
>
> 1) mouth open, right hand gestures palm-up
> 2) raises right index finger as a point
> 3) pats chest with right hand
> 4) taps temple with right finger
> 5) gestures wide with right palm
> 6) nods firmly
> 7) chuckles (eyes squint, mustache lifts)
> 8) returns to neutral

---

### 7.5 Lucien — young man in gray turtleneck

**Canvas (each):** 800×170. **Grid:** 8×1. **Cell:** 100×170.
**Paths (TWO files):**
- `assets/images/locations/paris/npc/coffee/cafe_patron_lucien_idle.png`
- `assets/images/locations/paris/npc/coffee/cafe_patron_lucien_talk.png`

> [style lock — character identity, paste into BOTH idle and talk
> prompts]
>
> Late 20s man, olive skin, dark wavy black hair just past ears,
> faint stubble, thoughtful brown eyes. Wears a **slate-gray
> turtleneck sweater `#5A6473`** with rolled cuffs at the wrists.
> Both hands cradle a small cream espresso cup `#E5DDC8` (NOT white)
> at chest height. Seated, chest-up framing. Black `#1C1C1C` ink
> outline ~3 px. **No pure white on the character or cup.** Pure
> white `#FFFFFF` background, no scenery. Every cell exactly
> 100×170, chair baseline locked across all 8 cells.

> **`cafe_patron_lucien_idle.png` — 8-frame idle loop, mouth closed.**
>
> 1) cup at chest, eyes lowered
> 2) raises cup to mouth
> 3) sips with eyes closed
> 4) lowers cup to chest
> 5) drums fingers of right hand once at chest height (cup in left)
> 6) glances left out of frame
> 7) head tilts back as he reflects
> 8) returns to neutral matching frame 1 silhouette

> **`cafe_patron_lucien_talk.png` — 8-frame talk loop, mouth open.**
>
> Cup in left hand at chest for all 8 frames so right hand is free.
>
> 1) mouth open, right hand palm-up
> 2) finger-point forward
> 3) palm pat at chest height (no table drawn)
> 4) shrugs
> 5) looks aside (sarcastic)
> 6) leans forward intently
> 7) eyebrows raise
> 8) returns to neutral with both hands around cup

---

### 7.6 Madame Élise — red-haired woman in autumn scarf

**Canvas (each):** 800×170. **Grid:** 8×1. **Cell:** 100×170.
**Paths (TWO files):**
- `assets/images/locations/paris/npc/coffee/cafe_patron_elise_idle.png`
- `assets/images/locations/paris/npc/coffee/cafe_patron_elise_talk.png`

> [style lock — character identity, paste into BOTH idle and talk
> prompts]
>
> Mid-40s woman, fair skin, **shoulder-length wavy auburn-red hair
> `#B8542D`**, green eyes, soft smile. Wears a **chunky cream
> cable-knit sweater `#E5DDC8`** (NOT white) with a **floral
> autumn-print scarf** (orange `#D9762A`, mustard `#C9A878`, brick
> red `#9E2A1F`) wrapped twice around the neck. Holds a small cream
> cup `#E5DDC8` (NOT white) at chin height with both hands. Seated,
> chest-up framing. Black `#1C1C1C` ink outline ~3 px. **No pure
> white on the character, sweater, or cup.** Pure white `#FFFFFF`
> background, no scenery. Every cell exactly 100×170, chair baseline
> locked across all 8 cells.

> **`cafe_patron_elise_idle.png` — 8-frame idle loop, mouth closed.**
>
> 1) cup at chin, eyes lowered
> 2) sips
> 3) lowers cup to chest, gentle exhale (pale grey `#C4C4C4` steam
>    wisp above the cup)
> 4) adjusts scarf with right hand (cup in left)
> 5) tucks a curl behind ear
> 6) lifts cup again
> 7) glances right
> 8) returns to neutral matching frame 1 silhouette

> **`cafe_patron_elise_talk.png` — 8-frame talk loop, mouth open.**
>
> Cup in left hand at chest height for all 8 frames so right hand is
> free.
>
> 1) mouth open, right hand gestures softly
> 2) hand to chest (sincere)
> 3) palm up
> 4) nods, eyes close briefly
> 5) leans forward
> 6) gestures outward
> 7) laughs silently (eyes crinkle)
> 8) returns to neutral with both hands on cup

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

---

## §NEW: PP full-set regen — restore saturation, kill background bleed

**Why:** PP sheets have been losing color and picking up cream-white halos
across regen passes. This prompt locks the look so every PP anim renders
with the same Pink Panther silhouette, palette, and outline weight.

**Per-anim canvas:** 1376 × 384 (8 frames in one row, cell 172 × 384).
For 2-row sheets (e.g. PP idle front 8×2) double the canvas height to
768 px and use the second row for the same character in a held breath /
slight wiggle variation.

**Prompt body** (copy verbatim, swap the `[POSE]` token per anim):

> Generate a clean Pink Panther sprite sheet, 8 frames in one row, 172×384
> per cell, total canvas 1376×384. **Pure white background**: RGB
> (255,255,255), flat, no gradient, no shadow, no painted halo, no rim
> light. Character is the classic Pink Panther: bright bubblegum pink
> body (#F58FB7), crisp dark pink shading on belly + tail (#D14E92),
> thin black ink outline (1-2 px). Yellow eyes (#FFE15A) with black
> pupils. Loose comedic posture, slightly off-balance.
>
> **Pose:** [POSE].
>
> Cartoon line art, NOT pixel art. Every frame must read as the same
> character — same pose anchor, same shading style, same outline
> thickness. No painted backdrop, no drop shadow, no decorative elements
> beside or behind the character. Hokus Pokus Pink reference style.

**`[POSE]` substitutions per sheet** (matches files in
`assets/images/player/`):

| Sheet | POSE token |
|---|---|
| `pp_idle_front.png` | "facing the camera, arms relaxed at sides, weight shifting between feet" |
| `pp_idle_side.png` | "in profile facing right, tail trailing behind, hands at hips" |
| `pp_idle_back.png` | "back to camera, looking over shoulder, tail flick" |
| `pp_walk_front.png` | "walking toward camera, alternating leg swings, tail counter-swinging" |
| `PP walk back.png` | "walking away from camera (back to viewer), alternating leg swings, tail counter-swinging" |
| `pp_walk_left.png` | "walking left in profile, full stride cycle" |
| `pp_talk_front.png` | "facing camera, mouth open in conversation, hands gesturing softly" |
| `pp_talk_side.png` | "in profile facing right, talking, one hand raised palm-up" |
| `PP grab flower.png` | "starting upright, crouching down to pick a flower from the ground, then rising holding the flower in front" |
| `PP receive map.png` | "facing right (toward Higgins), hand extended to receive a folded map, then bringing it to chest level and slipping it into pocket" |
| `pp_celebrate.png` | "joyful jump, arms up, tail curled exclamation" |
| `pp_sneak_examine.png` | "crouched, holding a magnifying glass to a small object" |
| `pp_sneak_use.png` | "crouched, applying an item with focused expression" |
| `pp_sleeping.png` | "lying on side, head on rolled jacket, slow rise-and-fall breathing" |
| `pp_waking.png` | "lifting head from sleep, stretch, rub eye, sit up" |
| `pp_airplane.png` | "seated inside a small biplane cockpit visible through the side window, head turning side to side" |

**After-generation pipeline:**

1. Run `python tools/clean_generated_sheet.py <path>` to strip the black
   frame + grid lines the generator bakes in.
2. Drop the cleaned PNG at the target path (overwrite the existing
   sheet).
3. Re-launch the game — runtime loaders use `gridFrames` /
   `SpriteGridFromPNGCleanAggressive` (tol 24, inset 4) which key out
   any residual white halo. **No code changes required**: paths and grid
   args are unchanged.

If a frame looks chopped, the generator gave back a different cell
count — verify the canvas is exactly 1376 × 384 (or 1376 × 768 for
2-row sheets) before slicing.

---

## §NEW: Paris Bakery Interior — café-corner regen

**Why:** First pass produced a counter-only bakery (no seating). To host
the six café patrons from §7 we need three small bistro tables in the
left foreground. The counter shifts right to make room; Madame Poulain
relocates accordingly; the rolling-pin floor item moves to the corridor
between the seating area and the counter.

**Path:** `assets/images/locations/paris/background/paris_bakery.png`
(overwrite).

**Canvas:** 1400 × 800.

### Spatial constraints (do not deviate)

- Madame Poulain renders at bounds `(820, 440, 140, 240)` (shifted right
  from `(540, 440)`) — feet at `y≈680`. She must read as standing behind
  a counter with the counter top edge at `y=620–650`.
- Six patron NPCs render at the bounds in the table below (matches §7
  cell size 100 × 170). Each pair faces each other across a small
  bistro table painted into the BG.

  | Patron | Table | Bounds (x, y, w, h) |
  |---|---|---|
  | Yvette  | A (left chair)   | `(170, 540, 90, 160)` |
  | Bernard | A (right chair)  | `(270, 540, 90, 160)` |
  | Camille | B (left chair)   | `(380, 555, 90, 160)` |
  | Henri   | B (right chair)  | `(480, 555, 90, 160)` |
  | Lucien  | C (left chair)   | `(560, 540, 90, 160)` |
  | Élise   | C (right chair)  | `(660, 540, 90, 160)` |

- Rolling-pin floor item moves to `(740, 720)` — clean floor patch
  between the café area and the counter.
- Walking corridor for PP: `x=200–1100, y=730–790` must be visually
  clean (no rugs, no debris).
- Left blocker `x=0–180, y=0–500` (curtained street exit) and right
  blocker `x=1080–1400, y=0–600` (shelving) are unchanged.

### Prompt

> **Style:** hand-drawn 1990s Saturday-morning cartoon, Pink Panther
> *Hokus Pokus Pink* (1997) / *Passport to Peril* (1996). Confident
> black ink linework ~2-3 px around major shapes, flat saturated fills,
> two cel tones max per material region, no airbrush, no gradients, no
> photorealism, no pixel art. Warm cozy golden-hour interior glow.
>
> **Scene:** small Parisian boulangerie that doubles as a tiny café.
> View from inside facing the back wall. **No people, no animals.**
> Three empty bistro tables on the left, an empty counter on the right,
> ready for the engine to drop sprites in.
>
> **Composition (left → right across the 1400 × 800 canvas):**
>
> - **x=0–180, y=200–760:** open arched doorway with a **half-drawn
>   red-and-white striped curtain** pulled aside. Sliver of cobblestone
>   street and a Paris lamppost peek through. Exit hotspot — must read
>   as "way out."
> - **x=180–700, y=540–790 — CAFÉ SEATING AREA:** three small **round
>   dark-oak bistro tables**, each ~120 px diameter, with **two bentwood
>   Thonet chairs** facing each other across the table. Tables are
>   empty: Table A has a small white espresso cup + saucer, Table B has
>   a folded *Le Figaro* newspaper resting on the edge, Table C has a
>   tiny vase with a single daisy. Approximate centers:
>   - Table A: center `(220, 680)` — chairs at `(170, 680)` and `(270, 680)`.
>   - Table B: center `(430, 695)` — chairs at `(380, 695)` and `(480, 695)`.
>   - Table C: center `(610, 680)` — chairs at `(560, 680)` and `(660, 680)`.
>   Chair seats at `y≈700`, chair backs rising to `y≈620`. Tables read
>   slightly behind the chairs so a seated patron sprite anchored at
>   the chair overlaps the table edge naturally.
> - **x=180–700, y=100–540 — café back wall:** cream plaster wall with
>   a small **black framed chalkboard menu** ("CAFÉ — 2F EXPRESSO,
>   5F CROISSANT") above Table B; a **vintage Paris poster** (Eiffel
>   Tower silhouette in red and black) above Table A; a small **brass
>   wall sconce** above Table C casting a warm pool of light.
> - **x=720–1080, y=580–700 — COUNTER (shifted right):** wooden bakery
>   counter, dark stained oak, top edge at `y=620` so an NPC anchored
>   at `(820, 440)` size `140×240` reads as standing behind it from
>   the waist up. **Brass scale** on the left end, **wicker basket of
>   golden baguettes** standing upright on the right end, **glass-domed
>   pastry display** in the center showing croissants, pain au chocolat,
>   and macarons.
> - **x=720–1080, y=100–580 — back wall behind counter:** **brick
>   wood-fire oven** built into the wall, dome-shaped iron door open
>   with a warm **orange glow** inside (`#F4A23C` core, `#C25A2C` rim),
>   three loaves of bread on the brick lip cooling. **Stack of split
>   firewood** at the oven's base on the left side of the counter.
> - **x=1080–1400, y=100–600 — right wall:** floor-to-ceiling **wooden
>   shelving** with stacked round country loaves on every shelf, two
>   sacks of flour stenciled "FARINE" leaning against the base, a
>   **chalkboard sign** ("BOULANGERIE POULAIN") hanging on the right
>   side with prices in francs scrawled in white chalk.
> - **Above all (y=0–180):** wooden rafter beams across the ceiling
>   with **dried wheat sheaves** hanging at each end and a single
>   **filament bulb** centered over the counter.
> - **Floor (y=620–800):** wide-plank wooden floor, warm honey-tan
>   (`#C8965A`), visible plank seams running toward the camera. Keep
>   the strip `x=200–1100, y=730–790` visually clean. Small clear
>   floor patch near `(740, 720)` for the rolling-pin floor item.
>
> **Palette:** warm cream walls (`#F5E6C8`), oak counter + tables
> (`#8B5A2B` outline, `#A87044` mid, `#5A3A1E` shadow), oven brick
> (`#A24A2C`), oven glow (`#F4A23C` → `#FFD074` core), bread crusts
> (`#D6A45C` with darker `#8B5A2B` scoring lines), red-white awning
> stripes (`#C4412A` + `#FFFFFF`), bentwood chairs (`#6E4A28`).
>
> **Hard rules:** No characters. No floating props. The walking
> corridor `y=730–790, x=200–1100` must be clear. The counter top edge
> is a clean horizontal line at `y=620`. Tables and chairs stay in
> `x=180–700, y=540–790`. No painted shadows on the floor (the engine
> adds those per-actor at runtime).

### Companion item — `assets/images/items/rolling_pin.png` (64 × 64)

> A simple wooden rolling pin viewed from the side, warm tan wood
> (`#D2A877`), darker grain lines, two small handles. Centered on a
> transparent background. Cartoon line art, no shadow.

After the BG and patron sheets land, the JSON wiring follow-up: update
`assets/data/scenes/paris_bakery.json` to (a) add the six patron NPC
ids to `npcs`, (b) move Madame Poulain's bounds to `(820, 440, 140,
240)`, and (c) move the floor-item drop to `(740, 720)`. (Tracked in
FIXME.md.)

---

## §NEW: Paris Clouds — airplane flight sky

**Why:** Current `paris_clouds.png` is a transparent canvas with a static
row of cloud puffs frozen near the top. Used as the airplane flight
cutscene background, it reads as cardboard cutouts in checkerboard void.
Needs to be a proper full-canvas sky so PP's biplane sprite has
something to fly through.

**Path:** `assets/images/locations/paris/background/paris_clouds.png`
(overwrite).

**Canvas:** 1400 × 800.

### Spatial constraints

- Background must fill the entire 1400 × 800 frame (no transparency).
- The `flight_cutscene` renderer parallax-scrolls this image
  horizontally; the painting must read as continuous flight, not
  floating-in-place clouds.
- PP's biplane sprite renders centered around `(700, 400)` at flight
  altitude — keep clouds away from a soft-focus zone of roughly
  `(550–850, 320–480)` so the biplane silhouette stays readable
  against the sky.

### Prompt

> **Style:** hand-drawn 1990s Saturday-morning cartoon, Pink Panther
> *Hokus Pokus Pink* (1997) / *Passport to Peril* (1996). Confident
> black ink linework ~2 px on cloud silhouettes, flat saturated fills,
> two cel tones max per element, no airbrush, no gradients (band the
> sky as discrete cel tones, not a smooth gradient), no pixel art.
>
> **Scene:** sunny daytime sky seen from cruising altitude. **No
> aircraft, no characters.** The painting is a backdrop that the engine
> will scroll horizontally behind PP's biplane.
>
> **Sky bands (top → bottom, flat cel tones, no gradient):**
> - `y=0–250`: bright sky blue (`#7ED4F2`) — the upper sky.
> - `y=250–520`: mid sky (`#B8E5F8`) — where the biplane will fly.
> - `y=520–740`: pale haze (`#E6F4FA`) — the air just above the
>   horizon.
> - `y=740–800`: a thin band of distant landscape — soft blue-gray
>   silhouettes (`#9FB6C8`) of rolling hills, a hint of a small
>   château or church spire in profile, no sharp detail. Gives the
>   biplane scale.
>
> **Clouds (8–12 total, NOT a single neat row):** cartoon
> Pink-Panther-era cumulus puffs, **off-white** (`#FAFAFA`) bodies with
> a single **soft gray shadow** (`#D8D8D8`) tucked under each cloud's
> belly. Vary the cloud sizes, Y positions, and shapes so the parallax
> scroll reads as forward flight.
>
> Suggested placement:
> - 3 large foreground clouds (~300 × 110 px) at `y≈200`, `y≈610`, and
>   `y≈100` — these scroll fastest in parallax.
> - 4 medium midground clouds (~180 × 70 px) scattered across `y=80`,
>   `y=350`, `y=480`, `y=620` — slightly lighter outline weight.
> - 3–5 small distant clouds (~100 × 40 px) drifting near `y=550–700`
>   — very faint outlines, mostly silhouette, suggest distance.
>
> Each cloud is **horizontally tilt-stretched ~5–10°** (subtly slanted
> trailing edges, as if motion-streaked) so the parallax scroll
> visually reinforces forward flight rather than sideways drift.
>
> **Avoid zone:** keep `(550–850, 320–480)` clear of large clouds so
> PP's biplane sprite (which the engine renders centered there) reads
> cleanly against the mid-sky band.
>
> **Hard rules:** Solid sky fills the entire canvas — no transparency,
> no checkerboard. No sun, no rainbow, no birds, no painted shadow on
> the clouds beyond the single soft underside tone. Cloud outlines are
> consistent ink weight (~2 px). The whole image must tile reasonably
> well horizontally (left edge sky color matches right edge sky color)
> so the parallax loop doesn't seam-pop.

After the PNG lands, the airplane flight cutscene auto-uses it via
`assets/data/scenes/airplane_flight.json` — no code changes.

---

## §NEW: Kid walk-away sheets (Tommy + Jake)

After PP finishes talking to a Day-1 camp kid, the kid walks out of the
scene so PP doesn't end up surrounded by a static cluster. **This pass:
Tommy walks left (out toward the lake path), Jake walks up-into-frame
toward his cabin door.** Danny right-bottom is queued for a later pass.

Both sheets are full body, side / 3⁄4 angle, **no white on the
character** (matches the `feedback_no_white_in_prompts.md` rule — the
engine's color-key strips white and would erase any white shirt detail
along with the BG). Background is plain pure-white `#FFFFFF` (engine
color-keys it on load).

Style lock: cartoon outlines, flat fills, **no pixel art**, no
gradients. Match the existing Tommy / Jake idle sheets at
`assets/images/locations/camp/npc/kids/{tommy,jake}/*_idle.png` —
same head shape, same outline weight, same palette.

### §20 Tommy walk-left

**Path:** `assets/images/locations/camp/npc/kids/tommy/npc_tommy_walk_left.png`
**Grid:** 8 columns × 1 row, **each cell 172 wide × 384 tall**
(sheet 1376×384).

> Cartoon child, exactly the same character as in
> `npc_tommy_idle.png` — tan skin, tousled medium-brown hair, GREEN
> short-sleeved t-shirt with a small pine-tree silhouette on the
> chest, NAVY-blue rolled-cuff jeans, BAREFOOT. Slim 8-year-old build,
> friendly excitable face, dimples. **Do not redesign him — match the
> idle sheet exactly.** Cartoon outlines, flat fills, NO pixel art,
> NO gradients.
>
> 8-frame walk-cycle strip, view from the SIDE, **facing LEFT** (he
> is leaving the scene toward the left edge). Each cell shows the
> same character in a different step of a natural walk cycle:
> 1. Right foot forward, left foot back (contact pose)
> 2. Right knee bent, body rising over right foot (down pose)
> 3. Right leg straight, body at peak height (passing pose)
> 4. Left foot lifting, right pushing off (high pose)
> 5. Left foot forward, right foot back (contact mirror)
> 6. Left knee bent, body rising over left foot
> 7. Left leg straight, body at peak height
> 8. Right foot lifting, left pushing off
>
> Arms swing in opposition to legs. Hair has subtle motion — slight
> trailing tilt opposite to the walk direction. Eyes looking forward
> (to the left). Slight forward lean to read as "walking away."
>
> Strict colors matching the idle sheet — green tee `#5BA84A`, navy
> jeans `#2B3A5C`, tan skin `#E0B894`, BLACK outlines. NO white pixels
> on the character (the eye whites in idle are off-white `#F3E9DE`
> already — preserve that). Pine-tree chest motif is a small darker-
> green silhouette. Background is solid pure white `#FFFFFF` (engine
> strips it).

### §21 Jake walk-up (back-walk toward his cabin)

**Path:** `assets/images/locations/camp/npc/kids/jake/npc_jake_walk_back.png`
**Grid:** 8 columns × 1 row, **each cell 172 wide × 384 tall**
(sheet 1376×384).

> Cartoon child, exactly the same character as in `npc_jake_idle.png`
> — peach-tan skin with light freckles, RED baseball cap worn
> straight, NAVY t-shirt, mid-blue rolled-cuff jeans, scuffed brown
> ankle sneakers. Stocky 9-year-old build, slightly tough/proud
> expression. **Do not redesign — match idle exactly.** Cartoon
> outlines, flat fills, NO pixel art, NO gradients.
>
> 8-frame walk-cycle strip, view from BEHIND (back of head and
> shoulders facing the camera). Jake is walking **UP into the frame**
> (away from camera, toward his cabin door). Cell sequence is a
> standard back-walk cycle — same contact / down / passing / high
> beats as a side walk, but the body silhouette is now back-of-head
> + back-of-shoulders + heels lifting alternately. Hat brim visible
> in profile silhouette. Shoulders shift left/right with each step.
>
> Strict colors matching the idle sheet — red cap `#C53A2C`, navy
> tee `#1F2D4E`, mid-blue jeans `#4A6FA0`, brown sneakers `#7A4E2B`,
> BLACK outlines. NO white pixels on the character. Background is
> solid pure white `#FFFFFF`.

### Wiring after the sheets land

Both sheets are loaded as `oneShotAnims["walk_away"]` on Tommy / Jake
respectively. The kid factory in [game/npc.go](game/npc.go) registers
them on construction; the `onDialogEnd` callback in
`setupCampCallbacks` ([game/game.go](game/game.go)) plays the one-shot
and translates the NPC's `bounds.X / bounds.Y` over the anim duration
so the kid visibly walks off:

- **Tommy** — one-shot for ~2.5 s while `bounds.X` lerps from his
  authored 130 to −180 (off-screen left). At the end, set
  `n.hidden = true` so he doesn't pop back on re-entry.
- **Jake** — one-shot for ~2.0 s while `bounds.Y` lerps from 405 to
  220 (up into the cabin doorway band). Set `n.hidden = true` at
  the end.

Both `hidden` flips can be reverted by the Day 2 trigger if we want
them back in their rooms; track in FIXME for the wire-up PR after the
PNGs land.

---

## NEW 2026-05-20 — FIXME sweep prompts (sizes, animations, items, story)

A batched set of regenerations and net-new sprites for the 2026-05-20 polish
pass. Code paths to drop each PNG into are listed inline. Save with the exact
filenames; the engine picks them up on next launch.

### A. PP talk-front canvas re-lock — regen

**Canvas:** 1600×540. **Grid:** 8×2. **Cell:** 200×270.
**Path:** `assets/images/player/PP talk front.png`.

Existing sheet has uneven per-frame canvases — PP's body shifts left/right
between cells and one frame renders taller than the others, producing the
"two frames at one" and "PP jumps from his place" report. Regen with a
locked baseline.

> [style lock]
>
> Pink Panther standing facing the camera, talking. Long lean cat-like
> silhouette, classic Pink Panther proportions, curved tail trailing
> behind him on the floor, yellow gloves, no shoes (bare panther feet),
> ivory off-white belly patch, pale grey eye sclera with black pupils,
> small smirk-capable mouth. Pink body fill `#E88BB5` with `#C4548A` cel
> shadow, black `#1C1C1C` ink outline ~3 px. **No pure white anywhere on
> the panther.**
>
> **Animation:** 16-frame talking loop (8 columns × 2 rows). Mouth open
> with subtle shape changes between frames, eyes alive. Half-body gestures
> only — feet stay planted, body shifts at hips only.
>
> Row 0 (frames 0–7): 0 neutral mouth open, 1 right paw raised palm-up,
> 2 shrug both shoulders, 3 right paw points forward, 4 sly grin, 5 both
> paws up, 6 left paw finger to chin (thinking), 7 confident closed-mouth
> smile.
>
> Row 1 (frames 8–15): same eight gestures repeated with eyes blinked
> half-closed and mouth half-shape for natural variation as the loop
> continues.
>
> **Critical canvas rules:** Every single cell is exactly 200 px wide ×
> 270 px tall. PP's center-of-body is on the cell's vertical centerline
> in every frame. PP's foot pixels rest on row pixel 268 in every frame —
> no frame is taller or shorter than another. Tail stays inside the cell;
> raised paws stay inside the cell. No bleed across cell borders. Pure
> white `#FFFFFF` background, zero scenery.

### B. PP talk-side canvas re-lock — regen

**Canvas:** 1600×540. **Grid:** 8×2. **Cell:** 200×270.
**Path:** `assets/images/player/PP talk side.png`.

Side-profile companion to §A. Same canvas locking issue, same fix.

> [style lock]
>
> Pink Panther in **side profile facing right**, talking. Long cat-like
> silhouette in three-quarter side view, tail curving behind him along
> the floor, yellow gloves, pink body `#E88BB5` with `#C4548A` cel
> shadow, ivory off-white belly, pale grey eye sclera, black `#1C1C1C`
> ink outline ~3 px. **No pure white on the panther** — use off-white
> for the belly and pale grey for the eye whites.
>
> **Animation:** 16-frame talking loop (8×2). Mouth open with subtle
> shape changes. Same eight gesture beats as PP-talk-front but performed
> in profile — gestures use the near-side arm with the far-side arm only
> partially visible. Body shifts at the hips, feet planted.
>
> Row 0 (frames 0–7): 0 neutral mouth open looking forward, 1 near paw
> raised palm-up, 2 shrug with shoulders up, 3 near paw points forward,
> 4 sly grin, 5 both paws up, 6 near paw finger to chin, 7 confident
> closed-mouth smile.
>
> Row 1 (frames 8–15): same eight gestures with eyes half-closed / mouth
> half-shape.
>
> **Critical canvas rules:** Every cell is exactly 200×270. PP's
> standing-foot baseline locked to row pixel 268 in every frame. Tail
> stays inside the cell. No frame is taller or wider than another. Pure
> white `#FFFFFF` background.

### C. PP grab-flower regen

**Canvas:** 900×235. **Grid:** 6×1. **Cell:** 150×235.
**Path:** `assets/images/player/PP grab flower.png`.

Current sheet shows the same uneven-cell problem as the talk sheets
and the flower itself is barely visible. Regen with the action staged
clearly.

> [style lock]
>
> Pink Panther bending down to pick a daisy from the ground, then
> standing back up holding it. Side three-quarter view facing slightly
> right. Long cat-like silhouette, yellow gloves, pink body `#E88BB5`
> with `#C4548A` cel shadow, ivory off-white belly patch, black
> `#1C1C1C` ink outline ~3 px. **No pure white on the panther.**
>
> **Animation:** 6-frame one-shot (no loop).
>
> - Frame 1: standing upright, neutral, arms relaxed at sides, looking
>   down toward a small daisy on the ground in front of him.
> - Frame 2: knees bend slightly, body lowers, near arm starts reaching
>   down. Mouth concentrating (closed).
> - Frame 3: deep crouch, near hand at ground level, fingertips just
>   above the daisy.
> - Frame 4: hand closes around the daisy stem. Daisy clearly visible
>   between his fingers — small ivory `#F2EFE5` petals (NOT white) with
>   a bright yellow `#F4C842` center, dark green `#3A5A3A` stem.
> - Frame 5: rising, knees half-extended, daisy held forward at chest
>   level. Small smile begins.
> - Frame 6: fully upright, daisy held close to chest in both hands,
>   pleased smile, eyes warm.
>
> **Critical canvas rules:** Every cell exactly 150×235. PP's foot
> baseline locked to row pixel 232 across all six frames — he does not
> change vertical position, only bends and rises in place. Pure white
> `#FFFFFF` background.

### D. Marcus talk sheet — regen at 7×2 to match idle

**Canvas:** 1372×768. **Grid:** 7×2. **Cell:** 196×384.
**Path:** `assets/images/locations/camp/npc/kids/marcus/npc_marcus_talk.png`.

Current talk sheet is 8 columns while idle is 7 — different cell width
makes talk render visibly bigger and shifts the body left/right between
idle and talk. Regen at 7×2 so cells match idle exactly.

> [style lock]
>
> Same Marcus as `npc_marcus_idle.png`. 10-year-old boy: messy
> light-brown hair `#4A2E1B` with stray strands, fair skin, round
> wire-rim glasses `#1C1C1C` frames, brown eyes, small white
> nose-bandage strip `#C4412A` (cartoon "boo-boo" plaster), light tan
> short-sleeve polo `#D8B47A`, olive-green cargo shorts `#4F7A3A`, white
> socks, brown lace-up shoes `#5A3A1E`. Always holds a small spiral
> sketchbook in his left hand. Black `#1C1C1C` ink outline ~3 px. **No
> pure white on the character** — sketchbook page uses off-white
> `#F2EFE5`, socks use cream `#E5DDC8`.
>
> **Animation:** 14-frame talking loop (7 columns × 2 rows). Mouth open
> with shape changes between frames, eager expression — Marcus is a
> know-it-all who loves explaining things.
>
> Row 0 (frames 0–6): 0 neutral mouth open looking at camera, 1 right
> hand raised explaining (one finger up), 2 small shrug, 3 both arms
> spread outward enthusiastically, 4 points at the sketchbook page,
> 5 pushes glasses up the bridge of his nose, 6 head-tilt explaining.
>
> Row 1 (frames 7–13): 7 small head-shake, 8 thinks (finger to chin),
> 9 closed-mouth grin between sentences, 10 both palms up "see what I
> mean?", 11 hand on hip explaining, 12 surprised "ah!" eyebrows-up,
> 13 returns to neutral matching frame 0 silhouette so the loop seams.
>
> **Critical canvas rules:** Every cell exactly 196×384. Marcus's foot
> baseline locked to row pixel 380 in every frame. Body centerline
> matches the idle sheet's body centerline. Sketchbook visible in every
> frame, held at hip height except frame 4. Pure white `#FFFFFF`
> background.

### E. Tommy walk-left sheet — new

**Canvas:** 1160×175. **Grid:** 8×1. **Cell:** 145×175.
**Path:** `assets/images/locations/camp/npc/kids/tommy/npc_tommy_walk.png`.

For the post-dialog beat where Tommy walks off the camp grounds to the
left. Style must match the existing `npc_tommy_idle.png` exactly so the
swap reads as the same kid.

> [style lock]
>
> Same Tommy as `npc_tommy_idle.png`. 9-year-old boy: floppy
> sandy-blond hair `#D8B47A`, fair skin with light freckles across the
> nose-bridge `#B06A3A`, pastel-blue short-sleeve t-shirt `#4A7AA8`,
> brown cargo shorts `#8B5A2B`, cream sneakers `#E5DDC8` (NOT white),
> small navy backpack `#3A5A7A` over both shoulders. Black `#1C1C1C`
> ink outline ~3 px. **No pure white on the character.**
>
> **Animation:** Side-view 8-frame walk cycle facing LEFT (engine
> mirrors for right). Standard cartoon gait — contact, down, passing,
> high — repeated for left and right legs. Natural arm swing,
> opposite-to-leading-leg. Backpack sways with the shoulders. Mouth
> closed, neutral expression, eyes forward in the direction of travel.
>
> - Frame 1: left-foot contact (heel down), right-arm forward.
> - Frame 2: left-foot down weight, body lowers ~3 px.
> - Frame 3: passing pose, right foot lifting past left.
> - Frame 4: right-leg high (knee up), left-arm forward.
> - Frame 5: right-foot contact, left-arm forward.
> - Frame 6: right-foot down weight, body lowers ~3 px.
> - Frame 7: passing pose, left foot lifting past right.
> - Frame 8: left-leg high (knee up), right-arm forward. Loops back to 1.
>
> **Critical canvas rules:** Every cell exactly 145×175. Tommy's
> contact-frame foot baseline locked to row pixel 172. Body proportions
> identical to the idle sheet (same head size, same shirt fit). Pure
> white `#FFFFFF` background.

### F. Jake walk-left sheet — new (must match Jake idle design)

**Canvas:** 1160×175. **Grid:** 8×1. **Cell:** 145×175.
**Path:** `assets/images/locations/camp/npc/kids/jake/npc_jake_walk.png`.

For the post-dialog beat where Jake walks back into his cabin. The
current walk PNG (if any) reads as a different character — regenerate
so it looks like the same Jake from idle.

> [style lock]
>
> **Match Jake's identity from `npc_jake_idle.png` exactly.** 9-year-old
> tough-kid: red baseball cap `#C4412A` worn backwards (brim points
> away from his face), dark-brown short hair peeking from the cap edges
> `#4A2E1B`, fair skin with light freckles `#B06A3A`, navy short-sleeve
> t-shirt `#2B3A5C`, denim shorts `#4A5A78`, cream sneakers `#E5DDC8`
> (NOT white). Black `#1C1C1C` ink outline ~3 px. **No pure white on
> the character.**
>
> **Animation:** Side-view 8-frame walk cycle facing LEFT. Standard
> cartoon gait — contact, down, passing, high — repeated for both legs.
> Jake walks with a slight forward lean and loose-fist hands (he's the
> tough kid, not the marching kid). Cap brim stays facing back through
> the whole cycle. Mouth closed, eyes forward, determined expression.
>
> - Frame 1: left-foot contact, right-arm forward in loose fist.
> - Frame 2: left-foot down weight, shoulders dip ~3 px.
> - Frame 3: passing, right foot lifting.
> - Frame 4: right-leg high, left-arm forward.
> - Frame 5: right-foot contact, left-arm forward.
> - Frame 6: right-foot down weight, shoulders dip.
> - Frame 7: passing, left foot lifting.
> - Frame 8: left-leg high, right-arm forward. Loops back to 1.
>
> **Critical canvas rules:** Every cell exactly 145×175. Contact-frame
> foot baseline locked to row pixel 172. Same proportions and ink
> weight as the idle sheet — same Jake, just walking. Pure white
> `#FFFFFF` background.

### G. Madame Colette talk sheet — regen-if-needed

**Canvas:** 1680×936. **Grid:** 8×2. **Cell:** 210×468.
**Path:** `assets/images/locations/paris/npc/outside/npc_french_guide_talk.png`.

The Go loader has already been corrected to read this sheet as 8×2
(was reading 8×1, skipping the bottom row, which produced the
"two frames at one" report). Only regenerate if the visual issue
persists after the data fix.

> [style lock]
>
> Same Madame Colette as `npc_french_guide_idle.png`. Mid-50s Parisian
> tour guide, warm and chatty. Short brown bob `#7B4A2E`, light olive
> skin, red lipstick `#C4412A`, navy blazer `#2B3A5C` over an
> off-white blouse `#F2EFE5` (NOT pure white), beige knee-length skirt
> `#C9A878`, dark brown low-heel shoes `#5A3A1E`. Holds a small folded
> Paris map in her left hand. Round tortoiseshell reading glasses
> `#5A3A1E` perched on top of her head. Black `#1C1C1C` ink outline
> ~3 px. **No pure white on the character.**
>
> **Animation:** 16-frame talking loop (8×2), mouth open with shape
> changes, warm and animated gestures.
>
> Row 0 (frames 0–7): 0 neutral mouth open, 1 right hand gestures
> outward palm-up, 2 left hand joins right (both palms up), 3 points
> off to the right (toward the museum), 4 brings folded map up to
> chest, 5 unfolds one corner showing it to camera, 6 head-tilt charm
> with closed-mouth smile, 7 returns to neutral.
>
> Row 1 (frames 8–15): 8 small shrug, 9 hand on hip, 10 laughs (eyes
> half-closed), 11 points at her own glasses on head, 12 both hands
> sweep wide, 13 finger-wag friendly correction, 14 head-tilt the
> other way, 15 loops back toward frame 0 silhouette.
>
> **Critical canvas rules:** Every cell exactly 210×468. Colette's foot
> baseline locked to row pixel 462 in every frame. Body centerline
> matches her idle sheet. Pure white `#FFFFFF` background.

### H. PP airplane sheet — regen for fuselage-locked centerline

**Canvas:** 1596×984. **Grid:** 6×2. **Cell:** 266×492.
**Path:** `assets/images/player/pp_airplane.png`.

The current sheet has the plane drawn at different vertical positions
in row 0 vs row 1, which made the cutscene bounce vertically. The
engine compensates with a per-row offset for now; once this regen
lands the offset can be deleted.

> [style lock]
>
> Pink Panther piloting a small yellow-and-blue biplane, viewed from the
> front three-quarter angle. Cartoon biplane with two stacked wings:
> bright sky-blue `#3A8AC0` fuselage and tail, sunny-yellow `#F4C842`
> upper and lower wings, red `#C4412A` trim stripe along the fuselage,
> mid-grey `#8A8A8A` propeller arc spinning at the nose, dark navy
> `#1F2D4E` cockpit rim. Pink Panther in the cockpit visible from chest
> up — pink body `#E88BB5`, aviator goggles with cream `#E5DDC8` straps
> (NOT white). Black `#1C1C1C` ink outline ~3 px. **No pure white on
> the plane or the panther.**
>
> **Animation:** 12-frame cycle (6 columns × 2 rows).
>
> Row 0 (frames 0–5): cruising level — slight wing tilt left, neutral,
> tilt right, neutral, with propeller arc rotating through 2 prop-blur
> variants. PP smiles confidently.
>
> Row 1 (frames 6–11): same cruising motion repeated for animation
> variety — slight wing wobble, prop variants, PP gives a thumbs-up on
> one frame.
>
> **Critical canvas rules — this is the whole point of the regen:**
> The biplane's **fuselage centerline must sit at exactly row pixel 220
> in every single cell**. Wings tilt up and down at the tips, but the
> fuselage horizontal axis does not move between frames. The plane's
> body is centered horizontally in each cell. No vertical drift between
> row 0 and row 1. Pure white `#FFFFFF` background.

### I. Higgins throw-map sheet — new

**Canvas:** 1092×235. **Grid:** 6×1. **Cell:** 182×235.
**Path:** `assets/images/locations/camp/npc/higgins/npc_director_higgins_throw_map.png`.

For the new "Higgins throws the map to PP" beat in the office. Higgins
is seated behind his desk and tosses the rolled map across the room.

> [style lock]
>
> Same Higgins as `npc_director_higgins_office_idle.png`. Seated behind
> a wooden desk, only torso and head visible from desk-edge up. Lanky
> ranger build, round wire-rim glasses `#1C1C1C`, small brown mustache,
> short side-parted brown hair `#4A2E1B`, khaki safari shirt `#A47148`
> with epaulets, brown leather belt visible, cream `#E5DDC8` neckerchief
> (NOT white). Black `#1C1C1C` ink outline ~3 px. **No pure white on
> the character.**
>
> **Animation:** 6-frame one-shot (no loop), the windup-throw-release
> beat.
>
> - Frame 1: holds a rolled-up parchment map at chest height in his
>   right hand. Looks at PP off-camera. Mouth open mid-sentence.
> - Frame 2: pulls the map back over his right shoulder, body rotates
>   slightly right, eyes track forward toward the throw target.
> - Frame 3: full wind-back peak — arm bent at the elbow behind the
>   shoulder, body coiled.
> - Frame 4: forward throw motion — arm comes through, body uncoils
>   toward the camera, map still in hand mid-arc.
> - Frame 5: release — map leaves the hand, arm fully extended forward
>   pointing toward the catch. The rolled map is visible just beyond
>   his fingertips heading off-frame to the right.
> - Frame 6: follow-through — empty hand still pointing forward, body
>   settled back to the seated posture, satisfied smile.
>
> **Critical canvas rules:** Every cell exactly 182×235. Higgins's
> seated-shoulder-line locked to row pixel 90 in every frame (no
> vertical drift). The map (cream parchment `#E5DDC8`, brown ribbon
> `#8B5A2B`, NOT white) is clearly visible in his hand on frames 1–4
> and just departing in frame 5. Pure white `#FFFFFF` background.

### J. PP catch-map sheet — new

**Canvas:** 1200×270. **Grid:** 6×1. **Cell:** 200×270.
**Path:** `assets/images/player/pp_catch_map.png`.

PP's half of the throw/catch choreography. Plays right after §I.

> [style lock]
>
> Pink Panther standing facing camera, catching a thrown rolled map
> from the right. Standard PP build, yellow gloves, pink body `#E88BB5`
> with `#C4548A` cel shadow, ivory off-white belly, pale grey eye
> sclera, black `#1C1C1C` ink outline ~3 px. **No pure white on the
> panther.**
>
> **Animation:** 6-frame one-shot (no loop).
>
> - Frame 1: relaxed neutral stance, arms at sides, looking off to the
>   right toward incoming map. Curious raised eyebrow.
> - Frame 2: arms start coming up, body shifts weight onto right foot
>   in anticipation.
> - Frame 3: both hands extended forward at chest height, palms open
>   and angled to receive the catch.
> - Frame 4: hands clasp around the rolled map — map visible between
>   his palms (cream parchment `#E5DDC8`, brown ribbon `#8B5A2B`, NOT
>   white). Body still leaning slightly right.
> - Frame 5: pulls the map back toward his chest, body returns to
>   center, looks down at it.
> - Frame 6: holds the map at chest height in both hands, examining it
>   with a small pleased smile. Loops/holds.
>
> **Critical canvas rules:** Every cell exactly 200×270. PP's foot
> baseline locked to row pixel 268 in every frame — PP does not bob or
> drift vertically during the catch, only his arms move. Body
> centerline same as `PP idle front.png`. Pure white `#FFFFFF`
> background.

### K. Thrown-map projectile — new

**Canvas:** 80×40. **Grid:** 1×1.
**Path:** `assets/images/items/inv_travel_map_throw.png`.

The rolled-map sprite that tweens across the screen between §I and §J.

> [style lock]
>
> A single rolled-up parchment map / scroll viewed at a slight diagonal
> tilt, as if frozen mid-throw. Cream parchment body `#E5DDC8` (NOT
> white) with darker `#B89868` shadow along one edge to give it form,
> two visible roll-edges showing the scroll layers, a brown ribbon
> `#8B5A2B` wrapped around the middle with the ends fluttering slightly
> back from the motion. Black `#1C1C1C` ink outline ~3 px. Pure white
> `#FFFFFF` background.

### L. Travel-map inventory icon — regen for clarity

**Canvas:** 96×96. **Grid:** 1×1.
**Path:** `assets/images/items/travel_map_icon.png`.

Current icon reads as abstract — needs to clearly read as a travel map
in the inventory bar.

> [style lock]
>
> A rolled parchment map, partially unrolled to show a fragment of its
> contents, viewed from a three-quarter angle (top-down-ish). The
> unrolled flap shows simplified continent outlines in navy `#2B3A5C`
> with a dotted flight-path line `#C4412A` looping across them and a
> small biplane silhouette on the path. Cream parchment `#E5DDC8` body
> (NOT white) with darker `#B89868` shadow under the curl, brown
> ribbon `#8B5A2B` around the rolled portion. Black `#1C1C1C` ink
> outline ~3 px. Pure white `#FFFFFF` background.

### M. Unified action cursor — pointing-hand, cursor_grab style  *(DONE — landed 2026-05-21)*

**Canvas:** 1677×938 (landscape, matches `cursor_grab.png` aspect).
**Grid:** 1×1 (single frame, no animation).
**Path:** `assets/images/ui/cursors/cursor_point.png`.

One sprite reused for two purposes: hovers above PP to indicate
"click to open inventory", and shown on travel-map pins to indicate
"this city is the current quest target". Engine tints it at draw
time (yellow for travel-relevant, cream for info / inventory open).
The landed sprite matches the `cursor_grab.png` pixel-art hand family
exactly. Prompt below preserved for future regen / variants.

> [style lock — pixel-art cursor family]
>
> A chunky cartoon hand, viewed from the back of the hand, **index
> finger extended pointing toward the upper-left of the canvas** (as
> if the cursor is pointing at something above and to the left of the
> wrist). The hand reads as a classic adventure-game "look at this"
> pointer. Match the style of `cursor_grab.png` exactly:
>
> - Pixel-art rendering with visibly stepped (staircased) edges, NOT
>   smooth anti-aliased curves
> - Bright magenta-pink fill `#ED1A8F` filling the whole hand silhouette
> - Thick black `#1C1C1C` pixel-art outline, chunky weight (roughly
>   16–20 source pixels at the 1677×938 canvas so it reads as a clear
>   ~2–3 px outline once downscaled to cursor size)
> - Subtle darker `#2A2A2A` drop shadow offset down-right ~14 px,
>   creating a slight floating feel
> - No inner shading, no gradients, no skin tone, no fingernails
> - Solid silhouette only — the hand is a flat magenta shape with a
>   black outline, nothing more
>
> **Pose details (matches the landed file):**
>
> - Index finger fully extended, pointing toward the upper-left at
>   roughly 30–35° above horizontal (fingertip is the highest and
>   leftmost point of the silhouette)
> - Middle, ring, and pinky fingers are curled into the palm — their
>   knuckles read as a stack of small bumps along the lower edge of
>   the hand
> - Thumb tucked along the side of the hand (small visible bump on
>   the upper-right of the palm silhouette)
> - Back of the hand facing the viewer; wrist visible at the lower-
>   right of the canvas as a stubby flat cuff
> - Overall silhouette reads instantly as "hand pointing up and to
>   the left"
>
> **Canvas rules:**
>
> - 1677×938 source canvas (landscape, same aspect as
>   `cursor_grab.png` so the cursor family scales uniformly)
> - Hand fills roughly 70% of the canvas area, biased toward the
>   left two-thirds
> - **Fingertip (the action hotspot) sits at roughly (canvas 15%,
>   15%)** — upper-left of the canvas — so the natural mouse hotspot
>   is the tip of the index finger
> - Pure white `#FFFFFF` background, zero scenery, single hand only
> - **No pure white anywhere on the hand itself.** The engine
>   chroma-keys white as transparent, so any white pixels inside the
>   silhouette will become holes. The hand is solid magenta-pink
>   `#ED1A8F` with a black `#1C1C1C` outline and a `#2A2A2A` drop
>   shadow — no highlights, no inner white
>
> Reference: open `assets/images/ui/cursors/cursor_grab.png` for the
> exact pink shade, outline weight, and drop-shadow treatment, and
> `assets/images/ui/cursors/cursor_point.png` for the exact pointing
> pose. Same hand identity across the whole cursor family.

### N. Item sprites — 8 new (story anchors + Paris quest props)

**Canvas:** 96×96 each. **Grid:** 1×1.
**Paths:** in `assets/images/items/` (see each prompt for filename).

Each item is currently declared in `items.json` but renders as a wrong
placeholder. Eight quick prompts to generate them all.

> [style lock]
>
> A single wooden rolling pin lying at a slight diagonal across the
> frame, as if a baker just set it down. Pale cream wood body
> `#E5C895` (NOT white) with darker `#B89868` cel shadow along the
> underside, brown `#8B5A2B` handles at both ends with visible turning
> rings. A faint dusting of flour particles `#D8C8B0` underneath.
> Black `#1C1C1C` ink outline ~3 px. Pure white `#FFFFFF` background.
> Save as `rolling_pin.png`.

> [style lock]
>
> A single crusty Parisian baguette lying at a slight diagonal across
> the frame, still warm and fresh from the oven. Golden-brown crust
> `#B97E3F` with darker `#8A5828` cel shadow along the underside, four
> or five diagonal slash marks across the top in slightly lighter
> `#D9A26A`, one end torn open showing pale `#E5C895` interior crumb.
> Two tiny steam wisps in pale grey `#C4C4C4` rising from the open end.
> Black `#1C1C1C` ink outline ~3 px. Pure white `#FFFFFF` background.
> Save as `baguette.png`.

> [style lock]
>
> A laminated photographer's press pass hanging from a red lanyard.
> Cream `#E5DDC8` card body (NOT white) with a thin black border, a
> small black-and-white head-shot photo in the top-left, the word
> "PRESS" hand-lettered in bold dark-red `#9E2A1F` across the middle.
> Red `#C4412A` woven lanyard looped through a silver-grey `#B0B0B0`
> clip at the top of the card, lanyard ends visible draping back. Card
> tilted slightly. Black `#1C1C1C` ink outline ~3 px. Pure white
> `#FFFFFF` background. Save as `press_pass.png`.

> [style lock]
>
> A Louvre postcard viewed from the front. Cream `#E5DDC8` postcard
> body (NOT white) with a small darker `#B89868` shadow on the
> back-right corner showing the card has a slight curl. Front face
> shows a simplified illustration of the Louvre glass pyramid: pale
> sky-blue `#9FC9E8` upper sky, the pyramid in navy `#2B3A5C` line-
> work with pale grey `#C4D8E8` glass panels, a hint of the museum
> building behind it in warm beige `#C9A878`. "PARIS" hand-lettered
> across the bottom of the postcard in bold dark-red `#9E2A1F`. Black
> `#1C1C1C` ink outline ~3 px. Pure white `#FFFFFF` background. Save
> as `postcard.png`.

> [style lock]
>
> A folded piece of parchment with a coin-rubbing imprint visible on
> its surface. Cream `#E5DDC8` paper (NOT white) with a soft darker
> `#B89868` shadow under the fold-crease, one corner of the paper
> torn jagged at the top-right. The rubbing itself is a roughly
> circular graphite-grey `#444444` smudge in the center of the
> visible face, with faint Hebrew-style markings inside the circle.
> Black `#1C1C1C` ink outline ~3 px around the paper. Pure white
> `#FFFFFF` background. Save as `coin_rubbing.png`.

> [style lock]
>
> A single pressed cherry-blossom flower flattened between two strips
> of pale washi paper. Five soft pink petals `#F4B6C2` with darker
> pink veins `#C97186` radiating from a small yellow `#F4C842` center,
> a thin dark-green `#3A5A3A` stem trailing off one side. Two cream
> washi paper strips `#E5DDC8` (NOT white) cross over the blossom
> diagonally, with rough deckled edges and faint fiber lines visible.
> Black `#1C1C1C` ink outline ~3 px. Pure white `#FFFFFF` background.
> Save as `pressed_sakura.png`.

> [style lock]
>
> Half of a carnival dance card, torn vertically on its right edge so
> only the left half is visible. Sunny yellow `#F4C842` card body
> with a bold red `#C4412A` decorative border, music notes
> `#1C1C1C` printed across the surface, a single feathered samba
> plume `#5A8A4A` sprouting from the top-left corner, a small
> gold-brown `#B98A28` tassel hanging from the bottom-left. The torn
> right edge has a jagged tear pattern. Black `#1C1C1C` ink outline
> ~3 px. Pure white `#FFFFFF` background. Save as `dance_card.png`.

> [style lock]
>
> A folded piece of cream parchment paper with a Roman-style
> inscription rubbing visible on its surface. Cream `#E5DDC8` paper
> (NOT white) with a darker `#B89868` shadow at the fold-crease, one
> corner curled up. The rubbing shows blocky Roman capitals (S P Q R
> or similar) in graphite grey `#444444`, with a faint chalk smudge
> at the bottom edge of the rubbing. Black `#1C1C1C` ink outline ~3 px
> around the paper. Pure white `#FFFFFF` background. Save as
> `inscription_rubbing.png`.
