package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jaytaph/project92/game"
)

const (
	maxX = 256
	maxY = 256
)

func main() {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	// Initialize screen
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()

	// Create new terrain map
	rand.Seed(time.Now().UnixNano())
	gm := game.New(maxX, maxY)
	gm.Regenerate(.1, .1, 3, rand.Int63())

	// Start displaying map at top left corner
	xOff := 0
	yOff := 0

	// Event loop
	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	// Display map

	for {
		// Poll event
		if !s.HasPendingEvent() {
			refresh(gm, s, xOff, yOff)
			time.Sleep(28 * time.Millisecond)
			continue
		}

		ev := s.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			// In case the terminal resized
			refresh(gm, s, xOff, yOff)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyRune {
				switch ev.Rune() {
				case 'r':
					// 'r' regenerates a new game map
					gm.Regenerate(.1, .1, 3, rand.Int63())
					// refresh(gm, s, xOff, yOff)
				case 'p':
					game.Ping(gm, 100, 20, 15)
					// refresh(gm, s, xOff, yOff)
				}
			}
			if ev.Key() == tcell.KeyUp {
				yOff--
				if yOff < 0 {
					yOff = 0
				}
				// refresh(gm, s, xOff, yOff)
			}
			if ev.Key() == tcell.KeyDown {
				yOff++
				if yOff > maxY {
					yOff = maxY
				}
				// refresh(gm, s, xOff, yOff)
			}
			if ev.Key() == tcell.KeyLeft {
				xOff--
				if xOff < 0 {
					xOff = 0
				}
				// refresh(gm, s, xOff, yOff)
			}
			if ev.Key() == tcell.KeyRight {
				xOff++
				if xOff > maxX {
					xOff = maxX
				}
				// refresh(gm, s, xOff, yOff)
			}

			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			}
		}
	}
}

func refresh(gm *game.GameMap, s tcell.Screen, xOff, yOff int) {
	// s.Clear()
	gm.Draw(s, xOff, yOff)
	s.Show()
}
