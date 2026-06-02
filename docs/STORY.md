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
| Mademoiselle Camille | paris_bakery (mid table) | Flavor — Pierre's pink-painting hint |
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
  - Gives postcard of the restored painting (ANCHOR OBJECT)
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

### Return-to-Paris beat (post-Marcus-healed)

User 2026-05-20: keep the city alive after PP heals Marcus. The bakery
gains a new dialog when PP comes back to Paris:

- Madame Poulain (`bakeryWomanLouvreSouvenirDialog`): asks PP to bring
  her another Louvre postcard for her grandson in Lyon as a thank-you
- Doesn't gate further quests — pure flavor / world-warmth beat
- Triggered automatically when `marcusHealed == true` on next bakery
  visit (see `setupParisCallbacks` in `game/game.go`)

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
| Postcard | Curator Beaumont | Marcus (back at camp) | Heals Marcus |

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

### Jerusalem (Jake's visions — tunnels, coins, echoes)
- [ ] Jake goes strange: rubbing surfaces, hearing echoes, "tunnels under old city"
- [ ] Travel to Jerusalem
- [ ] Visit Western Wall area
- [ ] Explore ancient tunnels
- [ ] Learn about Old City, Western Wall history
- [ ] Find anchor object: coin rubbing or tunnel map
- [ ] Return to camp, give to Jake

### Tokyo (Lily's visions — glowing garden, temple bells)
- [ ] Lily goes strange: arranging petals in symbols, hearing bells
- [ ] Travel to Tokyo
- [ ] Visit Senso-ji temple
- [ ] Explore cherry blossom garden
- [ ] Learn about temples, cherry blossoms, torii gates
- [ ] Find anchor object: pressed flower or charm card
- [ ] Return to camp, give to Lily

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
| Jerusalem | Jake | Not yet implemented |
| Tokyo | Lily | Not yet implemented |
| Rome | Danny | Not yet implemented |
| Rio de Janeiro | Tommy | Not yet implemented |
| Buenos Aires | Tommy | Not yet implemented |
| Mexico City | Secret | Not yet implemented |
| Stonehenge / London countryside | TBD | **Planned** — inspired by the PtP "part 3" clip ending (PP flies from London to a stone circle for a druid puzzle). Use as a side destination or a step inside the Jerusalem/Tokyo chain. Authoring workflow lives at `docs/SKILL.md §4a`. |
