package game

import (
	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

// devMenu is a developer overlay that lets the user click a scenario name to
// jump straight into that chapter / scene with the right story flags pre-set.
// Toggle with F1. Always available — shipping builds can gate this behind a
// build tag if needed; for now it's part of the QA + dev workflow the user
// asked for ("create a test file so we will be able to start chapters by
// click after we will be happy").
type devMenu struct {
	visible bool
	rows    []devMenuRow
	rowH    int32
	panelX  int32
	panelY  int32
	panelW  int32
}

type devMenuRow struct {
	label string
	jump  func(g *Game)
}

func newDevMenu() *devMenu {
	dm := &devMenu{
		rowH:   34,
		panelX: 80,
		panelY: 80,
		panelW: 540,
	}
	dm.rows = []devMenuRow{
		{"Day 1 — Camp Entrance (fresh start)", jumpDay1Entrance},
		{"Day 1 — Camp Grounds (no kids met)", jumpDay1Grounds},
		{"Day 1 — Flower in pocket (give to Lily)", jumpDay1FlowerInPocket},
		{"Day 1 — Bedtime ready (5 kids met + flower given)", jumpDay1Bedtime},
		{"Day 1 — Night campfire scene", jumpDay1Night},
		{"Day 2 — Marcus's Cabin (strange)", jumpDay2MarcusRoom},
		{"Day 2 — Higgins Office (give map)", jumpDay2Office},
		{"Paris — Street (fresh arrival)", jumpParisStreet},
		{"Paris — Bakery (rolling pin puzzle)", jumpParisBakery},
		{"Paris — Holding Press Pass", jumpParisPressPass},
		{"Paris — Louvre interior (with ticket)", jumpParisLouvre},
	}
	return dm
}

func (dm *devMenu) toggle()        { dm.visible = !dm.visible }
func (dm *devMenu) Visible() bool  { return dm.visible }

// handleClick returns true if the click was inside the panel (consumed).
func (dm *devMenu) handleClick(x, y int32, g *Game) bool {
	if !dm.visible {
		return false
	}
	for i, row := range dm.rows {
		ry := dm.panelY + 60 + int32(i)*dm.rowH
		rect := sdl.Rect{X: dm.panelX, Y: ry, W: dm.panelW, H: dm.rowH - 4}
		pt := sdl.Point{X: x, Y: y}
		if pt.InRect(&rect) {
			if row.jump != nil {
				row.jump(g)
			}
			dm.visible = false
			return true
		}
	}
	// Click outside any row but inside the panel — eat it.
	panelH := int32(60) + int32(len(dm.rows))*dm.rowH + 20
	panelRect := sdl.Rect{X: dm.panelX, Y: dm.panelY, W: dm.panelW, H: panelH}
	pt := sdl.Point{X: x, Y: y}
	if pt.InRect(&panelRect) {
		return true
	}
	// Click outside the panel closes the menu.
	dm.visible = false
	return false
}

func (dm *devMenu) draw(renderer *sdl.Renderer, font *engine.BitmapFont) {
	if !dm.visible {
		return
	}
	panelH := int32(60) + int32(len(dm.rows))*dm.rowH + 20

	// Backdrop dim
	renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND)
	renderer.SetDrawColor(0, 0, 0, 140)
	renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})

	// Panel
	renderer.SetDrawColor(20, 20, 30, 230)
	renderer.FillRect(&sdl.Rect{X: dm.panelX, Y: dm.panelY, W: dm.panelW, H: panelH})
	renderer.SetDrawColor(255, 200, 120, 255)
	renderer.DrawRect(&sdl.Rect{X: dm.panelX, Y: dm.panelY, W: dm.panelW, H: panelH})

	font.DrawText(renderer, "DEV MENU  (F1 to close)", dm.panelX+16, dm.panelY+16, 3,
		sdl.Color{R: 255, G: 220, B: 140, A: 255})

	for i, row := range dm.rows {
		ry := dm.panelY + 60 + int32(i)*dm.rowH
		// Row background
		renderer.SetDrawColor(40, 40, 60, 220)
		renderer.FillRect(&sdl.Rect{X: dm.panelX + 8, Y: ry, W: dm.panelW - 16, H: dm.rowH - 4})
		font.DrawText(renderer, row.label, dm.panelX+18, ry+6, 2,
			sdl.Color{R: 230, G: 230, B: 240, A: 255})
	}
}

// --- Jump helpers --------------------------------------------------------
//
// Each helper mutates Game state then transitions PP into the target scene.
// They are deliberately blunt — flipping flat flags directly rather than
// going through the proper story callbacks — so the dev can land at any
// point in the story without playing through. Don't use these from
// production code paths.

func jumpDay1Entrance(g *Game) {
	g.day = 1
	g.day2Started = false
	g.metKids = 0
	g.day1BedtimeStarted = false
	g.parisUnlocked = false
	g.marcusHealed = false
	g.monologuePlayed = true // skip intro monologue when jumping
	resetCampGroundsKids(g)
	g.sceneMgr.transitionTo("camp_entrance", g.player)
}

func jumpDay1Grounds(g *Game) {
	g.day = 1
	g.day2Started = false
	g.metKids = 0
	g.day1BedtimeStarted = false
	g.monologuePlayed = true
	resetCampGroundsKids(g)
	g.sceneMgr.transitionTo("camp_grounds", g.player)
}

func jumpDay1FlowerInPocket(g *Game) {
	g.day = 1
	g.day2Started = false
	g.metKids = 4 // pretend the other 4 kids have been met
	g.day1BedtimeStarted = false
	g.monologuePlayed = true
	if !g.inv.hasItem("Flower") {
		if item := g.items.createItem("flower"); item != nil {
			g.inv.addItem(item)
		}
	}
	// Lily must already have heard the shy beat so altDialog can fire.
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
		for _, n := range grounds.npcs {
			if n.name == "Lily" {
				n.hintState = 1
				n.altDialogRequiresHeld = false
				n.altDialogRequiresItem = "Flower"
			}
		}
	}
	g.sceneMgr.transitionTo("camp_grounds", g.player)
}

func jumpDay1Bedtime(g *Game) {
	g.day = 1
	g.metKids = 5
	g.day1BedtimeStarted = false
	g.monologuePlayed = true
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
		for _, n := range grounds.npcs {
			if n.name == "Lily" {
				n.hintState = 2
			}
		}
	}
	g.sceneMgr.transitionTo("camp_grounds", g.player)
	g.checkDay1Complete()
}

func jumpDay1Night(g *Game) {
	g.day = 1
	g.day1BedtimeStarted = true
	g.monologuePlayed = true
	g.sceneMgr.transitionTo("camp_night", g.player)
}

func jumpDay2MarcusRoom(g *Game) {
	g.startDay2()
	g.day2Started = true
	g.monologuePlayed = true
	g.sceneMgr.transitionTo("marcus_room", g.player)
}

func jumpDay2Office(g *Game) {
	g.startDay2()
	g.day2Started = true
	g.parisUnlocked = false
	g.monologuePlayed = true
	g.sceneMgr.transitionTo("camp_office", g.player)
}

func jumpParisStreet(g *Game) {
	g.startDay2()
	g.day2Started = true
	g.parisUnlocked = true
	g.travelMap.setUnlocked("paris_street", true)
	if !g.inv.hasItem("Travel Map") {
		if item := g.items.createItem("travel_map"); item != nil {
			g.inv.addItem(item)
		}
	}
	g.sceneMgr.transitionTo("paris_street", g.player)
}

func jumpParisBakery(g *Game) {
	jumpParisStreet(g)
	g.sceneMgr.transitionTo("paris_bakery", g.player)
}

func jumpParisPressPass(g *Game) {
	jumpParisStreet(g)
	if !g.inv.hasItem("Press Pass") {
		if item := g.items.createItem("press_pass"); item != nil {
			g.inv.addItem(item)
		}
	}
	g.sceneMgr.transitionTo("paris_street", g.player)
}

func jumpParisLouvre(g *Game) {
	jumpParisStreet(g)
	if !g.inv.hasItem("Museum Ticket") {
		if item := g.items.createItem("museum_ticket"); item != nil {
			g.inv.addItem(item)
		}
	}
	g.sceneMgr.transitionTo("paris_louvre", g.player)
}

// resetCampGroundsKids zeroes the per-NPC dialogDone / hintState flags so the
// kid intros replay cleanly. Used by the Day 1 jumps.
func resetCampGroundsKids(g *Game) {
	if grounds, ok := g.sceneMgr.scenes["camp_grounds"]; ok {
		for _, n := range grounds.npcs {
			n.dialogDone = false
			n.hintState = 0
		}
	}
}
