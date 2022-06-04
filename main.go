package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/jaytaph/project92/game"
	"github.com/jaytaph/project92/screen"
	"github.com/mattn/go-runewidth"
)

const (
	maxX = 256
	maxY = 256
)

type MoveMode int

const (
	MoveModePlayer = iota
	MoveModeMap
	MoveModeMenu
	moveModeLen
)

var moveMode MoveMode = MoveModePlayer

func drawStr(s tcell.Screen, x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.SetContent(x, y, c, comb, style)
		x += w
	}
}

var mainScreen *screen.Screen

func main() {
	scr := initPhysicalScreen()
	initGameScreens(scr)
	gm := initGame(maxX, maxY)

	// Event loop
	quit := func() {
		scr.Fini()
		os.Exit(0)
	}

	for {
		// Refresh the screen every 28ms
		if !scr.HasPendingEvent() {
			refresh(gm, scr)
			time.Sleep(28 * time.Millisecond)
			continue
		}

		// Poll event
		ev := scr.PollEvent()

		// Process event
		switch ev := ev.(type) {
		case *tcell.EventResize:
			// In case the terminal resized
			refresh(gm, scr)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyRune {
				switch ev.Rune() {
				case 'r':
					// 'r' regenerates a new game map
					gm.Regenerate(.1, .1, 3, rand.Int63())
				case ' ':
					game.Ping(gm, gm.P.X, gm.P.Y, 15)
				}
			}
			if ev.Key() == tcell.KeyTab {
				moveMode++
				moveMode %= moveModeLen
			}

			if ev.Key() == tcell.KeyUp {
				move(gm, moveMode, 0, -1)
			}
			if ev.Key() == tcell.KeyDown {
				move(gm, moveMode, 0, 1)
			}
			if ev.Key() == tcell.KeyLeft {
				move(gm, moveMode, -1, 0)
			}
			if ev.Key() == tcell.KeyRight {
				move(gm, moveMode, 1, 0)
			}

			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			}
		}
	}
}

func initGame(maxX, maxY int) *game.GameMap {
	// Create new terrain map
	rand.Seed(time.Now().UnixNano())
	gm := game.New(maxX, maxY)
	gm.Regenerate(.1, .1, 3, rand.Int63())

	return gm
}

func initPhysicalScreen() tcell.Screen {
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)

	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	s.SetStyle(defStyle)
	s.Clear()

	return s
}

func initGameScreens(scr tcell.Screen) {
	// Generate
	w, h := scr.Size()
	vp, err := screen.New(0, 0, w, h, nil)
	if err != nil {
		log.Fatal("%+v", err)
	}
	vp.SetBordered(true)
	vp.SetActive(false)
	vp.SetTitle("Project92")

	mapVp, _ := screen.New(1, 1, vp.W-2, vp.H-10, vp)
	mapVp.SetBordered(true)
	mapVp.SetActive(true)
	mapVp.SetTitle("Map")

	menuVp, _ := screen.New(1, vp.H-9, vp.W-2, 9, vp)
	menuVp.SetBordered(true)
	menuVp.SetTitle("Menu")

	mainScreen = &screen.Screen{
		Scr:    scr,
		MainVP: vp,
		MapVp:  mapVp,
	}

}

func move(gm *game.GameMap, mode MoveMode, x, y int) {
	switch mode {
	case MoveModePlayer:
		gm.P.X = ModLimit(gm.P.X, x, 0, gm.W)
		gm.P.Y = ModLimit(gm.P.Y, y, 0, gm.H)
	case MoveModeMap:
		gm.MapXOff = ModLimit(gm.MapXOff, x, 0, gm.W)
		gm.MapYOff = ModLimit(gm.MapYOff, y, 0, gm.H)
	}
}

// ModLimit increases or decreases v and makes sure that v stays between min and max-1
func ModLimit(v int, delta int, min int, max int) int {
	v += delta

	if v < min {
		v = min
	}
	if v > max-1 {
		v = max - 1
	}

	return v
}

func refresh(gm *game.GameMap, s tcell.Screen) {
	gm.Draw(s)

	// Draw player
	st := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGreen)
	s.SetContent(gm.P.X+3, gm.P.Y+2, 'P', nil, st)

	white := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)
	red := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)

	drawStr(s, 0, 0, white, "[      ]")
	switch moveMode {
	case MoveModePlayer:
		drawStr(s, 1, 0, red, "PLAYER")
	case MoveModeMap:
		drawStr(s, 1, 0, red, "  MAP ")
	case MoveModeMenu:
		drawStr(s, 1, 0, red, " MENU ")
	}

	mainScreen.Draw()

	s.Show()
}
