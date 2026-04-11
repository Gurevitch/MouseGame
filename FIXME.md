# Pink Panther: Camp Chilly Wa Wa — Known Issues & Fixes

> **Reference:** See [STORY.md](STORY.md) for full story flow and design.
> **Progress:** See [STATUS.md](STATUS.md) for implementation status.

---

## How to Use

Add issues below as you find them. Mark priority: `[P0]` critical, `[P1]` important, `[P2]` minor.
When fixed, move to the **Resolved** section with the date.

---

## Open Issues

### Story / Flow

- [ ] `[P1]` Postcard not added to inventory after Curator dialog (item pickup missing)
- [ ] `[P1]` No Marcus healing flow — giving postcard to Marcus should cure him
- [ ] `[P2]` Higgins bedtime dialog uses existing entrance sprite — needs new camp sprite (see STATUS.md assets)
- [ ] `[P2]` Airplane cutscene uses PP standing idle — needs sitting-in-plane sprite

### Scenes / Navigation

- [ ] `[P2]` Kid rooms (Tommy, Jake, Lily, Danny) have no NPCs inside — should have kids in their rooms on Day 2

### Travel Map

- [ ] `[P2]` Buenos Aires and Rio pins are close together — landmarks may overlap visually

### Assets

- [ ] `[P1]` Need Higgins camp sprite (idle + talk) for bedtime scene
- [ ] `[P1]` Need PP airplane idle sprite for flight cutscene
- [ ] `[P2]` No airplane background — currently using paris_clouds.png as fallback

---

## Resolved

| Issue | Date | Notes |
|-------|------|-------|
| Map used dots instead of landmarks | 2026-04-08 | Replaced with landmark images from ui/landmarks/ |
| Night scene too simple (instant Marcus room) | 2026-04-08 | Reworked to multi-phase: campfire sleep → Marcus freakout → wakeup |
| No airplane transition before cities | 2026-04-08 | Added airplane_flight scene with 4-second cutscene |

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
- [ ] `[P1]`  lily dialog is still wrong. the first time i click on her is like i gave the flower 
- [ ] `[P1]` in the camp scnece i cant reach the inventory normaly, and if i do cant press on the righ kid
- [ ] `[P1]` in the night senarion, we need to remove the bg from both sleep, wake up and fire. put the fire in (622,573) and the pp in (335,591). when goin to marcus room, the pp is inside there for some reason
- [ ] `[P0]` we need to make the other kid ot be like lily is displayed. the frames of them are not good and loose colors like tommy and danny
- [ ] `[P1]` in marcus room, he need to be display in (666,561) and bigger in the size, same with the pink 
- [ ] `[P1]` in higgins office, he need to be in (1059,370) and no bg
- [ ] `[P1]` still no location in the map!! and travel to paris isnt working