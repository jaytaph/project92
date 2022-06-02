package game

import (
	"math"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jaytaph/project92/terrain"
)

// Ping will trigger a ping sourced from x/y with a strength s
func Ping(gm *terrain.GameMap, x, y, s int) {

	// Ping inside a separate go routine
	go func(gm *terrain.GameMap, sourceX, sourceY, s int) {

		// Items will store the element below a ping asterisk. We should figure out a better way to deal with this
		items := make([]*terrain.TerrainItem, 0, 360)

		// Strength is basically the max radius of the ping
		for i := 0; i != s; i++ {

			// Remove old pings
			for _, item := range items {
				gm.SetTile(item.X, item.Y, item.S, item.C)
			}
			items = items[:0]

			// Number of "points", basically, on each of these angles
			for a := 0; a != 360; a += 18 {
				x = sourceX + int(math.Cos(float64(a)*math.Pi/180)*float64(i))
				y = sourceY + int(math.Sin(float64(a)*math.Pi/180)*float64(i))

				// Get tile on position, store it when it's not a '*' ping
				it, _ := gm.GetTile(x, y)
				if it.C != '*' {
					items = append(items, it)
				}

				gm.SetTile(x, y, tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite), '*')
			}

			time.Sleep(50 * time.Millisecond)
		}

		// Remove last ping elements
		for _, item := range items {
			gm.SetTile(item.X, item.Y, item.S, item.C)
		}

	}(gm, x, y, s)
}
