# Pink Panther: Camp Chilly Wa Wa — Known Issues & Fixes

> **Reference:** See [STORY.md](STORY.md) for full story flow and design.
> **Progress:** See [STATUS.md](STATUS.md) for implementation status.

---

## How to Use

Add issues below as you find them. Mark priority: `[P0]` critical, `[P1]` important, `[P2]` minor.
When fixed, move to the **Resolved** section with the date.

---

## Open Issues

### Reported (2026-06-21 — biker/pigeon/worshippers white-box bg, REAL fix)

- [x] Earlier diagnosis was WRONG: `biker.png`, `npc_pierre_pigeon_lands.png` and
  `people_pray.png` are NOT transparent - they ship opaque near-white (every pixel
  alpha=255, corners ~233-253), so the RAW ambient load drew a white box around each.
  FIX: their loaders (newAmbientBiker / newAmbientPigeonFlyUp / newAmbientWorshippers)
  now use `loadAmbientStripKeyedTol(..., 40)` - the EDGE-CONNECTED white key strips the
  background while protecting enclosed whites (the biker's striped shirt). Verify the
  biker's shirt survives in-game; if the key leaks into it, re-export that sheet with
  true alpha transparency.

### Reported (2026-06-21 — opening flow now starts at the airstrip + day2 art)

- [x] Game now STARTS at camp_landing (was camp_entrance): scene.go start scene + the opening
  monologue/walk-in moved to camp_landing. New day-1 flow: landing → monologue → up-arrow to the
  ENTRANCE (Higgins) → grounds. The landing's up-arrow target is conditional - camp_entrance on
  the first arrival (!paris_done), camp_grounds on city returns (the dark-landing flow).
- [x] day2 room art LANDED (user) — renamed to the folder scheme (day2/jake_room.png etc., same
  filename as day1) so moodBG finds them; day2/camp_grounds.png + day2/camp_landing.png + the 4
  rooms wired. Marcus's room added to applyCampMood (gated !marcus_healed, so the heal still
  brightens it). day3/camp_landing.png renamed in too. NOTE: a stray day3/camp_camp_dark.png is
  unused (day3/camp_grounds.png is the canonical full-dark grounds) - delete or rename if intended.

### Reported (2026-06-21 — Jerusalem art structure refinement)

- [x] Shimon art LANDED (`npc/wall/npc_shimon.png`, full-body 6×2); loader repointed.
- [x] Reorganised Jerusalem NPC art into `npc/wall/` (plaza+Wall) + `npc/market/` (souk).
- [x] Spice + coffee seller: now load SEPARATE idle/talk sheets and render FULL BODY (bounds
  150→230 tall), per user; kid also SEPARATE idle/talk. Fallback to Paris placeholders until art.
- [x] GIVE one-shots wired on every giving NPC (Shimon/spice/coffee/bagel/praying) + PP take
  beats in each trade callback (§8b both-sides); separation-fence prop wired in the plaza
  (`props/fence.png`, no-ops until art). All queued at EXTRA_PROMPTS §JERUSALEM.

### Reported (2026-06-21 — PR Step C: Jerusalem chapter, item 26)

- [x] Built the full Jerusalem daisy-chain (game/jerusalem.go rewrite), retiring the trivial
  "Miriam hands the coin rubbing" stub and the old sardine/cat/jar design:
  Shimon (plaza fence, directs up→Wall / left→market) → spice seller gives **Cardamom** →
  coffee seller (souk centre) trades it for **Coffee** (sits + teaches Jerusalem) → bagel
  seller trades for a **Bagel/ka'ak** → praying man at the Wall (idle=praying, turns to talk)
  takes the bagel + gives a **Note Paper** → Shimon gives a **Pen** → the Wall-crack hotspot
  writes + places the note (`jer_note_placed`) → Shimon gives the **Coin**.
- [x] Coin is Jake's anchor now — Jake's heal repointed from "Coin Rubbing" to "Coin" (game.go).
- [x] Return flight gated on `jer_note_placed` (travel_map.json camp pin) — can't leave Jerusalem
  until the note is in the Wall.
- [x] Scenes given Paris-style walk lines (#23); market uses the user's #25 coords (entry far→
  centre, up-square exit). Worshippers MULTIPLIED at the plaza + Wall (#22).
- [x] New items (Cardamom/Coffee/Bagel/Note Paper/Pen/Coin) in items.json; PP write_note/put_note
  one-shots wired with grab fallback. All Jerusalem NPC/item/one-shot art queued at
  EXTRA_PROMPTS §JERUSALEM (NPCs borrow Paris/camp sheets, icons/one-shots no-op, until it lands).

### Reported (2026-06-21 — PR Step B: dark camp + Marcus arc, items 20-23)

- [x] #20 Marcus room entry higher — marcus_room spawnY 360→330.
- [x] #21 Graded dark camp — `applyCampMood` now grades (campMoodLevel 0/1/2): mid-dark
  post-Paris, fully-dark from the Jerusalem leg. Swaps camp_grounds + camp_landing + the
  4 non-Marcus cabin interiors, each falling back through available art. Marcus room darkens
  via its existing day/night bg.
- [x] BG folder reorg (user 2026-06-21) — camp_grounds + camp_landing + the 5 rooms moved into
  `camp/background/day1/` (normal); `day2/` = mid-dark, `day3/` = full dark, SAME filenames in
  each. `day3/camp_grounds.png` is the moved `camp_dark.png`. moodBG picks the folder by grade
  (day1→day2→day3) and falls back down. Scene JSONs + sceneAltBGs repointed to day1/. Art queued
  at EXTRA_PROMPTS §DARK-CAMP/§DARK-ROOMS (new folder naming).
- [x] #22 Marcus rude pre-heal dialog — marcusPostStrangeDialog now irritable ("What do you
  want NOW?"). Postcard hand-over plays a `receive_postcard` NPC one-shot (art §MARCUS-POSTCARD,
  no-ops until it lands).
- [x] #23 Sleep now STICKS — the strange-alt "freakout" punctuation (altIdleGrid/altIdleAfterSec)
  was still firing after the heal; the heal callback now disables it + clears any in-flight alt
  cycle, so the sleeping idle persists. Sleepy dialog reworded ("I'm so tired... maybe tomorrow").

### Reported (2026-06-21 — PR Step A: Paris/camp fixes, items 1-19,25)

- [x] #1 Lake dock still low — raised walkSegments another ~40px (camp_lake.json). F3-verify.
- [x] #2 Office Higgins talk BLINKS — root cause: per-frame foot anchoring jittered on his
  waist-up bust (FOOT drift 121px). New `fixedFootAnchor` (npc.go drawScaled) pins content
  center-X + bottom for seated NPCs; set on office Higgins.
- [x] #3 Office Higgins shrunk (H 200→185, foot kept).
- [x] #5 Biker — lane is y=750 (foot anchor) on the cobbles; biker.png IS transparent + loads
  RAW (gap-detect confirms transparent gaps), so no bg box from the asset. Verify in-game.
- [x] #7 Colette approached from the other side (approachRight→approachLeft); she's not
  fixedFacing so she turns to face PP.
- [x] #8 PP shrinks at Pierre again — restored Pierre's recede onClickOverride (the clean
  version with the recedeHeld in-place-talk guard from #11); Margaux stays standard so no
  inter-NPC size pop. depthScale can't do it (PP's y is clamped on the street).
- [x] #9 Bakery walk lines added (user coords). NOTE center/foot+clamp — F3-tune (paris_bakery.json).
- [~] #11 Rolling-pin pickup "jump" — couldn't pin the exact cause from code (stand-offset vs the
  +60 grab draw-offset snap); needs F3 repro in-game before a fix to avoid regressing the pickup.
- [~] #12 Pierre white pocket between board/arm — the connected key reads it right; the pocket is
  ENCLOSED (edge-flood can't reach it). Needs the §PIERRE-BOARD re-roll (cream canvas), not a code fix.
- [x] #13 Henri coffee→confiture — reordered to clean BRING-then-PICK-UP: coffee handed in the
  pre-dialog handoff, then Henri give_jam → PP get_jam (chained) → Confiture added on completion.
- [x] #14 Louvre arrival monologue — PP faces FRONT (dirRight→dirDown).
- [x] #10 Poulain: PP showed his back — new `ppFacePlayer` (she's behind the back counter) makes
  PP face the camera for her dialog/receive.
- [x] #15 PP give heel blinks — §GIVE-HEEL re-roll queued (arm+heel cross cells). Code reads it right.
- [x] #16 Pigeon fly-up — npc_pierre_pigeon_lands.png IS transparent + loads RAW; no bg box from
  the asset (verify in-game).
- [x] #17/#18 PP give/get item sheets — verified: give/receive sheets load keyed (white stripped);
  `PP get bagguette/jam.png` gap-detect 1×8 clean (8 frames, correct grid). The "broken" look was
  the facing (fixed via #10). No re-roll needed.
- [x] #19 Paris pin off after the postcard (relevantWhen → paris_done==0); #25 the map never
  offers PP's current region (travel_map.go hitTest skips it via travelRegionOf).

### Reported (2026-06-20 — bug-sweep PR, Step 1, items 1-20)

- [x] #1 Jake talks too fast — `newJake` talkFrameSpeed 0.10→0.18 (room Jake inherits). npc.go.
- [x] #2 Lake dock: PP's foot floated at (844,608) below the planks — raised the dock
  walkSegments ~40px (camp_lake.json). F3-verify the foot lands on the planks.
- [x] #3 Office: PP stood on the trash bin — new per-NPC `approachGapX` (npc.go/player.go);
  office Higgins gap 280 so PP stops left of the bin, facing right.
- [x] #5 Thrown map missed PP — new `toPlayer` flag on `tween_item` retargets the projectile
  to PP's runtime paw (sequence.go/sequence_loader.go + higgins_give_map.json).
- [x] #6 Removed the "m"-opens-map shortcut (game.go HandleKey).
- [x] #7 Margaux foot → (656,639); bounds {617,494,78,145}. npc.go + hit-test.
- [x] #8 Flower pot shrunk + moved beside Pierre (bounds {868,562,82,88}); fly-up spawn moved
  to (909,565). game.go.
- [x] #9 Biker lane y 735→750 (scene_ambient.go). NOTE: biker.png is already transparent and
  loads RAW, so no background box comes from the asset — couldn't reproduce a "bg box."
- [x] #10 Margaux talk speed 0.13→0.22. npc.go.
- [x] #11 Pierre "jump-back" — removed Pierre's recede onClickOverride; he now uses the
  standard walk-up-and-talk (depthScale handles perspective). game.go.
- [x] #12 Camille moved up (bounds Y 384→360). npc.go + hit-test.
- [x] #13 PP "disappeared" after EVERY bakery NPC — the talk-stand row put his body box into
  the full-width top-wall blocker {0,0,1400,280}, which shoved him off-screen. paris_bakery
  minY 200→290 so his top-left can't enter the blocker; Margaux's recede override removed too.
- [x] #14 Giving Pierre an item played his "portrait" sheet — Pierre's `give` one-shot
  (used by playHandOff as the take-fallback) was the painting-display art; removed it. npc.go.
- [~] #15 "PP give heel" blinks — the extended arm+heel cross cell borders (jitter_audit:
  CONTENT CROSSES 18-30px). Layout re-roll queued at EXTRA_PROMPTS §GIVE-HEEL. Functional.
- [x] #16 Flying-pigeon "bg" — `npc_pierre_pigeon_lands.png` is transparent and loads RAW;
  no background box comes from the asset. (No change needed; flag if it still shows in-game.)
- [x] #17 Pencil inventory icon too big — new per-item `iconScale` (items.json/inventory.go/
  item_registry.go); charcoal pencil iconScale 0.6.
- [x] #18 camp_landing exit → arrow UP, hotspot center (1194,303); added road walkSegments +
  a camp_landing→camp_grounds waypoint walk so PP follows the road, not the side. game.go/json.
- [x] #19 Marcus post-heal — fixed the dialog revert (onDialogEnd now guards on marcusHealed)
  and wired a go-to-sleep one-shot + sleeping-idle so Higgins's "sleeping soundly" line is
  true. ART LANDED 2026-06-20: `npc_marcus_going_to_sleep.png` + `npc_marcus_sleeping.png`
  (sitting-doze among his drawings), both GAP-DETECTED 1×8 clean; loader points at those names.
- [x] #20 Room Jake + room Marcus too big — shrunk (Jake H 245→200, Marcus H 205→185, feet
  kept). npc.go + hit-test.
- [ ] #4 Office Higgins talk last frame "disappears" — DEFERRED to verify the office_talk
  sheet's trailing cell in a playtest before trimming the loaded frame count.

### Reported (2026-06-15 — playtest batch, 20 items; grouped fix sweep)

**Group 1 — camp office / Higgins:**
- [x] #2/#3 office double-click + stuck at entrance + "standing on air" — the
  scene had no walk-Y range (min/max 0 → engine defaults 265-395, foot max
  665), so walking up to Higgins raised PP's 270px body-box into the top wall
  blocker {500,0,900,360}; that shoved him back to the entrance and fired a
  false blocker-arrival (eating the 1st click), and kept his feet floating.
  FIX: `camp_office.json` minY/maxY → 420/520 (foot 690-790, the floor strip),
  so PP stays below the wall blocker and the desk blocker gives a clean
  first-click arrival on the floor.

**Group 2 — bakery PP "disappears":**
- [x] #7/#10 PP vanished after talking to ANY bakery NPC — PR#12 moved PP's
  stand row up into the seated patrons' bust band (foot ~470-480), and the old
  #27 hack forced every patron to draw IN FRONT of PP (drawFootY=900), so a
  patron bust swallowed him. FIX (game.go bakery setup): patrons' sort-foot
  pinned BELOW PP's min foot (900→400) so the roaming PP always renders on top
  of the seated regulars; he's always above the tablecloth line so this doesn't
  put him "on the cloths." Poulain unchanged (renders behind the counter).

**Group 3 — Paris give / trade flow:**
- [x] #11/#17 baguette / heel "bring broken" — Pierre & Margaux register only a
  `give` reach one-shot (no `receive_*`/`receive_item`), so `playHandOff`'s NPC
  half was an instant no-op. FIX (player.playHandOff): fall back to the NPC's
  `give` reach when no receive anim exists, so they visibly take the item.
- [x] #12 confiture → PP "jumped back to the main road" — Pierre's
  onClickOverride re-ran walkToAndDo(690,510)+recede on every click; on the
  stage-2 click PP (already recede-held at Pierre, logical y drifted up by the
  recede) walked back down to the road then re-receded. FIX: factored the
  conversation into `talk()`; when `recedeHeld` is already set, talk in place.
- [x] #13 press pass → Claude played a "pick" anim + PP stuck shrunk — (a)
  `giveAnimKeyForItem` had no "Press Pass" → fell to the generic grab/"pick"
  frames; mapped it to the flat-paper `postcard` give sheet. (b) recedeHeld
  (carried from Pierre) was only released by setTarget, so walking to Claude
  left PP shrunk; `walkToAndInteract` + `walkToTalkPos` now releaseRecedeSmooth
  when recedeHeld, so walking to any other NPC grows PP back.

**Group 4 — pencil pickup + inventory softlock:**
- [x] #19/#20 pencil never entered the inventory + PP stuck — the pot pencil
  used `playAction(stateGrabbing, cb)`, but player.update force-resets any
  non-talking state to idle every frame while !moving, killing the grab before
  its add-item callback fired. FIX: register a generic `grab` one-shot and use
  the guaranteed `playOneShot` (like the rolling pin) so the item is always
  added and PP returns to idle.
- [x] #19 inventory L/R arrow bands opened the (single) Travel Map — with one
  item the paging bands were skipped and any in-oval click selected it. FIX
  (inventory.handleClick): clicks in the arrow bands are consumed without
  selecting even when len(items)==1.
- [x] #9 rolling-pin pickup — dropped the grab pose further (offset 30→60).

**Group 5 — positioning (bottom-centre dots, Poulain convention):**
- [x] #5 Margaux → (559,593) bottom-centre → bounds {520,448,78,145}; given
  Pierre-style walk-to-the-line + recede onClickOverride ("act the same way").
- [x] #8 Camille nudged left + down → {470,384,...} (test bound synced).
- [x] #14 flower pot → bottom-centre (947,745) → {882,600,130,145} (bigger,
  front of scene); pigeon fly-up spawn moved to match (1118,590 → 947,610).
- [x] #15/#16 Beaumont brought forward + bigger ({546,450,150,290} →
  {520,490,165,315}) and gallery characterScale 0.7→0.85 so the meeting reads
  closer; PP foot-aligns beside him. (TUNE in playtest.)
- [x] #4 biker lane restored "as was before" (755 → 735).

**Group 6 — sprite backgrounds:**
- [x] #4-bg biker + #18 flying pigeon — both sheets already ship TRANSPARENT
  and load raw (SpriteGridFromPNGRaw), so the in-game box is gone; verify.
- [x] #1 PP pick-flower "looks at the other side" — the grab_flower sheet leans
  LEFT toward a daisy on PP's left, but walkToFloorItem stood PP to the item's
  LEFT (daisy on his right) so he bent away. FIX: flower floor item now
  `standRight: true` (PP stands to the daisy's right, daisy on his left).
- [x] #1 PP pick-flower outer white halo — `gridFramesConnected` keyed at tol 8
  (soft fringe); added `gridFramesConnectedTol` and load grab_flower at tol 36.
- [~] #1 PP pick-flower BLINKING — re-generated 2026-06-20, STILL broken (same
  layout flaws). `go test ./engine -run ContentGrid` gap-detects the new sheet
  1×6 but cell 0 is an 82px sliver = the daisy lying DETACHED at the far left of
  frame 1, so frame 0 is "just the daisy, no panther" and every PP pose shifts a
  slot → the blink. `jitter_audit` also shows the 6 PP poses TOUCHING (cross
  borders 13-16px, "2 cells touch both edges"), so a fixed grid can't cut it
  either. Engine can't recover 6 clean frames from a detached object + figures
  with no gaps. Pickup still works functionally (item is added) but looks wrong.
  Re-roll sharpened at EXTRA_PROMPTS §FLOWER-PICK: daisy in PP's paw EVERY frame
  (never a separate ground object at the cell edge) + ≥15px gap between the 6
  poses so none touches its neighbour.
- [ ] #6 Pierre "missing board" frame — white-on-white chroma-key (his easel
  canvas is pure white → the edge key eats it where it abuts the bg). Re-roll
  queued at EXTRA_PROMPTS §PIERRE-BOARD (cream canvas).

### SOFTLOCK FIX (2026-06-12) — couldn't travel to Jerusalem (or any city) after a heal

- [x] Root cause: a travel pin must be **unlocked AND relevant** to be a
  travel target; each city's `relevantWhen` reads `vars.game.<id>_unlocked`,
  but those vars (jerusalem/tokyo/rome/rio/mexico) were **defined and never
  written** - heal callbacks only called `travelMap.setUnlocked(scene)` (the
  pin's bool), not the var. So after healing Marcus the Jerusalem pin lit up
  but couldn't be clicked to travel. (Paris worked only because
  `paris_unlocked` is synced from a Go bool.) FIX: `setUnlocked` now also
  mirrors `vars.game.<id>_unlocked`, so every existing unlock call fixes its
  city at once. Live flow Paris→camp→Marcus→Jerusalem verified reachable.
  NOTE: save/load only RE-restores the paris/jerusalem/camp pins on load
  (saveload.go) - tokyo+ pins aren't re-applied after a reload yet (pre-
  existing gap, separate from this live-flow fix).

### Reported (2026-06-12 — Marcus strange idle: day + night variants)

- [x] User split the strange idle into `npc_marcus_strange_idle_day.png` +
  `_night.png`. Wired: newRoomMarcus loads both (`strangeIdleDay/Night`);
  new `npc.setStrangeVariant(night)` picks the one matching the cabin bg.
  Hooked into `setSceneAltBG` (the single bg chokepoint) so the night
  cutscene → night variant and Day-2 → day variant (Day-2's bg swap now
  routes through setSceneAltBG too). Both sheets gap-detect clean; manifests
  + content-grid test updated. strange_alt stays the periodic fidget.

### Reported (2026-06-12 — pigeon lady becomes a quest NPC)

- [x] Street felt too crowded + the pot pigeon needed an owner → the
  background "crumb lady" ambient is promoted to a real NPC, **Madame
  Margaux**, on the LEFT of paris_street (Pierre-sized, x=230 foot 645). PP
  now brings the day-old Baguette Heel to HER (not Pierre); she lures the
  flower-pot pigeon off (shared `clearPotPigeon`: swaps the pot to the
  pencil + flaps the pigeon up via the existing fly-up sheet). Pierre is
  done questing after the press pass. The two street-density ambient stubs
  (accordion + crumb lady) were dropped. Camille's hint + the pot dialog now
  point to Madame Margaux. Art queued at §PIGEON-LADY (idle required, give
  optional) — wired + clickable now, invisible until the sheet lands.

### Reported (2026-06-12 — 30-item playtest batch; plan replicated-dazzling-parrot)

**Batch A — scene/position/speed:**
- [x] #4 office blocker added {278,588,131,101}.
- [x] #11 Poulain → bottom-center dot (726,318): bounds {641,173,170,145}.
- [x] #12 patrons "stand on the table" — approachYOverride 405→210 + bakery
  minY 400→200 so PP stands in the aisle BEHIND the front tables. F3-verify.
- [x] #14 Camille talk 0.10→0.22, bounds up (370→352); sketch one-shot now
  uses a HOLD (new npc.playOneShotAnimHold) so the reveal lingers ~1.2s.
- [x] #15 Bernard/Yvette talk 0.10→0.20; Bernard down (355→372).
- [x] #22 Beaumont too small under the 0.7 louvre scale → bounds H 205→290
  (foot kept), W→150.
- [x] #10 Pierre shrink (84×156→78×145) + CONNECTED key (was losing colour).

**Batch B — anim mechanics:**
- [x] #1 Lily flower → PP faced front: re-assert side-facing `dir` at the
  moment talk begins (give one-shot was leaving it stale). Fixes all gives.
- [x] #3 wake-up idle jump: resume-snap aligned to the wake render spot
  (foot 565, was 650).
- [x] #16 rolling-pin pickup: face PP left (toward basket) + drop the grab
  pose 30px so the reach lands in the basket.
- [x] #17/#19 give "broken"/coffee "too fast" — were TIMING (sheets are
  fine): playHandOff give 0.8→1.3 + per-give override.
- [x] #20 get confiture "not working" — get_jam/give_jam 1.0→1.6/1.5.
- [x] #21 give-to-Pierre: held item rode the cursor through the dialog
  (onClickOverride skipped the heldItem clear) → cleared at hand-off start;
  PP walked to foot 790 not Pierre's 645 → walk to centre-Y 510 (foot 645).
- [x] #25 walk-left stayed shrunk: a hotspot click now releases the recede
  so PP grows heading to the exit.
- [x] #18 Poulain handed items away from PP → face her toward PP on give.

**Batch C — Paris quest-chain reorder (full, softlock-proof):**
- [x] #24 Press pass is no longer a silent key — PP HANDS it to Claude, who
  consumes it and sets `louvreUnlocked`; the Louvre hotspot checks the flag.
- [x] #27/#29 Pot pigeon: leaves when PP gives Pierre the day-old Baguette
  Heel (Camille's hint → Poulain gives heel → Pierre shoos). Pot swaps to
  the pencil; an optional fly-up ambient (new ambientFlyOff kind) plays when
  its art lands (§AMB-PIGEON). Clicking the pot while holding an item no
  longer pockets it ("disappeared").
- [x] #30 softlock: chain is now strictly linear (rolling pin→Poulain→Henri→
  Pierre→press pass→Claude→Beaumont→Camille pencil via heel→sketch→postcard);
  no out-of-order pencil/sketch/portrait. The old easel pigeon-critic +
  mini_portrait beat is removed (the heel's job is the pot pigeon now).

**Batch D — color-key:**
- [x] #13 Henri tol 8→16 (clears his fringe; measured +3.4k bg px).
- [x] #7 biker tol→40 (shaves the edge halo). Interior white pockets inside
  the bike frame can't be reached by the edge key → transparent-bg re-roll
  queued (§BK2).

**Batch E — art:**
- [x] #9 PP jump-back — `PP jump back.png` LANDED, fitted (gap-detect OK),
  wired into the biker bump.
- [x] #23 Camille lost-pencil — `cafe_patron_camille_lostpencil.png` LANDED,
  wired as her `lost_pencil` one-shot (plays on the ask).
- [x] #26 Pierre idle/talk — split sheets LANDED, loader prefers them.
- [x] #18 Poulain give — re-rolled to a baguette, LANDED + code flip.
- [x] #7 biker transparent bg — LANDED, loader switched to raw.
- [ ] #2 Marcus room strange idle — standalone re-roll prompt now at
  §MARCUS-STRANGE-IDLE (no smiling, no touching figures). Interim still in
  place (strange_alt frantic frames + connected key) until it lands.
- [x] "Marcus too scary" pullback (2026-06-12) — dialled the freak-out back
  to eerie-sad, not horror: tremble amplitude/frequency halved (1.1/0.9 at
  ~20Hz, was 2.2/1.8 at ~43Hz), strange-idle cadence 0.14→0.26, punctuation
  2s→4.5s, and §MARCUS-STRANGE-IDLE prompt rewritten to "a little off /
  faraway" (no bloodshot eyes, sweat, manic shaking).
- [x] "Marcus MORE strange/freak out" (2026-06-12) — code intensity boost,
  art-independent: (1) NERVOUS TREMBLE on all strange NPCs (buzzing summed
  sines on the drawn rect; hit-test stays steady); (2) Marcus's strange idle
  IS the frantic scribbling loop now (not the smiling sheet); (3) faster
  manic cadence (new per-NPC strangeIdleFrameSpeed=0.14) + freakout
  punctuation every 2s (was 5). §MARCUS-STRANGE-IDLE prompt rewritten to a
  genuinely disturbing freak-out (wide bloodshot eyes, trembling, manic
  scribbling, jolting frames).
- [x] #29 pigeon fly-up — DONE: reuses the existing transparent
  `npc_pierre_pigeon_lands.png` (a perched→takeoff strip) via the
  newAmbientPigeonFlyUp ambient. Plays when Pierre shoos the pot pigeon.
- [x] #5 map pull — already wired (§PM1, art pending).

### Reported (2026-06-12 — travel-map pocket beat)

- [ ] Pull-map-from-pocket sprite before the map screen — WIRED, art
  pending: every map-open path (Travel Map inventory click, held-item
  drop, all six city street hotspots) now routes through
  `Game.openTravelMap`, which plays a `pull_map` one-shot (~0.9s) and
  opens the globe when it finishes. Until the §PM1 sheet lands on
  `assets/images/player/PP pull map.png`, the one-shot no-ops and the map
  opens immediately as before. Prompt queued at EXTRA_PROMPTS §PM1.

### Retro plan items 3-5 (2026-06-12, from the retro_frames design review)

- [x] #3 Character presence — paris_street `characterScale` 1.0 → 1.15
  (PTP reference shows PP at ~35-40% of screen height; ours was ~26%).
  Scales PP + NPCs together so the tuned relative sizes survive; feet stay
  planted (both draw paths anchor at the foot line). EXPERIMENT on this
  scene only — if it reads right in playtest, roll out per scene.
- [x] #4 Iris wipe transition — scene changes now use a retro iris instead
  of the plain alpha fade: opaque black with a soft-edged circular hole
  that closes onto PP on the way out and reopens from his spawn in the new
  scene (scene.go: newIrisMask + drawTransition; fadeAlpha doubles as wipe
  progress, plain fade kept as fallback if the mask texture fails).
- [ ] #5 Street density — two flavor ambients pre-wired crow-style on
  paris_street (accordion player at x≈120, pigeon lady at x≈1080): they
  auto-appear when their art lands, prompts queued at §AMB5/§AMB6 (renamed
  from AMB3/4 — those numbers were already the camp crow / biker). New
  generic `newAmbientSway` constructor for any future walk-line flavor
  figure.

### Reported (2026-06-12 — night playtest: Marcus night, walk front, Higgins hand, bakery)

- [x] Marcus strange idle (night) smiling + losing colors — FIXED both
  sides: strange_alt now loads with the CONNECTED key (the global key ate
  ~7.5k px of eye whites/teeth/highlights), and the relaxed/smiling poses
  (frames 5-7 per row) are filtered out so only the frantic scribbling
  beats loop. Art-side fix stays §JIT-MARCUS (prompt updated: "never
  smiling in any frame").
- [x] Marcus talk blinks one frame — root cause: the gap slicer gave two
  stray specks their own cells (52×39 and 49×45 vs ~150×365 real frames)
  and the size normalizer blew them up. New `dropMalformedFrames` filter in
  framesFromGrid drops RUNT frames (<40% of median content height) and
  DOUBLE-FIGURE frames (>1.9× median width) on every NPC sheet load — this
  also covers Bernard/Camille below.
- [ ] PP walk front is not a walk cycle (16 near-identical standing poses,
  PP "glides" when walking toward camera) — regen queued at §JIT-WALKFRONT.
- [x] Office Higgins loses color in his hand — his pale hand/skin
  highlights sit inside the connected key's tol-8 band, so the background
  flood bled into them. Office idle/talk/give-map now load at tol 4
  (measured +12k opaque px on idle; background still keys cleanly).
- [x] Biker — user: "brilliant". No action.
- [x] Poulain reposition — Y=319 per playtest (waist-cut foot 464).
- [x] Poulain work alt-idle never played while idle — ENGINE BUG: the
  alt-idle restore check (frame 0 + tiny timer) was also true on the tick
  right after the swap-in, so every alt cycle ended after ONE invisible
  tick (this also silently killed Marcus's freakout punctuation). Restore
  now waits for the cycle to actually advance past frame 0 and wrap.
- [x] Camille + Bernard idle sprites broken — interim fix shipped via the
  `dropMalformedFrames` filter (kills the 3px/7px sliver "blink" frames and
  the two-patrons-in-one-cell frames). Real fix stays the §JIT-PATRONS
  re-roll.

### Reported (2026-06-12 — PR batch: give ordering + give flower + Marcus strange idle + office spawn)

- [x] PR#1 give-item order (EVERY give in the game) — FIXED: the alt-dialog
  contract now returns an optional `handOff` third value; every dispatcher
  (held-item click, walk-up startNPCDialog, Pierre's onClickOverride) plays
  the hand-over BEFORE the text: PP's `give_<item>` one-shot → the NPC's
  receive one-shot (`npcAnim` override / `receive_<key>` / `receive_item`,
  skipped if absent) → talk anim + dialog. The end-callbacks keep only the
  hand-BACK anims (NPC give + PP receive) and state flips. New plumbing:
  `npc.playOneShotAnimThen` (one-shot with completion callback),
  `player.playHandOff`, `giveAnimKeyForItem`. Wired for: Lily's flower, all
  five kid anchor heals, Pierre baguette/confiture/heel, Poulain rolling
  pin + signed postcard, Henri coffee, Camille pencil, Beaumont sketch.
- [x] PR#2 give flower not smooth — TWO causes. Engine side FIXED: the
  player one-shot loader used the GLOBAL white color key, which erased the
  white parts of the handed items (daisy petals, coffee cup, postcard,
  sketch page) - the give/receive/grab-flower one-shots now load with the
  edge-connected key (recovers ~4.6k petal px on the give sheet, ~7.7k on
  grab). Art side QUEUED: frames 1-2 and 4-6 of the sheet are
  near-duplicates (a 4-pose hand-over) - re-roll at §JIT-GIVEFLOWER.
- [x] PR#3 Marcus strange idle swapping frames — root cause: it's the ONE
  Marcus sheet that fails gap detection (figures touch; content crosses the
  proportional borders by up to 335px), so the fallback slices two
  Marcuses into each frame. INTERIM FIX: strange idle plays the clean
  strange_alt sheet (same freakout-sketching vibe). Real fix stays the
  §JIT-MARCUS re-roll.
- [x] PR#4 PP office spawn — moved down (spawnY 480 → 505, feet ~775) so he
  stands on the open floor instead of floating at the door/radiator line.

### Reported (2026-06-12 — follow-up: bakery sizes + PP vs Paris adults)

- [x] Bakery patrons huge / only heads visible — FIXED (root cause of #16's
  bad tune): measured the sheets — every cafe-patron sheet is a WAIST-UP
  BUST (~190×300px opaque content, no legs, no chair), not full-body art.
  Filling a 240px full-body box and then cropping to the top 55% rendered
  a ~2×-life-size head and nothing else. Fix (npc.go): `srcCropBottomFrac`
  dropped from all 5 patrons; the whole bust renders at 110×135 with the
  waist cut (bounds.Y+H) anchored at each table's cloth-top edge (left
  ~y490, middle ~y505, right ~y500). Heads now read at standing-NPC scale.
- [x] Madame Poulain oversized — FIXED: her sheets are the same bust
  framing (192×227 content); a 180px bust implied a ~330px standing person
  vs PP's rendered 211px. Bounds 170×180 → 170×145, waist-cut foot kept at
  y=508.
- [x] PP smaller than the Paris adults — FIXED: PP renders at
  playerDstH 270 × fillFrac 0.78 ≈ 211px, but Colette/Nicolas/Claude filled
  H=235 and Beaumont 240. All four → H=205 (feet kept; Y shifted down by
  the delta), so PP now reads a touch taller than the humans, per the
  classic look. Pierre (156, mid-distance) untouched.
- [ ] Bernard + Camille idle sheets: two figures TOUCH → the gap slicer
  yields a 3px sliver frame (blink) + a double-figure frame in each.
  Re-roll queued at EXTRA_PROMPTS §JIT-PATRONS. Their talking sheets and
  the other patrons split clean.

### Reported (2026-06-12 — playtest batch #3, 17 items; plan velvety-exploring-pnueli)

**Regression fixes:**

- [x] #1/#3 Higgins/Jake talk ≠ idle size — FIXED: reverted the idle-grid
  size-reference change (last PR's #30 fix); `drawScaled` normalizes each
  state by its OWN sheet's `maxOpaqueH` again (idle/talk sheets ship at
  different resolutions). The original #30 was a crop-path bug, fixed for
  real below.
- [x] #4 Marcus talk vanishing frames — FIXED: `sheet_clean` had run on the
  RESTORED old talk sheet (art strays outside fixed cells → real body parts
  erased). PNG restored from HEAD; all four Marcus sheets REMOVED from the
  cleaner allowlist. Real fix stays §JIT-MARCUS (art).
- [x] #16 bakery NPC sizes ("you must fix this!") — FIXED: the legacy seated
  crop aspect-fit the FULL 192×1024 cell into 80×140 bounds (≈26px-wide
  patrons). Seated crop now runs through the opaque-box pipeline: scale from
  content height, crop the CONTENT box to its top 55%, anchor at bounds.Y.
  Patron bounds → full-body 110×240 (visible ≈132px) — tune in playtest.
  **SUPERSEDED same day (playtest: "huge, only heads")** — see the
  2026-06-12 follow-up batch below.
- [x] #7 sleeping PP bigger than idle — FIXED: `sleepStandH` now multiplies
  by `playerRenderFillFrac` (0.78), matching the rendered idle height ~211.
- [x] #5 Lily pre-trigger + missing give anim — FIXED: the hover probe CALLS
  `altDialogFunc` every frame, and Lily's receive one-shot sat in the func
  BODY → fired on hover. Moved into the returned callback (give anim chain
  `playGive("flower")` → `receive_flower` now runs on the real handoff).
  Camille's hover-firing debug print removed. New SKILL.md §8a rule:
  altDialogFunc must be PURE.

**Diagnosis tooling:**

- [x] #2 walk line (3rd report) — built the F3 walk-debug overlay
  (`game/walk_debug.go`): yellow walkSegments with endpoint coords, green
  PP-foot crosshair, red last-snap marker. Walk the path with F3 on,
  screenshot, and we set the segment coords once with real numbers.
- [x] #6 Higgins shout (again) — ROOT CAUSE FOUND: the bellow the player
  actually WATCHES is the pre-transition beat at camp_grounds
  (`checkDay1Complete`) — delivered by the GROUNDS Higgins, who had no shout
  frames, in plain talk. The camp_night shout wiring was always fine (the
  scene + sheet verified: 8×2 gap-detects, 16 frames). FIXED: grounds
  Higgins registers the shout sheet; the beat swaps idle+talk to shout for
  the dialog and restores after. Duplicated "It's getting very late / NOW"
  lines removed from night_bedtime.json (kept the campfire-only lines).

**Camp / office:**

- [x] #9 standing over Marcus's desk — FIXED: spawn 700→500 (spawn placement
  bypasses blockers — that was the actual overlap) + desk blocker extended
  left/down {850,460,380,240} → {700,430,530,270}.
- [x] #10 office exit walked off-screen — FIXED: camp_grounds → camp_office
  now chains waypoints (mid-path 900,483 → down-right trail 1050,640 →
  transition) instead of `walkToExit(downRight)`.
- [x] #11 office Higgins "a tiny up" — Y 300→290. (Speed untouched.)

**Paris street:**

- [x] #12 biker v2 — FIXED: halo keyed at tol 24 (new engine
  `SpriteGridFromPNGCleanConnectedTol`); faster vx 120→190; lane 735→755;
  encounter REDESIGNED: click → PP walks into the lane ahead, biker keeps
  riding, brakes when he reaches PP (±50px), PP flinches (stateReacting
  holds through the dialog), apology plays, rides on at close.
- [x] #13 Colette talk spot still on the bike rack — FIXED: `approachRight`
  (PP stands street-side). The previous X-nudge alone wasn't enough.
- [x] #14/#15 floor-item stand points — FIXED: new `walkToFloorItem` helper
  (HandleClick pickup path): PP stands BESIDE the item, feet on its base
  line, left by default / right via new `floorItem.standRight`; faces FRONT
  on arrival so blocked beats (pigeon pot) play to camera. Rolling pin sets
  standRight so the reach hand lands in the basket. SKILL.md §8c addendum.

**Art (user generation pass — prompts already live):**

- [ ] #8 Marcus strange idle swiping — §JIT-MARCUS (straddled cells). ART.
- [ ] #17 remaining broken last-PR sprites — the live §JIT batch: Marcus set,
  §JIT-POULAIN idle/talk, §JIT-PP2 walk back, §JIT-JAKE strange talk,
  §JIT-LILY, §JIT-FLOWER, §OD office talk, §NIC1-v2b Nicolas talk. After the
  #1/#3 engine fix, re-judge in-game which still look broken — the size
  mismatches were engine, not art.

### Reported (2026-06-11 — NEW PR playtest, 39 items; plan velvety-exploring-pnueli)

**Systemic fixes:**

- [x] `[P0]` NPC one-shots never terminated when played standalone → NPCs
  froze on the last give frame forever (#26 Poulain, #29 Henri, Camille
  mid-sketch). FIXED: `npc.update` auto-ends at `oneShotDuration` elapsed.
- [x] `[P0]` Held-item drops bypassed `onClickOverride` (#33: no shrink when
  giving to Pierre) → routed through the override; `giveItemTo` now clears a
  matching held item (no ghost cursor item).
- [x] `[P0]` #38 pencil→Camille dead end: gate loosened to
  `sketchAsked || camilleAsked`, plus never-silent branches (holding pencil
  pre-quest = polite line; pencil in bag = "hand it here!" nudge) + a
  `[camille]` debug print for the playtest.
- [x] `[P1]` #28 "Caf??": the bitmap font is ASCII-only — renamed
  `Cafe au Lait`/`Elise`, ASCII-swept ALL strings (accents + em-dashes →
  ASCII) across game/*.go and assets/data JSONs.
- [x] `[P1]` #12/#27/#35 talk speed vs text: office Higgins 0.20, Henri 0.20,
  Beaumont 0.22.
- [x] `[P1]` #30 NPC one-shot size: drawScaled's size reference now comes
  from the IDLE grid, so talk + one-shots render at idle's size.
- [x] `[P1]` #37 inventory icons centered by CONTENT box (engine
  ContentBoxKeyed; bag + held cursor draw by it).
- [x] `[P1]` #4 cursors: floor items show the GRAB (open hand); open bag +
  PP click show the pink POINT.

**Camp / office / flight:**

- [x] #1 walk line rerouted through (755,470) — TUNE IN PLAYTEST.
- [x] #2 Jake approached from his LEFT (new `approachLeft`); #3 Lily Y 440→425.
- [x] #6/#11 pp_sleeping/pp_waking + office Higgins: Aggressive(32)→Connected
  key (face/eye colors survive); sleeping height matched to idle (270).
- [x] #10/#13 office Higgins re-mirrored (flipped=false) with the give-map
  throw kept at its old orientation (new per-one-shot `oneShotFlip`); Y→300.
- [x] #14 airplane drawn by per-frame OPAQUE BOX — no more two-row jumping.
- [ ] #5 Higgins shout: wiring verified correct; the SHEET is one pose ×14 —
  re-roll with the existing storyboarded prompt. ART-ONLY.
- [ ] #7 Marcus room strange idle: known straddled sheet (§JIT-MARCUS). ART.
- [x] #8/#24/#31 last-frame blink: trailing BLANK frames now trimmed at load
  (`trimBlankTail`). If Bernard idle still looks broken it's art (§BRN1 then).

**Paris:**

- [x] #16 biker (user choice: INTERACTIVE): keyed load (no more white box),
  rides the MAIN street (foot 735, scale 0.85), clickable → brakes +
  "Pardon! Sorry monsieur, but you are blocking ze way!" → rides on.
  §BK1 art landed 2026-06-11 with real alpha.
- [x] #17 Colette X 300→335 (PP off the bike rack — tune).
- [x] #19 Pierre walk-away: two-stage exit (down to full size, THEN to the
  clicked point).
- [x] #37 pigeon: now an ambient perch critter at the easel (returns every
  morning) — playing it as Pierre's one-shot literally turned HIM into a
  pigeon.
- [x] #18/#33 generic pickup lines (`genericPickupDialog`, SKILL §8c) + PP
  `give_item` one-shot wired at every hand-over (SKILL §8b, §PG1 art landed
  2026-06-11).
- [x] #34 museum: single walk-in via entryWalkPending (no double-spawn);
  Beaumont moved onto PP's line (Y 359→500 — tune W/H if he reads big).
- [x] #22/#23/#25/#32 bakery: patrons 80×140 + seated crop (0.55) restored +
  PP stands in the counter/tables aisle (band minY 400, approachY 405);
  Poulain → (733,328).
- [~] #20 Nicolas regen: split IDLE sheet landed + verified 2026-06-12
  (gap-detected, camera routine, mouth closed). TALK sheet still pending
  (§NIC1-v2b); the loader reuses the old combined sheet's talk row meanwhile.
- [ ] #26 Poulain outfit mismatch → §JIT-POULAIN amended (match WORK sheet). ART.
- [x] #36 flower pot re-roll → §PA2-v2. ART landed 2026-06-11.
- [x] `[P2]` Biker §BK1, PP give §PG1 — sheets landed 2026-06-11.

**Out of scope:** #39 Jerusalem (next batch).

### Reported (2026-06-10 — automated sprite jitter audit, camp + Paris)

Ran `go run ./tools/jitter_audit` (new tool: measures per-frame content
bounding boxes inside each sheet cell; FOOT drift = character bottom moves
between frames → vertical jumping, CENTER-X drift = horizontal sliding).

**Code fixes landed this pass:**

- [x] `[P1]` Poulain "give baguette" one-shot loaded 0 frames — `npc.go`
  still pointed at `outside/npc_madame_poulain_give.png` after the sheet
  moved to `coffee/`. FIXED: path updated; the hand-over animation plays
  again on the rolling-pin trade.
- [x] `[P1]` Night Higgins idle loaded as 7×1 but the sheet is 6 frames
  (2304 = 6×384; 2304/7 isn't whole) → sliced mid-character, horizontal
  sliding. FIXED: `newNightHiggins` loads 6×1 (matches entrance Higgins).
- [x] `[P2]` `assets/data/npc/{kids,higgins,paris}.json` grids/paths were
  stale (pre-regen art): Tommy talk 4×2, Jake 5×2, Marcus idle 7×2,
  Higgins idle 7×1, curator old path + 8×2/5×2, Colette old french_guide
  sheets. SYNCED to what the engine actually loads, so a future
  JSON-driven loader won't regress.
- [x] `[P1]` Curator Beaumont minted a DUPLICATE postcard (and replayed the
  "head back to camp" monologue) on every repeat conversation — onDialogEnd
  fires after every chat and had no guard. FIXED: one-shot `gaveCard` guard.

**Paris quests wired this pass (see STORY.md):**

- [x] `[P1]` "Camille and the Sold-Out Postcard" — MAIN-CHAIN gate (user
  rework same day: the first draft put the pencil under a sleeping Lucien
  → too dark; user wanted outside/inside/museum back-and-forth on the way
  to the postcard). New flow: Beaumont's postcards SOLD OUT → asks for
  Camille's replica sketch → Camille lost her pencil at sunrise → Nicolas
  saw it roll into the flower pot by the Louvre steps (hidden floor item,
  generic grab anim) → pencil to Camille → sketch one-shot → sketch to
  Beaumont → Postcard + paris_done flags. Lucien reverted to awake flavor;
  Yvette foreshadows the sell-out. User additions same day: (a) the
  pigeons BLOCK the flower pot until Pierre repays his baguette+confiture
  debt — favor beat where he whistles them off (seeds the Pigeon Critic
  gag); (b) Camille plays her sketching one-shot at the end of her first
  regular chat (npc_camille_sketching.png already on disk).
- [x] `[P1]` "The Pigeon Critic" (optional) — post-press-pass Pierre asks
  for crumbs → Poulain donates the Baguette Heel → pigeon lands (dialog
  beat) → "Mini Portrait" keepsake + plein-air/Monet fact.
- [x] `[P1]` Grandson souvenir loop CLOSED — Poulain asks (existing beat) →
  Beaumont signs a second postcard from the new print run (altDialog) →
  hand-in at the bakery → "Le Panthère Rose" éclair reward.
- [x] `[P1]` Poulain "counter service": after the rolling-pin trade she
  refills the Café au Lait while Henri's trade is pending, hands out the
  heel, and accepts the Signed Postcard — the chain can't soft-lock.
- [ ] `[P2]` Art for the new items: 4 icons queued at EXTRA_PROMPTS §PI1
  (charcoal_pencil, camille_sketch, baguette_heel, mini_portrait — items
  work now but show blank icons until the PNGs land). Optional pigeon
  one-shot at §PA1. Signed Postcard reuses postcard.png (no art needed).
- [ ] `[P2]` Pencil flower-pot pickup coords (1085, 615, 70, 50) on
  paris_street are a starting guess — tune against the BG in a playtest.
- [ ] `[P1]` The paris_street BG has NO flower pot near the Louvre exit —
  the pencil spot is invisible (cursor-only). Two-state prop prompt landed
  at EXTRA_PROMPTS §PA2 (pigeon perched → pencil revealed); wire the prop
  swap after the PNGs land.
- [x] `[P1]` NEW STANDING RULE (user 2026-06-10, documented SKILL.md §8b +
  memory): every collected item must be visibly acquired — PP plays a
  pickup/receive one-shot AND the giving NPC plays a give one-shot. Applied
  to all six new Paris beats: heel (Poulain give + PP get_baguette), coffee
  refill (Poulain give + generic grab), Camille sketch (her sketch one-shot
  + generic grab), mini portrait / postcard trade / signed postcard (generic
  grab; postcard monologue now plays AFTER the grab completes).
- [x] `[P1]` §PR2 (Pierre give) + §PR3 (Beaumont give) GENERATED + WIRED
  2026-06-10; pigeon one-shot sequenced before Pierre's give (back-to-back
  playOneShotAnim calls cancel each other). Flower-pot prop wired with the
  pigeon→pencil texture swap on Pierre's favor.
- [x] `[P1]` §PR1 `PP receive.png` GENERATED + WIRED 2026-06-10: registered
  as player one-shot `receive_item` (8×1); all five call-sites in game.go
  swapped from generic grab (Pierre portrait, Beaumont postcard ×2, Poulain
  coffee refill, Camille sketch).
- [ ] `[P1]` Jerusalem give one-shots queued at EXTRA_PROMPTS §JG1 (Gary,
  Eli, Dov, Miriam) — required by SKILL.md §8b before the daisy-chain is
  wired.

**Art regens still needed — ranked by measured drift (worst first).**
All are art-only; the fix is the same for every sheet: keep the
character's feet and horizontal center locked in the SAME pixel position
in every cell:

Regen batch #1 generated + wired 2026-06-10 (user). Audit re-run results:

**FIXED (audit-clean, prompts retired to the Done log):** PP talk front,
PP talk side, PP grab, PP receive map, Marcus idle, Higgins shout, Higgins
give-map (seated 6×2), Poulain work, Colette talk grid (clean 8×2, loader
switched).

**Ghost-limb sweep (user 2026-06-10 — floating hand visible in-game on PP
talk front):** the audit gained two detectors — GHOST PIECES (detached limb
painted inside a cell) and CONTENT CROSSES (figure straddling a cell border).
New `tools/sheet_clean` erases ghosts on prop-free sheets (keeps only the
largest connected piece per cell; NEVER run it on sheets with legit separate
objects — thrown map, handed items, pigeon). Cleaned + visually verified:
PP talk front (the reported bug — FIXED), PP idle side, PP celebrate,
Colette talk, Higgins office idle + talk. Marcus talk/strange sheets turned
out to have figures STRADDLING borders (regen unusable) → reverted to the
pre-regen art from HEAD; re-roll queued.

**ENGINE STABILIZATION (user 2026-06-10: "read the sprite properly + the
object stays in the same spot" — no padding/splitting hacks):**

- [x] `[P0]` Proportional cell slicing (`engine.gridCellRect`): cell boundary
  i sits at floor(i*W/cols), so sheets whose dims don't divide by the grid
  (1535-wide talk front, 1672/1685-wide Paris sheets) load correctly with the
  remainder distributed — no more truncated strip on every frame. Applied to
  all six grid loaders + eraseGridLines; jitter_audit mirrors it (and no
  longer flags non-divisible dims).
- [x] `[P0]` Median anchor stabilization: the renderer pins every character
  to ONE spot regardless of art drift inside the cells.
  - Player: `stabilizeFootCX` also computes the sheet-median foot ROW;
    drawScaled anchors Y by it (a dipped tail extends past the line instead
    of lifting the body). X was already median-anchored.
  - NPCs: new `stableAnchors` (median box-center column + bottom row per
    animation); drawScaled anchors X/Y by the medians instead of the frame's
    own cell position/bottom, with proper flip mirroring.
  - CONSEQUENCE: the audit's FOOT/CENTER-X drift numbers now measure ART
    quality only — the renderer cancels them on screen. Re-rolls below are
    still worthwhile (cleaner limb framing) but no longer urgent.
- [x] `[P1]` sheet_clean v2: only erases pieces OUTSIDE the body's bbox
  (v1 erased interior belly details → see-through holes, user-reported).
  New tools/sheet_repair refills enclosed pure-white holes on global-key
  player sheets with the surrounding color. Re-cleaned + visually verified:
  PP talk front (ghost hand gone, colors intact), the restored Marcus sheets
  (neighbor-spill erased — the old art now reads clean), Colette talk.

**Still drifting / ghosted — live re-roll prompts in EXTRA_PROMPTS §JIT
(now ART-QUALITY only; the renderer cancels positional drift):**

> 2026-06-12 prune: the drift-only re-rolls below were RETIRED from
> EXTRA_PROMPTS (user: file too long; the feet-anchoring renderer cancels
> drift on screen, so nothing is visibly wrong in-game). Their lines are
> checked off here; re-open from git history if one ever shows on screen.

- [ ] `[P1]` Marcus talk + strange idle/talk/alt: regen #1 straddled cell
  borders; reverted to old art (old drift numbers back) (§JIT-MARCUS).
- [x] ~~Poulain idle/talk~~ — RETIRED 2026-06-12 (renders correctly at bust
  scale; re-open only for the #26 outfit mismatch).
- [ ] `[P1]` PP walk back: regen #1 made it WORSE (60→97px FOOT) (§JIT-PP2).
- [x] ~~Jake strange talk / Lily idle+talk / PP grab flower / PP idle front /
  Colette talk~~ — RETIRED 2026-06-12 (drift-only; renderer compensates).
- [x] ~~Higgins office talk (§OD)~~ — RETIRED 2026-06-12: accepted in every
  playtest since 06-05; got the tol-4 color-key fix instead.
- [ ] `[P2]` Ghosted but NOT auto-cleanable (legit props in frame — re-roll
  or hand-edit): PP get baguette / get jam / grab rolling pin / receive map
  (incoming map), Lily receive flower, Poulain give + bring baguette, pigeon
  lands, curator idle/talk, café patrons (Bernard/Camille/Henri/Lucien
  talking), Higgins give-map (thrown map = legit piece).

Measured CLEAN (no regen needed): PP walk side/front, PP idle back,
Higgins entrance idle + talk + walk back, Tommy idle, Jake idle,
Danny idle/talk, Claude both rows, Curator both 8×1 strips.

### Reported (2026-06-05 — playtest pass, museum + Paris flow, 32 items)

Engine/JSON fixes landed this pass; art-bound items are queued in
`EXTRA_PROMPTS.md` under "Playtest pass — museum".

- [~] `1.` PP walk-side not a full sprite + idle-front jitter. **Art** — regen
  EXTRA_PROMPTS §AA (idle front, eyes open + feet locked). §AB walk-side DONE
  2026-06-10 (8×1 regular cycle, anchor-clean).
- [~] `2.` PP talk-front not cut normally. **Art** — §AC (match idle dims 1536×1024).
- [x] `3.` PP walks into camp at (755,533) on arrival.
- [~] `4.` New Lily idle/talk/get-flower/post-flower. **Art** — §LL (keep design).
- [~] `5.` Marcus talk not smooth. **Art** — §MM (unify idle+talk to clean 8×2).
- [x] `6.` Active pointing cursor while carrying an item + on the relevant travel pin.
- [~] `7.` Higgins shout sprite. **Wiring verified** (camp_night→night_higgins); the
  sheet's right half is blank → **art** regen §SH.
- [x] `8.` Marcus room frames swiping — engine: room idle/talk already load at the
  correct grid; remaining slide is the strange/talk cell-count mismatch → §MM art.
- [~] `9.` PP sleeping/waking + first idle frame open eyes. **Art** — §AD + §AA.
  (Wake dialog already plays from night_bedtime.json.)
- [x] `10.` Room Marcus shrunk (150×205) so he reads shorter than PP.
- [x] `11.` Higgins-office arrow moved to ~(1186,692).
- [~] `12.` Higgins office frames swiping → **art** §OF; **flip done** (office Higgins
  now faces PP); throw-map uses the give_map sheet (loads 6×2 correctly).
- [~] `13.` Airplane not cut well. **Art** — §AP.
- [x] `14.` Rolling pin hidden in bike basket (~539,644): cursor reveals it, grab
  one-shot plays on pickup.
- [~] `15.` Colette talk not smooth / right-side gap / last frame blank. **Art** — §CO.
- [x] `16.` Pierre: talks to the side; eases back to size after dialog (no pop).
- [x] `17.` Multi-item inventory left/right chevron logos removed.
- [x] `18.` Inventory bag oval enlarged (816×680).
- [x] `19.` Claude talk cadence slowed (0.10 → 0.16).
- [~] `20.` Poulain repositioned (~605,318); talk frame spacing → **art** §PO.
- [~] `21.` Lucien/Henri/Yvette bg + talk spacing. **Art** — §CF.
- [x] `22.` Inventory only opens when clicking on PP (tightened hit test).
- [x] `23.` Camille nudged right (legs tuck behind table).
- [~] `24.` Bernard talking tiny. Root cause: talk sheet framing ≠ idle → **art** §CF.
- [x] `25.` Item trades require handing the item over (held), not bag-only. Rule
  documented in `SKILL.md` §8a.
- [x] `26.` Bakery exit: PP walks to the door (~1261,422) then walks back through it.
- [x] `27./28.` Museum first arrival: one-time arrival monologue.
- [x] `29.` Beaumont flipped + repositioned (~546,599); **new talk sprite** → §BE art.
- [x] `30.` PP walks in from the left tunnel (381,481); scene scale 0.7 shrinks both.
- [x] `31.` Removed the bottom-right travel-map button in the museum.
- [x] `32.` Travel-map "fly back to Camp" pin unlocks after the postcard
  (relevantWhen `paris_done==1 && marcus_healed==0`).

### Reported (2026-05-24 — playtest pass 6, 7 items)

- [~] `[P1]` 1. PP talk-front cut wrong + not matching idle + two rows of frames visible. **Art-only** — EXTRA_PROMPTS §A (talk-front to 1376×768, 8×2 matching idle) and §B (talk-side to 1672×941, 8×2 matching idle-side) still queued. Current PNGs are wrong dims, hence the two-row strip showing.
- [x] `[P1]` 2. PP still not walking over the painted path. ROOT CAUSE FOUND: `camp_grounds.json` had `minY: 0, maxY: 0` → engine defaulted to `playerMinY=265, playerMaxY=395`. WalkSegments were authored at y=483-653, but `setTarget` clamped Y back to 395 → PP could never reach the segments and ended up ~90px above the painted path. FIXED: scene JSON now sets `minY: 300, maxY: 560`. Segments are now reachable; foot lands ON the dirt path. spawnY raised 455 → 483 to match the front-of-camp segment. Verified with `go build` clean + click hit-box test (27 NPCs, 135 probes) still passes.
- [x] `[P1]` 3. Need to click RIGHT of every NPC to start dialog. FIXED: `npc.containsPoint` now uses the live `lastDrawRect` (where the sprite actually rendered last frame) with ±25 X / ±15 Y padding, clamped to the authored `bounds`. Click-test now matches visible body for any aspect-preserve fit, regardless of how narrow the rendered sprite is inside the wider bounds rect.
- [x] `[P1]` 4. Jake walk_back has head cut. FIXED: `npc_jake_walk_back.png` (941 tall, kid content y=231-660) was being loaded as 8×3 row=1 → cell y=313-627 chopped 82px off Jake's head. Loader reverted to 8×1 (full canvas cell) — narrower render but no head crop.
- [~] `[P1]` 5. Marcus talking and idle not in same size. **Art-only** — EXTRA_PROMPTS §D regen still pending (talk 7×2 197×384 to match idle dims). Engine has no fix until the PNG lands.
- [x] `[P1]` 6. Can't move item to Lily after picking from inventory. FIXED: `HandleClick` PP-toggle branch was firing even when PP had a held item on the cursor, swallowing the click before it could route to Lily's drop handler. Condition tightened to `if g.player.containsPoint(x, y) && g.inv.heldItem == nil` — clicks fall through to NPC drop logic when an item is on the cursor.
- [x] `[P1]` 7. Game gets stuck after item #6. FIXED as side-effect of #6 — the swallowed click was leaving the inventory + held-item state in an inconsistent half-modal, which blocked further input until restart. With the PP-toggle fix, the held-item drop flow completes cleanly and control returns to the world.

### Reported (2026-05-23 — playtest pass 5, 12 items)

- [x] `[P1]` 1. PP click goes to location behind. FIXED: `HandleClick` no longer "checks if anything is under PP and falls through" — clicking on PP ALWAYS toggles inventory (or eats the click if inventory is empty). The fall-through was confusing the user; per their feedback, PP click is sacred and should never route to a hotspot behind him.
- [~] `[P1]` 2. PP idle and talk not same size. **Art-only** — EXTRA_PROMPTS §A (talk-front match idle-front dims 1376×768) and §B (talk-side match idle-side dims 1672×941) already targeted at PRIORITY regen. The on-disk PNGs are still the wrong dims; once regenerated to match idle exactly, talk + idle will render at identical visual size.
- [x] `[P1]` 3. Jake walks to the sky. FIXED: `jake_exit.json` target Y 200 → **340** (cabin door threshold). Y=200 was above the cabin into the sky; Y=340 lands his foot at the visible door area. Move duration 2.0 → 1.6s for snappier exit.
- [~] `[P1]` 4. Marcus talking two-frames + not same size as idle. **Art-only** — EXTRA_PROMPTS §D (Marcus talk regen to 7×2 matching idle dims) still pending. Engine has no fix until regen lands.
- [x] `[P1]` 5. Danny click misses. FIXED: bounds widened to **(1090, 405, 180, 175)** — covers 1090-1270 (vs the visible Danny at ~1135-1224). Generous margins on both sides ensure clicks register no matter which animation frame is showing. Also reverted Tommy/Jake/Lily/Marcus to original 145-wide rects (the 100-wide tighten introduced misses for them too).
- [x] `[P1]` 6. Lily flower dialog without giving the flower. FIXED: `kid.altDialogRequiresHeld = true` on Lily (was false). Now the flower beat only fires when PP has actively PULLED the flower from inventory (it's on the cursor) and clicked Lily. Just having flower in the bag is no longer enough. Strict-missing hint updated to nudge the user toward the drag-onto-Lily motion.
- [x] `[P1]` 7. Lily/Danny/Higgins outside on Day 2. FIXED: `startDay2()` now sets `hidden=true` on ALL camp_grounds NPCs (Marcus, Tommy, Jake, Lily, Danny, grounds-Higgins). They're all "in their cabins / at the office" per the Day-2 story. Tommy + Jake were already hidden via their Day-1 exit sequences; Lily/Danny/Higgins were leaking through because my previous pass only hid Marcus.
- [~] `[P1]` 8. Higgins office: position + frame smoothness + facing + throw blink. PARTIAL: bounds Y 310 → 300 (head a bit higher above desk); `talkFrameSpeed` 0.10 → 0.08 (smoother cadence); **`flipped: true`** so Higgins faces LEFT toward PP entering from spawn (was looking away). The "two lines" jitter is from the office idle PNG (1748×900) and office talk PNG (1694×928) having different cell dimensions — that's art-only, regen prompt landed at EXTRA_PROMPTS §W (both sheets to 1376×768, 6×2, 229×384 matching cells).
- [~] `[P1]` 10. Colette vanilla shirt still chroma-keyed. PARTIAL: standing-rule updated in EXTRA_PROMPTS — `#F2EFE5` (vanilla) is NOT safe for large fabric areas; use `#E5DDC8` (cream) instead. New Colette regen prompt added at EXTRA_PROMPTS §G with the corrected color. Engine has no fix until the regenerated PNG lands.
- [x] `[P1]` 11. PP standing on air in Paris. FIXED: `paris_street.json` spawnY 470 → **510**, walkSegments Y also 510. Foot now lands at y=780 (well into the cobblestone band).
- [x] `[P1]` 12. Pierre stuck post-dialog. FIXED: Pierre's `onClickOverride` now calls `clearRecede()` after dialog ends — the recede tween was leaving `recedeActive=true` so PP was frozen at the receded position (couldn't move). `clearRecede()` releases the freeze; `depthScale` (Y-based) keeps PP visually small at Pierre's depth until the next click moves him forward in the scene.

### Reported (2026-05-22 — playtest pass 4, 18 items)

- [x] `[P1]` 1. Danny click hit-region offset. FIXED: tightened bounds W 145→100 + shifted X right 22px so click rect hugs the visible ~89px sprite (was 56px of empty bounds to the right where user-clicks misfired). Applied to all 5 camp kids since they share idle cell dimensions.
- [x] `[P1]` 2. Inventory should freeze the world. FIXED: added `g.inv.open` gate at the top of `Update()` that early-returns BEFORE scene NPC update, floor-item ticks, and trigger checks fire. Dim overlay alpha 140→190 for better visual freeze. Clicks behind already blocked at HandleClick:1235.
- [x] `[P1]` 3. Cursor on pickup. FIXED: reordered hover loop in `ui.go` so floor items check BEFORE npcs (grab cursor wins). Also reordered click handler in `game.go:HandleClick` so floor-item click routes BEFORE npc-click.
- [~] `[P1]` 4. Higgins shout sprite. ENGINE FIXES APPLIED: switched shout PNG loader from `loadNPCGridClean` (tight tolerance) → `loadNPCGrid` (lenient), added debug `fmt.Printf` showing `len(shoutFrames)` at load. Also added `[SeqNPCAnim]` log if NPC lookup or frame-map lookup fails. Run game and check console output to diagnose if still failing.
- [x] `[P1]` 5. Wake-up dialog plays twice. FIXED: removed the duplicate `day2Monologue` trigger at `game/game.go:1527-1530`. The night_bedtime.json sequence dialog at the wake-up is now the single canonical source.
- [x] `[P1]` 6a. Marcus alt sprite. FIXED: added `altIdleGrid` / `idleAccumSec` / `altIdleAfterSec` fields on `npc` struct + inactivity swap logic in `npc.update`. Wired on `newRoomMarcus` with `npc_marcus_strange_alt.png` and threshold 5.0 s.
- [~] `[P1]` 6b. Marcus strange-idle BG fringes. PROMPT UPDATED: EXTRA_PROMPTS §D2 now spells out "Background must be PURE WHITE `#FFFFFF` — no off-white, no light-grey shading at edges; near-white survives the chroma-key as halos around the silhouette." Art regen required.
- [x] `[P1]` 7. PP walks side for right-down arrow. FIXED: `game/player.go:752` threshold lowered from 1.2 → 0.8. Now any `|dy| >= 0.8*|dx|` picks vertical motion — down-right clicks show the down-walk sprite (front-facing).
- [x] `[P1]` 8a. Higgins office position. FIXED: bounds Y nudged 290 → 310 so head sits slightly lower / more natural above desk.
- [x] `[P1]` 8b. Map throw blinking + arc. FIXED: added parabolic Y arc to `SeqTweenItem` draw (`arcHeight=200` lifts the projectile at midpoint t=0.5). Tween duration extended 0.8s → 1.5s so the arc plays slowly. The "blink" should be gone because the projectile now flies for longer than the brief gap between sequence steps.
- [x] `[P2]` 9. White-color rule with vanilla/cream. DOCUMENTED: added concrete color recipe block (vanilla `#F2EFE5`, cream `#E5DDC8`, pale grey `#C4C4C4`, bone `#EDE5D3`, silver `#C0C0C0`) at the top of `docs/EXTRA_PROMPTS.md` AND in `memory/feedback_no_white_in_prompts.md`.
- [~] `[P1]` 10. Madame Colette white shirt + frames not smooth. ENGINE FIX: `talkFrameSpeed` 0.10 → 0.08 for smoother cadence. White shirt → art regen needed (use vanilla `#F2EFE5` per the new rule); update queued in EXTRA_PROMPTS.
- [x] `[P1]` 11. PP on air in Paris. FIXED: `paris_street.json` spawnY 450 → 470 + added walkSegments along the cobblestone band (y=470 horizontally) so PP foot lands on the floor consistently.
- [x] `[P1]` 12. Pierre walk-back-and-shrink choreography. FIXED: added `onClickOverride` field on `npc` struct. Wired on Pierre in `setupParisCallbacks`: PP walks to middle of road (700, 700) → `playRecede(1.0s, scale=0.65, dyUp=50)` → dialog plays (uses the existing 2-stage altDialogFunc). After dialog, PP stays at the receded scale until next walk.
- [x] `[P1]` 13. Travel map back-button huge. FIXED: special-cased return-scene pin in `pinHitRect` — 40×40 hit rect (vs the standard 60×70 for travel pins).
- [x] `[P1]` 14. Paris NPC sizes not unified. FIXED: standardized Colette / Nicolas / Claude bounds to **120×235** (was 135×230, 86×230, 115×240 respectively). Pierre stays at 95×175 (back-of-line by design).
- [x] `[P1]` 15a. Café patrons sitting on chairs (position). FIXED: repositioned to match the BG chair layout — Yvette (80,555), Bernard (240,555), Camille (420,555), Henri (580,555), Lucien (920,555). Élise removed from `paris_bakery.json` (no 6th chair in the BG). Y unified at 555.
- [x] `[P1]` 15b. Café patrons full-body vs upper-body. FIXED: added `srcCropBottomFrac` field on npc + drawScaled trim logic. Set to 0.55 on each café patron — engine renders only the top 55% of each cell so the chair art from the BG fills the lower half. Avoids double-drawn legs/feet. Art regen still queued in EXTRA_PROMPTS §7.x for chest-up-only sprites.
- [x] `[P1]` 16. Café patron chroma-key fringe. FIXED: switched `loadCafePatronGrids` from `loadNPCGridRow` (tolerance 8) → `loadNPCGridRowClean` (tolerance 16) so off-white edges around clothing/cups chroma-key cleanly.
- [x] `[P1]` 17. Madame Poulain behind the desk, upper body. FIXED: bounds (780,380,135,230) → **(717,250,135,160)**. Top at Y=250 (head clearly above counter), foot at Y=410 (sprite bottom anchors at counter surface). Visible 71×160 — classic behind-the-counter pose.
- [x] `[P1]` 18. NPC talking sprites missing. FIXED: `loadCafePatronGrids` now (a) uses cleaner chroma-key, (b) logs `len(idle)` and `len(talk)` per patron at load, (c) falls back to `talk = idle` if the talk row comes back empty so the NPC still renders during dialog.
- [x] `[P2]` 19. Plan-mode. Done — plan file approved.
- [x] `[P0]` Day-2 silencing removed. Side-fix: dropped the blanket `silent=true` on all `camp_grounds` NPCs at `startDay2()`. Only Marcus still gets `hidden=true` (cabin freakout). Danny + others stay clickable.

### Reported (2026-05-21 — playtest pass 3, 5 items)

- [~] `[P1]` 1. PP talk side cut wrong + not same size as idle. FIXED PROMPT: EXTRA_PROMPTS §B retargeted to **1672×941, 8×2, 209×470 cells** — matching `PP idle side.png` exactly. Also removed the stale "yellow gloves" reference. Current PNG (1536×1024) is wrong dims — once regenerated to the new spec, talk + idle will visually match. Engine wiring unchanged.
- [x] `[P1]` 2. PP walks LEFT of the path again. FIXED: audited `PP idle front.png` — PP's body silhouette sits 8 source pixels LEFT of the cell horizontal center in every frame (tail trails right, body shifts left). Engine was centering the CELL on the walk-segment point, so the visible body always ended up left of the path. Added `bodyOffsetPx := int32(8.0 * frameScale)` shift in `player.drawScaled` (game/player.go) to compensate; mirrors when sprite is flipped for left-facing motion. Works at any depth scale.
- [x] `[P1]` 3. Tommy walking VERY small vs idle. FIXED: the on-disk `npc_tommy_walk_left.png` (1536×1024) has kid content in only the MIDDLE band (rows 324–678, ~35% of canvas height). Loading it as `8×1` gave cells of 192×1024 where the kid took only ~35% of vertical space, so aspect-preserve rendered him at ~37 px wide. Changed `newTommy` to `loadNPCGridRow(..., 8, 3, 1)` — picks the middle row of an 8×3 grid (cells 192×341, kid fills the cell) → renders at ~112 px wide, much closer to idle size. Full art regen still tracked in EXTRA_PROMPTS §E.
- [x] `[P1]` 4. Jake walks toward the fire instead of his cabin. FIXED: `jake_exit.json` target Y changed from 441 (his cabin hotspot center, but Y=441 ≈ his current Y=405 → he moved DOWN/toward camera/fire) → **200** (well above the cabin footprint, so the back-view walk-anim reads as 'walking AWAY from camera INTO the cabin'). Also fixed the same "kid in middle band" art issue as Tommy — `newJake` walk_back load changed to `loadNPCGridRow(..., 8, 3, 1)`.
- [~] `[P1]` 5. Marcus talking not same size as idle + two frames. **Art-only** — Marcus talk PNG is 1672×941 (8×2 cells 209×470) vs idle 1376×768 (7×2 cells 197×384), so talk renders wider+taller than idle at the same bounds. EXTRA_PROMPTS §D already prompts a regen to 1372×768 (7×2, 196×384) matching idle. Engine has nothing to fix until the regen lands.

### Reported (2026-05-21 follow-ups)

- [x] `[P1]` Café patrons not animating. FIXED: the on-disk sheets are single 1376×768 PNGs in `paris/npc/coffee/` with row 0 = idle and row 1 = talk (8×2 grid, same as kids). `loadCafePatronGrids` was looking for `_idle.png`/`_talk.png` pair files that don't exist; updated to load the single PNG with `loadNPCGridRow(sheet, 8, 2, row)` for each row.
- [x] `[P1]` Remove "Camp Chilly Wa Wa Air" left exit from `camp_entrance` (the first scene). FIXED: removed the obsolete `arrowLeft` hotspot in `game/game.go:setupTravelHotspots` that was adding a "Camp Chilly Wa Wa Air" airstrip to `camp_entrance`. Travel is already opened by clicking the Travel Map item in the inventory (see `g.inv.onSelectItem`), so this scene-edge hotspot was a pre-2026-04-26 leftover. **NOTE:** user clarified — the camp_grounds → camp_lake left exit was correct and was un-touched on the second pass (the previous "remove lake exit" change has been reverted: lake hotspot restored in `camp_grounds.json`, flower restored to `camp_lake`, Higgins+Lily hints restored to "the lake").
- [ ] `[P1]` Henri give-jam one-shot sprite. Prompt landed at EXTRA_PROMPTS §V (`npc_henri_give_jam.png`, 6×1, 800×170). Engine wiring will hook the one-shot into Henri's coffee-trade altDialog once the PNG is on disk.
- [x] `[P1]` Camille sketch should be shown from her hands, not as a separate image. FIXED in prompt: EXTRA_PROMPTS §T rewritten so the existing 8-frame sketching one-shot now ends with Camille turning her sketchpad TOWARD the camera in frames 7-8, revealing the drawing ON the page. §U (separate sketch card) deleted — no separate overlay needed. Existing `npc_camille_sketching.png` may need a regen with the 7-8 reveal frames added.

### Reported (2026-05-21) — playtest regression sweep + new asks (24 items)

- [x] `[P1]` 1. PP walks LEFT of the camp path. FIXED: `camp_grounds.json` walkSegments shifted +3 Y so visible foot lands ON the path.
- [x] `[P1]` 2. Tommy/Jake don't walk off after dialog. FIXED: added `tommy_exit.json` + `jake_exit.json` sequences, registered `walk_left`/`walk_back` one-shots on the kid factories (`game/npc.go`), wired exit-on-first-Day-1-dialog in `setupCampCallbacks` (`game/game.go`). Walk PNGs already exist on disk (`npc_tommy_walk_left.png`, `npc_jake_walk_back.png`) so the swap fires immediately. Optional character-identity polish regens for both sheets tracked in EXTRA_PROMPTS §E/§F — current sheets play but read as a generic camp kid (green shirt) rather than the idle Jake/Tommy palette.
- [ ] `[P2]` 3. Windows taskbar / .exe icon. Prompt landed at EXTRA_PROMPTS §O. Wire-up via `rsrc -ico` + `.syso` deferred until user generates the icon PNG.
- [~] `[P1]` 4+5. PP idle vs talk size + PP talk side regen. **Art-only**, prompts in EXTRA_PROMPTS §A/§B already; no engine change needed once the PNGs land.
- [~] `[P1]` 6. PP grab flower size. **Art-only**, EXTRA_PROMPTS §C.
- [x] `[P1]` 7. Lily flower-give gate (general rule). FIXED: new `altDialogStrictMissingHint` field on `npc` struct (`game/npc.go`); engine plays the hint when alt-dialog gate fails (`game/player.go:startNPCDialog`); wired on Lily in `setupCampCallbacks`. Generic mechanism — reusable for Marcus postcard, Jake coin-rubbing, etc.
- [~] `[P1]` 8. Higgins shout sprite still not playing. **DEBUG NEEDED**: code path looks correct on inspection (shout PNG loaded as `oneShotAnims["shout"]`, sequence step uses `"anim": "shout"`, default case in `SeqNPCAnim` calls `swapIdleForOneShot`). Suspect silent-NPC flag may interact. Verify with playtest first; if still failing add log to `SeqNPCAnim` default case.
- [~] `[P1]` 9. Marcus strange idle two-frames-at-one. **Art-only**, EXTRA_PROMPTS §D2 prompt landed for regen with silhouettes confined to cell width.
- [x] `[P1]` 10. PP idle visible in Marcus room last second. FIXED: reordered `night_bedtime.json` so transition to camp_night happens BEFORE the `hide_player false` step — no more leak frame.
- [x] `[P1]` 11. PP walks through Marcus room table. FIXED: added blocker rect `{x: 850, y: 460, w: 380, h: 240}` to `marcus_room.json`.
- [x] `[P1]` 12. Higgins office "sitting behind the desk" + smooth frames. FIXED: bounds (1091,365,182,235) → (990,290,220,200) so head clearly above desk and sprite bottom rests at desk surface y=490; `talkFrameSpeed` 0.18 → 0.10 so idle cycles at 0.25s not 0.45s.
- [x] `[P1]` 13. Throw/catch map bad aim direction. FIXED: new `SeqTweenItem` step type in `game/sequence.go` lerps a sprite across screen; `higgins_give_map.json` rewritten to throw the travel-map projectile right-to-left from Higgins's hand at (1100, 380) to PP's hands at (320, 520) over 0.8s before PP's receive_map plays. Engine has new `SequencePlayer.Draw` hook called from main draw loop.
- [~] `[P1]` 14. Travel map clicking wrong city. PARTIAL FIX: shrunk `pinHitRect` from 90×110 to 60×70 so adjacent pins can't share a rect; distance tie-break already in place. Verify with playtest — if still mis-clicking, additional debug needed.
- [x] `[P1]` 15. Airplane two-row vertical jump. FIXED: measured actual per-cell plane positions in the PNG — row 0 plane at cell-Y 212, row 1 at cell-Y 68. Last pass had the offset sign FLIPPED. Now `flight_cutscene.go` ADDS 216 px to row-1 frame Y (was subtracting 128) so row 1 planes align with row 0.
- [x] `[P1]` 17. Pierre on air in Paris. FIXED: bounds (820, 390) → (780, 470) — moved left 40 + down 80 so feet land on the mid-distance cobblestones.
- [x] `[P1]` 18. PP on air in Paris. FIXED: `paris_street.json` spawnY 365 → 450 (foot at 720, grounded).
- [x] `[P1]` 19. Madame Colette sprite not smooth. FIXED: all Paris NPC `talkFrameSpeed` 0.12 → 0.10 (idle cadence 0.30s → 0.25s).
- [x] `[P1]` 20. Rolling pin moves OUTSIDE the bakery. FIXED: removed floor-item from `paris_bakery` setup; added new floor-item to `paris_street`. User 2026-05-21 refinement: pin is now in the **bicycle basket** on the cobblestones at (320, 490, 80, 55) — previous café-table placement was rejected and the matching `PP grab rolling pin.png` sprite read as awkward (table appearing/disappearing in cells). EXTRA_PROMPTS §P rewritten for the new pose: PP reaches FORWARD at chest height into the bike basket (no crouch, no table in cell), pulls out the pin, holds it overhead. Pickup dialog updated.
- [x] `[P1]` 21. Inventory left/right arrows. FIXED: positions hardcoded to (453, 546) and (929, 547) in `game/inventory.go` per user spec.
- [x] `[P1]` 22. Madame Poulain to oven. FIXED: bounds (540, 490) → (780, 380) — behind the counter near the oven. She still trades the rolling-pin for baguette (and now also for the new Café au Lait).
- [x] `[P1]` 23. 6 café patrons + story item. FIXED: added 6 new NPC factories (`newCafePatronYvette/Bernard/Camille/Henri/Lucien/Elise`) in `game/npc.go`, registered in `game/npc_factory.go`, wired into `paris_bakery.json` npcs list. **NEW QUEST CHAIN**: Henri carries the coffee-jam beat — PP brings Café au Lait → Henri gives Confiture → Pierre's 2-stage trade (Baguette stage 1, Confiture stage 2) → Press Pass. Two new items in `items.json` (`cafe_au_lait`, `confiture`). Camille has a sketch beat (no inventory exchange). Patron art still needs the 12 sheets from EXTRA_PROMPTS §7.1-7.6.
- [x] `[P1]` 24. Can't go right after museum ticket. FIXED: trade now consumes the Press Pass (`inv.giveItemTo("Press Pass", "claude")`) AND adds the Museum Ticket, so the gate's `hasItem("Museum Ticket")` check succeeds cleanly.
- [x] `[P2]` 25. Plan-mode before acting. Done — plan file approved, all 24 items addressed.

### Reported (2026-05-20) - sizes, sprites, story progression sweep

- [~] `[P1]` Size re-audit: PP idle vs talk drawn at different visual sizes. FIXED engine-side compensations land via cell-size standardization in EXTRA_PROMPTS §A-§D regens. JSON bounds adjusted for Paris NPCs (Poulain Y 400→490, Nicolas 950/400→950/490 + W 106→86, Colette Y 400→490, Pierre to back-of-line, Claude Y 390→480, Marcus room Y 350→385, office Higgins X 1106→1091 / Y 480→365). PP talk/side/grab-flower canvas regen pending (art-only).
- [ ] `[P1]` PP walk-the-line: PP foot is NEAR the camp path but not on it. **File:** `assets/data/scenes/camp_grounds.json` walkSegments + `game/player.go` foot anchor.
- [ ] `[P1]` After PP talks to Tommy, Tommy should walk left off-scene. **File:** `game/game.go` Tommy post-dialog callback (line ~396). Prompt for `npc_tommy_walk.png` landed in EXTRA_PROMPTS §E; wire-up after PNG.
- [ ] `[P1]` After PP talks to Jake, Jake should walk into his room. **File:** `game/game.go` Jake post-dialog callback (line ~297). Prompt for `npc_jake_walk.png` landed in EXTRA_PROMPTS §F; wire-up after PNG.
- [ ] `[P1]` PP talk sprite shows two frames at once + one frame makes PP "jump". Regen prompts landed in EXTRA_PROMPTS §A (`PP talk front`) and §B (`PP talk side`). Art-only.
- [x] `[P1]` Higgins walk-back too slow / not smooth. FIXED: swapped-idle frames now cycle at 0.10s (was 0.45s) in `game/npc.go:update`; `higgins_walk_in.json` move duration shortened 3.0→1.8s.
- [~] `[P1]` Marcus talking and idle different sizes; talk two-frames-at-one. Engine reads `kids.json` 8×2 talk grid correctly; sheet regen to 7×2 (matching idle) tracked in EXTRA_PROMPTS §D. Art-only.
- [ ] `[P1]` PP pick-flower sheet: two frames at one. Regen tracked in EXTRA_PROMPTS §C. Art-only.
- [ ] `[P1]` After PP has flower, talking to Lily doesn't auto-give it. Needs `altDialogRequiresItem` wiring on Lily + `inv.giveItemTo` on dialog end. **File:** `game/npc.go` newLily + `game/game.go` setupCampCallbacks. Tracked for next pass.
- [x] `[P1]` Higgins night "go to sleep" plays normal talk sprite instead of shout. FIXED: `npc_director_higgins_shout.png` now loaded as `oneShotAnims["shout"]` on newNightHiggins; `night_bedtime.json` swapped `anim: "talk"` → `anim: "shout"`.
- [~] `[P1]` Marcus room cleanup. FIXED: Marcus moved down (Y 350→385); sleeping PP overlay gated to `camp_night` only (was leaking into marcus_room during transition). `_strange_alt` ambient hook tracked for next pass (`game/atlas.go` + npc inactivity timer).
- [ ] `[P1]` PP exits his room walking through the table — add blocker. **File:** PP-bedroom scene json `blockers`. Need to identify which cabin scene the user means; deferred pending in-game check.
- [x] `[P1]` Higgins office Higgins position wrong + PP on air. FIXED: office Higgins bounds (1091, 365) in `game/npc.go`; `camp_office.json` spawnY 365→480 + spawnX 150→220 so PP foot lands on the office floor.
- [x] `[P1]` Camp-return arrow in office points down-right; should point left. FIXED: `camp_office.json` hotspot arrow `downRight` → `left`.
- [ ] `[P1]` PP should talk to Higgins from farther; Higgins should throw map and PP should catch it. Plan + art prompts landed in EXTRA_PROMPTS §I (`npc_director_higgins_throw_map.png` + `pp_catch_map.png` + `inv_travel_map_throw.png`). New sequence + tween_item step type pending after PNGs. Current give_map sequence still plays as fallback.
- [ ] `[P1]` Map item sprite not clear at all — regen. Tracked in EXTRA_PROMPTS §J (`travel_map_icon.png` regen). Art-only.
- [ ] `[P2]` Unified action arrow design — single sprite used on PP-hover (inventory) and travel-map relevant pin. Plan + prompt landed in EXTRA_PROMPTS §K. Engine wiring pending after PNG.
- [~] `[P1]` Airplane bounce. STOPGAP FIX: row-1 frames now lift -127px in `game/flight_cutscene.go` to compensate for in-cell Y drift; sheet regen with locked fuselage centerline tracked in EXTRA_PROMPTS §H. Once art lands, drop the offset block.
- [x] `[P1]` Paris NPCs floating at y≈627. FIXED: Poulain, Nicolas, Colette, Claude all moved Y +90 in `game/npc.go` so feet land on the street floor (~y=720). Pierre placed back-of-line per Step 6 (Y=390 W=95 H=175).
- [x] `[P1]` Madame Colette talk sheet two-frames-at-one. FIXED: `game/npc.go:1034` now loads talk as 8×2 (was 8×1); `paris.json` already declared 8×2 → engine + data now agree. PNG is genuinely 8 cells wide (1686/8=210).
- [~] `[P2]` Pierre shrink + back-of-line. PARTIAL: Pierre bounds shrunk + moved back (820,390,95,175). PP's existing `depthScale` (driven by player.y) auto-shrinks PP when he walks up to Pierre. Smooth tween between scales is the engine default. Verify in playtest.
- [~] `[P1]` Nicolas click opens map. FIXED: Nicolas hit-rect W shrunk 106→86 in `game/npc.go` so it no longer overlaps with the Louvre exit hotspot (x≥1300). Verify in playtest.
- [~] `[P1]` Story progression: bakery lady + officer dialog stale. FIXED text + wiring: `higginsPostMarcusHealedDialog` plays in office once `marcusHealed=true` (points PP at Lily/Tokyo); `bakeryWomanLouvreSouvenirDialog` armed in Paris bakery callback once `marcusHealed=true` (asks for Louvre postcard for grandson). 8 missing item PNGs (postcard, baguette, rolling_pin, press_pass, coin_rubbing, pressed_sakura, dance_card, inscription_rubbing) tracked in EXTRA_PROMPTS §L.
- [ ] `[P1]` Wire the 6 café patrons (Yvette, Bernard, Camille, Henri, Lucien, Élise) into `paris_bakery` as clickable flavor NPCs with idle+talk anims and short dialogs (each gives a hint as described in `docs/STORY.md` Paris bakery flow). Needs the 12 new patron sheets from EXTRA_PROMPTS §7.1–§7.6 first, then new NPC factories in `game/npc.go` + bounds in `paris_bakery.json`'s `npcs` list.

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
- [ ] `[P2]` No airplane background — currently using `paris_clouds.png` as fallback. Regen prompt now landed at `EXTRA_PROMPTS.md` §NEW Paris Clouds (full 1400×800 sky, parallax-friendly clouds, biplane avoid-zone). PNG generation pending.
- [ ] `[P1]` Paris bakery café-corner rework — six ambient patron sheets (`EXTRA_PROMPTS.md` §7.1–§7.6) and BG regen (`§NEW Paris Bakery`) drafted. After PNGs land, wire `paris_bakery.json`: add the six patron NPC ids to `npcs`, move Madame Poulain bounds to `(820, 440, 140, 240)`, move the rolling-pin floor item to `(740, 720)`. Also write the ambient renderer hookup that was deferred when `cafe_patrons.png` first landed.

---

## Resolved

### Reported (2026-04-19 pass 2) - FIXME sweep (map + Paris NPCs + fire + cabin doors + Higgins post-map + Marcus freakout + Paris quest chain)

- [x] `[P1]` "_map_ location in the map got bg around them" — FIXED: ran `tools/clean_landmarks.py` (one-shot flood-fill-from-edges color-key pass) over every PNG in `assets/images/ui/landmarks/`. Christ Redeemer shed 84% baked-in bg; every other landmark 3–5%. Runtime loader already uses `SafeTextureFromPNGRaw` so the cleaned alpha renders cleanly.
- [x] `[P1]` "i try to click on brazil spot to get info and it took me to paris" — FIXED: map hit-rect shrunk from 110×140 to 90×110 (the label box no longer bleeds into the adjacent pin) AND when two rects overlap the closest-pin-center wins via `distanceSqFromPin` tie-break. See `game/travel_map.go:pinHitRect` / `hitTest` / `hitTestAny`.
- [x] `[P1]` "i want the info to stay in the map screen and not jump to the pp location back every time" — FIXED: new `game/travel_map_panel.go` renders a 720×400 card overlay on the globe (landmark image on the left, bulleted facts on the right). Map stays visible behind. Click-anywhere or Esc dismisses the panel, map stays open.
- [x] `[P1]` "for each location add at least 3 infoes and the famous location" — FIXED: `assets/data/travel_map.json` schema extended with `facts: []` (a list of paragraph-style strings). Every city now has 3 facts; legacy single-line `info` is kept as a fallback for backwards compat.
- [x] `[P1]` "paris people standing on air y~585" — FIXED: Madame Colette / Pierre / Claude bounds Y moved from 340–360 down to 430–440 so feet land at y≈680 on the street line.
- [x] `[P1]` "fire animation is huge... around (577,591)-(700,590)" — FIXED: day-grounds + night fire particles, smoke, and glow centers shifted from (620, 520) to (622, 573). Glow rect resized to `{x: 560, y: 555, w: 130, h: 45}` so the visible flame falls roughly inside the user's target band.
- [x] `[P1]` "i already place the right points where are the doors of each cabinet. fix it" — FIXED: `assets/data/scenes/camp_grounds.json` cabin hotspots swapped from 240×200 blanket rects to 120×90 zones centered on user-specified coords: Tommy (195,479), Jake (441,441), Marcus (820,435), Lily (1077,403), Danny (1243,503). All with `arrow: "up"`.
- [x] `[P1]` "walking in the camp should also have a logical routes" (partial) — FIXED: 5 new vertical walk-segments branching from the main path at y=480/500 up to each cabin door coord, so PP's snap-to-path lands on the door instead of cutting across bushes.
- [x] `[P1]` "higgins office... after the talking is finished, you can change to text to something like: i already gave you the map, comeon panther we need to fix this up" — FIXED: `higginsPostWorriedDialog` rewritten: _"I already gave you the map, Panther." / "Come on — we need to fix this up. The kids are counting on us." / "Marcus is in the camp grounds. Start there."_
- [x] `[P1]` "marcus freak out sprite is too fast" — FIXED: new `strangeTalkFrameSpeed` field on `npc` (0 = inherit talkFrameSpeed). Marcus overrides to 0.16 (was 0.10 default) — 60% slower.
- [x] `[P1]` "...add another freakout sprite that will run if we dont click on him for a few seconds" — SHEET LANDED (see Resolved 2026-04-19 campaign): `npc_marcus_strange_alt.png` generated, yaml extended, atlas repacked. Inactivity-trigger code hook still pending (tracked in Deferred below).
- [x] `[P0]` "i want to create now a story and object we need to collect before we enter the louver. so we need a bagguete" — FIXED: Paris pre-Louvre quest chain landed. New NPC `Madame Poulain` (Bakery Woman, placeholder sprite pointing at the french_guide sheet). 3 new items (`baguette`, `press_pass`, `museum_ticket`). Quest flow: bakery → Pierre (baguette trade for press pass) → Claude (press pass for museum ticket) → Louvre entrance now gated on `Museum Ticket` via `setupParisCallbacks` hotspot override.
- [x] `[P1]` "remove left arrow to the map" — verified: only "To the Lake" / "Higgins' Office" / cabin entries remain on camp_grounds; the travel-map opener on camp_entrance is a gated bus-stop hotspot (requires map item), not a free left-arrow.
- [x] `[P2]` "Buenos Aires and Rio pins are close together — landmarks may overlap visually" — FIXED indirectly: 90×110 hit rect + nearest-pin tie-break means overlapping pins now route clicks to the intended pin.

**Deferred to EXTRA_PROMPTS.md (pass-2 snapshot — most items landed in the
2026-04-19 campaign below; see that Resolved block for status per item):**

- §1 Higgins entrance idle redesign — LANDED (canonical).
- §2 Higgins walk-back — LANDED (canonical).
- §3 PP walk-back for "leaving camp" transition — still deferred.
- §4 Marcus strange_alt — LANDED.
- §5 PP flower pickup — still deferred.
- §6 Campfire small sheet — LANDED.
- §7 Paris cafe ambient patrons — LANDED (asset); renderer hookup deferred.
- §8 Bakery Woman dedicated sheet — LANDED + wired.
- §9 Press Photographer sheet — LANDED + wired.

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
- [x] `[P2]` Migrate Paris NPCs (`french_guide`, `museum_curator`, `pierre_artist`, `gendarme_claude`, `bakery_woman`, `press_photographer`) to the atlas pipeline — LANDED 2026-04-19 (see Resolved campaign block below). Manifests live under `tools/characters/paris/`; packed atlases at `assets/sprites/paris/<name>.(png|json)`. Factories now prefer atlases via `applyNPCAtlas` and fall back to legacy `loadNPCGrid*` PNG slicing when the atlas is missing.
- [ ] `[P1]` Higgins walk-in cutscene from office when Lily shy dialog ends — needs an `npc_move_to` sequence step. Tracked as task #11.
- [ ] `[P1]` Story-flag collapse from Game struct into VarStore — ~10 fields (marcusHealed, metKids, talkedToMarcus, parisUnlocked, nightSceneDone, day1BedtimeStarted, day2Started, monologuePlayed, parisMonologuePlayed). Blocked on setupCampCallbacks closure migration (see above).
- [ ] `[P2]` Re-introduce Camp Chilly Wa Wa as a travel-map destination when post-Paris quests require going back (retro _Hokus Pokus Pink_ style). Currently omitted from `assets/data/travel_map.json` so the player can't loop to camp_entrance and retrigger Higgins's introduction dialog. When re-adding: give Camp a `relevantWhen` expression (e.g. `vars.chapter.paris.return_to_camp == 1`) AND gate entrance-Higgins's `dialog` field on a VarStore key so returning visits show a "welcome back" variant instead of the initial greeting.
- [ ] `[P2]` Drop location voice clips into `assets/audio/locations/<id>.wav` and wire the paths into the `audio` field of each location in `assets/data/travel_map.json`. Playback is already hooked up in `audioManager.playSFX` — the popup just won't speak until the files exist.

| Issue                                                                 | Date       | Notes                                                                                                                                                                |
| --------------------------------------------------------------------- | ---------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Map used dots instead of landmarks                                    | 2026-04-08 | Replaced with landmark images from ui/landmarks/                                                                                                                     |
| Night scene too simple (instant Marcus room)                          | 2026-04-08 | Reworked to multi-phase: campfire sleep → Marcus freakout → wakeup                                                                                                   |
| No airplane transition before cities                                  | 2026-04-08 | Added airplane_flight scene with 4-second cutscene                                                                                                                   |
| Kid rooms have no NPCs inside                                         | 2026-04-16 | `tommy_room`, `jake_room`, `lily_room`, `danny_room` in `game/scene.go` now host the kid NPC; silent until healing chain activates them from `game/game.go` ~260-470 |
| Postcard not added to inventory after Curator dialog                  | 2026-04-16 | Paris handoff + Marcus heal at `game/game.go:461` give the Postcard to PP and consume it on Marcus                                                                   |
| No Marcus healing flow                                                | 2026-04-16 | Marcus `altDialogFunc`, `VarMarcusHealed`, Jerusalem unlock, day-bg restore (`game/game.go:449`)                                                                     |
| Night scene complete rework                                           | 2026-04-16 | 5-phase sequence in `nightSceneArrival` / `nightSceneUpdate` (Higgins speech, PP sleeping/waking, Marcus freakout, Day 2 transition)                                 |
| Map landmark positions                                                | 2026-04-16 | Travel map pins placed at user coords; Rome added; BA and Mexico pins in place; Paris click-to-fly working                                                           |
| Airplane animation 3-row sheet                                        | 2026-04-16 | Loaded as 4x3 via `SpriteGridFromPNGClean`, drawn in `airplane_flight` with bob                                                                                      |
| Flying / map / Paris travel broken                                    | 2026-04-16 | `pinHitRect` 110x140 + `HandleClick` rework in `game/travel_map.go` and `game/game.go`                                                                               |
| Rooms dots too small / exit radius too tight                          | 2026-04-16 | Cabin hotspots enlarged to 240x200 in `game/scene.go`                                                                                                                |
| Can't move between inventory items / yellow arrows / circle too small | 2026-04-16 | 720x600 oval + chevrons + 0.20 click zones in `game/inventory.go`                                                                                                    |
| Hard to find talk spot on kids                                        | 2026-04-16 | Partial: `npc.containsPoint` padX=70, padY=50 in `game/npc.go`; cursor/hover alignment still open                                                                    |
| Can't exit Higgins office / can't reach it                            | 2026-04-16 | Office hotspot + exit bounds landed, walk paths extended                                                                                                             |

---

_When adding issues, check STORY.md to verify expected behavior._

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
      assets\images\player\PP grab.png. so we need to modifty it to be with the same name. i want to see the flower in my inventory as it should be - in addition we need to generate pp bring the flower - lily geting the flower
      -fit the flower idle to the one we picking up and change the name from grab to taking flower or something.
- [ ] `[P1]` change the talking with danny, hes not behind tree in the first conversation
- [ ] `[P1]` pp id is one frame

- [ ] `[P0]` night schen again! focus. 1. higgns(need to be already in the frame in the right bottom corner) need to say that it become late. 2. we see the pp in the middle of the camp with fire turn on assets\images\locations\camp\campfire_idle.png(for some reason its four rows of frames)
      3.we only hearing marcus freak out
      4.then!! moving to his room and see his freak out assets\images\locations\camp\npc\kids\marcus\npc_marcus_strange_talk copy.png over and over 5. morning, pp is waking up . same spot as the sleeping around (298,582) assets\images\player\pp_sleeping.png,assets\images\player\pp_waking.png after finish the senario he need to speak front and said he heard something wired... 6. serching marcus 7. speak to him and he moving between idle and talkin freak out. 8. goin to higgins office, speak about the 9. instruction about the map 10. using the map
- [ ] `[P1]` i added bottom right arrow.
- [ ] `[P1]` higgins office when speak to him we need to change to position of him to (1065,413)
- [ ] `[P1]` i want to generate a new animation, higgins give us the map, then we need to walk to him,generate a new idle of us taking it and put in pocket
- [ ] `[P1]` when getting out of higgins office, we need to go out from the botton right corrner with walking back animation
- [ ] `[P1]` you didnt use the location in the map!! and the map isnt working so does the airplane animation assets\images\player\pp_airplane.png its 3 lines - location:
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

| What Retro Does            | What We Have Now                       | What We Need                      |
| -------------------------- | -------------------------------------- | --------------------------------- |
| Handler+Condition          | `setupCampCallbacks()` closures        | Declarative handler registry      |
| Sequences                  | `nightSceneArrival()` nested callbacks | Sequence player with steps        |
| Game/Module/Page Variables | 15 flat `bool`/`int` fields on Game    | `VarStore` with 3 scopes          |
| Item Ownership             | `inv.hasItem("name")`                  | `item.owner` field                |
| Scene Data Files           | 800+ lines hardcoded in scene.go       | JSON scene files                  |
| Dialog Data Files          | 500+ lines hardcoded in npc.go         | JSON dialog files                 |
| NPC State Machine          | `onDialogEnd` + manual swap            | Named states with auto-transition |
| Walk Locations             | `walkSegments` line pairs              | Polygon walk zones                |
| Save/Load                  | not implemented                        | VarStore serialization            |
| PDA                        | simple map overlay                     | Multi-page UI system              |

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
- [x] `[P1]` lily dialog is still wrong. the first time i click on her is like i gave the flower — FIXED: item-in-bag gate + cursor hint in pass 2 (see Reported 2026-04-16 pass 2 above).
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

### Resolved (2026-04-19) - Hokus-Pokus style unification campaign

**[P1] regen-campaign checklist** — COMPLETED this session. Every sheet
called out in the prior P2 items below and in the 2026-04-19 style-sweep
(see `docs/EXTRA_PROMPTS.md` "Character sheet regen status — updated
2026-04-19" table) is now landed. Style lock anchored on Hokus Pokus Pink
cartoon (canonical refs = new Higgins idle, Danny, Lily, PP idle_front).

Pipeline: `GenerateImage` prompt → `tools/clean_generated_sheet.py` (strips
the baked-in black frame + grid lines that the generator bakes in) → drop
at the target path → `python tools/pack_atlas.py tools/characters/<name>.yaml`
for the atlas-loaded NPCs (Higgins / kids) or straight `loadNPCGrid*` for
legacy loaders (Paris NPCs).

**Camp — Higgins:**

- [x] `[P1]` Regenerated `npc_director_higgins_talk.png` per `§10`. Matches new idle (silver hair, red lanyard, khaki trousers) — dropped the ruddier-face / olive-pants drift. Atlas repacked.
- [x] `[P1]` Regenerated `npc_director_higgins_office_idle.png` per `§18`. Pixel art → cartoon.
- [x] `[P1]` Regenerated `npc_director_higgins_office_talk.png` per `§18`. Pixel art → cartoon.
- [x] `[P1]` Regenerated `npc_director_higgins_give_map.png` per `§19`. Pixel-leaning → cartoon.

**Camp — kids:**

- [x] `[P1]` Regenerated Tommy 4-sheet set per `§11` (idle + talk aligned to canonical strange_idle / strange_talk — tan skin, tousled brown hair, green pine-tree tee, navy jeans, barefoot). Atlas repacked.
- [x] `[P1]` Regenerated Jake 4-sheet set per `§12`. All four sheets on-style; yaml `take_row: 1` dropped from idle so both rows render.
- [x] `[P1]` Regenerated Marcus idle + talk + strange_talk per `§13` (canonical identity anchored on strange_idle: spiky brown hair, yellow polo, khaki cargo shorts, brown ankle shoes). YAML grid normalized to 8×2 across all four animations.
- [x] `[P1]` Generated Marcus `strange_alt` per `§4`. Sheet landed in yaml as a 5th anim. Inactivity-trigger code hook is still a future enhancement but the asset is production-ready.

**Paris:**

- [x] `[P1]` Regenerated French Guide (`npc_french_guide_idle.png` 8×2 + `npc_french_guide_talk.png` 8×1) per `§14`. Pure pixel → cartoon. Also used as the back-sprite placeholder by Nonna/Obachan/Abuela/Marisa/Lucia/Miriam so every chapter benefits.
- [x] `[P1]` Regenerated Museum Curator (`npc_museum_curator_idle.png` 8×1 + `npc_museum_curator_talk.png` 4×2) per `§15`. Updated the §15 canvas/cell dims from the authored 1376×384 / 688×768 target to the actual 1376×768 canvas that the generator produces, with cells 172×768 (idle) and 344×384 (talk) — matches the live `loadNPCGrid` args in `newMuseumCurator`, so no constructor change.
- [x] `[P1]` Generated Bakery Woman (`npc_bakery_woman.png` 8×2) per `§8`. Switched `newBakeryWoman` from the french_guide fallback to `loadNPCGridRow(sheet, 8, 2, 0)` / `(…, 1)` and flipped `flipped: false` since the new sheet draws her right-facing already.
- [x] `[P1]` Generated Press Photographer (`npc_press_photographer.png` 8×2) per `§9`. Added `newPressPhotographer` factory (rows 0/1 = idle/talk) + registered `press_photographer` in `npcFactories` + listed in `paris_street.json` NPCs (X=1010 between Pierre and Claude, so the "photographer near ze museum" breadcrumb in Madame Poulain's post-dialog actually points at someone on stage).

**Environmental / ambient:**

- [x] `[P1]` Generated `campfire_small.png` per `§6`. Sized down to the spec'd 1032×172 / 6×1 / cell 172×172 so the visible flame lands inside the (581,592)-(702,594) band when drawn at 1× scale.
- [x] `[P2]` Generated `paris/ambient/cafe_patrons.png` per `§7`. Updated §7 canvas to the actual 1376×768 / 4×2 grid (8 distinct patrons) that the generator produces. Folder `assets/images/locations/paris/ambient/` created. Paris ambient renderer hookup still deferred — the asset is ready for the next ambient sweep.

**Player (PP):**

- [ ] `[P1]` Regenerate `PP walk back.png` per `§3` (still queued — clearer walk cycle for leaving-camp transition). Deferred to the next PP-focused sweep; out of scope for this campaign.
- [ ] `[P1]` Regenerate `PP grab flower.png` per `§5` (still queued — visible crouch+grab+rise). Deferred.

**Verified KEEP (not regen'd this pass):** Higgins idle / walk / walk_back / shout / desk; Danny all sheets; Lily all sheets; Tommy strange_idle + strange_talk; Marcus strange_idle; PP idle_front / talk_front / idle_side / idle_back / walk_front / walk_left / talk_side / grab / receive_map / celebrate / sneak_examine / sneak_use / pp_sleeping / pp_waking / pp_airplane; Paris art_vendor (Pierre) / security_guard (Claude) / mystery_figure / suspicious_dealer.

**Retired:** `docs/PIXEL_PROMPTS.md` (opposite direction — style direction consolidated onto cartoon only).

**Tooling added this session:**

- `tools/clean_generated_sheet.py` — strips the black frame + grid lines the image generator bakes in, so `pack_atlas.apply_color_key` sees a clean white background and its flood-fill works without chewing character outlines.
- `tools/clean_landmarks.py` — one-shot flood-fill-from-edges color-key pass for the travel-map landmarks (already used by the pass-2 map fix).

**Paris atlas migration (follow-up to the campaign):** Paris NPCs now follow
the same atlas pipeline as the camp kids. Manifests live under
`tools/characters/paris/` (`french_guide.yaml`, `museum_curator.yaml`,
`bakery_woman.yaml`, `press_photographer.yaml`, `pierre_artist.yaml`,
`gendarme_claude.yaml`) with `subfolder: paris` so atlases land at
`assets/sprites/paris/<name>.(png|json)`. `game/atlas.go` got a new
`applyNPCAtlas(renderer, n, "paris/<name>")` helper — the 2-animation
(idle / talk) counterpart to `applyKidAtlas`. Each Paris factory in
`game/npc.go` tries `applyNPCAtlas` first and falls back to the legacy
`loadNPCGrid*` PNG slicing if the atlas is missing, so the fleet still
boots even if `tools/pack_atlas.py` hasn't been run.

### Deferred (after 2026-04-19 campaign)

**Still-open regens (PP only):**

- [ ] `[P1]` `PP walk back.png` per `§3` — queued for next PP sweep.
- [ ] `[P1]` `PP grab flower.png` per `§5` — queued for next PP sweep.

**Code hooks to land separately:**

- [ ] `[P2]` Marcus inactivity timer: after N seconds of no interaction on `camp_grounds`, swap `npc.strangeIdle` for the `strange_alt` atlas anim once, then fall back to `strange_idle`. Asset is in place (`assets/sprites/marcus.json` now has a `strange_alt` entry); the trigger is the only missing piece.
- [ ] `[P2]` Paris `cafe_patrons.png` ambient renderer: draw the 8 patrons as a parallax band at the back of the `paris_street` background, with a gentle per-patron Y bob so they read as alive.

---

**Prior P2 kid-sheet notes (superseded by the campaign above):**

- [x] `[P2]` Regenerate `npc_director_higgins_idle.png` as a clean single-row sheet so we can drop the "row 0 only" / mismatched-grid workaround. — FIXED: regenerated as a clean 7x1 strip at 172x384 per cell (matches `_talk.png` geometry so idle and talk render at the same on-screen size); `newDirectorHiggins` + `newNightHiggins` still call `loadNPCGrid(..., 7, 1)` unchanged. Built via `tools/stitch_higgins_idle.py` from 7 per-pose generations + edge-flood color-key. Stale duplicate at `assets/images/locations/camp/npc/npc_director_higgins_idle.png` deleted.

### Reported (2026-04-16) - this session (pass 3)

- [x] `[P1]` Higgins idle still reads as double on `camp_entrance`. Root cause: the idle sheet was a 2-row mess (row 0 = 8 frames, row 1 = 6) but `newDirectorHiggins` loaded it as `loadNPCGrid(..., 7, 1)`, so cell slicing divided the sheet into 7 wide-and-tall cells that each straddled a pose and its neighbor. — FIXED: regenerated the idle PNG as a clean `7x1` strip at 172x384 per cell via `tools/stitch_higgins_idle.py` (matches talk-sheet cell geometry so both animations render at the same size). `newDirectorHiggins` now loads it cleanly as `loadNPCGrid(..., 7, 1)`; `newNightHiggins` was already on the same path.
- [x] `[P1]` Lily flow still inconsistent — shy dialog and flower dialog could fire in the wrong order depending on scene re-entry because `lilyHinted` was a closure-local that reset whenever `setupCampCallbacks` ran. — FIXED: promoted to `npc.hintState` (per-NPC, survives reloads). Lily's altDialog is armed at setup but gated on `hintState == 1`, so first click is always shy, and flower handoff fires exactly once after the shy beat finishes.
- [x] `[P1]` PP visibly walked around inside Marcus's cabin during the night freakout beat. Root cause: phase 3 cleared `playerSleeping` before transitioning, so `Draw` fell through to `scene.drawActors(..., g.player)`. — FIXED: added `nightHidePlayer` bool set true on phase-3 entry and false on phase-4 exit; `Draw` calls `scene.drawActorsNoPlayer` while the flag is on.
- [x] `[P1]` Warm orange tint bled across the sleep/wake sprites making them look bronzed. — FIXED: `drawWarmTint` now skipped while `camp_night` + `playerSleeping`.
- [x] `[P1]` Cursor had no held-item state — pointer stayed on `cursorNormal` over empty space even while PP was carrying something. — FIXED: `updateHover` now defaults to `cursorGrab` whenever PP is carrying an item. No new cursor asset added — only the existing cursor PNGs are used.

- [] `[P1]` match every higgins design to the last idle.
- [x] `[P1]` when the icon change to click it got white bg. — FIXED 2026-05-09: cursor PNGs regenerated as single-frame portraits + loader now uses `engine.TextureFromPNGAggressive` (tol=16); old 2-frame idle|click split removed in `game/ui.go::drawCursor` and replaced with a triangle scale pulse so click feedback survives without the second frame.
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
- [x] `[P1]` when the icon change to click it got white bg. — FIXED 2026-05-09: cursor PNGs regenerated as single-frame portraits + loader now uses `engine.TextureFromPNGAggressive` (tol=16); old 2-frame idle|click split removed in `game/ui.go::drawCursor` and replaced with a triangle scale pulse so click feedback survives without the second frame.
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

- [ ] `[P1]` logic walking, in the first screen you can walk freely but as a result you walking on unlogic places
- [ ] `[P1]` higgings idle is not the same design we said to use. as a result there is a huge different between the idle and talking — PROMPT WRITTEN: `docs/EXTRA_PROMPTS.md §1` (regen + drop PNG in place; constructor already points there)
- [ ] `[P1]` talking logo got white square around him
- [ ] `[P1]` when we goin the the camp, pp is just walking left. i want him to walk back for a few seconds(make him srink for a seconds ) and then change screen — PROMPT WRITTEN: `docs/EXTRA_PROMPTS.md §3` for PP walk-back art; "shrink-then-transition" hook still needs to land after the sheet drops in
- [ ] `[P1]` tommy framing change is not good
- [ ] `[P1]` jake got white square near to his leg,same as marcus
- [ ] `[P1]` tommy talking sprite is in different color then the other sprites
- [ ] `[P2]` can we make the talkin to work with both click and fluent talking. like if we will add a talking to every npc i should be just like we hear them
- [ ] `[P1]` finally higgins is arrived when lily is she. we need to generate a walking back sprite for this. — PROMPT WRITTEN: `docs/EXTRA_PROMPTS.md §2`; sequence already calls `npc_move` so swapping to walk_back animation is a one-line change once the sheet lands
- [ ] `[P0]` isnt all object are huge? we bearly we see other objects
- [ ] `[P1]` i want to change danny talking sprtie. he got a wired stuff in hes legs
- [ ] `[P1]` when picking the flower i want the picking animation we already got.and also its very hard to find the right place to click in order to pick it up — PROMPT WRITTEN: `docs/EXTRA_PROMPTS.md §5`; hit-zone widening + animation hook still pending
- [x] `[P1]` fire animation is huge, it need to be from like (577,591)-(700,590) — FIXED in pass 2: flame + smoke + glow re-centered at (622, 573) with glow rect `{560, 555, 130, 45}` covering the user's band in both day + night scenes (`game/scene_ambient.go`)
- [ ] `[P1]` pp need to walk to the sleeping point and then the animation need to start(sleeping). also the sleeping got a white bg. also the wake up
- [ ] `[P0]` marcus got a bg aroung him and hes huge in his room. we dont see the text for the entire senario.
- [x] `[P1]` i already place the right points where are the doors of each cabinet. fix it — FIXED in pass 2: `camp_grounds.json` cabin hotspots now 120×90 zones centered at Tommy (195,479), Jake (441,441), Marcus (820,435), Lily (1077,403), Danny (1243,503) with `arrow: "up"`
- [x] `[P1]` marcus freak out sprite is too fast — FIXED in pass 2: new `strangeTalkFrameSpeed` field on `npc`; Marcus set to 0.16 (60% slower than normal 0.10)
- [ ] `[P1]` add another freakout sprite that runs if we dont click on him for a few seconds — PROMPT WRITTEN: `docs/EXTRA_PROMPTS.md §4` for Marcus strange_alt; inactivity-timer code hook still pending the PNG
- [ ] `[P1]` check the size in the room between the pp and marcus
- [x] `[P1]` walking in the camp should also have a logical routes. i made it already check the cords — PARTIAL FIX in pass 2: 5 vertical walk-segments added from the main path up to each cabin door coord; main perimeter still original, more tuning may be needed
- [ ] `[P1]` line to walk to higgins office is not easy to find the right dot
- [x] `[P1]` _higgins in his office_ — after the talking is finished, change to text like: i already gave you the map, comeon panther we need to fix this up — FIXED in pass 2: `higginsPostWorriedDialog` rewritten to this exact message
- [ ] `[P1]` _higgins in his office_ — his not in a correct place to sit / in order to talk with him we standing on the table / giving the map animation isnt implemented — office bounds may need another tweak; give-map animation still pending
- [x] `[P1]` _map_ location in the map got bg around them — FIXED in pass 2: `tools/clean_landmarks.py` stripped baked-in backgrounds from every landmark PNG (3–84% alpha coverage added per file)
- [x] `[P1]` _map_ i try to click on brazil spot to get info and it took me to paris — FIXED in pass 2: pin hit-rect shrunk from 110×140 → 90×110; overlapping rects now resolve to the closest-pin-center via `distanceSqFromPin` tie-break
- [x] `[P1]` _map_ i want the info to stay in the map screen and not jump to the pp location back every time — FIXED in pass 2: new `game/travel_map_panel.go` overlays a 720×400 card (landmark + facts) on the map; map stays visible underneath
- [x] `[P1]` _map_ for each location add at least 3 infoes and the famous location — FIXED in pass 2: `assets/data/travel_map.json` schema extended with `facts: []`; every city has 3 facts
- [ ] `[P1]` traveling got pp for somereason over there and a gray bg
- [x] `[P1]` _paris_ people standing on air y~585 — FIXED in pass 2: Madame Colette / Pierre / Claude bounds Y moved from 340–360 down to 430–440 so feet land at y≈680
- [x] `[P1]` remove left arrow to the map. when we use it in the right time no need arrow to active it — VERIFIED in pass 2: no standalone left-arrow hotspot exists; travel map opens from the gated bus-stop on camp_entrance (requires map item)
- [ ] `[P1]` npc cot lines around them as a result of remove bg
- [ ] `[P1]`

new PR 26/04

- [x] `[P1]` still when the pp is walking to the camp he walk a little left and then we move scene. i want him to walk in the "same place" and srink and then move. — FIXED: new `player.playRecede` tween in `game/player.go`; "Enter Camp" hotspot in `game/game.go::setupTravelHotspots` swapped from `walkToAndDo(599,200)` to `playRecede(1.6, 0.35, 80)`. PP holds X, plays back-walk frames, drifts up 80 px and shrinks 1.0→0.35 over 1.6 s, then transitions.
- [x] `[P1]` lily scenario is broken. i spoke with the kids and then the night scene just started. also the text isnt working at all. — FIXED (root cause): kid Day-1 callbacks (Marcus/Tommy/Jake/Danny) were unguarded so re-clicking after the post-dialog swap kept incrementing `metKids`. Added closure-local `XMet` latches in `game/game.go::setupCampCallbacks`. Belt-and-braces: `checkDay1Complete` now also requires `lily.hintState == 2` so night can't fire before the flower handoff. Dialog text was never broken in code — the bedtime beat just couldn't render because night transition fired during it. NEEDS RUNTIME VERIFY.
- [ ] `[P0]` i left in the commit changes some frames i didnt like because it cut the idle in the middle. most of the because of the same reason. — PARTIAL: documented the loader/grid args per modified PNG in plan. User needs to a) confirm the new PNG dims match the loader's `(cols, rows)`, and b) run `python tools/pack_atlas.py tools/characters/<name>.yaml` for the atlas-loaded NPCs (higgins/tommy/marcus/jake) so the packed atlas catches up to the new sheets.
- [x] `[P1]` i cant click normal on npcs. i starting to guess where is the right spot to talk with them. — FIXED: `npc.containsPoint` now hit-tests against `lastDrawRect` (the actual on-screen rect post-`characterScale` aspect-preserve) instead of the looser `bounds`. Set every frame in `drawScaled`. Falls back to bounds on the first tick. `game/npc.go`.
- [x] `[P1]` when i get the map i dont want the inventory to open. i want the higgings give map sprite and pp get map sprite to work. also i want you to scan the frames because higgins isnt sitting the a logical place. — FIXED for the sequence: new `assets/data/sequences/higgins_give_map.json` plays Higgins's `give_map` one-shot anim, then PP's `receive_map` one-shot, then drops the map silently into inventory. Three new sequence step types (`give_item`, `player_anim`, `npc_one_shot`). Higgins office position SCAN still pending — bounds unchanged this pass; needs a follow-up that reads `camp_office.png` to align torso behind the desk.
- [x] `[P1]` marcus in the room is still huge and the fire is not the new sprite we created. — FIXED: `newRoomMarcus` bounds dropped from 280×380 to 200×300 (kid-sized at scene scale 0.85). Campfire renderer in `game/game.go` now points at `assets/images/locations/camp/campfire_small.png` (6×1) instead of the bulky `campfire_idle.png` (8×4).
- [x] `[P1]` i want a new md file that will keep the design of each npc so we will always have the same design and more messing around with it. — FIXED: per user direction (don't split docs), the per-NPC info table was merged INTO `docs/CHARACTERS.md` rather than created as a sibling file. Each NPC entry covers role, instances, sprite paths, atlas yaml, grid, factory bounds, characterScale overrides, animation speeds, dialog handles, and a regen-prompt cross-link.
- [ ] `[P1]` paris npc still standing too high — NOT DONE in code yet. Plan calls for inspecting `paris_street.png` paving line and updating Y on all 5 paris NPCs. Pending — needs the bg-image inspection pass.
- [x] `[P1]` remove left arrow to open map — FIXED: deleted the `Travel Map` left-arrow hotspot append in `game/game.go::setupParisCallbacks`. Map now opens via the new `inventory.onSelectItem` hook: clicking the Travel Map item in the inventory immediately calls `travelMap.Show(currentScene)`.
- [x] `[P1]` i want to move some npc and generate a backary when click left arrow. — FIXED: new `assets/data/scenes/paris_bakery.json` interior scene, registered in `scene.newSceneManager`. `bakery_woman` moved off `paris_street` into `paris_bakery`. `paris_street.json` left-arrow now opens the bakery. NEEDS ASSET: `assets/images/locations/paris/background/paris_bakery.png` (regen prompt landed in `EXTRA_PROMPTS.md` §NEW: Paris Bakery Interior).
- [x] `[P1]` what will be the story over there — DESIGNED: rolling-pin retro-style intro puzzle. New `Rolling Pin` item (added to `items.json`); `bakery_woman` opens with `bakeryWomanLostPinDialog` (lost the pin), altDialogFunc requires `Rolling Pin` in inventory and trades for the baguette via `bakeryWomanPinTradeDialog`. Floor item registered in `setupParisCallbacks` so PP can find the pin on the bakery floor before handing it back. Existing baguette → press pass → museum ticket chain unchanged downstream.

bonus from this pass (not in original list but landed):

- [x] `[P1]` Dev / chapter-jump menu (F1) — `game/dev_menu.go` adds an in-game overlay listing 11 scenarios (Day 1 fresh, Lily-flower-in-pocket, bedtime, night, Day 2 Marcus room, Day 2 office, Paris fresh, Paris bakery, Paris with press pass, Paris Louvre). Click a row to jump straight there with the right flags pre-set. Toggle with F1.
- [x] `[P1]` PP full-set regen prompt — appended to `docs/EXTRA_PROMPTS.md` (`§NEW: PP full-set regen`) covering all 16 PP sheets in one consistent style lock to fix the color-loss / bg-bleed drift. No code changes — drop the new PNGs in place after generation.

new PR

- [x] `[P1]` higgins idle and talking not cut well. same as the mouse talk image, it got a white bg around. — FIXED: added `engine.TextureFromPNGAggressive` (tol=16 vs default 8) and routed all 8 cursor PNGs through it in `game/ui.go`. Higgins entrance/office/night factories now use new `loadNPCGridClean` / `loadNPCGridRowClean` (kid-grade tol=16 color-key) in `game/npc.go`.
- [~] `[P1]` kids idle and talking need to loan from kids folder. for example tommy is old idle — CODE PATH FIXED: new `applyKidAtlasOrFallback` in `game/atlas.go` falls through to `kids/<name>/npc_<name>_*.png` when the packed atlas is missing, so a kid can never silently render as a frameless ghost. NEEDS USER STEP: re-run `python tools/pack_atlas.py tools/characters/<name>.yaml` for tommy/jake/marcus/lily/danny so packed atlases catch up to the source PNGs (this is what fixes the literal "Tommy is old idle"), and delete the 20 stale duplicate PNGs at `assets/images/locations/camp/npc/npc_<kid>_*.png`.
- [x] `[P1]` when enter to rooms, the pp is goin to the sky insted of shrink or just stand on the door and moving to inside — FIXED in `game/game.go` hotspot handler: `arrow=up` hotspots whose `targetScene` ends in `_room` now `walkToAndDo(door)` then `playRecede(0.7s, scale 0.45, dy 60)` before the scene transition. Reuses the existing camp-entrance recede tween.
- [x] `[P1]` in the first day marcus is both inside and outside his room. — FIXED: `newRoomMarcus` now `hidden = true` by default; `startDay2` unhides room Marcus AND hides ground Marcus. Night-bedtime sequence (`assets/data/sequences/night_bedtime.json`) adds an `npc_hidden hide:false` step before transitioning to `marcus_room` so Marcus appears for the freakout cutscene.
- [x] `[P1]` cant click the cabin exit (and similar hotspots) when PP is standing on top of them — FIXED 2026-05-10 in `game/game.go::HandleClick`: the "click PP to open inventory" branch now defers to any hotspot / NPC / floor-item under the click. Previously PP's 170×235 click rect at (700, 550) overlapped the cabin exit hotspot bounds (200,500,1000,260) and inventory toggling hijacked the exit click. Now inventory only opens when nothing actionable sits under the cursor.
- [x] `[P1]` i cant talk with the kids, they are not clickable. fix the talking spot to be on the npc itself after cut — FIXED 2026-05-10 (same root cause as P0 below): `walkToAndInteract` now snaps + opens dialog immediately when PP is already within 80 px of the talk-target, and stray floor clicks while a pending NPC interaction is in flight are now ignored so the second click doesn't nuke `interactTarget`. `npc.containsPoint` already used `lastDrawRect` so hover and click stay in lockstep.
- [x] `[P1]` when higgins shown up hes not walking, hes talking spirte is displayed. — FIXED 2026-05-10: registered `npc_director_higgins_walk_back.png` as a `walk_back` one-shot anim on `newGroundsHiggins` (`game/npc.go`); `assets/data/sequences/higgins_walk_in.json` now sets anim `walk_back` before the move and `idle` after, with a longer 2.5s arc so the back-walk frames cycle clearly. `SeqNPCAnim` extended in `game/sequence.go` to dispatch named animations to `playOneShotAnim` (treating the registered name as a long-running cycle that the next idle/talk anim ends).

PR 10/05

- [x] `[P1]` tommy talking loosing color after cut, same as jack idle. tommy idle frame change is not cut smooth — FIXED 2026-05-10: `applyKidAtlasOrFallback` in `game/atlas.go` now uses `loadNPCGrid` (default tol=8) instead of `loadNPCGridClean` (tol=16). The aggressive tol=16 was eating pastel shirt pixels on the new 2026-04-29 kid sheets where the BG is already clean enough for default tolerance.
- [x] `[P1]` higgins display after lily first shy dialog is not working well. npc_director_higgins_walk_back.png is not happening and i want it to come from right bottom to the middle of the camp.also his dialog not display at all — FIXED 2026-05-10: `higgins_walk_in.json` retargets the move from (1010, 612) to (700, 612) so he ends in the middle of camp; he now teleports in at (1380, 612) and the back-walk anim plays for the whole 2.5 s move; dialog step (`higginsLilyHintDialog`) plays after the wait. Same change as the previous bullet.
- [x] `[P0]` still cant click right on npcs. for example i cant talk with marcus from (964,417)-(985,571) as the talking icon is displayed. that what i want to happen with every npc. if hte mouse icon change to talking, the dialog need to be able [to start] — FIXED 2026-05-10 in `game/player.go::walkToAndInteract`: when PP is already close (`< 80 px`) to the resolved talk-target, snap there and call `startNPCDialog` immediately instead of running a 1–2 s walk that the user couldn't see fire. Also: `Game.HandleClick` now drops stray floor clicks while PP is mid-walk to an NPC so the user's second click doesn't cancel the pending dialog. Dialog still only ever opens AS A RESULT of a click — no auto-fire from hover.
- [x] `[P1]` it happen multi time. when we go bakc and shrink as it should be, the last frame is display in normal size and i dont want it. so from full size->shrink well-> again full side for a second — FIXED 2026-05-10 in `game/player.go`: the recede-completion path now clamps `recedeScale = recedeEndScale` and KEEPS `recedeActive = true` until the scene transition lands and `sceneMgr.transitionTo` repositions PP at the new spawn (it now calls `player.clearRecede()`). Previously the tween reset to scale 1.0 + idle state immediately, leaving one full-size render frame between the recede and the fade-out → the flash the user saw.
- [x] `[P1]` i cant talk with marcus and danny at all — FIXED 2026-05-10: same fix as the P0 above. The walk-target Marcus/Danny resolved to was further away than other kids (their bounds are at the corners of the camp), so the long walk made it feel like the click was eaten. The 80-px snap now lets the dialog open promptly.
- [x] `[P1]` when i got the flower, the picking animation is not happening, i cant give it to her, i just clicked on her and the dialog start.. — FIXED 2026-05-10: registered `PP grab flower.png` as the `grab_flower` one-shot anim on player (`game/player.go`); flower floor item's `onPickup` callback now plays it for 0.9 s before the inventory pulse + monologue. Lily's `altDialogFunc` kicks off her `receive_flower` one-shot (`npc_lily_receive_flower.png`, 1.4 s) in parallel with the flower-handoff dialog so the give/receive is visible.
- [x] `[P1]` i want to remove unused old photoes. for now not items, only from camp — FIXED 2026-05-10: deleted 16 unreferenced kid PNGs at `assets/images/locations/camp/npc/` (4 each for danny/jake/lily, 3 for marcus, 1 for tommy_strange_talk). Kept the 3 stale Higgins PNGs at the same folder per user direction.

- [x] `[P1]` higgins walikng is showing only the first frame. also i want him to walk from 1266,750->1058,640 — FIXED 2026-05-12: `higgins_walk_in.json` retargeted to teleport (1266,750) → move to (1058,640) over 3.0 s. Frame-cycling bug fixed in `game/sequence.go::SeqNPCAnim`: named anims (e.g. `walk_back`) now SWAP `idleGrid` for the one-shot frames (looping via the normal idle cycler) instead of calling `playOneShotAnim(name, 60.0)` which froze each frame for 7.5 s. New `npc.swapIdleForOneShot` / `restoreSwappedIdle` helpers (`game/npc.go`).
- [x] `[P1]` in the rooms, pp is too small — FIXED 2026-05-12 via the global size rebalance: PP bounds bumped 170×235 → 245×340, all room `characterScale` 0.85 → 1.0, room spawnY shifted −105 to keep PP's foot on the floor. PP now fills the cabin at retro proportions.
- [x] `[P1]` it still hard to talk with danny and marcus — FIXED 2026-05-12: kid grounds bounds widened 150×180 → 183×220 as part of the rebalance. Visible post-aspect-preserve hit rect is now ~100 px wide instead of 81 px, giving a more generous click target.
- [x] `[P1]` pp pick up the flower sprite is not in the same scale as his idle. and he stack after pick up — FIXED 2026-05-12: (1) loader swapped from the SQUARE-cell `PP grab flower.png` (1024×256, 128×128 cells) to the canonical portrait `PP grab.png` (1376×768, 172×384 cells) — matches idle aspect so PP no longer renders shorter during pickup. (2) `playOneShot` completion path in `game/player.go` now resets `state = stateIdle` and `moving = false` BEFORE firing onDone, so PP doesn't stay frozen in the walk pose after the grab finishes.
- [x] `[P1]` flower pickup STILL stuck + white box around PP — FIXED 2026-05-12 (follow-up): (1) `player.playOneShot` completion math used `totalElapsed = stepLen*(idx+1) - (stepLen - timer)`, but `idx` was capped at `len-1` while `timer` kept decrementing — total stuck at `0.9 - stepLen` and never crossed the duration threshold, so the one-shot never completed. Replaced with explicit `oneShotElapsed += dt` wall-clock tracking that always advances. (2) `PP grab.png` had a baked-in gray frame border around a white interior — engine's color-key sampled the gray corners and made gray transparent but left the white interior opaque, producing a white box around PP. Ran `tools/clean_generated_sheet.py` to overwrite the gray border with white; engine now color-keys the whole white BG correctly.

new PR 12/05 — visual cleanup + position drift

- [x] `[P1]` Fire / PP sleeping / PP waking sprites show white BG halo — FIXED 2026-05-12: ran `tools/clean_generated_sheet.py` over `campfire_small.png` (6×1), `pp_sleeping.png` (8×2), `pp_waking.png` (8×2). Same gray-border-around-white-interior pattern as `PP grab.png`.
- [x] `[P1]` Travel map info popup shows "???" before every fact line + landmarks have white BG — FIXED 2026-05-12: (1) bullet "• " → "- " in `travel_map_panel.go` (bitmap font has no U+2022 glyph; renderer mis-treated UTF-8 bytes as three chars → "???"). (2) ran `tools/clean_landmarks.py` over all 12 landmark PNGs.
- [x] `[P1]` PP renders below screen at night + dialogs not shown during cutscenes — FIXED 2026-05-12: (1) `camp_night.json` spawnY 556 → 395 (the +70 shift from the rebalance had drifted it past playerMaxY). (2) `sceneMgr.transitionTo` now clamps `spawnY` to the player's Y range AFTER setting sceneMinY/Max, so any future scene with a drifted spawnY can't drop PP below-screen.
- [x] `[P1]` Higgins office position wrong (head visible in corner only) — FIXED 2026-05-12: bounds Y 402 → 360 (foot 595 — sits properly behind the desk instead of low on the floor).
- [x] `[P1]` Higgins night standing on air — FIXED 2026-05-12: bounds Y 470 → 430 (foot 650, at campfire ground level).
- [x] `[P1]` Paris NPCs standing on air + overlapping each other — FIXED 2026-05-12: all five Paris adults' Y shifted from 440-450 down to 390-400 (foot 620-630 lands on the paving). X spread: French Guide 300 / Pierre 720 (was 880) / Press Photographer 950 (was 1010) / Gendarme Claude 1180 (was 1120). Madame Poulain at X=540 (paris_bakery scene). Museum Curator Y 330 → 320.
- [x] `[P1]` Can't walk LEFT to bakery (Paris street walking stuck) — FIXED 2026-05-12: `paris_street.json` left blocker height 500 → 200. The blocker was covering the entire left edge including the bakery hotspot rect at y=200-700, so PP couldn't pathfind into it.
- [x] `[P1]` Airplane sprite huge and mis-aligned — FIXED 2026-05-12: (1) grid args 4×3 → **6×2** in `loadAirplaneFrames` (real layout per user; previously cropped half-of-frame per cell). (2) render scale 3.0 → 1.5 so the plane fits the 1400×800 background instead of dwarfing it.
- [x] `[P1]` Hard to click Marcus / Danny — visible sprite is narrow column — FIXED 2026-05-12: `npc.lastDrawRect` now padded horizontally by 25 px each side in `drawScaled`. Click + hover hit-test feels generous over the character's actual silhouette without falling off into empty space.
- [x] `[P1]` Walking animation not smooth (PP walk_back, etc.) — FIXED 2026-05-12: `walkFrameTime` 0.12 (8.3 fps) → 0.08 (12.5 fps). Foot Y drift across walk_back frames was already 0-1 px, so the choppiness was purely frame rate. Both PP and any NPC walks now cycle smoother.

new PR 19/05 — playtest fixes (cursor==click, fire pos, Higgins office, landmark BG)

- [x] `[P1]` Marcus / Danny still hard to talk — only right-edge clicks worked — FIXED 2026-05-19: `npc.containsPoint` reverted to bounds-based hit-test (was using `lastDrawRect`, which after my +25 pad was 128 wide inside 145-wide bounds — narrower than expected). Post-rebalance bounds are tight to the visible character (145×175 kids, 200×270 PP-class), so bounds give a natural ~30 px forgiveness margin. `lastDrawRect` still set in `drawScaled` for click-probe diagnostics, just not used for hits.
- [x] `[P1]` Flower pickup hard to find — FIXED 2026-05-19: bounds widened 50×50 → 100×100 in `game/game.go::setupCampCallbacks` (flower floor item). Cursor change in `updateHover` + click in `checkFloorItemClick` both read the same bounds → cursor==click across the full patch.
- [x] `[P1]` Night fire huge and mis-positioned — FIXED 2026-05-19 in `game/game.go::Draw`: scale 2.5 → 1.5; anchor (622, 573) → (646, 598) to sit on the actual fire-pit in `camp_night.png`.
- [x] `[P1]` Higgins office position wrong (head in corner) — FIXED 2026-05-19: `newOfficeHiggins` bounds (1062, 360, 182, 235) → (1106, 480, 182, 235). Higgins's torso/head visible at top-left anchor (1106, 480); PP stands in front of the desk.
- [x] `[P1]` Travel-map landmark images show white BG — FIXED 2026-05-19: user hand-edited PNGs to have white background; `travel_map.go` loader swapped from `SafeTextureFromPNGRaw` (no key) to `SafeTextureFromPNGKeyed` so the corner-sample color-key strips white at load. No PNG re-edits needed.
- [x] `[P1]` Marcus in his Day-2 room standing too high — FIXED 2026-05-19: `newRoomMarcus` bounds Y 290 → 350. Foot at 620 (cabin floor) instead of 560 (mid-room).

new PR 17/05 — playtest follow-up: positions, talk speeds, click areas, sequence ghost
- [x] `[P1]` jake talking sprite moving too fast — FIXED 2026-05-17: all kid `talkFrameSpeed` bumped 0.10 → 0.14 across Tommy/Jake/Lily/Marcus/Danny.
- [x] `[P1]` Higgins shows up after Lily shy dialog isn't displayed — best-effort 2026-05-17: walk-in endpoint Y 640 → 580 so he ends at foot ~790 instead of overlapping the dialog panel band; sequence dialog step plays normally. Verify in-game; if the text is STILL invisible, the next pass needs the dialog system to draw above any NPC overlay.
- [x] `[P1]` Higgins + Lily standing on each other in camp scene — FIXED 2026-05-17: Higgins endpoint Y 640 → 580 (up the camp), Lily bounds Y 400 → 440 (down toward the path). No more clustering.
- [~] `[P1]` After talking, Tommy/Jake/Danny walk away to cabins/offscreen — DEFERRED: needs new walk-away sprite generation (Tommy left, Jake to cabin, Danny down-right) + a `kid_walk_away` sequence. Track for the next art pass.
- [~] `[P1]` Inventory open icon (finger / pickup) — DEFERRED: needs new cursor PNG.
- [x] `[P1]` Camp fire size + dialog + PP sleep/wake position + idle stays at wake spot — FIXED 2026-05-17: fire anchor Y 598 → 615 (sits lower on the fire pit). PP sleep/wake render anchor (335, 591) → (335, 615) so both poses share the same spot. `SeqPlayerSleep` with `hide:false` now snaps `player.x/y` to the sleep anchor + sets `state = stateIdle`, so PP's idle appears at the wake location instead of his pre-cutscene coords. (Dialog visibility ride-along on the Higgins fix above.)
- [x] `[P1]` PP sleeping ghost visible during Marcus freakout — FIXED 2026-05-17: in `assets/data/sequences/night_bedtime.json`, swapped the order so `hide_player:true` runs BEFORE `player_sleep:false`. Eliminates the one-frame walk-pose ghost.
- [x] `[P1]` PP walking back too fast — FIXED 2026-05-17 in `game/player.go`: walk-frame tick is now direction-aware. Forward / side / down walks stay at 0.08 s/frame; `dirUp` (back-walk) is 0.12 s/frame.
- [~] `[P1]` Marcus strange_idle alt variation — DEFERRED: needs anim-swap-with-random-timer logic. Sheet exists but not wired.
- [x] `[P1]` PP walks through desk leaving camp_office — FIXED 2026-05-17: added a second blocker `{x:950, y:480, w:450, h:200}` to `camp_office.json` covering the desk surface.
- [x] `[P0]` Higgins office position wrong — every sprite — FIXED 2026-05-17: `newOfficeHiggins` bounds (1106, 480, 182, 235) → (1113, 382, 182, 235). PP foot max is 665; Higgins foot is 617 — PP stands ~50 px in front of the desk.
- [x] `[P1]` Map western wall opens wrong info (Rome instead) — FIXED 2026-05-17: Rome `pinY` 330 → 280 (clears the 90×110 overlap with Jerusalem at 782, 349). `pinHitRect` + `distanceSqFromPin` tie-break still in place as backup.
- [~] `[P1]` Add traditional dish + population per landmark — DEFERRED: content addition (JSON-only). Track for a later content pass.
- [x] `[P1]` Paris NPCs + PP still standing on air — FIXED 2026-05-17: French Guide / Pierre / Press Photographer / Madame Poulain / Bakery Woman Y 400 → 430. Gendarme Claude Y 390 → 420. Foot lines now ~660 (paving).
- [x] `[P1]` Madame Colette PNG path uses `french_guide` — FIXED 2026-05-17: renamed `npc_french_guide_idle/talk.png` → `npc_madame_colette_idle/talk.png`; updated factory loader paths. (Atlas registration stays at "paris/french_guide" since the packed atlas doesn't exist yet.) Underlying frame-bleed concern (10-frame idle loaded as 8×2) tracked for a separate re-pack pass.
- [x] `[P1]` NPC talking too fast — FIXED 2026-05-17: kid `talkFrameSpeed` 0.10 → 0.14; Paris adults 0.12 → 0.16.
- [x] `[P1]` Ambient cafe NPCs — move from `paris/ambient/` to `paris/npc/coffee/` — DONE 2026-05-17: 6 patron PNGs moved; empty `ambient/` dir removed. No code refs updated (renderer hookup is still deferred).
- [x] `[P1]` Inventory item navigation hard — FIXED 2026-05-17: `handleClick` cut ratio 0.20 → 0.10 in `game/inventory.go`. Prev region widens to x < 628 (was 556), next to x > 772 (was 844). Center "pick" zone shrinks; nav is generous.
- [x] `[P1]` Stuck inside the bakery, can't exit — FIXED 2026-05-17: `paris_bakery.json` left blocker height 500 → 200, freeing the bakery-exit hotspot at y ≥ 200 (same fix as `paris_street` last pass).

new PR 18/05 — opening polish + new-destination workflow
- [x] `[P1]` PP should walk in from off-screen-left at game start, then start the monologue — DONE 2026-05-17 in `game/game.go` opening trigger: PP parks at `x = -200`, walks right to the scene's spawn coords (`allowOffscreen=true` so the boundary blocker doesn't bounce him back), then `onArrival` fires the existing monologue with `state=stateTalking`. `player.update` skips blocker collision while `allowOffscreen=true`.
- [x] `[P1]` Higgins on camp_entrance — PP stands on the fence at his resting spawn — DONE 2026-05-18: Higgins entrance bounds X 660 → 760 (+100 px); camp_entrance `spawnX` 500 → 580 (+80 px). PP lands clear of the left gate post.
- [~] `[P1]` Walk-away sprites for Tommy / Jake after talking — PROMPTS WRITTEN 2026-05-18: `docs/EXTRA_PROMPTS.md` §20 (Tommy walk-left) and §21 (Jake walk-back-up-to-cabin). Both authored as 8×1 portrait strips matching idle palette. Wiring plan documented in the prompts: register as `walk_away` one-shots, lerp bounds X/Y over 2–2.5 s on `onDialogEnd`, then `hidden = true`. Code-side wiring waits on the PNGs landing.
- [ ] `[P1]` New destination — **Stonehenge** (per the PtP "part 3" clip ending: PP flies from London to Stonehenge for a druid puzzle). When the art lands, follow the 5-step workflow now documented in `docs/SKILL.md §4a`. Needs: BG art, scene JSON, travel-map pin + landmark, NPCs (likely 1 druid + a stone-circle puzzle), `setupStonehengeCallbacks` for the story chain.

new PR — JSON dialog parse + opening walk-in tween + Higgins click-to-talk
- [x] `[P0]` Marcus freakout dialog at night not showing AND Higgins post-Lily-shy dialog never appeared — ROOT CAUSE FOUND + FIXED: `dialogEntry`'s fields (`speaker`, `text`, `audio`) were lowercase / unexported, so `json.Unmarshal` silently skipped every `"speaker"` / `"text"` JSON key — every sequence-loaded dialog produced entries with empty strings. The dialog DID start (len(entries) > 0) but rendered a blank panel. FIX: added `UnmarshalJSON` on `dialogEntry` that maps JSON `speaker`/`text`/`audio` keys to the unexported fields. All 117 in-code `dialogEntry{}` literals continue to work. Affects every JSON-loaded sequence dialog (night_bedtime, higgins_walk_in, higgins_give_map, etc.).
- [x] `[P1]` PP not walking in from left at game start — DONE 2026-05-18: replaced moving/allowOffscreen approach with dedicated `playWalkIn` tween (mirrors `playRecede`). Drives p.x directly via lerp, short-circuits the normal moving/clamp/blocker pipeline. Opening trigger calls `playWalkIn(-200, spawnX-100, y, 2.5, monologueStart)`. Side-walk frames cycle during the walk. On arrival fires the monologue.
- [x] `[P1]` Higgins shows up after Lily-shy but text never appeared — DONE 2026-05-18: removed the auto-dialog step from `higgins_walk_in.json`. Higgins now walks in, plays idle, and waits silently. PP must CLICK on him to trigger `higginsLilyHintDialog` (already set on the `newGroundsHiggins` factory). The `npc_hidden hide:false` step un-silents him so the click registers. Also fixes the underlying root cause via the dialog JSON parse fix above.
- [x] `[P1]` Marcus talk / strange_idle + PP talk side regenerated — VERIFIED 2026-05-18: new PNG dimensions still load correctly with the existing `8×2` loader args (Marcus talk 1672×941 → cells 209×470 with 1 px slack; Marcus strange_idle 1648×954 → cells 206×477; PP talk side 1536×1024 → cells 192×512). No code changes needed — the renderer's aspect-preserve scaling adapts to the new cell aspect automatically.
