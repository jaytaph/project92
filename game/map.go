package game

import (
	"errors"
	"math/rand"
	"sync"

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
	S tcell.Style // Style (color) of the element
	C rune        // actual character
	X int
	Y int
	//  here will be other stuff about the actual element. Like a player, enemy, building, flag etc
}

type Player struct {
	X int
	Y int
}

type GameMap struct {
	mu sync.Mutex      // Mutex so we can make sure only one thing can edit the map
	m  [][]TerrainItem // Map of the terrain
	H  int             // Height of the map
	W  int             // Width of the map

	MapXOff int
	MapYOff int

	P Player
}

// New will create a new map based on width and height. It will be empty
func New(w, h int) *GameMap {
	gameMap := GameMap{
		mu: sync.Mutex{},
		H:  h,
		W:  w,
	}

	// Initialize multidimensional array
	gameMap.m = make([][]TerrainItem, h)
	for i := 0; i != h; i++ {
		gameMap.m[i] = make([]TerrainItem, w)
	}

	return &gameMap
}

// Regenerate generates a new gamemap terrain based on the params given
func (gm *GameMap) Regenerate(a float64, b float64, n int32, seed int64) {
	p := perlin.NewPerlinRandSource(a, b, n, rand.NewSource(seed))

	for x := 0; x != gm.W; x++ {
		for y := 0; y != gm.H; y++ {

			// Get perlin noise for x/y coordinate
			f := int(p.Noise2D(float64(x), float64(y)))

			// Cap between 0 and len(terrain colors)
			if f < 0 {
				f = 0
			}
			if f >= len(terrainColors)-1 {
				f = len(terrainColors) - 1
			}

			// Display empty char
			c := ' '

			// Set color
			st := tcell.StyleDefault.Background(terrainColors[f] | tcell.ColorValid).Foreground(tcell.ColorGreen)

			gm.m[x][y] = TerrainItem{
				S: st,
				C: c,
				X: x,
				Y: y,
			}
		}
	}
}

// Draw will draw the gamemap onto the given screen. xOff and yOff are the top left corner of the map to display
func (gm *GameMap) Draw(s tcell.Screen) {
	w, h := s.Size()

	// This is the viewport which we draw
	viewportHeight := h - 10
	viewportWidth := w - 2
	viewportX := 1
	viewportY := 1

	emptyStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	for y := 0; y != viewportHeight; y++ {
		for x := 0; x != viewportWidth; x++ {
			// find screenX and screenY
			// find mapX and mapY

			screenX := viewportX + x
			screenY := viewportY + y
			mapX := gm.MapXOff + x
			mapY := gm.MapYOff + y

			var c rune
			var st tcell.Style
			if mapX >= gm.W || mapY >= gm.H {
				// outside boundaries
				c = ' '
				st = emptyStyle
			} else {
				c = gm.m[mapX][mapY].C
				st = gm.m[mapX][mapY].S
			}

			s.SetContent(screenX, screenY, c, nil, st)

		}
	}
}

// GetTile will fetch a specific x/y coordinate. It's ok to be out of the gamespace
func (gm *GameMap) GetTile(x, y int) (*TerrainItem, error) {
	if x < 0 || y < 0 {
		return nil, errors.New("out of bounds")
	}
	if x >= gm.W || y >= gm.H {
		return nil, errors.New("out of bounds")
	}

	return &TerrainItem{
		S: gm.m[x][y].S,
		C: gm.m[x][y].C,
		X: gm.m[x][y].X,
		Y: gm.m[x][y].Y,
	}, nil
}

// SetTile will set a specific x/y coordinate to the element. It's ok to be out of the gamespace
func (gm *GameMap) SetTile(x, y int, s tcell.Style, c rune) {
	if x < 0 || y < 0 {
		return
	}
	if x >= gm.W || y >= gm.H {
		return
	}

	gm.mu.Lock() // Needed?
	gm.m[x][y].C = c
	gm.m[x][y].S = s
	gm.m[x][y].X = x
	gm.m[x][y].Y = y
	gm.mu.Unlock() // Needed?
}
