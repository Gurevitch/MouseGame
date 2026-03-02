package engine

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

func TextureFromBMP(renderer *sdl.Renderer, filename string) *sdl.Texture {
	img, err := sdl.LoadBMP(filename)
	if err != nil {
		panic(fmt.Errorf("loading BMP %s: %v", filename, err))
	}
	defer img.Free()

	key := GetPixelColor(img, 0, 0)
	img.SetColorKey(true, key)

	tex, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		panic(fmt.Errorf("creating texture from %s: %v", filename, err))
	}
	tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	return tex
}

func TextureFromBMPRaw(renderer *sdl.Renderer, filename string) *sdl.Texture {
	img, err := sdl.LoadBMP(filename)
	if err != nil {
		panic(fmt.Errorf("loading BMP %s: %v", filename, err))
	}
	defer img.Free()

	tex, err := renderer.CreateTextureFromSurface(img)
	if err != nil {
		panic(fmt.Errorf("creating texture from %s: %v", filename, err))
	}
	return tex
}

func GetPixelColor(s *sdl.Surface, x, y int32) uint32 {
	bpp := int(s.Format.BytesPerPixel)
	px := s.Pixels()
	off := int(y)*int(s.Pitch) + int(x)*bpp
	if off+bpp > len(px) {
		return 0
	}
	switch bpp {
	case 1:
		return uint32(px[off])
	case 2:
		return uint32(px[off]) | uint32(px[off+1])<<8
	case 3:
		return uint32(px[off]) | uint32(px[off+1])<<8 | uint32(px[off+2])<<16
	case 4:
		return uint32(px[off]) | uint32(px[off+1])<<8 | uint32(px[off+2])<<16 | uint32(px[off+3])<<24
	}
	return 0
}
