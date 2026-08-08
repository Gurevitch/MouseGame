# Retro Pink Panther Games — Architecture Analysis

> Analysis of the original PP games to guide our game development.
> Source: `C:\Users\Roii\Documents\PP HP\HokusPP` and `C:\Users\Roii\Documents\PP P2P\ppp2p`

---

## Games Analyzed

### 1. Pink Panther: Hokus Pokus Pink (1997)
- **Engine:** Wanderlust Interactive proprietary C++/MFC
- **Executable:** `hpp.exe` (627 KB)
- **Game Data:** `HPP.ORB` (480 MB — all assets in one container)
- **Audio:** Miles Sound System (`mss32.dll`)
- **Resolution:** 640x480, 256-color palette
- **Story:** 3 modules, 39 steps — Periowinkle mansion, Israel, Siberia, Kenya, Indonesia, Greece

### 2. Pink Panther: Passport to Peril (1997)
- **Engine:** Same Wanderlust engine
- **Executable:** `pptp.exe` (592 KB)
- **Game Data:** `PPTP.ORB` (616 MB) + `pptp.bro` (9 MB)
- **Script Files:** `game.ptp`, `allsongs.ptp` (text-based property format)
- **Locations:** England, Egypt, India, China, Australia

---

## Core Architecture Patterns

### 1. Scene/Page System
- Scenes called "Pages", grouped into "Modules"
- Each page loads independently, can be purged from memory
- Background + foreground actors + walkable zones + hotspots
- Pages triggered on entry via `CHandlerStartPage`

### 2. Actor System (3 types)
| Type | Class | Purpose |
|------|-------|---------|
| Lead Actor | `CLeadActor` | Player (Pink Panther) — full control |
| Supporting Actor | `CSupportingActor` | NPCs — scripted behavior |
| Cursor Actor | `CCursorActor` | Interactive cursor with states |

### 3. Action System (modular, reusable)
| Action | Purpose |
|--------|---------|
| `CActionPlay` | Play animation |
| `CActionLoop` | Loop animation |
| `CActionStill` | Hold static frame |
| `CActionTalk` | Dialog with lip-sync |
| `CActionSound` | Sound effect |
| `CActionPlayWithSfx` | Animation + sound synced |
| `CActionHide` | Hide actor |
| `CActionText` | Display text |
| `CActionCEL` | Frame-by-frame cel animation |

### 4. Handler/Event System (data-driven)
| Handler | Trigger |
|---------|---------|
| `CHandlerLeftClick` | Left mouse click on object |
| `CHandlerUseClick` | Use inventory item on object |
| `CHandlerStartPage` | Page/scene loaded |
| `CHandlerTimer` | Timed events |
| `CHandlerSequences` | Complex choreographed sequences |

### 5. Condition System (gates all logic)
| Condition | Checks |
|-----------|--------|
| `CConditionGameVariable` | Global game state |
| `CConditionPageVariable` | Per-scene state |
| `CConditionModuleVariable` | Per-chapter state |
| `CConditionInventoryItemOwner` | Who has an item |
| `CConditionNotInventoryItemOwner` | Who doesn't have an item |

### 6. Variable Scoping (3 levels)
| Scope | Persists | Example |
|-------|----------|---------|
| Game Variable | Entire game | `JacksonsNameKnown` |
| Module Variable | Within chapter | `FoxEnabled` |
| Page Variable | Within scene | `DidTea`, `BoyBlocked` |

### 7. Sequence/Choreography System
- `CSequencer` plays complex multi-actor cutscenes
- `CSequence` = choreographed multi-step interaction
- `CSequenceItem` = individual step (actor + action + timing)
- `CSequenceItemLeaderAudio` = audio-timed leader actions
- All actors must reach READY state after sequence completes
- Side effects trigger state changes on completion

### 8. Inventory System
- `CInventoryMgr` manages all items
- Items track `InitialOwner` and `CurrentOwner`
- Items can be given to NPCs (ownership transfer)
- UI: scrollable window with left/right arrows
- Conditions can check item ownership

### 9. Movement & Pathfinding
- `CWalkMgr` manages character movement
- `CWalkLocation` defines walkable positions/areas
- `CWalkShortestPath` for automatic pathfinding
- 3-axis movement: MOVEX, MOVEY, MOVEZ (depth sorting)
- Click-to-move with smooth walking animation

### 10. Side Effects (state changes on actions)
| Side Effect | Does |
|-------------|------|
| `CSideEffectVariable` | Change a variable value |
| `CSideEffectGameVariable` | Change global state |
| `CSideEffectPageVariable` | Change scene state |
| `CSideEffectLocation` | Move actor to location |
| `CSideEffectInventoryItemOwner` | Transfer item ownership |
| `CSideEffectRandomPageVariable` | Set random value |

---

## What We Should Adopt (Priority Order)

### HIGH PRIORITY

**1. Sequence System for Cutscenes**
Our night scene and dialog chains use fragile callback nesting. A proper sequence player would:
- Define cutscenes as data (list of steps)
- Each step: actor, action, timing, conditions
- Auto-advance through steps
- Handle actor states (talking, idle) automatically
- Support audio sync

**2. Variable Scoping**
Replace flat Game struct fields with scoped variables:
```
GameVars:   parisUnlocked, nightSceneDone, talkedToMarcus
SceneVars:  metHiggins, flowerPickedUp (reset per scene visit)
ChapterVars: day, metKids (persist within chapter)
```

**3. Ownership-Based Inventory**
Track who has each item. Enable "give item to NPC" cleanly:
- Flower: Player → Lily
- Postcard: Curator → Player → Marcus
- Map: Higgins → Player

### MEDIUM PRIORITY

**4. Handler + Condition System**
Move from hardcoded `onDialogEnd` closures to declarative handlers:
```
Handler: LeftClick on Lily
Condition: Player has "Flower"
Action: Play lilyFlowerSequence
SideEffect: Set metKids++, Transfer Flower to Lily
```

**5. Walk Locations (not segments)**
Define walkable areas as polygons/zones instead of line segments. More intuitive and easier to tune.

**6. PDA/Map as Full UI**
Their PDA has pages, buttons, and navigation. Our map should evolve into a proper travel UI with:
- City info pages
- Travel history
- Clue tracking

### LOWER PRIORITY

**7. Asset Bundling (ORB format)**
They pack everything into one file. Good for distribution but not critical for development.

**8. Save/Load System**
Serialize entire game state. Important but can wait until gameplay is solid.

**9. Audio Sync**
Time animations to voice/sound clips. Nice to have when we add audio.

---

## Game Design Patterns

### Puzzle Flow (from both games)
1. **Encounter** — Meet character/see object
2. **Block** — Something prevents progress (NPC won't talk, door locked)
3. **Hint** — Another character hints at solution
4. **Collect** — Find required item in different location
5. **Use** — Apply item to resolve the block
6. **Reward** — New dialog, new area unlocked, story advances

### Character State Machine
```
IDLE → TALKING → IDLE (dialog complete)
IDLE → PLAYING → IDLE (animation complete)
IDLE → HIDDEN (removed from scene)
ANY → READY (sequence complete, waiting)
```

### Scene Lifecycle
```
1. CHandlerStartPage fires
2. Conditions evaluated
3. Initial animations/dialogs play
4. Player can interact
5. Handlers respond to clicks
6. Side effects modify state
7. Page exit → transition
```

---

## Music System (from Passport to Peril)

| Song | Region |
|------|--------|
| GInTimeSong | Main theme |
| GChinaSong | China |
| GFinaleSong | Ending |
| GMummySong | Egypt |
| GDreamTimeSong | Dream sequences |
| GTajmahalSong | India |
| GGuyFoxSong | Guy Fawkes |
| GCasteSong | India (caste) |

Audio controlled via PDA jukebox with per-region enable/disable.

---

---

# Footage Analysis — Presentation & Feel (2026-08-08)

> The sections above were reverse-engineered from the game DATA files
> (handlers, sequences, variables). This section comes from frame-level
> analysis of actual gameplay footage: four walkthrough clips in
> `~/Documents/pp retro frames and clips/` —
> **PtP part 2** (London park + the Mucky Duck pub + manor),
> **PtP part 6** (camp hub, India, China, Bhutan),
> **HPP part 6** (Africa, reef, jungle),
> **PtP Stonehenge** (manor library → silhouette travel → Stonehenge, Hebrew dub).
> Timestamps below are seconds/minutes into the named clip. Extracted frame
> sheets live in the session scratchpad (`ptp2/ ptp6/ hpp6/ stonehenge/`).

## A. PP movement (how the original actually moves)

- **The walk is a strut, and it's SLOW.** Upright, chest out, huge stride
  arcs, arms swinging; ~3–4 s to cross half a 640-wide screen in every clip.
  The game never hurries PP; the leisure IS the character.
- **Two gaits.** On a far click PP breaks into a stretched-body dash
  (Stonehenge 345s). There is also a comedy **all-fours sneak-prowl** used as
  flavor inside the pub (PtP2 ~352s). We currently have one walk speed only.
- **Turn-around = corkscrew.** PP turns via a 1–2 frame leg-cross "twist"
  pose, then unwinds facing the other way (PtP6 976s) — a snap with one
  in-between frame, not a smooth multi-frame turn.
- **Depth scaling is used selectively.** Interiors (pub, restaurant, theater)
  stage PP on a single flat depth band — no scaling at all. Scaling appears
  only on explicit into-the-screen paths: the HPP beach (PP grows ~3–4× walking
  from waterline to camera, 600–610s) and Stonehenge's diagonal ledge→circle
  wedge (~1.05 → ~0.7). Validates our narrowed `depthScale`; interiors should
  pin it to ~1.0.
- **Crowded-room trick: the raised catwalk.** In the Chinese restaurant PP
  walks a raised sideboard BEHIND the crowded table (PtP6 830–877s) so he stays
  full-size and unoccluded. Give crowded interiors a rear walk band instead of
  shrinking PP.
- **Every scene entry/exit is in-fiction.** Walk in from frame edge (our
  `entryWalk` — confirmed correct), climb down stairs, dive off the boat rail,
  climb back over it **in a distinct soaked/slicked-down wet sprite** (HPP
  437–441s), squeeze out along the manor mantel, climb a dropped rope ladder.
  Never a pop.
- **Idle fidgets are personality punctuation, not tight loops**: hands-on-hips
  (the signature listening pose), weight shifts, shrugs with splayed fingers,
  an unprompted sparkle-dance. In a pitch-dark cellar PP is a silhouette with
  **glowing yellow eyes only** (PtP6 ~3:17) — cheap and charming for any night
  interior.
- **Talking to seated NPCs**: PP stays standing and **bends deeply at the
  waist** to get face-level (pub ladies' table, the dog), props a foot on a
  stool rung, leans over the bar. He picks tiny NPCs UP and talks to them in
  his palm (PtP6 24–31s).

## B. Item transfer — the oval "loupe" inset is the whole grammar

Neither game has a visible inventory strip, text, or HUD. Every item beat runs
through one reusable device: a **pink/white-ringed oval inset** that grows out
of PP's position (~0.5 s), fills with a slate-blue field showing **PP's
disembodied paws holding the item large**, holds 2–3 s, and shrinks back.

- **Pickup**: kneel/reach one-shot in scene → oval close-up = the "you got it"
  feedback (PtP6 worm dig 1286–1296s; PtP2 computer-mouse gag ~54s).
  A pocket beat exists too (HPP 53–55s, book slides into hip fur) — our
  "always pocket" rule is era-accurate.
- **Examine/read**: documents and props get the same close-up, with gags baked
  in (the "A.S.BESTOS chief architect" blueprint credits; PP's arm wiping the
  mirror glass). Close-ups have an explicit exit-arrow button.
- **Give**: talk loop → oval close-up showing the item in paws → NPC plays a
  bespoke receive one-shot (druid caveman paints a heart on the sarsen). Our
  two-stage `handOff` is already RICHER than the original — keep it, add the
  close-up layer on top.
- **Inventory browsing is diegetic** (HPP 798–802s): a small circle next to PP
  cycles carried items one per click, each drawn held by his paws. No grid.
- **Consumed items dissolve in a white sparkle cloud** (HPP magic note 56–60s);
  used items can **persist as scene props** (the mirror lands propped by a
  stone and later becomes a light source).
- **Pacing**: every reward moment is held ~1 s longer than modern instinct —
  kneel, oval grows, beat, oval shrinks. `playHandOff` should keep that
  deliberate slowness.

## C. Dialog presentation

- **Zero on-screen text in ~61 minutes of footage.** Voice + gesture only. We
  keep our text/dialog system (no VO budget), but copy the staging:
- Both parties stay full-body in the wide shot — no camera moves, portraits,
  or letterboxing. NPCs have clearly separate idle vs talk loops (our rule,
  confirmed); talk = mouth flaps + whole-body gesture loop.
- **The world never pauses for dialog.** Every ambient loop keeps running while
  NPCs talk; only the player is locked out. (Worth auditing our scenes for
  ambients that freeze during dialog.)
- **Click feedback**: a small white X-shaped sparkle burst at the click point
  (Stonehenge 218s).
- **Thought bubbles as wordless hints**: PP gets a cloud bubble picturing an
  item, then one picturing the destination (PtP2 506–510s — the black duck
  sign = "go to the Mucky Duck").
- **Videophone calls** (PtP6 22:32, Stonehenge 5:50): full-screen pink PDA held
  by PP's thumbs at the frame edges; the caller is a **giant face over a live
  scene BG with rich lip-sync + eyebrow acting** (a new expression every
  1–2 s); calls open/close with ~1 s of TV static that also masks the loop
  seam. **This is exactly how Danny's Kyoto phone call should be staged.**

## D. Transitions & travel grammar

Fixed sequence at every scene change:
1. Cursor becomes a **white directional arrow at the screen edge** (the only
   exit affordance — no permanent on-screen arrows); PP walks off.
2. **Iris/oval wipe closes** centered on PP (~1–1.5 s) → beat of black —
   during which a **tiny white mouse scurries across the darkness** (PtP2
   178s, 528s, 674s): the transition itself carries a running gag.
3. For travel: a full-screen painted map where visited locations are landmark
   icons ON the map; a **dotted trail draws itself segment-by-segment** while a
   small vehicle icon crawls along it (~4–8 s); **new icons appear only after
   PP learns of the place in dialog** (the manor pin appears at PtP2 780s only
   after the pub story).
4. Iris opens on the new scene's focal point.

Alternates: HPP's purple vortex tumble (chapter travel), and the Stonehenge
**silhouette travel vignette** — one black hill + tree against a gradient sky,
PP and the tailed pair crossing as flat silhouettes while the sky ramps
day → dusk → starry night. One background + existing walk cycles rendered
black = an entire travel cutscene. Perfect for short hops (London → Stonehenge).

Scripted **whole-scene lighting progression as a quest clock**: Stonehenge
starts night-blue and a yellow dawn band behind the trilithons widens from the
first puzzle beat to the end (~210s→410s). The scene tells time passing.

## E. Ambient life — the room-density formula

Every strong scene follows the same recipe (best seen in the pub and the
Chinese restaurant):

1. **3–6 independent loops on independent timers** (never synchronized).
2. **One mobile critter** stitching the room together (the pub dog wanders,
   begs, hops on the bar, follows PP; the white park mouse; the escaped
   parrot flying laps of the restaurant).
3. **One loud periodic beat** readable across the room (the red fat man's huge
   toast-cheer every ~10–15 s).
4. **Walk-on extras that change between visits** (a boy on a stool, a
   brown-coat man drifting by the dartboard — population differs visit to visit).
5. **Persistent micro-state**: props move after interactions and STAY moved
   (the newspaper migrates to the ladies' table; a fish-and-chips plate
   appears on the bar after the pump gag; the painted heart stays on the
   sarsen). Consequence without dialog.

## F. The Mucky Duck pub — deep dive (London clip)

Layout: bar across the back wall, arched window center, chalkboard specials
top-right, TWO dartboards (set dressing only — no darts minigame in the
footage), ladies' tea table left, policemen's table center-right, newspaper
reader far right, one clear mid-floor corridor for PP.

Cast loops (all continuous, independent): barmaid polish/pour + arm-raise
greet; red fat man's periodic toast; dozing flat-cap man who sleeps through
everything; old ladies sip + knit + one stands to gossip; two policemen
sip/wipe-mustache; newspaper man raise/lower/turn-page; the dog (mobile).

Click toys: **beer pump** → PP yanks the handle, a ~1 s foam burst (~15 bubble
sprites) douses him; **sleeping man** → PP flips him fully upside-down in
midair, drops him back, he never wakes (~3 s); the sausage-on-a-fork pickup →
immediately fed to the dog.

**Silhouette-audience flashback** (10:36–12:00): the barmaid's Guy Fawkes
story plays as flat limited-palette panels while the CURRENT pub patrons
watch as dark-blue foreground silhouettes raising teacups. Dirt-cheap way to
stage "an NPC tells a story" — earmark for Rome/Rio storytelling beats.

**What the pub gives OUR scenes right now**: the density formula (§E) applied
to the Paris bakery, the Kyoto ramen street, and the Jerusalem market — each
currently has loops but no mobile critter, no loud periodic beat, and no
walk-on variation. The Jerusalem **alley cat** should be the market's wandering
critter, not a static prop. And if the planned Stonehenge/London side
destination lands, a Mucky Duck-style pub is the obvious hub scene
(`docs/SKILL.md §4a` authoring flow; STORY.md travel table already seeds it).

## G. Stonehenge scene — build sketch (from the Hebrew clip)

Beats observed: arrive by silhouette vignette → iris-in, walk in from left
ledge → two druids (granny knitting a pink panther doll + skinny caveman) →
give the hand mirror (matchmaker beat) → caveman paints a **heart on the
standing stone that persists** → scripted dawn glow begins → place/throw the
mirror at a stone, it catches light as a moon-disk beam → the couple flees in
a smoke-puff dash → PDA call → rope ladder drops, PP climbs out in sparkles.

Clone mapping: one night BG + glow-band decorator (`scene_ambient.go`) driven
by a `sh_heart_painted` flag; diagonal walk wedge with depthScale 1.05→0.7;
two NPCs (idle+talk each, house rules); heart decal + propped-mirror props as
persistent state; PP one-shots `place mirror` (starts pocket-out per GIVE
rule) + `climb rope ladder`; exit gated on both flags. Retro-faithful and
almost entirely made of systems we already have.

## H. Adoption plan — prioritized (engine + animation)

### Tier 1 — cheap engine wins, big retro payoff
1. **Iris in/out scene transition** centered on PP (we hard-cut today).
   One shader-less mask circle over the frame; ~1–1.5 s each way.
2. **Oval examine/close-up inset** — one overlay ring texture + per-item
   "paws holding it" close-up art (graceful fallback: big item icon in the
   oval). Hook: pickups, "look" clicks, and the show-item beat inside
   `playHandOff`. This is the single most identity-defining device in the
   footage.
3. **Travel-map trail animation**: dotted trail draws segment-by-segment +
   small plane icon crawl + iris bookends; pins appear only when learned
   (our `travel_map.json` display state already supports show/hide).
4. **Edge-exit arrow cursor states** replacing/augmenting always-visible
   arrows; pointing-finger cursor over interactables (cursor change is the
   ONLY affordance the originals need).
5. **Ambient density pass** on bakery / ramen street / market using the §E
   formula: add a `wander` mover to `ambient_sprite.go` (mobile critter —
   the Jerusalem alley cat first), one loud periodic beat per room, walk-on
   extras varied by visit count.

### Tier 2 — feel & staging
6. **Danny's Kyoto call as a full-screen PDA videophone**: giant lip-synced
   face sheet (idle+talk+3–4 expression frames), PP's thumbs at frame edges,
   static-burst open/close. Reusable for any camp-calls-PP beat.
7. **One long-gag puzzle payoff per city** (the HPP shark ballroom-dance is
   45 s): stage each anchor-object reveal as an extended comedy one-shot,
   not just a handOff.
8. **Scene lighting progression decorator**: crossfading glow band driven by
   VarStore flags (quest clock; needed for Stonehenge dawn anyway).
9. **Consumed-item sparkle dissolve** for anchor handovers; **persistent
   micro-state props** formalized in scene JSON (prop position/visibility
   keyed to flags).
10. **Thought-bubble wordless hints** (bubble with item icon → bubble with
    place icon) as a hint layer that needs no dialog lines.

### Tier 3 — animation/art additions (queue prompts in EXTRA_PROMPTS.md when scheduled)
11. **Run gait** for far clicks + **corkscrew turn** in-between frame +
    **bend-at-waist talk** variant for seated NPCs (bakery tables, pub).
12. **Wet/soaked PP** entrance sprite (any water exit) and **glowing-eyes
    dark idle** (night interiors) — both single sheets.
13. **Silhouette travel vignette** (one hill BG + sky palette ramp + existing
    walks drawn black) for short hops; **silhouette-audience flashback**
    framing for storytelling NPCs (Rome/Rio).
14. **Costume state variants** (PP wore a tweed detective set for the whole
    England chapter) — big art cost, only if a chapter demands it.

*This analysis is referenced in development. See FIXME.md and STATUS.md for current implementation progress.*
