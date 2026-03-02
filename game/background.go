package game

import (
	"math"
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type background struct {
	tex  *sdl.Texture
	srcW int32
	srcH int32
}

func newBackground(renderer *sdl.Renderer, path string) *background {
	tex := engine.TextureFromBMPRaw(renderer, path)
	return &background{tex: tex, srcW: 626, srcH: 626}
}

func newSolidBackground(renderer *sdl.Renderer, r, g, b uint8) *background {
	surface, err := sdl.CreateRGBSurface(0, engine.ScreenWidth, engine.ScreenHeight, 32, 0, 0, 0, 0)
	if err != nil {
		panic(err)
	}
	defer surface.Free()
	surface.FillRect(nil, sdl.MapRGB(surface.Format, r, g, b))
	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		panic(err)
	}
	return &background{tex: tex, srcW: engine.ScreenWidth, srcH: engine.ScreenHeight}
}

func newPNGBackground(renderer *sdl.Renderer, path string) *background {
	tex, w, h := engine.TextureFromPNGRaw(renderer, path)
	return &background{tex: tex, srcW: w, srcH: h}
}

// ---------- London Night Background ----------

func newLondonBackground(renderer *sdl.Renderer) *background {
	s, err := sdl.CreateRGBSurface(0, engine.ScreenWidth, engine.ScreenHeight, 32,
		0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
	if err != nil {
		panic(err)
	}
	defer s.Free()

	f := s.Format
	horizonY := int32(440)
	parkBottom := int32(610)
	pathBottom := int32(720)

	// Sky gradient
	for y := int32(0); y < horizonY; y++ {
		t := float64(y) / float64(horizonY)
		r := uint8(8 + t*20)
		g := uint8(12 + t*30)
		b := uint8(35 + t*45)
		s.FillRect(&sdl.Rect{X: 0, Y: y, W: engine.ScreenWidth, H: 1}, sdl.MapRGBA(f, r, g, b, 255))
	}

	// Horizon city glow
	for y := horizonY - 60; y < horizonY+20; y++ {
		dist := math.Abs(float64(y) - float64(horizonY))
		a := uint8(math.Max(0, 35-dist*0.8))
		s.FillRect(&sdl.Rect{X: 0, Y: y, W: engine.ScreenWidth, H: 1}, sdl.MapRGBA(f, 60, 55, 40, a))
	}

	// Stars
	for i := 0; i < 150; i++ {
		sx := int32(rand.Intn(int(engine.ScreenWidth)))
		sy := int32(rand.Intn(int(horizonY - 80)))
		sz := int32(rand.Intn(2) + 1)
		a := uint8(rand.Intn(180) + 75)
		s.FillRect(&sdl.Rect{X: sx, Y: sy, W: sz, H: sz}, sdl.MapRGBA(f, 255, 255, 240, a))
	}

	// Moon
	moonX, moonY, moonR := int32(980), int32(90), int32(35)
	for dr := moonR + 60; dr > moonR; dr-- {
		a := uint8(float64(3) * float64(moonR+60-dr) / 60.0)
		fillCircleSurface(s, moonX, moonY, dr, sdl.MapRGBA(f, 180, 190, 220, a))
	}
	fillCircleSurface(s, moonX, moonY, moonR, sdl.MapRGBA(f, 220, 225, 240, 255))
	// Craters
	fillCircleSurface(s, moonX-12, moonY-8, 6, sdl.MapRGBA(f, 195, 200, 215, 255))
	fillCircleSurface(s, moonX+10, moonY+5, 4, sdl.MapRGBA(f, 200, 205, 220, 255))
	fillCircleSurface(s, moonX+3, moonY-15, 3, sdl.MapRGBA(f, 205, 208, 222, 255))

	bFar := sdl.MapRGBA(f, 12, 15, 28, 255)
	bMid := sdl.MapRGBA(f, 18, 22, 38, 255)
	bNear := sdl.MapRGBA(f, 25, 30, 48, 255)
	winLit := sdl.MapRGBA(f, 255, 200, 80, 255)
	winDim := sdl.MapRGBA(f, 200, 155, 55, 255)

	// ----- Far building silhouettes -----

	// Big Ben tower (left area)
	s.FillRect(&sdl.Rect{X: 155, Y: 120, W: 40, H: horizonY - 120}, bFar)
	fillTriangleSurface(s, 175, 70, 150, 200, 120, bFar)
	// Clock face
	fillCircleSurface(s, 175, 160, 14, sdl.MapRGBA(f, 50, 55, 70, 255))
	fillCircleSurface(s, 175, 160, 12, sdl.MapRGBA(f, 180, 175, 140, 180))

	// Parliament building (left-center)
	s.FillRect(&sdl.Rect{X: 50, Y: 300, W: 260, H: horizonY - 300}, bFar)
	// Roofline crenellations
	for x := int32(50); x < 310; x += 20 {
		s.FillRect(&sdl.Rect{X: x, Y: 290, W: 12, H: 10}, bFar)
	}
	// Small towers
	s.FillRect(&sdl.Rect{X: 70, Y: 260, W: 18, H: 40}, bFar)
	fillTriangleSurface(s, 79, 245, 68, 90, 260, bFar)
	s.FillRect(&sdl.Rect{X: 280, Y: 270, W: 16, H: 30}, bFar)
	fillTriangleSurface(s, 288, 255, 278, 298, 270, bFar)

	// Center buildings
	s.FillRect(&sdl.Rect{X: 350, Y: 320, W: 120, H: horizonY - 320}, bMid)
	s.FillRect(&sdl.Rect{X: 480, Y: 290, W: 80, H: horizonY - 290}, bMid)
	s.FillRect(&sdl.Rect{X: 420, Y: 280, W: 50, H: horizonY - 280}, bFar)
	fillTriangleSurface(s, 445, 260, 420, 470, 280, bFar)

	// Westminster Abbey (right side, gothic with spires)
	s.FillRect(&sdl.Rect{X: 680, Y: 200, W: 250, H: horizonY - 200}, bFar)
	// Central tall spire
	s.FillRect(&sdl.Rect{X: 795, Y: 100, W: 22, H: 100}, bFar)
	fillTriangleSurface(s, 806, 60, 790, 822, 100, bFar)
	// Side spires
	s.FillRect(&sdl.Rect{X: 710, Y: 150, W: 16, H: 50}, bFar)
	fillTriangleSurface(s, 718, 125, 707, 729, 150, bFar)
	s.FillRect(&sdl.Rect{X: 890, Y: 155, W: 16, H: 45}, bFar)
	fillTriangleSurface(s, 898, 130, 887, 909, 155, bFar)
	// Gothic window arches
	for x := int32(700); x < 920; x += 35 {
		fillTriangleSurface(s, x+10, 205, x, x+20, 230, sdl.MapRGBA(f, 20, 25, 45, 255))
	}

	// Far right buildings
	s.FillRect(&sdl.Rect{X: 1000, Y: 310, W: 200, H: horizonY - 310}, bMid)
	s.FillRect(&sdl.Rect{X: 1020, Y: 280, W: 60, H: 30}, bMid)

	// ----- Near foreground buildings with windows -----

	// Left building
	s.FillRect(&sdl.Rect{X: 0, Y: 360, W: 110, H: horizonY - 360}, bNear)
	for wy := int32(375); wy < horizonY-20; wy += 30 {
		for wx := int32(12); wx < 100; wx += 28 {
			c := winLit
			if rand.Float64() < 0.3 {
				c = winDim
			}
			if rand.Float64() < 0.15 {
				continue
			}
			s.FillRect(&sdl.Rect{X: wx, Y: wy, W: 14, H: 18}, c)
		}
	}

	// Right building (newsstand/shop area)
	s.FillRect(&sdl.Rect{X: 1060, Y: 370, W: 140, H: horizonY - 370}, bNear)
	for wy := int32(385); wy < horizonY-20; wy += 30 {
		for wx := int32(1075); wx < 1190; wx += 30 {
			c := winLit
			if rand.Float64() < 0.4 {
				c = winDim
			}
			if rand.Float64() < 0.2 {
				continue
			}
			s.FillRect(&sdl.Rect{X: wx, Y: wy, W: 16, H: 18}, c)
		}
	}

	// ----- Green park -----
	for y := horizonY; y < parkBottom; y++ {
		t := float64(y-horizonY) / float64(parkBottom-horizonY)
		r := uint8(30 + t*15)
		g := uint8(80 - t*25)
		b := uint8(25 + t*5)
		s.FillRect(&sdl.Rect{X: 0, Y: y, W: engine.ScreenWidth, H: 1}, sdl.MapRGBA(f, r, g, b, 255))
	}
	// Grass tufts
	for i := 0; i < 200; i++ {
		gx := int32(rand.Intn(int(engine.ScreenWidth)))
		gy := horizonY + int32(rand.Intn(int(parkBottom-horizonY)))
		gw := int32(rand.Intn(4) + 2)
		gh := int32(rand.Intn(3) + 1)
		brightness := uint8(rand.Intn(30))
		s.FillRect(&sdl.Rect{X: gx, Y: gy, W: gw, H: gh},
			sdl.MapRGBA(f, 25+brightness, 70+brightness, 20+brightness, 255))
	}

	// Trees
	drawTree(s, f, 250, horizonY-5, 45)
	drawTree(s, f, 560, horizonY+10, 40)
	drawTree(s, f, 850, horizonY, 38)
	drawTree(s, f, 1140, horizonY+5, 35)

	// ----- Cobblestone path -----
	pathBase := sdl.MapRGBA(f, 110, 95, 70, 255)
	s.FillRect(&sdl.Rect{X: 0, Y: parkBottom, W: engine.ScreenWidth, H: pathBottom - parkBottom}, pathBase)
	// Cobblestone pattern
	for y := parkBottom + 3; y < pathBottom-3; y += 12 {
		offset := int32(0)
		if ((y-parkBottom)/12)%2 == 1 {
			offset = 12
		}
		for x := offset; x < engine.ScreenWidth; x += 24 {
			shade := uint8(rand.Intn(25))
			s.FillRect(&sdl.Rect{X: x + 1, Y: y + 1, W: 21, H: 9},
				sdl.MapRGBA(f, 95+shade, 80+shade, 60+shade, 255))
		}
	}
	// Path edges
	s.FillRect(&sdl.Rect{X: 0, Y: parkBottom, W: engine.ScreenWidth, H: 3},
		sdl.MapRGBA(f, 80, 70, 50, 255))
	s.FillRect(&sdl.Rect{X: 0, Y: pathBottom - 2, W: engine.ScreenWidth, H: 3},
		sdl.MapRGBA(f, 80, 70, 50, 255))

	// ----- Dark street -----
	for y := pathBottom; y < engine.ScreenHeight; y++ {
		t := float64(y-pathBottom) / float64(engine.ScreenHeight-pathBottom)
		r := uint8(45 - t*15)
		g := uint8(40 - t*12)
		b := uint8(38 - t*10)
		s.FillRect(&sdl.Rect{X: 0, Y: y, W: engine.ScreenWidth, H: 1}, sdl.MapRGBA(f, r, g, b, 255))
	}

	// ----- Street lamp -----
	lampX := int32(130)
	lampTop := int32(250)
	lampBase := parkBottom + 60
	poleColor := sdl.MapRGBA(f, 35, 35, 40, 255)
	s.FillRect(&sdl.Rect{X: lampX - 2, Y: lampTop, W: 5, H: lampBase - lampTop}, poleColor)
	// Lamp arm
	s.FillRect(&sdl.Rect{X: lampX - 15, Y: lampTop, W: 30, H: 4}, poleColor)
	// Lamp housing
	s.FillRect(&sdl.Rect{X: lampX - 10, Y: lampTop - 20, W: 20, H: 22},
		sdl.MapRGBA(f, 40, 40, 45, 255))
	// Lamp glass (lit)
	s.FillRect(&sdl.Rect{X: lampX - 7, Y: lampTop - 17, W: 14, H: 14},
		sdl.MapRGBA(f, 255, 220, 150, 200))
	// Lamp base
	s.FillRect(&sdl.Rect{X: lampX - 8, Y: lampBase - 5, W: 16, H: 8}, poleColor)

	// ----- Park bench -----
	benchX := int32(400)
	benchY := parkBottom - 30
	benchColor := sdl.MapRGBA(f, 55, 40, 25, 255)
	s.FillRect(&sdl.Rect{X: benchX, Y: benchY, W: 60, H: 5}, benchColor)
	s.FillRect(&sdl.Rect{X: benchX, Y: benchY - 20, W: 60, H: 5}, benchColor)
	s.FillRect(&sdl.Rect{X: benchX + 3, Y: benchY - 20, W: 4, H: 25}, benchColor)
	s.FillRect(&sdl.Rect{X: benchX + 53, Y: benchY - 20, W: 4, H: 25}, benchColor)
	// Legs
	s.FillRect(&sdl.Rect{X: benchX + 5, Y: benchY + 5, W: 4, H: 15}, benchColor)
	s.FillRect(&sdl.Rect{X: benchX + 51, Y: benchY + 5, W: 4, H: 15}, benchColor)

	tex, err := renderer.CreateTextureFromSurface(s)
	if err != nil {
		panic(err)
	}
	return &background{tex: tex, srcW: engine.ScreenWidth, srcH: engine.ScreenHeight}
}

// ---------- Interior Background ----------

func newInteriorBackground(renderer *sdl.Renderer) *background {
	s, err := sdl.CreateRGBSurface(0, engine.ScreenWidth, engine.ScreenHeight, 32,
		0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
	if err != nil {
		panic(err)
	}
	defer s.Free()
	f := s.Format

	wall := sdl.MapRGBA(f, 120, 85, 60, 255)
	s.FillRect(nil, wall)

	floor := sdl.MapRGBA(f, 80, 55, 35, 255)
	s.FillRect(&sdl.Rect{X: 0, Y: 550, W: engine.ScreenWidth, H: 250}, floor)

	// Wood grain on walls
	for i := 0; i < 80; i++ {
		gx := int32(rand.Intn(int(engine.ScreenWidth)))
		gy := int32(rand.Intn(540))
		s.FillRect(&sdl.Rect{X: gx, Y: gy, W: int32(rand.Intn(40) + 10), H: 1},
			sdl.MapRGBA(f, 115, 80, 55, 255))
	}

	windowColor := sdl.MapRGBA(f, 150, 180, 220, 255)
	s.FillRect(&sdl.Rect{X: 500, Y: 100, W: 200, H: 250}, windowColor)
	windowFrame := sdl.MapRGBA(f, 60, 40, 25, 255)
	s.FillRect(&sdl.Rect{X: 500, Y: 100, W: 200, H: 8}, windowFrame)
	s.FillRect(&sdl.Rect{X: 500, Y: 342, W: 200, H: 8}, windowFrame)
	s.FillRect(&sdl.Rect{X: 500, Y: 100, W: 8, H: 250}, windowFrame)
	s.FillRect(&sdl.Rect{X: 692, Y: 100, W: 8, H: 250}, windowFrame)
	s.FillRect(&sdl.Rect{X: 596, Y: 100, W: 8, H: 250}, windowFrame)

	shelf := sdl.MapRGBA(f, 90, 60, 35, 255)
	s.FillRect(&sdl.Rect{X: 800, Y: 300, W: 250, H: 12}, shelf)
	s.FillRect(&sdl.Rect{X: 800, Y: 400, W: 250, H: 12}, shelf)

	baseboard := sdl.MapRGBA(f, 70, 45, 25, 255)
	s.FillRect(&sdl.Rect{X: 0, Y: 540, W: engine.ScreenWidth, H: 15}, baseboard)

	// Carpet on floor
	carpet := sdl.MapRGBA(f, 120, 40, 35, 255)
	s.FillRect(&sdl.Rect{X: 300, Y: 600, W: 600, H: 150}, carpet)
	carpetBorder := sdl.MapRGBA(f, 150, 55, 40, 255)
	s.FillRect(&sdl.Rect{X: 300, Y: 600, W: 600, H: 6}, carpetBorder)
	s.FillRect(&sdl.Rect{X: 300, Y: 744, W: 600, H: 6}, carpetBorder)
	s.FillRect(&sdl.Rect{X: 300, Y: 600, W: 6, H: 150}, carpetBorder)
	s.FillRect(&sdl.Rect{X: 894, Y: 600, W: 6, H: 150}, carpetBorder)

	tex, err := renderer.CreateTextureFromSurface(s)
	if err != nil {
		panic(err)
	}
	return &background{tex: tex, srcW: engine.ScreenWidth, srcH: engine.ScreenHeight}
}

func (b *background) draw(renderer *sdl.Renderer) {
	renderer.Copy(
		b.tex,
		&sdl.Rect{X: 0, Y: 0, W: b.srcW, H: b.srcH},
		&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight},
	)
}

// ---------- Drawing helpers ----------

func fillCircleSurface(surface *sdl.Surface, cx, cy, radius int32, color uint32) {
	for dy := -radius; dy <= radius; dy++ {
		halfW := int32(math.Sqrt(float64(radius*radius - dy*dy)))
		if halfW > 0 {
			surface.FillRect(&sdl.Rect{X: cx - halfW, Y: cy + dy, W: halfW * 2, H: 1}, color)
		}
	}
}

func fillTriangleSurface(surface *sdl.Surface, tipX, tipY, baseLeft, baseRight, baseY int32, color uint32) {
	height := baseY - tipY
	if height <= 0 {
		return
	}
	baseW := float64(baseRight - baseLeft)
	centerX := float64(tipX)
	for dy := int32(0); dy <= height; dy++ {
		t := float64(dy) / float64(height)
		w := int32(baseW * t)
		if w < 1 {
			w = 1
		}
		surface.FillRect(&sdl.Rect{X: int32(centerX) - w/2, Y: tipY + dy, W: w, H: 1}, color)
	}
}

func drawTree(surface *sdl.Surface, format *sdl.PixelFormat, x, y, size int32) {
	trunk := sdl.MapRGBA(format, 40, 30, 18, 255)
	tw := size / 6
	if tw < 3 {
		tw = 3
	}
	surface.FillRect(&sdl.Rect{X: x - tw/2, Y: y, W: tw, H: size * 2 / 3}, trunk)

	leaves := sdl.MapRGBA(format, 22, 55, 18, 255)
	leafR := size * 2 / 3
	fillCircleSurface(surface, x, y-leafR/2, leafR, leaves)

	lighter := sdl.MapRGBA(format, 30, 70, 25, 200)
	fillCircleSurface(surface, x-leafR/3, y-leafR/2-leafR/4, leafR*2/3, lighter)
}
