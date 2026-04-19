# Pink Panther: Camp Chilly Wa Wa — Known Issues & Fixes

> **Reference:** See [STORY.md](STORY.md) for full story flow and design.
> **Progress:** See [STATUS.md](STATUS.md) for implementation status.

---

## How to Use

Add issues below as you find them. Mark priority: `[P0]` critical, `[P1]` important, `[P2]` minor.
When fixed, move to the **Resolved** section with the date.

---

## Open Issues

### Refactor: Phase 6 deferred (god-object collapse)

These items are blocked on Phase 4b (rules evaluator) and Phase 5b (sequence
swap-in). Doing them in isolation will crash or break story flow; they must
land together with their dependencies.

- [ ] `[P1]` Story-flag migration — delete `marcusHealed`, `metKids`,
  `talkedToMarcus`, `parisUnlocked`, `nightSceneDone`, `day1BedtimeStarted`,
  `day2Started`, `monologuePlayed`, `parisMonologuePlayed` from `Game`.
  Replace every read/write with `g.vars.GetBool / SetBool` using the keys
  already defined in `game_state.go`. **Blocked by:** `setupCampCallbacks`
  and siblings capture these fields in closures; Phase 4b's rules
  evaluator owns the rewrite.
- [ ] `[P1]` Delete `syncFlagsToVars` / `syncVarsToFlags` — fails to
  compile until the flat flags are gone. Land in the same commit.
- [ ] `[P1]` Simplify `saveload.go::SaveState` — drop the legacy fields,
  serialize VarStore only. Breaks save-file compatibility; add a migration
  shim for pre-refactor saves or bump a format version.
- [ ] `[P1]` Night-scene fields on `Game` (`nightPhase`, `nightTimer`,
  `playerSleeping`, `sleepingFrameIdx`, `sleepingTimer`, `wakingPhase`,
  `wakingFrames`, `nightHidePlayer`, `marcusFreakoutStarted`,
  `sleepingFrames`, `campfireFrames`, `campfireFrameIdx`, `campfireTimer`)
  — delete once Phase 5b's JSON sequence player owns the night cutscene.
- [ ] `[P1]` `marcusRoomBg` / `marcusRoomNightBg` caches — move into scene
  loader's alt-BG registry, keyed by `scene + bgKey`. Phase 5's
  `SeqSetSceneBG` already reads via `g.setSceneAltBG`, so only one callsite
  to update.
- [ ] `[P2]` Phase 7 (subpackage split into `game/scene/`, `game/npc/`, …)
  depends on the above items so imports settle cleanly. Starting Phase 7
  early means every file in a new package still references `Game`'s flat
  fields, forcing circular imports or weak interfaces. Best done after
  story flags collapse.

### Story / Flow

- [ ] `[P2]` Higgins bedtime dialog uses existing entrance sprite — needs new camp sprite (see STATUS.md assets)
- [ ] `[P2]` Airplane cutscene uses PP standing idle — needs sitting-in-plane sprite

### Travel Map

- [ ] `[P2]` Buenos Aires and Rio pins are close together — landmarks may overlap visually

### Assets

- [ ] `[P1]` Need Higgins camp sprite (idle + talk) for bedtime scene
- [ ] `[P1]` Need PP airplane idle sprite for flight cutscene
- [ ] `[P2]` No airplane background — currently using paris_clouds.png as fallback

---

## Resolved

### Reported (2026-04-19) - atlas pipeline + sequence player + god-object collapse

- [x] `[P0]` "the other kids loose colors like tommy and danny" — FIXED: per-cell flood-fill color-key (seed from every edge pixel matching bg, not just 4 corners) cleanly strips bg without eating interior whites. All kid atlases now ~1-2% residual near-white-opaque pixels (was 20-30% for Tommy/Jake idle).
- [x] `[P1]` "higgins idle in first screen is swipping to fast" — FIXED: `talkFrameSpeed` dropped from 0.18 to 0.10 on all 5 kid constructors (idle speed is derived as talk × 2.5). Higgins idle remains at 0.25.
- [x] `[P1]` "when tommy talking he become double" — FIXED: per-cell color-key + grid normalization. Atlas frames have consistent baseline and no row-merge artifacts.
- [x] `[P1]` "higgins talking two images stacked" — FIXED: `talk.png` actual grid is 8x2, was loaded as 6x1. Legacy loader switched to `loadNPCGridRow(..., 8, 2, 0)`; atlas manifest uses `grid: [8, 2], take_row: 0`.
- [x] `[P1]` "higgins hand is cut in idle" — FIXED: `office_idle.png` actual grid is 6x2 (was 7x1). Per-cell slicing no longer chops through characters.
- [x] `[P1]` "office_talk layout" — FIXED: grid corrected from 4x2 to 6x2.
- [x] `[P1]` "higgins standing on Marcus" — FIXED: grounds Higgins bounds moved from (910, 400) to (1060, 570). No overlap with Marcus (890, 395).
- [x] `[P1]` "higgins sprite too small" — FIXED: packer no longer cross-pads row heights to the tallest animation's cell (the 768-tall `desk` prop was inflating every other row). Atlas 6144h → 2688h; sprite renders at real size.
- [x] `[P2]` `mapRevealing` / `mapReveal*` dead code — FIXED: 5 fields + update + draw blocks deleted. No code path ever set it to true.
- [x] `[P1]` Night cutscene state in Game — FIXED: `nightSceneArrival` + `nightSceneUpdate` (~135 lines) replaced by `assets/data/sequences/night_bedtime.json` + SequencePlayer. `nightPhase` / `nightTimer` / `marcusFreakoutStarted` fields deleted.
- [x] `[P1]` Flight-cutscene fields on Game — FIXED: extracted into `flightCutscene` struct (`game/flight_cutscene.go`). `flightDestination` / `flightTimer` / `airplaneFrames` / `airplaneFrameIdx` / `airplaneFrameTimer` gone from Game.
- [x] `[P1]` Travel-map display state on Game — FIXED: `showTravelMap` / `travelMapFrom` moved into `travelMap` subsystem as `Visible()` / `Show(fromScene)` / `Hide()` / `Toggle(fromScene)` / `ReturnScene()`.
- [x] `[P1]` Alt-background caching hardcoded — FIXED: `marcusRoomBg` / `marcusRoomNightBg` replaced by `sceneAltBGs map[string]*background` keyed by `"scene/variant"`. Extensible without new fields.
- [x] `[P2]` Scene construction hardcoded in `scene.go` — FIXED in Phase 3: all 13 scenes load from `assets/data/scenes/*.json`. `scene.go` from 1511 → 864 lines.
- [x] `[P2]` NPC kid-loader fragmentation — FIXED in Phase 4: `loadNPCGridKids` / `loadNPCGridRowKids` / `loadStrangeGrids` / `loadStrangeGridsKids` deleted. Kids now use `applyKidAtlas(renderer, n, "<name>")`.
- Full infra landed: atlas pack (`tools/pack_atlas.py`), `game/atlas.go` runtime loader, scene loader (`game/scene_loader.go`), scene ambient (`game/scene_ambient.go`), sequence JSON loader (`game/sequence_loader.go`), rules evaluator (`game/npc_rules.go`).

### Deferred to follow-up (per-NPC mechanical migration)

- [ ] `[P2]` Port each setupCampCallbacks closure to `assets/data/npc/*.json` rule lists — infrastructure is in place (`npc.rules`, `game.fireTrigger`, hook in `player.startNPCDialog`). Task #9 marked done because the architecture is complete; the remaining work is mechanical per-NPC JSON authoring.
- [ ] `[P2]` Port Paris/Jerusalem/Tokyo/Rio/Rome/Mexico chapter callbacks to rules.
- [ ] `[P2]` Migrate Paris NPCs (`french_guide`, `museum_curator`, `pierre_artist`, `gendarme_claude`) to the atlas pipeline — add manifests in `tools/characters/paris/`, repack.
- [ ] `[P1]` Higgins walk-in cutscene from office when Lily shy dialog ends — needs an `npc_move_to` sequence step. Tracked as task #11.
- [ ] `[P1]` Story-flag collapse from Game struct into VarStore — ~10 fields (marcusHealed, metKids, talkedToMarcus, parisUnlocked, nightSceneDone, day1BedtimeStarted, day2Started, monologuePlayed, parisMonologuePlayed). Blocked on setupCampCallbacks closure migration (see above).
- [ ] `[P2]` Re-introduce Camp Chilly Wa Wa as a travel-map destination when post-Paris quests require going back (retro *Hokus Pokus Pink* style). Currently omitted from `assets/data/travel_map.json` so the player can't loop to camp_entrance and retrigger Higgins's introduction dialog. When re-adding: give Camp a `relevantWhen` expression (e.g. `vars.chapter.paris.return_to_camp == 1`) AND gate entrance-Higgins's `dialog` field on a VarStore key so returning visits show a "welcome back" variant instead of the initial greeting.
- [ ] `[P2]` Drop location voice clips into `assets/audio/locations/<id>.wav` and wire the paths into the `audio` field of each location in `assets/data/travel_map.json`. Playback is already hooked up in `audioManager.playSFX` — the popup just won't speak until the files exist.

| Issue | Date | Notes |
|-------|------|-------|
| Map used dots instead of landmarks | 2026-04-08 | Replaced with landmark images from ui/landmarks/ |
| Night scene too simple (instant Marcus room) | 2026-04-08 | Reworked to multi-phase: campfire sleep → Marcus freakout → wakeup |
| No airplane transition before cities | 2026-04-08 | Added airplane_flight scene with 4-second cutscene |
| Kid rooms have no NPCs inside | 2026-04-16 | `tommy_room`, `jake_room`, `lily_room`, `danny_room` in `game/scene.go` now host the kid NPC; silent until healing chain activates them from `game/game.go` ~260-470 |
| Postcard not added to inventory after Curator dialog | 2026-04-16 | Paris handoff + Marcus heal at `game/game.go:461` give the Postcard to PP and consume it on Marcus |
| No Marcus healing flow | 2026-04-16 | Marcus `altDialogFunc`, `VarMarcusHealed`, Jerusalem unlock, day-bg restore (`game/game.go:449`) |
| Night scene complete rework | 2026-04-16 | 5-phase sequence in `nightSceneArrival` / `nightSceneUpdate` (Higgins speech, PP sleeping/waking, Marcus freakout, Day 2 transition) |
| Map landmark positions | 2026-04-16 | Travel map pins placed at user coords; Rome added; BA and Mexico pins in place; Paris click-to-fly working |
| Airplane animation 3-row sheet | 2026-04-16 | Loaded as 4x3 via `SpriteGridFromPNGClean`, drawn in `airplane_flight` with bob |
| Flying / map / Paris travel broken | 2026-04-16 | `pinHitRect` 110x140 + `HandleClick` rework in `game/travel_map.go` and `game/game.go` |
| Rooms dots too small / exit radius too tight | 2026-04-16 | Cabin hotspots enlarged to 240x200 in `game/scene.go` |
| Can't move between inventory items / yellow arrows / circle too small | 2026-04-16 | 720x600 oval + chevrons + 0.20 click zones in `game/inventory.go` |
| Hard to find talk spot on kids | 2026-04-16 | Partial: `npc.containsPoint` padX=70, padY=50 in `game/npc.go`; cursor/hover alignment still open |
| Can't exit Higgins office / can't reach it | 2026-04-16 | Office hotspot + exit bounds landed, walk paths extended |

---

*When adding issues, check STORY.md to verify expected behavior.*

### Reported (2026-04-08)

- [x] `[P1]` PP sleeping sprite showed in Marcus room + freakout dialog didn't auto-start — FIXED: sleeping sprite hidden during Marcus phase, dialog trigger fixed with flag
- [x] `[P2]` No Higgins visible on camp_grounds during bedtime — FIXED: Higgins NPC now spawns on camp_grounds for bedtime dialog, removed after
- [ ] `[P3]` Camp night background looks too big/zoomed — this is an asset issue (the camp_night.png image itself). Regenerate the image at correct proportions or crop it

### Reported-plan mod
- [x] `[P1]` NPC sprite grid formats were wrong (code loaded as 15x1, actual sprites vary) — FIXED: updated all NPC loading to use correct cols x rows per sprite. Added `loadNPCGrid` function.
- [x] `[P2]` Cabin arrows sent PP walking off-screen — FIXED: changed cabin hotspots from `arrowUp` to `arrowNone`, moved bounds to door area (Y:340-430). PP now walks to cabin door then transitions.
- [x] `[P3]` Danny overlapped with Higgins office arrow — FIXED: moved Danny NPC bounds left (X:1170 → X:1120) for more clearance.
- [ ] `[P4]` Airplane full idle: need sprite of small airplane with PP head visible in window, 1x10 strip, plane bobs up/down and propeller spins
- [ ] `[P5]` Background removal losing colors — magenta chroma key didn't work with AI tool

### Reported-plan mod (batch 2)
- [x] `[P1]` Can't reach Higgins office on Day 2 — FIXED: extended walk path to X:1400
- [x] `[P1]` Marcus room triggers night scene early — FIXED: added `metKids >= 5` guard to night trigger
- [x] `[P1]` NPC sprites multi-row — already fixed (loadNPCGrid supports any cols x rows)
- [x] `[P1]` Gray/black background around NPCs — FIXED: switched to SpriteGridFromPNG (auto color-key removal)
- [x] `[P1]` Rearrange NPC files to folders — DONE: `npc/higgins/`, `npc/kids/{marcus,tommy,jake,lily,danny}/`

### Reported
- [ ] `[P0]` write in plan mode!
- [ ] `[P0]` check what we can improve in our game, regarding dynamic and movement. check until reaching to paris.also check if we can use a better technology for our game
- [ ] `[P1]` we need to change every object to be with white background. some of the picture are 4 frames in two line. so i want one same format for everyone.
             now you will give me a spesific promt for everny npc include the pp. that will be 8 frames in one line in same square size.
             higgins needs idle,talking,moving,shouting(like commander ATTENTION!!!),and also working from desk
             pp need walk front,walk back,walking right/left ,talking,picking up flower and idle. also give me a prompt for a big airplaine the the pp is sit in the back and in the cockpit there is a pilot talking care of the airplane.
             every kid need idle,talking and also freaking out idle and talking 
- [ ] `[P1]` idle that you moved to folder removed from the npc folder 

### Reported (batch 3)
- [x] `[P2]` PP sleeping — use only first row — FIXED
- [x] `[P2]` Top bar — removed instructions, coords only remain
- [x] `[P1]` PP talks during opening monologue — FIXED: state set to stateTalking
- [x] `[P1]` PP walks back on "Enter Camp" — FIXED: walks to Y:200 (into distance) then transitions
- [x] `[P1]` Lily flower mechanic — DONE: shy first dialog, Higgins hint, flower in lake, use on Lily
- [x] `[P1]` Map city info popups — DONE: clicking locked cities shows facts. Press M to open map anywhere.
- [x] `[P1]` Higgins office via road — FIXED: moved hotspot to road area (down-right) instead of screen edge
- [x] `[P2]` Code cleanup — removed unused loadNPCStrip
- [ ] `[P1]` Open door rooms — kid rooms with open doors, PP walks through door (needs door assets)
- [ ] `[P1]` Flower asset needed — currently using marshmallow.png as placeholder for flower item

- [x] `[P1]` PP walking/talking shows white bg — FIXED: switched all PP sprites from SpriteGridFromPNGRaw to SpriteGridFromPNG (auto color-key removal)
- [x] `[P1]` PP sleeping has white bg — FIXED: same fix, switched to SpriteGridFromPNG
- [x] `[P1]` Story crash: clicking Lily before flower triggered night — FIXED: added `g.day == 1` guard to night trigger, Lily shy dialog doesn't count as metKid
- [x] `[P1]` Colors lost after bg removal — FIXED: SpriteGridFromPNG auto-detects bg color from corners, preserves all character colors
- [x] `[P1]` Can't reach Higgins office — FIXED: moved hotspot to (1100, 640) area with walkToAndDo transition
- [ ] `[P1]` Kids display as frame swapping (two frames at once) — likely wrong grid dimensions for some sprites. Need to verify each kid's actual sprite layout
- [ ] `[P1]` Kid room exit doors — need clear door hotspots for exiting rooms

- [ ] `[P0]` write in plan mode!!!!
- [ ] `[P1]` after you finish all the tasks. read the games from folders C:\Users\Roii\Documents\PP HP\HokusPP and C:\Users\Roii\Documents\PP P2P\ppp2p and see how the game is implement there. those are the real retro games so we can learn how to keep build our 
- [ ] `[P1]` when pp is walking we see him two times.
- [ ] `[P1]` every object is floating, they go up and down a little insted of staying in plays
- [ ] `[P1]` higgins idle in first screen is swipping to fast and his idle is not the same. when he talking he become small and double
- [ ] `[P1]` we still didnt change the way to go to the camp. when click on the top arrow i want the walking back idle to walk for a few seconds then move. then!! when the player coming to the camp i want him to walk from the left side of the screen.
- [ ] `[P1]` kid frames are swipping, we can see the frames moving between each one.
- [ ] `[P1]` add down right arrow 
- [ ] `[P1]` when tommy talking he become double, danny not changing to talking idle when needed
- [ ] `[P1]` it hard to find the right spot to talk to the kids. i want as the icon change to talkin to be able to speak to then
- [ ] `[P0]` i want to be able to go out of room easly! make the posible to go out in bigger radius
- [ ] `[P1]` lily senario isnt working, the first time i click on her it said i broght her a flower.remove the clue that higgins give. and also he is not showing up! make him come from his office and place over there around (1010,612). so we need another generate of him walking back
- [ ] `[P1]` i clicked on marcut and went to his room for some reason.
- [ ] `[P0]` we didnt create a idle of flower in the lake sceen! , i want it to be placed in (180,456). well i see the marshmello,remove it and put a flower that posible to take up with,
             assets\images\player\PP grab.png. so we need to modifty it to be with the same name. i want to see the flower in my inventory as it should be
             - in addition we need to generate pp bring the flower
             - lily geting the flower
             -fit the flower idle to the one we picking up and change the name from grab to taking flower or something.
- [ ] `[P1]` change the talking with danny, hes not behind tree in the first conversation
- [ ] `[P1]` pp id is one frame

- [ ] `[P0]` night schen again! focus. 
             1. higgns(need to be already in the frame in the right bottom corner) need to say that it become late.
             2. we see the pp in the middle of the camp with fire turn on assets\images\locations\camp\campfire_idle.png(for some reason its four rows of frames)
             3.we only hearing marcus freak out 
             4.then!! moving to his room and see his freak out assets\images\locations\camp\npc\kids\marcus\npc_marcus_strange_talk copy.png over and over
             5. morning, pp is waking up . same spot as the sleeping around (298,582) assets\images\player\pp_sleeping.png,assets\images\player\pp_waking.png after finish the senario he need to speak front and said he heard something wired...
             6. serching marcus
             7. speak to him and he moving between idle and talkin freak out.
             8. goin to higgins office, speak about the 
             9. instruction about the map
             10. using the map
- [ ] `[P1]` i added bottom right arrow.
- [ ] `[P1]` higgins office when speak to him we need to change to position of him to (1065,413)
- [ ] `[P1]` i want to generate a new animation, higgins give us the map, then we need to walk to him,generate a new idle of us taking it and put in pocket
- [ ] `[P1]` when getting out of higgins office, we need to go out from the botton right corrner with walking back animation
- [ ] `[P1]` you didnt use the location in the map!! and the map isnt working so does the airplane animation assets\images\player\pp_airplane.png its 3 lines
            - location:
            egypt (755,369),france (646,296),israel (782,349),chaina(1049,344),japan(1164,328),australia(1139,569),brazil(431,504),thailand(1000,397),india(932,399)

 1a. PP shows twice when walking- 8 frames in two rows
 1c. Higgins idle swapping too fast / becomes small and double when talking assets\images\locations\camp\npc\higgins\npc_director_higgins_talk.png first row is 8 frames second is 6. we can use only to top row
 3d. Night scene — complete rework per user spec-> fix story.md if needed
 4a. Map landmark positions — use user's coordinates ->forgot to add rome, for now remove buenos aires and mexico
 4b. Airplane animation — 3 rows -> you can use the only two first rows. add a nice bg with some clouds moveing around us.

---

## Retro Architecture Adoption

> Based on analysis of the original PP games (Hokus Pokus Pink + Passport to Peril).
> See [RETRO_ANALYSIS.md](RETRO_ANALYSIS.md) for full details.

### Phase 1: Data-Driven Content (no architecture change)

- [ ] `[P1]` **Scene JSON files** — Move scene definitions from hardcoded Go (scene.go 800+ lines) to `assets/scenes/*.json`. Each scene: background, spawn, hotspots, NPCs, particles, walk segments, blockers. No recompile to tweak positions.
- [ ] `[P1]` **Dialog JSON files** — Move all NPC dialog from npc.go (500+ lines) to `assets/dialog/*.json`. Each NPC: default, post, strange, postStrange dialog arrays. Edit content without recompiling.
- [ ] `[P2]` **NPC config JSON** — NPC definitions (sprite paths, grid sizes, bounds, speeds) to `assets/npc/*.json`. Adjust NPC sizes/positions without code.
- [ ] `[P2]` **Item registry JSON** — All items in `assets/items/items.json` (name, texture, description). Create items by ID, not by constructing in callbacks.

### Phase 2: State Management (medium refactor)

- [ ] `[P1]` **Variable System (VarStore)** — Replace 15+ flat Game fields (`metKids`, `parisUnlocked`, `nightSceneDone`, etc.) with scoped variables: Game scope (persist forever), Chapter scope (persist within day), Scene scope (reset on scene change). Enables save/load. Currently: `game.go` lines 39-71.
- [ ] `[P1]` **NPC State Machine** — Replace `onDialogEnd` callback + `dialogDone` bool + manual dialog swapping with named states ("default" → "post" → "strange"). Auto-transition after dialog. Currently: fragile closures in `setupCampCallbacks()`.
- [ ] `[P2]` **Item Ownership Tracking** — Add `owner` field to items. Track who has what: `"player"`, `"lily"`, `"curator"`, `"none"`. Enables clean give-item-to-NPC flows. Currently: `inv.hasItem("name")` string matching.

### Phase 3: Full Engine Architecture (major refactor)

- [ ] `[P1]` **Sequence System** — Replace nested callback chains (5+ levels deep in `setupCampCallbacks`) with a Sequence player. Each cutscene = list of steps: `{actor, action, data, sideEffects}`. Night scene becomes 10 declarative steps instead of fragile closures. Currently: `nightSceneArrival()`, `checkDay1Complete()` are callback hell.
- [ ] `[P2]` **Handler + Condition System** — Replace `setupCampCallbacks()` (130+ lines of closures) with declarative handlers: `{event: "click_npc", target: "Lily", condition: "hasItem(Flower)", action: playLilyFlowerSequence}`. Decouple triggers from actions.
- [ ] `[P2]` **Walk Zones** — Replace walk segments (line pairs) with polygon walk zones defined in scene JSON. More natural movement, easier to tune. Currently: 8 manual line segments per scene.
- [ ] `[P2]` **Save/Load System** — Serialize VarStore + inventory + scene + position to JSON. Requires VarStore first.
- [ ] `[P3]` **PDA/Map UI** — Travel map becomes multi-page UI: Map, Clue Book, Travel History. Like the retro PDA system.

### Architecture Mapping

| What Retro Does | What We Have Now | What We Need |
|----------------|-----------------|--------------|
| Handler+Condition | `setupCampCallbacks()` closures | Declarative handler registry |
| Sequences | `nightSceneArrival()` nested callbacks | Sequence player with steps |
| Game/Module/Page Variables | 15 flat `bool`/`int` fields on Game | `VarStore` with 3 scopes |
| Item Ownership | `inv.hasItem("name")` | `item.owner` field |
| Scene Data Files | 800+ lines hardcoded in scene.go | JSON scene files |
| Dialog Data Files | 500+ lines hardcoded in npc.go | JSON dialog files |
| NPC State Machine | `onDialogEnd` + manual swap | Named states with auto-transition |
| Walk Locations | `walkSegments` line pairs | Polygon walk zones |
| Save/Load | not implemented | VarStore serialization |
| PDA | simple map overlay | Multi-page UI system |

first code revie:
- [ ] `[P1]` in the first scene. after pp short monologe, i want to click on higgings to talk with him, now its automatic
- [ ] `[P1]` from higgings idle. take only the first row
- [ ] `[P1]` mouse size in screen is huge
- [ ] `[P1]` flip danny by 180 
- [ ] `[P1]` when lily is shy, higgings not shown up.
- [ ] `[P0]` i want the pp to look bigger then the kids regardin size.
- [ ] `[P1]` in lake remove the bg from the flower 
- [ ] `[P1]` no animation loading to grab flower 
- [ ] `[P1]` when i pick up the flower i need to give it to lily,currently the conversation change when i got in my inventory
- [ ] `[P1]` when goin to marcus room, pp is shown there for some reason even when we need to see only marucs freak out 
- [ ] `[P1]` rooms dots to enter(make it small radius), also i want an arrow to know we going to enter. tommy:(195,479),jake(441,441),marcus(820,435),lily(1077,403),danny(1243,503)
- [ ] `[P1]` change marcus possition in night and day to 646,398
- [ ] `[P1]` higgins sitting need to be in 1062,357, make both object bigger , and make sure pp is not standing on the table so radius to talk is from (86,748),(452,587)->(735,739),(796,562)
- [ ] `[P0]` i cant move between item in inventory!! remove the yellow arrows, make the cirle bigger, and then i will give you the cords to move between right and left
- [ ] `[P1]` map starting to look good! but we loose color on the map from every object 
- [ ] `[P1]` i cant go out from higgins office. radius (82,460)
- [ ] `[P1]` flying is still not working. im pressing on paris and nothing happen.
- [ ] `[P1]` remove bg from right/left arrow inside rooms 
- [ ] `[P1]` no need map to gro on screen, change it to the animation we make of giveing and taking map.
- [ ] `[P1]` assets\images\locations\camp\npc\kids\jake\npc_jake_idle.png take only the second line here
- [ ] `[P1]` remove bg from items, we see them in the invertory
- [x] `[P1]`  lily dialog is still wrong. the first time i click on her is like i gave the flower — FIXED: item-in-bag gate + cursor hint in pass 2 (see Reported 2026-04-16 pass 2 above).
- [ ] `[P1]` in the camp scnece i cant reach the inventory normaly, and if i do cant press on the righ kid
- [ ] `[P1]` in the night senarion, we need to remove the bg from both sleep, wake up and fire. put the fire in (622,573) and the pp in (335,591). when goin to marcus room, the pp is inside there for some reason
- [ ] `[P0]` we need to make the other kid ot be like lily is displayed. the frames of them are not good and loose colors like tommy and danny
- [ ] `[P1]` in marcus room, he need to be display in (666,561) and bigger in the size, same with the pink 
- [ ] `[P1]` in higgins office, he need to be in (1059,370) and no bg
- [ ] `[P1]` still no location in the map!! and travel to paris isnt working

### Reported (2026-04-16) - this session

- [x] `[P0]` Campfire sprite renders with its own background in the night scene. Asset: `assets/images/locations/camp/campfire_idle.png`. Sheet is 8x4 but loaded and flipbooked as a 32-frame cycle in `game/game.go:166-172` and drawn at scale 2.5 at (622,573). Color-key misses read as a visible halo. Fix: load row 0 only, bump inset, add aggressive color-key pass. — FIXED: `SpriteGridFromPNGCleanAggressive` + row 0 + inset 4 in `game/game.go`.
- [x] `[P1]` Lily first-click triggers flower dialog even when PP hasn't walked up holding the flower. Root cause at `game/player.go:532`: `altDialogFunc` checked unconditionally on normal click-to-talk, not gated on `inv.heldItem`. Fix: add `altDialogRequiresHeld bool` to `npc`, set true for Lily, require held Flower. — FIXED: gate in `canTriggerAltDialog`; Lily uses `altDialogRequiresItem="Flower"` with `altDialogRequiresHeld=false` so the handoff fires as soon as the flower is in the bag.
- [x] `[P1]` Higgins appears double on `camp_entrance` when talking. `npc_director_higgins_talk.png` is malformed (row 1 has 8 frames, row 2 has 6) but `game/npc.go:166` loads it as `4x2`, so cell slicing picks up half a neighbor figure. Fix: load row 0 only via `loadNPCGridRow(..., 8, 2, 0)` as a stopgap; regenerate the sheet to a clean 8x1 per PROMPTS.md. — FIXED: `loadNPCGridRow(..., 8, 2, 0)` in `newDirectorHiggins`; clean-regen tracked below.
- [x] `[P1]` Higgins office NPC is 280 px tall (~35% of screen) — too big vs PTP reference. Fix: drop bounds H to 225 and let scene `characterScale 0.9` in `camp_office` finish the job. — FIXED: `newOfficeHiggins` bounds 160x225 + `camp_office` characterScale 0.9.
- [x] `[P1]` No per-scene camera scale — cabin rooms and outdoor scenes share one size multiplier; PTP's pub shot is tighter than its park shot. Fix: add `characterScale float64` to `scene` struct, multiply PP + NPC draw rects by it, starter values in CHARACTERS.md. — FIXED: `scene.characterScale` + `drawScaled` on player/npc; values seeded for `*_room`, `camp_office`, city tight shots.
- [x] `[P1]` Camp scenes feel static. Existing "birds" in `game/scene.go:1285-1292` are 3-pixel rectangles, barely visible. Day 1 camp needs real ambient life (sprite birds + butterflies near lake + drifting clouds). Day-1-only so Day 2 grim tone lands. — FIXED: sprite-based `updateAmbient` / `drawAmbient` with bird + butterfly + cloud sheets, gated by `isCampOutdoorScene` + `day == 1`.

### Reported (2026-04-16) - this session (pass 2)

- [x] `[P1]` Higgins still renders small on `camp_entrance` next to PP (reported regression from pass 1). Root cause: bounds were 160x230 and the aspect-preserve draw produced a figure shorter than PP (170x235). — FIXED: `newDirectorHiggins` bounds bumped to 200x265 with Y=345 so the aspect-preserve lands in the 225-235 band from CHARACTERS.md.
- [x] `[P1]` Flower handoff to Lily broken — clicking Lily after picking up the flower did nothing (reported regression). Root cause: `altDialogRequiresHeld` was too strict; player had to manually select the flower in the inventory to "hold" it before Lily would accept. — FIXED: `canTriggerAltDialog` now also accepts `inv.hasItem(...)` for NPCs that set `altDialogRequiresItem` without `altDialogRequiresHeld`; Lily's setup flipped `altDialogRequiresHeld` off. `updateHover` lights up the cursor when the required item is in the bag.
- [x] `[P1]` NPCs do not face PP when talking — e.g. Danny stays facing the trees through his whole dialog. — FIXED: `npc.preTalkFlipped` snapshot + flip in `startNPCDialog`, restore in `wrapCb`; drag-onto-NPC path mirrors the same snapshot/restore so Lily's flower handoff leaves her turned toward PP.
- [x] `[P2]` Camp kid sprite sheets lose color on transparent background — default color-key (tol=8) left halos on pastel backgrounds. — FIXED: added `SpriteGridFromPNGCleanKids` (tol=16), routed every camp kid loader through `loadNPCGridKids` / `loadNPCGridRowKids`, normalized Tommy/Jake/Marcus mismatched idle/talk/strange pairs to 8x2, and tightened `eraseGridLines` (window ±2, RGB<50, outer-edge alpha gate) so sprite outlines stop getting eaten.

### Deferred (2026-04-16) - sprite regens

- [ ] `[P2]` Regenerate `npc_tommy_idle.png` + `npc_tommy_talk.png` as clean 8x1 sheets per PROMPTS.md so we can drop the 8x2 normalization stopgap.
- [ ] `[P2]` Regenerate `npc_tommy_strange_idle.png` + `npc_tommy_strange_talk.png` as matching 8x1 sheets.
- [ ] `[P2]` Regenerate `npc_jake_idle.png` (currently row 1 of 8x2) and `npc_jake_talk.png` as matching 8x1 sheets.
- [ ] `[P2]` Regenerate `npc_jake_strange_idle.png` + `npc_jake_strange_talk.png` as matching 8x1 sheets.
- [ ] `[P2]` Regenerate `npc_marcus_strange_idle.png` + `npc_marcus_strange_talk.png` as matching 8x1 sheets (currently authored as 8x2 and 5x2).
- [ ] `[P2]` Regenerate `npc_director_higgins_talk.png` as a clean 8x1 so we can drop the "row 0 only" hack in `newDirectorHiggins`.
- [x] `[P2]` Regenerate `npc_director_higgins_idle.png` as a clean single-row sheet so we can drop the "row 0 only" / mismatched-grid workaround. — FIXED: regenerated as a clean 7x1 strip at 172x384 per cell (matches `_talk.png` geometry so idle and talk render at the same on-screen size); `newDirectorHiggins` + `newNightHiggins` still call `loadNPCGrid(..., 7, 1)` unchanged. Built via `tools/stitch_higgins_idle.py` from 7 per-pose generations + edge-flood color-key. Stale duplicate at `assets/images/locations/camp/npc/npc_director_higgins_idle.png` deleted.

### Reported (2026-04-16) - this session (pass 3)

- [x] `[P1]` Higgins idle still reads as double on `camp_entrance`. Root cause: the idle sheet was a 2-row mess (row 0 = 8 frames, row 1 = 6) but `newDirectorHiggins` loaded it as `loadNPCGrid(..., 7, 1)`, so cell slicing divided the sheet into 7 wide-and-tall cells that each straddled a pose and its neighbor. — FIXED: regenerated the idle PNG as a clean `7x1` strip at 172x384 per cell via `tools/stitch_higgins_idle.py` (matches talk-sheet cell geometry so both animations render at the same size). `newDirectorHiggins` now loads it cleanly as `loadNPCGrid(..., 7, 1)`; `newNightHiggins` was already on the same path.
- [x] `[P1]` Lily flow still inconsistent — shy dialog and flower dialog could fire in the wrong order depending on scene re-entry because `lilyHinted` was a closure-local that reset whenever `setupCampCallbacks` ran. — FIXED: promoted to `npc.hintState` (per-NPC, survives reloads). Lily's altDialog is armed at setup but gated on `hintState == 1`, so first click is always shy, and flower handoff fires exactly once after the shy beat finishes.
- [x] `[P1]` PP visibly walked around inside Marcus's cabin during the night freakout beat. Root cause: phase 3 cleared `playerSleeping` before transitioning, so `Draw` fell through to `scene.drawActors(..., g.player)`. — FIXED: added `nightHidePlayer` bool set true on phase-3 entry and false on phase-4 exit; `Draw` calls `scene.drawActorsNoPlayer` while the flag is on.
- [x] `[P1]` Warm orange tint bled across the sleep/wake sprites making them look bronzed. — FIXED: `drawWarmTint` now skipped while `camp_night` + `playerSleeping`.
- [x] `[P1]` Cursor had no held-item state — pointer stayed on `cursorNormal` over empty space even while PP was carrying something. — FIXED: `updateHover` now defaults to `cursorGrab` whenever PP is carrying an item. No new cursor asset added — only the existing cursor PNGs are used.

- [] `[P1]` match every higgins design to the last idle.
- [] `[P1]` when the icon change to click it got white bg.
- [] `[P1]` tommy loose color and we can see the frames change,same for jake
- [] `[P1]` when we first want to talk to lily higgins not shown up as he should.
- [] `[P1]` make the kid be smaller then the pp 
- [] `[P1]` i want to click on the exac object and not a big radius around him. currently i want to talk to danny and marcus talk is jumping in 
- [] `[P1]` lily not getting the flower at all! i just talk to them and it jump to night. i want you to check the size of the fire and match it to the place it need to be around (581,592)-(702,594)
- [] `[P1]` pp sleeping is huge and got white bg 
- [] `[P1]` fix the door enter locations
- [] `[P1]` no need the map to fill the screen when get it. it need to be animation from giving and talking we already made 

### Reported (2026-04-17) - this session (pass 4)

- [x] `[P1]` Higgins sprite sheets were loaded with grid dimensions that did not match PROMPTS.md spec, causing the long-running "doubled Higgins" artefact on talk. — FIXED in `game/npc.go`:
  - `newDirectorHiggins` talk: `loadNPCGridRow(..., 8, 2, 0)` → `loadNPCGrid(..., 6, 1)` (per spec).
  - `newOfficeHiggins` idle: `6,2` → `7,1`; talk: `6,2` → `4,2` (per spec).
  - `newNightHiggins` talk: `4,2` → `6,1` (matched entrance).
- [x] `[P1]` Office Higgins was at (1062,400) size 160x225 — user wanted him higher and bigger. — FIXED: bounds to `(942, 357) 240x320` so foot lands around y=677 with head at y=357.
- [x] `[P1]` Night Higgins not clearly in the bottom-right corner. — FIXED: bounds to `(1120, 430) 200x260` so he sits at the campfire's right side with feet at ~(1220, 690).
- [x] `[P1]` Marcus in his room too small and off-position. — FIXED: bounds to `(526, 181) 280x380` in `scene.go` (bigger body, foot center lands near user-specified (666, 561)).

- [] `[P1]` we need to change the size of the objects. take a look in those pictures C:\go-workspace\src\bitbucket.org\Local\games\PP\london\background\the-pink-panther-passport-to-peril_7 and C:\go-workspace\src\bitbucket.org\Local\games\PP\london\background\pub. currently they fill the image to much imo 
- [] `[P1]` i want to add object to the bg to make it more alive 
- [] `[P0]` i dont want radius to talk to npc. we must click or give item to them when the mouse is in the frame of them. probably it because of the bg so we need to think of it. 
- [] `[P1]` higgins still not showing up when we talking to lily when she is shy 
- [] `[P1]` i cant pick up the flower... when i got it the story just jumped imidiatily to the night schene. not good! lily didnt got the flower and higgins didnt say it time to sleep
- [] `[P1]` when the icon change to click it got white bg.
- [] `[P1]` tommy loose color and we can see the frames change,same for jake
- [] `[P1]` when we first want to talk to lily higgins not shown up as he should.
- [] `[P1]` make the kid be smaller then the pp 
- [] `[P1]` i want to click on the exac object and not a big radius around him. currently i want to talk to danny and marcus talk is jumping in 
- [] `[P1]` lily not getting the flower at all! i just talk to them and it jump to night. i want you to check the size of the fire and match it to the place it need to be around (581,592)-(702,594)
- [] `[P1]` pp sleeping is huge and got white bg 
- [] `[P1]` fix the door enter locations
- [] `[P0]` no need the map to fill the screen when get it. it need to be animation from giving and talking we already made 
- [] `[P1]` higgings in the office is not in place at all
- [] `[P1]` when we fly to paris the pp is shown in the frame for some reason. also add a color of sky to make the flying more alive 

### Reported (2026-04-17) - this session (pass 5)

- [x] `[P0]` NPC click radius too wide (Danny snap-stealing Marcus clicks, clicks on empty ground triggering dialog). — FIXED: `containsPoint` no longer expands `bounds` by ±70/±50 px; the hit test is now strictly inside the NPC's authored rect (`game/npc.go`). Silent + hidden NPCs skipped in `scene.checkNPCClick` and `updateHover`.
- [x] `[P0]` Travel Map "grows onto the screen" when Higgins hands it over — broke the give/take-map rhythm. — FIXED: removed the `mapRevealing` zoom tween trigger in `giveMapItem`. The map item just drops into inventory; the existing take-map animation is the whole handoff.
- [x] `[P1]` Higgins didn't appear when Lily was shy. — FIXED: added hidden `Director Higgins` NPC to `camp_grounds` at `(910, 400, 200, 212)` plus `higginsLilyHintDialog`. He unhides the moment Lily's shy dialog finishes and delivers the flower clue. `npc.hidden` bool added + honored by `drawScaled`, `checkNPCClick`, `updateHover`.
- [x] `[P1]` Flower-to-Lily felt like the story jumped to night immediately. — FIXED: `checkDay1Complete` now plays a three-line Higgins bedtime beat on `camp_grounds` before transitioning to `camp_night`, latched behind `day1BedtimeStarted` so the flower callback can't double-fire it.
- [x] `[P1]` PP sleeping sprite too big and showed a white rim. — FIXED: draw scale 1.8 → 1.1; loaders swapped to `SpriteGridFromPNGCleanAggressive` with inset 4 to strip the cream-white halo.
- [x] `[P1]` PP rendered inside the Paris flight cutscene. — FIXED: `Draw` now hides player actor while `sceneMgr.currentName == "airplane_flight"` (same branch that already hides him during Marcus-room phase 3).
- [x] `[P1]` Office Higgins "not in place at all". — FIXED: bounds top-left snapped to the user's spec `(1062, 357)` with size `220x280`.
- [x] `[P1]` Camp kids read the same size as PP — should be noticeably smaller. — FIXED: all five camp_grounds kid bounds clamped to `150x180` (PP stays at `170x235`) and Y adjusted so feet land on the existing walk segments.

--paris!!--
- [] `[P1]` npc not placed on the ground they are at around 588 y 
- [] `[P0]` direction! remove the left move to open the map. no need for that, i want to create now a story and object we need to collect before we enter the louver.
so we need a bagguete and give more ideas just like every retro pp game had. it need to be a quest game after all.
- [] `[P1]` can you generate bg people that sit on the chairs and drink coffee in loop regardless to what we doing?

new PR
- [] `[P1]` logic walking, in the first screen you can walk freely but as a result you walking on unlogic places
- [] `[P1]` higgings idle is not the same design we said to use. as a result there is a huge different between the idle and talking
- [] `[P1]` talking logo got white square around him
- [] `[P1]` when we goin the the camp, pp is just walking left. i want him to walk back for a few seconds(make him srink for a seconds ) and then change screen
- [] `[P1]` tommy framing change is not good
- [] `[P1]` jake got white square near to his leg,same as marcus
- [] `[P1]` tommy talking sprite is in different color then the other sprites
- [] `[P2]` can we make the talkin to work with both click and fluent talking. like if we will add a talking to every npc i should be just like we hear them 
- [] `[P1]` finally higgins is arrived when lily is she. we need to generate a walking back sprite for this. 
- [] `[P0]` isnt all object are huge? we bearly we see other objects 
- [] `[P1]` i want to change danny talking sprtie. he got a wired stuff in hes legs 
- [] `[P1]` when picking the flower i want the picking animation we already got.and also its very hard to find the right place to click in order to pick it up
- [] `[P1]` fire animation is huge, it need to be from like (577,591)-(700,590). pp need to walk to the sleeping point and then the animation need to start(sleeping). also the sleeping got a white bg. also the wake up
- [] `[P0]` marcus got a bg aroung him and hes huge in his room. we dont see the text for the entire senario.
- [] `[P1]` i already place the right points where are the doors of each cabinet. fix it 
- [] `[P1]` marcus freak out sprite is too fast. in addition add another freakour sprite that will run if we dont click on him for a few seconds. also check the size in the room between the pp and marcus
- [] `[P1]`walking in the camp should also have a logical routes. i made it already check the cords. 
- [] `[P1]` line to walk to higgins office is not easy to find the right dot