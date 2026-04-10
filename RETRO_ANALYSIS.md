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

*This analysis is referenced in development. See FIXME.md and STATUS.md for current implementation progress.*
