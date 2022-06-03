package screen

import (
	"github.com/gdamore/tcell/v2"
	"github.com/mattn/go-runewidth"
)

type Screen struct {
	Scr tcell.Screen // Screen to plot the stuff against

	MainVP *Viewport // Complete viewport of the screen with children
	MapVp  *Viewport // Map viewport
}

// Draw draws the main viewport (and consequtive child viewports)
func (s *Screen) Draw() {
	s.MainVP.Draw(s)
}

func (s *Screen) DrawText(x, y int, style tcell.Style, str string) {
	for _, c := range str {
		var comb []rune
		w := runewidth.RuneWidth(c)
		if w == 0 {
			comb = []rune{c}
			c = ' '
			w = 1
		}
		s.Scr.SetContent(x, y, c, comb, style)
		x += w
	}
}
