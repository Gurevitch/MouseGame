# Pink Panther: Camp Chilly Wa Wa — Story Flow

## Overview
Pink Panther arrives at Camp Chilly Wa Wa as a substitute counselor. The kids start experiencing strange visions of real-world cities they've never visited. PP must travel to each city, find the source of their visions, and bring back an anchor object to help each child.

---

## DAY 1 — Arrival (Normal)

```
Camp Entrance
    |
    v
[Opening monologue plays]
    |
    v
[Auto-walk to Director Higgins]
    |
    v
Higgins: "Welcome! Go meet the kids."
    |
    v (arrow UP — walk up the road)
    |
Camp Grounds -----> Higgins' Office (right arrow)
    |
    v
Meet all 5 kids (normal personalities):
  - Marcus: Know-it-all, loves sketching
  - Tommy: Storyteller, music lover  
  - Jake: Tough kid, coin collector
  - Lily: Shy, loves flowers
  - Danny: Prankster, treasure stories
    |
    v (after talking to all 5)
    |
[Higgins arrives: "It's getting late!"]
    |
    v
```

## NIGHT — The Campfire (Multi-phase cutscene)

```
[Higgins bedtime dialog at camp_grounds]
Higgins: "Attention everyone! It's getting very late."
Higgins: "All campers to their cabins. NOW."
    |
    v
[Transition to camp_night scene]
    |
    v
[PP falls asleep by campfire — sleeping animation + campfire_idle animation]
[~3.5 seconds of sleeping]
    |
    v
[Transition to marcus_room (night background)]
[Marcus freaks out for 6-8 seconds]
Marcus: "No no no... the lines won't stop..."
Marcus: "A woman's face... golden frames everywhere..."
Marcus: "A GLASS PYRAMID! And the building is ENORMOUS!"
Marcus: "The painting is WRONG! Something is hidden underneath!"
Marcus: "I have to draw it... I HAVE to draw it ALL!"
    |
    v
[Transition back to camp_night]
[PP waking animation plays ~2 seconds]
    |
    v
PP: "*yawn* What a night..."
PP: "I heard Marcus freaking out. Something about paintings and pyramids."
PP: "I need to check on him."
    |
    v
```

## DAY 2 — The Weirdness

```
[Transition to camp_grounds]
    |
    v
[Day 2 monologue: "I need to find Marcus..."]
    |
    v
Navigate to Marcus's Cabin (arrow UP hotspot)
    |
    v
Marcus Room (day background, Marcus has strange sprite)
    |
    v
Talk to Marcus (strange dialog):
  "A woman's face... ornate golden frames... a glass pyramid..."
  "The biggest museum in the world!"
    |
    v
Navigate to Higgins' Office (arrow RIGHT from camp_grounds)
    |
    v
Higgins: worried, gives Travel Map
    |
    v
[Paris unlocks on travel map]
    |
    v
Use Travel Map (from inventory, click anywhere) --> Globe opens
    |
    v
Click Paris --> [Airplane cutscene ~4 seconds, PP idle in plane] --> Fly to Paris
```

## PARIS — Marcus's Anchor (Louvre Postcard)

Paris is a multi-step pre-Louvre puzzle: PP can't just walk into the
museum, he has to assemble a press pass via three street NPCs first. This
matches the retro PP "collect props before the main door opens" rhythm
from *Hokus Pokus Pink* / *Passport to Peril* and gives the city its own
mini-quest before the anchor-object reveal.

### Cast (paris_street + paris_bakery + paris_louvre scenes)

| NPC | Scene | Role |
|-----|-------|------|
| Madame Poulain | paris_bakery (interior) | Baker; gives baguette in trade for her lost rolling pin |
| Pierre | paris_street (back of line, smaller in perspective) | Street artist; gives press pass in trade for baguette |
| Gendarme Claude | paris_street (right side, near Louvre steps) | Police officer; gives museum ticket in trade for press pass |
| Madame Colette | paris_street (left side) | Flavor guide; chats about Paris landmarks, points toward the museum |
| Nicolas | paris_street (mid) | Flavor NPC; chatty photographer near the Louvre, no items |
| Curator Beaumont | paris_louvre (interior) | Curator; identifies Marcus's painting, hands over the postcard |
| Madame Yvette | paris_bakery (left table) | Flavor — camp gossip about the Marcus drawings |
| Monsieur Bernard | paris_bakery (left table) | Flavor — newspaper headline about the museum restoration |
| Mademoiselle Camille | paris_bakery (mid table) | **MAIN QUEST** — "Camille and the Sold-Out Postcard": sketches the Room 7 replica Beaumont needs, once PP recovers her lucky pencil |
| Monsieur Henri | paris_bakery (mid table) | **QUEST** — asks PP for a Café au Lait; trades it for homemade Confiture from his bag |
| Lucien | paris_bakery (right table) | Flavor — foreshadows Lily / Tokyo arc |
| Madame Élise | paris_bakery (right table) | Flavor — warm encouragement, no quest info |

### Quest flow

```
[Airplane cutscene ~4 s, PP in biplane] -> Arrive paris_street
    |
    v
Talk to Madame Colette (flavor):
  - Eiffel Tower: built 1889, was temporary, Gustave Eiffel
  - Louvre: largest museum, 380,000 objects
  - Glass Pyramid: I.M. Pei, 1989
  - Points: "Ze museum is to ze right, monsieur"
    |
    v
Talk to Pierre (flavor pre-baguette):
  - Painting Pink Panthers on the easel
  - Hint: "Ze Curator knows every face in Paris. Ask her."
    |
    v
Talk to Nicolas (flavor):
  - Photographer near the museum, 20 years on the street
  - "Talk to Pierre ze painter and Claude ze gendarme"
    |
    v
Find the rolling pin OUTSIDE in the **bicycle basket** on paris_street
  (the black bike parked on the cobblestones — someone borrowed Madame
  Poulain's pin and dropped it in the wicker basket). Click the basket
  to grab it.
    |
    v
Try to enter Louvre (right arrow on paris_street)
  -> GATED. Need: Museum Ticket.
    |
    v
Enter paris_bakery (down/left arrow on paris_street)

Café patrons seated around the bakery — most are flavor only (clickable
in any order, no `|/v` chain). Henri carries a NEW quest beat:

  - Madame Yvette (beret + pearls): "Ze restoration is all anyone talks
    about. A hidden symbol under ze portrait — imagine!"
  - Monsieur Bernard (Le Figaro): rustles paper, repeats the Louvre
    restoration headline as confirmation of the museum rumor
  - Mademoiselle Camille (red beret, art student): SKETCH BEAT — asks
    to draw PP on the spot, shows him the finished sketch in-dialog,
    no inventory exchange (pure character beat)
  - **Monsieur Henri (silver mustache, croissant + bag):** asks PP to
    fetch a Café au Lait. Promises something nice from his bag in
    exchange. THIS IS THE NEW QUEST GATE for Pierre's press pass.
  - Lucien (gray turtleneck): cryptic Tokyo foreshadow — "a tower
    covered in flowers, bells ringing far away"
  - Madame Élise (auburn hair, autumn scarf): warm seasonal flavor
    line, no quest info

Talk to Madame Poulain (initial):
  - She's lost her rolling pin
  - "Find it and ze first baguette is yours"
    |
    v
Talk to Madame Poulain (alt dialog, holding Rolling Pin):
  -> Trade Rolling Pin → Baguette + **Café au Lait**
  -> Inventory: -Rolling Pin, +Baguette, +Café au Lait
  -> Poulain: "Take ze coffee to Henri, he's been waiting all morning."
    |
    v
Back to paris_street -> Talk to Pierre (alt dialog, holding Baguette):
  Stage 1 — Pierre: "Bread is good, but it is dry. Bring me a spread."
  -> Pierre TAKES the Baguette (inventory loses it) but does NOT hand
     over the press pass yet. Pierre's hintState advances to "waiting
     for spread".
    |
    v
Return to bakery -> Talk to Henri (alt dialog, holding Café au Lait):
  -> Henri remembers his promise: "Here, homemade strawberry confiture,
     made it zis morning."
  -> Trade Café au Lait → Confiture.
    |
    v
Back to paris_street -> Talk to Pierre (alt dialog, holding Confiture):
  Stage 2 — Pierre: "Strawberries from ze south — perfect."
  -> Trade Confiture → Press Pass.
    |
    v
Talk to Gendarme Claude (alt dialog, holding Press Pass):
  -> Trade Press Pass → Museum Ticket
  -> Press Pass is also consumed (consumed by Claude).
    |
    v
Right arrow on paris_street -> Enter paris_louvre (now ungated)
    |
    v
Talk to Curator Beaumont:
  - Mona Lisa: Leonardo da Vinci, 1503
  - Venus de Milo: Greek, 100 BC
  - PP describes Marcus's drawings (woman's face, golden frames,
    "something missing")
  - Curator: "Zat sounds like ze portrait in Room 7"
  - BUT the gift-shop postcards are SOLD OUT (the restoration news —
    ties into Yvette/Bernard's café gossip). Beaumont keeps the LAST
    one in the archive and will trade it for a replica sketch of the
    portrait for the archive wall: "Mademoiselle Camille at ze café —
    ze fastest charcoal in Paris."
    |
    v
Back OUTSIDE -> into the BAKERY -> Talk to Camille:
  - Thrilled Beaumont asked for HER — but she lost her lucky charcoal
    pencil sketching the museum at sunrise.
  - "Ask Nicolas, ze photographer by ze steps. Nothing happens on zat
    street without his lens seeing it."
    |
    v
Back OUTSIDE -> Talk to Nicolas (street):
  - "It rolled off ze curb — straight into ze flower pot by ze Louvre
    steps. Ze pigeons have been guarding it."
    |
    v
Try the flower pot -> BLOCKED: the pigeons defend it.
  PP: "Pierre speaks fluent pigeon... and after that baguette and
  confiture, he owes me a favor."
    |
    v
Talk to Pierre (altDialog favor beat — he owes PP after the 2-item
trade): he whistles + scatters crumbs by his easel, pigeons abandon
the pot. "Zey do ANYTHING for crumbs... except land on my canvas.
Critics!" (seeds the Pigeon Critic side quest gag)
    |
    v
Pick up the Charcoal Pencil (hidden floor item in the flower pot near
the Louvre entrance — cursor reveals it, PP grab animation plays)
    |
    v
Back INSIDE the bakery -> Hand the pencil to Camille:
  - Sketching one-shot plays (npc_camille_sketching.png, ends revealing
    the page) -> PP gets "Camille's Sketch" (Room 7 replica)
    |
    v
Back to the MUSEUM -> Hand the sketch to Beaumont:
  - "Straight to ze archive wall!" -> trades the last Postcard of the
    restored painting (ANCHOR OBJECT) + hints a commission for Camille
    |
    v
Use Travel Map -> Return to Camp Chilly Wa Wa
```

### Healing Marcus (back at camp)

```
Camp grounds, Day 2 (post-Paris)
    |
    v
Walk to Marcus's cabin (arrow UP)
    |
    v
Marcus in strange state, drawing endlessly
    |
    v
Give Postcard to Marcus (use item on NPC):
  - Marcus sees the complete painting
  - Realizes what was "missing" was a hidden symbol
  - Strange sprite swap OFF, normal sprite ON
  - VarStore: marcusHealed = true
  - Higgins office dialog pivots to "now Lily is the one I'm worried
    about" (postMarcusHealed branch)
    |
    v
Lily begins showing strange behavior (next chapter, see Tokyo)
```

### "Camille and the Sold-Out Postcard" (main-chain gate, 2026-06-10)

Reworked same day per user feedback — the original draft (pencil under a
sleeping Lucien) read too dark, and the user wanted the quest to ping-pong
between outside, inside and the museum on the way to the postcard. The
full flow is in the main chain diagram above: museum (sold out, Beaumont
asks) → bakery (Camille lost her pencil) → street (Nicolas saw it roll)
→ flower-pot pickup by the Louvre steps → bakery (sketch one-shot) →
museum (sketch → Postcard trade). Lucien reverted to his awake
Tokyo-foreshadow flavor dialog; Yvette now gossips that the postcards
sold out in a day, foreshadowing the gate.

Extras (user 2026-06-10): Pierre repays the baguette + confiture debt —
the pigeons guarding the flower pot only move when he whistles them off
(favor beat that also seeds the Pigeon Critic gag). And Camille plays
her sketching one-shot at the end of her FIRST chat, so the existing
npc_camille_sketching.png animation is visible from the start.

### Paris side quest: "The Pigeon Critic" (optional)

Retro pattern (encounter → block → hint → collect → use → reward); does
not gate the main chain. After the rolling-pin trade, Madame Poulain
runs "counter service": refills the Café au Lait while Henri's trade is
pending, donates the Baguette Heel, and accepts the Signed Postcard.

**"The Pigeon Critic" (paris_street, post-press-pass):**

```
After the Confiture → Press Pass trade, the next chat with Pierre
  becomes the ask: his masterpiece is done but no pigeon critic
  will land to approve it. "Crumbs, monsieur — Madame Poulain
  always has a stale heel for ze birds."
    |
    v
Ask Madame Poulain -> she donates the day-old Baguette Heel
  ("ze ends are for ze birds anyway").
    |
    v
Hand the Baguette Heel to Pierre -> crumbs scattered, a pigeon
  lands on the canvas, the painting is APPROVED. Pierre teaches
  the "plein air" / Monet beat and gives PP the "Mini Portrait"
  (keepsake; the pigeon posed for the background).
```

### Return-to-Paris souvenir loop (post-Marcus-healed, completed 2026-06-10)

Keeps the city alive after PP heals Marcus, and now closes end-to-end:

```
Return to the bakery (marcusHealed) -> Madame Poulain asks PP to
  bring a Louvre postcard for her grandson in Lyon
    |
    v
Back to the Louvre -> Curator Beaumont (altDialog) signs a SECOND
  postcard ("every collection needs a rare piece") -> "Signed Postcard"
    |
    v
Hand the Signed Postcard to Poulain -> she names the pink éclair in
  her window "Le Panthère Rose". Pure world-warmth; gates nothing.
```

### Paris perspective polish (2026-05-20)

- Pierre stands "back of line" in the street scene — smaller bounds
  (~95×175 vs the 135×230 of front-line NPCs) so the street reads as
  having depth
- PP's existing `depthScale` (driven by player.y) automatically shrinks
  PP as he walks up toward Pierre and restores his size as he walks
  back, matching the retro perspective rule
- All Paris NPCs feet anchored to street tile y≈720 (was floating
  at y≈630 before)

### Items used in this chapter

| Item | Source | Used on | Result |
|------|--------|---------|--------|
| Rolling Pin | Paris street (bicycle basket on the cobblestones) | Madame Poulain | Trade for Baguette + Café au Lait |
| Baguette | Madame Poulain | Pierre (stage 1) | Pierre keeps it, asks for a spread |
| Café au Lait | Madame Poulain | Monsieur Henri | Trade for Confiture |
| Confiture | Monsieur Henri | Pierre (stage 2) | Trade for Press Pass |
| Press Pass | Pierre | Gendarme Claude | Trade for Museum Ticket |
| Museum Ticket | Gendarme Claude | Louvre door (auto-consume) | Unlocks museum |
| Postcard | Curator Beaumont (trade for Camille's Sketch) | Marcus (back at camp) | Heals Marcus |
| Charcoal Pencil | Flower pot by the Louvre steps (hidden, Nicolas hints) | Mademoiselle Camille | Sketch one-shot → Camille's Sketch |
| Camille's Sketch | Mademoiselle Camille | Curator Beaumont | Trade for the Postcard (main chain) |
| Baguette Heel | Madame Poulain (counter service) | Pierre (post-press-pass) | Pigeon lands → Mini Portrait |
| Mini Portrait | Pierre | — (keepsake) | Side-quest reward |
| Signed Postcard | Curator Beaumont (post-heal, after Poulain asks) | Madame Poulain | Closes the grandson souvenir loop |

### Checklist (current status)

- [x] Paris street + bakery + Louvre scenes all rendering
- [x] Rolling-pin → baguette → press pass → museum ticket trade chain
- [x] Louvre entrance gated on Museum Ticket
- [x] Curator gives postcard
- [x] Pierre moved back-of-line with depth scaling
- [x] Madame Colette talk frames fixed (was reading 8×1 of an 8×2 sheet)
- [x] Paris NPCs grounded at street tile y≈720
- [x] Nicolas hit-radius shrunk so clicks don't bleed into Louvre exit
- [x] Post-Marcus-healed bakery beat (postcard-for-grandson) wired
- [ ] Give postcard to Marcus → Marcus calms down (Day 2 anchor handover wiring)
- [ ] Marcus auto-swaps to normal sprite after postcard
- [ ] Higgins office dialog confirmed to switch to postMarcusHealed branch in playtest
- [ ] Lily begins showing strange behavior (next chapter trigger)

---

## JERUSALEM — Jake's Anchor (Western Wall Coin Rubbing)

Jake's arc is **courage**. His nightmare is the face in the tunnels — a Roman
coin's portrait staring out of the Western Wall. Jerusalem is a multi-step retro
daisy-chain (like the Paris press-pass chain) of *tangential, slightly absurd*
errands — an alley cat, a tourist's sardines, a lost guidebook — that only pay
off in the anchor object (the Coin Rubbing) at the END. The note-in-the-wall
ritual is the FINAL task, and placing it is what unlocks the flight home.

**Status:** scenes + NPCs + dialog are built and wired; the trivial version
(Miriam hands the rubbing on first talk) is live. The full task chain below is
DESIGNED but not yet wired — see `docs/JERUSALEM_TASKS.md` for the build plan.

### Scenes & cast

Three scenes, entrance plaza is the hub:
- **jerusalem_entrance** (`wall_enterence.png`) — plaza where PP lands; LEFT arch
  → market, RIGHT → wall. Distant worshippers ambient.
- **jerusalem_wall** (`wall_close.png`) — up at the Western Wall; worshippers
  sway overlay; the coin/crack lives here.
- **jerusalem_market** (`market.png`) — the covered Old City souk.

| NPC | Scene | Role |
|-----|-------|------|
| Miriam | wall | **QUEST anchor** — archeologist; makes the Coin Rubbing once she has her charcoal; teaches the note custom |
| Dov | wall | **QUEST** — trades the charcoal stick for his glow-bug jar; scared kid, courage mirror of Jake |
| Eli (spice seller) | market | **QUEST** — alley-cat-at-the-stall errand; hands PP the Paper |
| Gary (tourist) | market | **QUEST** — trades the Sardine Tin (+ Pencil) for his lost Guidebook; comic relief |
| Worshippers | wall | ambient sway overlay (no quest) |
| Alley cat | market | obstacle/prop — guards the jar, lured off by sardines |

### Quest flow (designed)

```
Arrive jerusalem_entrance -> LEFT arch = market, RIGHT = wall
    |
    v
WALL — Miriam: "A rubbing of the Hadrian coin, for a frightened boy? I can — but
  my charcoal stick is gone; little Dov ran off with it. Get it back."
    |
    v
WALL — Dov: "Charcoal? I dropped my glow-bug jar and the alley cat batted it
  under the spice table. I'm not going near that cat. Bring my jar, get the charcoal."
    |
    v
MARKET — Eli (spice seller): "Jar's under my saffron table; the cat guards it.
  Only fish moves that cat — ask the tourist." (Presses a sheet of PAPER on PP.)
    -> PICKUP Paper  (PP: "...dunno what for, but oh well.")
    |
    v
MARKET — Gary (tourist): "Sardines? Take the tin — but find my Guidebook first."
  (Done sketching, also hands PP his PENCIL.)
    -> find Guidebook (floor item) -> trade -> Sardine Tin + PICKUP Pencil
       (PP: "...still no idea what for.")
    |
    v
MARKET — use Sardine Tin on the cat -> cat leaves -> grab Glow-bug Jar
    |
    v
WALL — Dov: return Glow-bug Jar -> get Charcoal Stick (Dov starts to lose his fear)
    |
    v
WALL — Miriam: give Charcoal Stick -> she makes the COIN RUBBING (Jake's key).
  Custom: "You never take from the Wall without leaving something. Write a note,
  tuck it in the stones, then go home."
    |
    v
WRITE the note (PP one-shot sprite) -> PUT it in the wall (PP one-shot sprite).
  Dov, brave now, places his own.
    |
    v
Note placed -> travel map UNLOCKS the flight to camp -> fly home
```

### Healing Jake (back at camp)

```
Camp grounds (darkened, post-France)
    |
    v
Jake's cabin -> Jake in strange state
    |
    v
Give Coin Rubbing to Jake (use item on NPC):
  - "That's HIM. That's the face in my head. You put him on PAPER!"
  - setStrange OFF, jake_healed = true
  - Unlocks Tokyo; Lily begins showing strange behavior (next chapter)
```

### Items used in this chapter

| Item | Source | Used on | Result |
|------|--------|---------|--------|
| Guidebook | Market floor item (Gary dropped it) | Gary | Trade for Sardine Tin (+ Pencil) |
| Sardine Tin | Gary | Alley cat | Lures the cat off the jar |
| Glow-bug Jar | Under Eli's stall (after cat leaves) | Dov | Trade for Charcoal Stick |
| Charcoal Stick | Dov (was Miriam's) | Miriam | She makes the Coin Rubbing |
| Paper | Eli (offhand) | written into the Note | — |
| Pencil | Gary (offhand) | writes the Note | — |
| Note (written) | PP (paper + pencil) | the wall crack | Placed → unlocks return flight |
| Coin Rubbing | Miriam | Jake (camp) | Heals Jake |

### Checklist (current status)

- [x] Three Jerusalem scenes (entrance / wall / market) rendering
- [x] NPCs wired: Miriam, Dov, Eli (spice seller), Gary the tourist
- [x] Worshippers ambient overlay at the wall
- [x] Trivial version live (Miriam hands the Coin Rubbing on first talk)
- [x] Coin Rubbing → Jake heal wiring (camp)
- [x] Camp-return pin relevant for the Jake leg
- [ ] Daisy-chain wired (guidebook → sardines → cat → jar → charcoal → rubbing)
- [ ] Note ritual (paper + pencil → write → put in wall) + two PP one-shot sprites
- [ ] Return flight gated on `jer_note_placed`
- [ ] New item art (guidebook, sardines, jar, charcoal, paper, pencil, note) + alley cat sprite

## Presentation & feel polish (2026-06-02)

Character-scale, camp-arrival and freakout-pacing pass. Code/data landed
this round; the three art regens are queued in `EXTRA_PROMPTS.md` (§M1–M3,
§P1).

- **PP one consistent size (#1/#2/#3/#7).** PP's draw scale now normalises
  by the tallest *opaque* pose in the active animation (the same method the
  NPCs use) instead of the raw cell height, so front- and side-facing
  walk/talk/idle all render at the same height. Killed the per-frame size
  jump that read as "two frames at once" on the front talk.
- **PP no longer shrinks mid-chat (#7).** `depthScale` range narrowed to
  0.95–1.05, so walking up-screen to a kid no longer visibly shrinks PP.
- **Camp arrival reads as an entrance (#4/#5).** New `entryWalk` scene flag:
  on reaching the camp grounds PP now walks in from off-screen left to his
  mark instead of popping into place; spawn moved right onto the path
  (x 80 → 210). This matches the retro "character strolls into the scene"
  arrival beat.
- **PP is the tallest again (#6).** The five grounds kids shrank
  (H 175 → 140, feet kept planted) so PP clearly towers over them, matching
  the original Camp Chilly Wa Wa wide shot (kids ≈ half PP's height).
- **Clicks on kids now animate (#8).** The talk-walk snap radius dropped
  80 → 30 px, so a short approach plays the walk instead of teleporting PP
  onto the talk spot (the "he doesn't move" pop at Marcus's room / Higgins).
- **Flower pickup lines up (#10).** PP's grab-flower pose lifts 38 px so his
  reach meets the flower on the ground.
- **Campfire sleep sits lower (#14).** Sleeping/waking PP anchored 615 → 650.
- **Marcus's freakout calmed down (#15).** The strange idle (and its periodic
  alt-idle "punctuation" beat) now cycle ~3.5× slower so the freakout reads
  as an uneasy fidget, not a flicker. The alt-idle still auto-fires after a
  few seconds of no interaction.

Open art (queued, not yet drawn):

- [ ] Marcus strange idle / talk / alt redrawn to his **canonical design** —
  round glasses + golden-yellow polo `#EEB421` + brown hair (the current
  strange sheets are the wrong kid: black hair, navy shirt, no glasses).
- [ ] PP walk-front redrawn to match the modern idle/talk-front PP (current
  sheet is an older design and a different size).

Still needs a live playtest pass (couldn't reproduce from assets alone):
action-button activation (#11), Higgins campfire shout firing (#13), Lily
"two rows" talk render (#7a), and whether giving the flower to Lily should
accept a plain click while holding it (#12).

---

## FUTURE CITIES (Not Yet Implemented)

### Jerusalem (Jake's visions — tunnels, coins, echoes) — PROMOTED

Moved up to its own chapter section above (**JERUSALEM — Jake's Anchor**). Scenes,
NPCs and the trivial heal path are built; the full retro task chain is designed in
`docs/JERUSALEM_TASKS.md`, not yet wired.

### Tokyo / Kyoto (Lily's visions) — OPENING BUILT (Batch 2, 2026-06-24)

**Build status (2026-06-24):** the chapter is wired and playable end-to-end on
placeholder/landed art. Built:
- Sad Lily revealed at the lake after Jake's heal (seen from behind); healed
  there with the Pressed Sakura (→ `lily_healed`, unlocks Tommy / Rio+BA).
- The rude-Higgins intercept (he strides halfway in via `npc_director_front_walk`)
  + PP's camera aside, which unlocks Tokyo.
- **Five Kyoto scenes:** `tokyo_torii` (arrival, Gary — his ramen tip opens the
  stall) → `tokyo_street` (ramen — Hiro + Kenji, falling leaves, the well) →
  `tokyo_temple` (flower store — Oba-chan + Kiku the dresser-geisha; matcha/bowl
  shelves) → `tokyo_teahouse` (UP from the store — the temple tea-house where the
  Tea Master hosts the matcha ceremony) → `tokyo_sakura` (hidden grove).
- **8-step offering chain + matcha gate** (see docs/JAPAN_TASKS.md): Gary →
  Hiro (needs his crow-stolen fire-striker for the blessed bowl) → Kenji (needs
  well-water) → well → Kenji (Voice Charm + clue) → Oba-chan (fire-striker) →
  Hiro (blessed Offering Bowl) → matcha ceremony (`jp_tea_done`) → Oba-chan
  "follow me" (`jp_grove_revealed`) → grove.
- **Grove payoff:** PP places the offering (+ Voice Charm) at the old tree →
  it blooms → **picks the Pressed Sakura** → Danny's foreshadowing phone call.
- All Kyoto items are consumed in-place (nothing lingers in the bag).

Still open (art only): the sakura-grove BG (`§JP-SAKURA-BG`), several talk
sheets, gap re-rolls, and BG edge-continuity — see EXTRA_PROMPTS §2026-06-24 Japan.


Lily's arc. Built after the Batch-1 bug sweep lands & playtests. Mirrors the
Paris/Jerusalem chapter pattern (chapter wiring in `game/tokyo.go` — a stub +
the Pressed Sakura heal handoff already exist in `game.go`; scenes in
`assets/data/scenes/`; NPCs in `assets/data/npc/`; items in `items.json`;
save-safe VarStore flags). All art ships behind graceful fallbacks; prompts go
to `docs/EXTRA_PROMPTS.md`.

**Opening — find Lily at the lake (not her cabin).**
- In `camp_lake`, Lily sits at the end of the dock seen from **behind, hugging
  her knees**, sad. (The old "strange Lily" art is retired.) New prompt:
  Lily-sad-from-behind idle.
- Talking to her starts her strange beat (petals, distant bells, "everything's
  the wrong colour").

**Higgins gets rude + the camera aside.**
- When PP heads to Higgins, Higgins walks **halfway down the line** to meet PP
  and is curt — the Lily conversation goes nowhere, he brushes PP off.
- After the dialog, PP turns to the **camera** (new aside beat): something like
  *"She did say she loves flowers… and it's the best season for them in Japan."*
  — the player's nudge toward Kyoto.

**Kyoto — three backgrounds** (prompts; pink/orange palettes, no pure white):
- **Torii path** — the orange wooden gate corridor (Fushimi-Inari style).
- **Kyoto city** — a machiya street.
- **Sakura grove** — pink cherry-blossom woods (revealed later, when a local
  lady tells PP to pick a blossom).

**Quest:** a Kyoto lady directs PP to the sakura grove to pick a single
blossom → that **Pressed Sakura** is Lily's anchor object (heals her back at
camp, kid-repair formula; the heal handoff `&handOff{item: "Pressed Sakura"}`
is already stubbed in `game.go`).

**Danny's phone call (sets up the NEXT destination).** Near the end, before
flying home, PP gets a call from Danny. Danny is NOT calling *from* a city — he
**wants an item / mentions something** that makes it make sense to head
somewhere immediately, teeing up the next chapter. PP joke:
*"I didn't know you guys have phones in the camp."*

**Return + heal Lily** at camp (give Pressed Sakura → `set_strange(false)` +
`lily_healed`), unlocking the next chapter.

### Kid → city map (locked 2026-06-24)

| Kid | City | Notes |
|-----|------|-------|
| Marcus | Paris | done |
| Jake | Jerusalem | done |
| Lily | Tokyo / Kyoto | Batch 2 (this chapter) |
| Tommy | **Rio** | confirmed |
| Danny | **Rome** | likely the **final** city — leave Danny's chapter for last |

### Endgame seed (far-off forward note, not built)

After the last kid, the **whole camp goes crazy** at once; PP calls on an old
friend for help — the **young magician from *Hokus Pokus Pink*** (confirm his
name). Capture only; not designed yet.

### Rome (Danny's visions — arena, gold paths, ruins)
- [ ] Danny goes strange: mapping camp like ruins, digging holes
- [ ] Travel to Rome
- [ ] Visit Colosseum
- [ ] Explore Roman ruins
- [ ] Learn about gladiators, Colosseum (72 AD), Roman engineering
- [ ] Find anchor object: seal sketch or inscription rubbing
- [ ] Return to camp, give to Danny

### Rio de Janeiro + Buenos Aires (Tommy's visions — music, statue, carnival, tango)
- [ ] Tommy goes strange: hearing music, dancing, mixing landmarks
- [ ] Travel to Rio de Janeiro
  - [ ] Christ the Redeemer statue
  - [ ] Copacabana beach
  - [ ] Carnival
- [ ] Travel to Buenos Aires
  - [ ] La Boca neighborhood
  - [ ] Obelisco
  - [ ] Tango
- [ ] Find anchor object: dance card or paired stamp
- [ ] Return to camp, give to Tommy

### Mexico City (Secret city — discovered during gameplay)
- [ ] Unlocked by clue found in another city (Rome?)
- [ ] Aztec pyramids, Teotihuacan
- [ ] Frida Kahlo museum
- [ ] Learn about Aztec civilization, pyramids

---

## Kid Repair Formula

Each kid follows the same pattern:
1. **Normal** — Kid has a natural personality trait (Day 1)
2. **Strange** — That trait becomes exaggerated/supernatural (Day 2+)
3. **Travel** — PP flies to the real city the kid is seeing
4. **Discover** — PP finds the source of the kid's visions, learns real history
5. **Anchor** — PP brings back a small real-world object from the place
6. **Repair** — Kid sees the object, their mind settles, they return to normal

---

## Travel System

- Camp Chilly Wa Wa Air (funny old propeller plane)
- Globe/world map with landmark icons per city
- Unlocked cities: full-color landmark image + yellow glow + label
- Locked cities: dimmed/gray landmark image + muted label
- Dotted flight path between camp and destination
- Airplane cutscene (~4 seconds, PP idle in plane) plays before arriving at any city

---

## Cities on Globe

| City | Kid | Status |
|------|-----|--------|
| Camp Chilly Wa Wa | — | Always unlocked |
| Paris | Marcus | Unlocked after Day 2 Marcus dialog |
| Jerusalem | Jake | Scenes + NPCs built; trivial heal live; full task chain designed (`docs/JERUSALEM_TASKS.md`), not yet wired |
| Tokyo / Kyoto | Lily | Designed (Batch 2, 2026-06-24); see the Tokyo/Kyoto chapter above |
| Rome | Danny | Not yet implemented |
| Rio de Janeiro | Tommy | Not yet implemented |
| Buenos Aires | Tommy | Not yet implemented |
| Mexico City | Secret | Not yet implemented |
| Stonehenge / London countryside | TBD | **Planned** — inspired by the PtP "part 3" clip ending (PP flies from London to a stone circle for a druid puzzle). Use as a side destination or a step inside the Jerusalem/Tokyo chain. Authoring workflow lives at `docs/SKILL.md §4a`. |
