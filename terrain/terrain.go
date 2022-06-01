package terrain

import (
	"fmt"
	"math/rand"

	"github.com/aquilax/go-perlin"
	"github.com/gdamore/tcell/v2"
)

// Set of colors for the terrain. If the perlin generator goes below 0 or after the length of this
// slice, it will be capped to those first or last color
var terrainColors = []tcell.Color{
	tcell.Color226,
	tcell.Color226,
	tcell.Color226,
	tcell.Color226,
	tcell.Color226,
	tcell.Color226,
	tcell.Color178,
	tcell.Color178,
	tcell.Color178,
	tcell.Color178,
	tcell.Color178,
	tcell.Color184,
	tcell.Color184,
	tcell.Color184,
	tcell.Color184,
	tcell.Color184,
	tcell.Color185,
	tcell.Color185,
	tcell.Color185,
	tcell.Color186,
	tcell.Color190,
	tcell.Color191,
	tcell.Color192,
	tcell.Color193,
	tcell.Color22,
	tcell.Color28,
	tcell.Color34,
	tcell.Color40,
	tcell.Color41,
	tcell.Color46,
	tcell.Color47,
	tcell.Color21,
	tcell.Color20,
	tcell.Color19,
	tcell.Color18,
}

type TerrainItem struct {
	s tcell.Style // Style (color) of the element
	c rune        // actual character
	//  here will be other stuff about the actual character
}

type GameMap struct {
	m [][]TerrainItem // Map of the terrain
	h int             // Height of the map
	w int             // Widht of the map
}

// New will create a new map based on width and height. It will be empty
func New(w, h int) *GameMap {
	gameMap := GameMap{
		h: h,
		w: w,
	}

	gameMap.m = make([][]TerrainItem, h)
	for i := 0; i != h; i++ {
		gameMap.m[i] = make([]TerrainItem, w)
	}

	return &gameMap
}

// Regenerate generates a new gamemap terrain based on the params given
func (gm *GameMap) Regenerate(a float64, b float64, n int32, seed int64) {
	p := perlin.NewPerlinRandSource(a, b, n, rand.NewSource(seed))

	for x := 0; x != gm.w; x++ {
		for y := 0; y != gm.h; y++ {
			c := '.'
			if y == 0 || x == 0 || y == gm.h-1 || x == gm.w-1 {
				c = 'X'
			}

			c = '.'

			f := int(p.Noise2D(float64(x), float64(y)))
			if f < 0 {
				f = 0
			}
			if f >= len(terrainColors)-1 {
				f = len(terrainColors) - 1
			}
			st := tcell.StyleDefault.Background(terrainColors[f] | tcell.ColorValid).Foreground(tcell.ColorGreen)
			// st := tcell.StyleDefault.Background(tcell.Color(p.Noise2D(float64(x), float64(y)) + 96) | tcell.ColorValid).Foreground(tcell.ColorWhite)

			gm.m[x][y] = TerrainItem{
				s: st,
				c: c,
			}
		}
	}
}

// Draw will draw the gamemap onto the given screen. xOff and yOff are the top left corner of the map to display
func (gm *GameMap) Draw(s tcell.Screen, xOff int, yOff int) {
	w, h := s.Size()

	// Leave a bit of room for game info
	h -= 10

	// Make sure we don't display outside of the map
	if w+xOff > gm.w {
		w = gm.w - xOff
	}
	if h+yOff > gm.h {
		h = gm.h - yOff
	}

	// Draw coordinates
	white := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	for y := 0; y != h; y++ {
		f := fmt.Sprintf("%02X", y+yOff)
		s.SetContent(0, y+2, rune(f[0]), nil, white)
		s.SetContent(1, y+2, rune(f[1]), nil, white)
		s.SetContent(2, y+2, ' ', nil, white)

		for x := 0; x != w; x++ {
			f := fmt.Sprintf("%02X", x+xOff)
			s.SetContent(x+3, 0, rune(f[0]), nil, white)
			s.SetContent(x+3, 1, rune(f[1]), nil, white)

			s.SetContent(x+3, y+2, gm.m[x+xOff][y+yOff].c, nil, gm.m[x+xOff][y+yOff].s)
		}
	}

	// draw actual map
	for y := 0; y != h; y++ {
		for x := 0; x != w; x++ {
			s.SetContent(x+3, y+2, gm.m[x+xOff][y+yOff].c, nil, gm.m[x+xOff][y+yOff].s)
		}
	}

}
