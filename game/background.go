package game

import (
	"math"
	"math/rand"

	"bitbucket.org/Local/games/PP/engine"
	"github.com/veandco/go-sdl2/sdl"
)

type bgLayer struct {
	tex      *sdl.Texture
	srcW     int32
	srcH     int32
	parallax float64
}

type background struct {
	tex    *sdl.Texture
	srcW   int32
	srcH   int32
	layers []bgLayer
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
	tex, w, h := engine.TextureFromPNGRawClean(renderer, path)
	return &background{tex: tex, srcW: w, srcH: h}
}

func (b *background) draw(renderer *sdl.Renderer, playerX float64) {
	if len(b.layers) == 0 {
		renderer.Copy(
			b.tex,
			&sdl.Rect{X: 0, Y: 0, W: b.srcW, H: b.srcH},
			&sdl.Rect{X: 0, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight},
		)
		return
	}
	screenCenter := float64(engine.ScreenWidth) / 2.0
	for _, l := range b.layers {
		offsetX := int32((playerX - screenCenter) * l.parallax)
		renderer.Copy(l.tex, nil,
			&sdl.Rect{X: -offsetX, Y: 0, W: engine.ScreenWidth, H: engine.ScreenHeight})
	}
}

// ---------- London Night Background (3-layer parallax) ----------

func makeSurface() (*sdl.Surface, *sdl.PixelFormat) {
	s, err := sdl.CreateRGBSurface(0, engine.ScreenWidth, engine.ScreenHeight, 32,
		0x000000FF, 0x0000FF00, 0x00FF0000, 0xFF000000)
	if err != nil {
		panic(err)
	}
	return s, s.Format
}

func surfaceToTexture(renderer *sdl.Renderer, s *sdl.Surface) *sdl.Texture {
	tex, err := renderer.CreateTextureFromSurface(s)
	if err != nil {
		panic(err)
	}
	tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	return tex
}

func newLondonBackground(renderer *sdl.Renderer) *background {
	horizonY := int32(440)
	parkBottom := int32(610)
	pathBottom := int32(720)

	skyTex := buildSkyLayer(renderer, horizonY)
	cityTex := buildCityLayer(renderer, horizonY)
	groundTex := buildGroundLayer(renderer, horizonY, parkBottom, pathBottom)

	return &background{
		layers: []bgLayer{
			{tex: skyTex, srcW: engine.ScreenWidth, srcH: engine.ScreenHeight, parallax: 0.02},
			{tex: cityTex, srcW: engine.ScreenWidth, srcH: engine.ScreenHeight, parallax: 0.05},
			{tex: groundTex, srcW: engine.ScreenWidth, srcH: engine.ScreenHeight, parallax: 0.0},
		},
	}
}

func buildSkyLayer(renderer *sdl.Renderer, horizonY int32) *sdl.Texture {
	s, f := makeSurface()
	defer s.Free()

	// Daytime sky gradient (bright blue)
	for y := int32(0); y < horizonY; y++ {
		t := float64(y) / float64(horizonY)
		r := uint8(100 + t*60)
		g := uint8(160 + t*40)
		b := uint8(235 - t*30)
		s.FillRect(&sdl.Rect{X: 0, Y: y, W: engine.ScreenWidth, H: 1}, sdl.MapRGBA(f, r, g, b, 255))
	}

	// Horizon haze
	for y := horizonY - 40; y < horizonY+10; y++ {
		dist := math.Abs(float64(y) - float64(horizonY))
		a := uint8(math.Max(0, 30-dist*0.8))
		s.FillRect(&sdl.Rect{X: 0, Y: y, W: engine.ScreenWidth, H: 1}, sdl.MapRGBA(f, 220, 210, 190, a))
	}

	// Sun
	sunX, sunY := int32(1300), int32(80)
	for dr := int32(80); dr > 20; dr-- {
		a := uint8(math.Min(255, float64(3)*float64(80-dr)/60.0))
		fillCircleSurface(s, sunX, sunY, dr, sdl.MapRGBA(f, 255, 245, 200, a))
	}
	fillCircleSurface(s, sunX, sunY, 20, sdl.MapRGBA(f, 255, 250, 220, 255))

	// Fluffy clouds
	drawCloud(s, f, 200, 100, 120, 35)
	drawCloud(s, f, 550, 60, 100, 30)
	drawCloud(s, f, 900, 130, 90, 25)
	drawCloud(s, f, 1200, 80, 110, 30)

	return surfaceToTexture(renderer, s)
}

func drawCloud(surface *sdl.Surface, f *sdl.PixelFormat, cx, cy, w, h int32) {
	col := sdl.MapRGBA(f, 255, 255, 255, 180)
	colLight := sdl.MapRGBA(f, 255, 255, 255, 120)
	fillEllipseSurface(surface, cx, cy, w/2, h/2, col)
	fillEllipseSurface(surface, cx-w/3, cy+h/6, w/3, h/3, col)
	fillEllipseSurface(surface, cx+w/3, cy+h/6, w/3, h/3, col)
	fillEllipseSurface(surface, cx-w/5, cy-h/4, w/3, h/3, colLight)
	fillEllipseSurface(surface, cx+w/6, cy-h/5, w/4, h/3, colLight)
}

func buildCityLayer(renderer *sdl.Renderer, horizonY int32) *sdl.Texture {
	s, f := makeSurface()
	defer s.Free()

	// Daytime building colors (muted blue-gray silhouettes)
	bFar := sdl.MapRGBA(f, 120, 130, 150, 255)
	bMid := sdl.MapRGBA(f, 100, 110, 130, 255)
	bNear := sdl.MapRGBA(f, 80, 85, 100, 255)
	winGlass := sdl.MapRGBA(f, 160, 190, 220, 255)
	winDark := sdl.MapRGBA(f, 70, 75, 90, 255)

	// Big Ben tower
	s.FillRect(&sdl.Rect{X: 200, Y: 140, W: 45, H: horizonY - 140}, bFar)
	fillTriangleSurface(s, 222, 80, 195, 250, 140, bFar)
	fillCircleSurface(s, 222, 180, 16, sdl.MapRGBA(f, 160, 155, 140, 255))
	fillCircleSurface(s, 222, 180, 14, sdl.MapRGBA(f, 230, 220, 190, 200))

	// Parliament building
	s.FillRect(&sdl.Rect{X: 60, Y: 330, W: 300, H: horizonY - 330}, bFar)
	for x := int32(60); x < 360; x += 22 {
		s.FillRect(&sdl.Rect{X: x, Y: 320, W: 14, H: 10}, bFar)
	}
	s.FillRect(&sdl.Rect{X: 80, Y: 290, W: 20, H: 40}, bFar)
	fillTriangleSurface(s, 90, 270, 78, 102, 290, bFar)
	s.FillRect(&sdl.Rect{X: 330, Y: 300, W: 18, H: 30}, bFar)
	fillTriangleSurface(s, 339, 282, 328, 350, 300, bFar)

	// Center buildings
	s.FillRect(&sdl.Rect{X: 420, Y: 350, W: 140, H: horizonY - 350}, bMid)
	s.FillRect(&sdl.Rect{X: 570, Y: 320, W: 90, H: horizonY - 320}, bMid)
	s.FillRect(&sdl.Rect{X: 500, Y: 310, W: 60, H: horizonY - 310}, bFar)
	fillTriangleSurface(s, 530, 285, 500, 560, 310, bFar)

	// Westminster Abbey
	s.FillRect(&sdl.Rect{X: 780, Y: 220, W: 300, H: horizonY - 220}, bFar)
	s.FillRect(&sdl.Rect{X: 920, Y: 120, W: 24, H: 100}, bFar)
	fillTriangleSurface(s, 932, 70, 915, 950, 120, bFar)
	s.FillRect(&sdl.Rect{X: 810, Y: 170, W: 18, H: 50}, bFar)
	fillTriangleSurface(s, 819, 145, 807, 831, 170, bFar)
	s.FillRect(&sdl.Rect{X: 1040, Y: 175, W: 18, H: 45}, bFar)
	fillTriangleSurface(s, 1049, 150, 1037, 1061, 175, bFar)
	for x := int32(800); x < 1070; x += 38 {
		fillTriangleSurface(s, x+12, 225, x, x+24, 255, sdl.MapRGBA(f, 105, 115, 135, 255))
	}

	// Far right buildings
	s.FillRect(&sdl.Rect{X: 1100, Y: 340, W: 150, H: horizonY - 340}, bMid)
	s.FillRect(&sdl.Rect{X: 1120, Y: 310, W: 60, H: 30}, bMid)
	s.FillRect(&sdl.Rect{X: 1270, Y: 360, W: 130, H: horizonY - 360}, bMid)

	// Near foreground buildings with windows
	s.FillRect(&sdl.Rect{X: 0, Y: 390, W: 120, H: horizonY - 390}, bNear)
	for wy := int32(405); wy < horizonY-20; wy += 30 {
		for wx := int32(14); wx < 110; wx += 28 {
			c := winGlass
			if rand.Float64() < 0.3 {
				c = winDark
			}
			if rand.Float64() < 0.15 {
				continue
			}
			s.FillRect(&sdl.Rect{X: wx, Y: wy, W: 14, H: 18}, c)
		}
	}

	s.FillRect(&sdl.Rect{X: 1250, Y: 400, W: 150, H: horizonY - 400}, bNear)
	for wy := int32(415); wy < horizonY-20; wy += 30 {
		for wx := int32(1265); wx < 1390; wx += 30 {
			c := winGlass
			if rand.Float64() < 0.4 {
				c = winDark
			}
			if rand.Float64() < 0.2 {
				continue
			}
			s.FillRect(&sdl.Rect{X: wx, Y: wy, W: 16, H: 18}, c)
		}
	}

	return surfaceToTexture(renderer, s)
}

func buildGroundLayer(renderer *sdl.Renderer, horizonY, parkBottom, pathBottom int32) *sdl.Texture {
	s, f := makeSurface()
	defer s.Free()

	// Bright green park
	for y := horizonY; y < parkBottom; y++ {
		t := float64(y-horizonY) / float64(parkBottom-horizonY)
		r := uint8(60 + t*20)
		g := uint8(140 - t*30)
		b := uint8(40 + t*10)
		s.FillRect(&sdl.Rect{X: 0, Y: y, W: engine.ScreenWidth, H: 1}, sdl.MapRGBA(f, r, g, b, 255))
	}
	// Grass tufts
	for i := 0; i < 300; i++ {
		gx := int32(rand.Intn(int(engine.ScreenWidth)))
		gy := horizonY + int32(rand.Intn(int(parkBottom-horizonY)))
		gw := int32(rand.Intn(4) + 2)
		gh := int32(rand.Intn(3) + 1)
		brightness := uint8(rand.Intn(40))
		s.FillRect(&sdl.Rect{X: gx, Y: gy, W: gw, H: gh},
			sdl.MapRGBA(f, 50+brightness, 120+brightness, 35+brightness, 255))
	}

	// Trees
	drawTree(s, f, 280, horizonY-5, 48)
	drawTree(s, f, 580, horizonY+10, 42)
	drawTree(s, f, 900, horizonY, 40)
	drawTree(s, f, 1200, horizonY+5, 38)

	// Cobblestone path
	pathBase := sdl.MapRGBA(f, 140, 125, 100, 255)
	s.FillRect(&sdl.Rect{X: 0, Y: parkBottom, W: engine.ScreenWidth, H: pathBottom - parkBottom}, pathBase)

	for y := parkBottom + 3; y < pathBottom-3; y += 12 {
		offset := int32(0)
		if ((y-parkBottom)/12)%2 == 1 {
			offset = 12
		}
		distFromTop := float64(y-parkBottom) / float64(pathBottom-parkBottom)
		edgeDarken := uint8(distFromTop * 10)
		for x := offset; x < engine.ScreenWidth; x += 24 {
			shade := uint8(rand.Intn(35))
			distFromEdge := float64(x) / float64(engine.ScreenWidth)
			if distFromEdge > 0.5 {
				distFromEdge = 1.0 - distFromEdge
			}
			edgeFade := uint8(math.Max(0, (0.5-distFromEdge)*20))
			r := uint8(math.Max(0, float64(120+shade-edgeDarken-edgeFade)))
			g := uint8(math.Max(0, float64(105+shade-edgeDarken-edgeFade)))
			b := uint8(math.Max(0, float64(80+shade-edgeDarken-edgeFade)))

			if rand.Float64() < 0.04 {
				r = uint8(math.Max(0, float64(r)-20))
				g = uint8(math.Max(0, float64(g)-18))
				b = uint8(math.Max(0, float64(b)-15))
			}
			s.FillRect(&sdl.Rect{X: x + 1, Y: y + 1, W: 21, H: 9},
				sdl.MapRGBA(f, r, g, b, 255))
		}
	}

	// Path edges
	s.FillRect(&sdl.Rect{X: 0, Y: parkBottom, W: engine.ScreenWidth, H: 3},
		sdl.MapRGBA(f, 100, 90, 70, 255))
	s.FillRect(&sdl.Rect{X: 0, Y: pathBottom - 2, W: engine.ScreenWidth, H: 3},
		sdl.MapRGBA(f, 100, 90, 70, 255))

	// Pavement / street
	for y := pathBottom; y < engine.ScreenHeight; y++ {
		t := float64(y-pathBottom) / float64(engine.ScreenHeight-pathBottom)
		r := uint8(130 - t*25)
		g := uint8(125 - t*22)
		b := uint8(115 - t*20)
		s.FillRect(&sdl.Rect{X: 0, Y: y, W: engine.ScreenWidth, H: 1}, sdl.MapRGBA(f, r, g, b, 255))
	}

	// Prop shadows
	treeShadow := sdl.MapRGBA(f, 30, 50, 25, 60)
	fillEllipseSurface(s, 280, horizonY+40, 35, 9, treeShadow)
	fillEllipseSurface(s, 580, horizonY+50, 30, 8, treeShadow)
	fillEllipseSurface(s, 900, horizonY+40, 28, 8, treeShadow)
	fillEllipseSurface(s, 1200, horizonY+45, 26, 7, treeShadow)

	// Street lamp (unlit in daytime)
	lampX := int32(170)
	lampTop := int32(320)
	lampBase := parkBottom + 40
	poleColor := sdl.MapRGBA(f, 50, 50, 55, 255)
	s.FillRect(&sdl.Rect{X: lampX - 3, Y: lampTop, W: 6, H: lampBase - lampTop}, poleColor)
	s.FillRect(&sdl.Rect{X: lampX - 18, Y: lampTop, W: 36, H: 5}, poleColor)
	s.FillRect(&sdl.Rect{X: lampX - 12, Y: lampTop - 24, W: 24, H: 26},
		sdl.MapRGBA(f, 55, 55, 60, 255))
	s.FillRect(&sdl.Rect{X: lampX - 9, Y: lampTop - 20, W: 18, H: 16},
		sdl.MapRGBA(f, 200, 210, 210, 200))
	s.FillRect(&sdl.Rect{X: lampX - 10, Y: lampBase - 6, W: 20, H: 10}, poleColor)
	fillEllipseSurface(s, lampX, lampBase+5, 14, 5, sdl.MapRGBA(f, 30, 30, 30, 60))

	// Park bench
	benchX := int32(500)
	benchY := parkBottom - 45
	benchColor := sdl.MapRGBA(f, 70, 50, 30, 255)
	benchW := int32(100)
	fillEllipseSurface(s, benchX+benchW/2, benchY+35, benchW/2+5, 6, sdl.MapRGBA(f, 30, 40, 25, 50))
	s.FillRect(&sdl.Rect{X: benchX, Y: benchY, W: benchW, H: 6}, benchColor)
	s.FillRect(&sdl.Rect{X: benchX, Y: benchY - 25, W: benchW, H: 6}, benchColor)
	s.FillRect(&sdl.Rect{X: benchX + 3, Y: benchY - 25, W: 5, H: 31}, benchColor)
	s.FillRect(&sdl.Rect{X: benchX + benchW - 8, Y: benchY - 25, W: 5, H: 31}, benchColor)
	s.FillRect(&sdl.Rect{X: benchX + 6, Y: benchY + 6, W: 5, H: 18}, benchColor)
	s.FillRect(&sdl.Rect{X: benchX + benchW - 11, Y: benchY + 6, W: 5, H: 18}, benchColor)

	return surfaceToTexture(renderer, s)
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

// ---------- Drawing helpers ----------

func fillEllipseSurface(surface *sdl.Surface, cx, cy, rx, ry int32, color uint32) {
	for dy := -ry; dy <= ry; dy++ {
		halfW := int32(float64(rx) * math.Sqrt(1.0-float64(dy*dy)/float64(ry*ry)))
		if halfW > 0 {
			surface.FillRect(&sdl.Rect{X: cx - halfW, Y: cy + dy, W: halfW * 2, H: 1}, color)
		}
	}
}

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
