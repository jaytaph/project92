package screen

import (
	"errors"
	"fmt"

	"github.com/gdamore/tcell/v2"
)

/*

	Viewports are boxes on a given position (x/y) inside
	another box. 0/0 in the viewport is the top left corner

	+------------------------------------------------+
	| viewport screen (0, 0, screen h/w)             |
	|    +-------------------+                       |
	|    | viewport map      |                       |
	|    |                   |                       |
	|    +-------------------+                       |
	|    +------------------------------------------+|
	|    | viewport menu                            ||
	|    +------------------------------------------+|
	|                                                |
	+------------------------------------------------+

*/

var (
	errOutOfViewportBounds = errors.New("coordinate out of bounds of viewport")
)

var (
	singleBorderRunes = [8]rune{tcell.RuneULCorner, tcell.RuneHLine, tcell.RuneURCorner, tcell.RuneVLine, tcell.RuneLRCorner, tcell.RuneHLine, tcell.RuneLLCorner, tcell.RuneVLine}
	doubleBorderRunes = [8]rune{tcell.RuneULCorner, tcell.RuneHLine, tcell.RuneURCorner, tcell.RuneVLine, tcell.RuneLRCorner, tcell.RuneHLine, tcell.RuneLLCorner, tcell.RuneVLine}
)

type DrawContentFunc func(vp *Viewport, s *Screen)

// Viewport defines a box on a position (in parent viewport)
type Viewport struct {
	H int // Height of the viewport
	W int // Width of the viewport
	X int // X offset
	Y int // Y offset

	title string

	bordered bool // When true, a border will be displayed
	active   bool // when active, the border will be different color?

	content *DrawContentFunc // Function to call to draw content in the viewport

	parent   *Viewport   // Parent viewport, to base the actual screen X/Y positions on. When nil, we use 0/0 as screen offsets
	children []*Viewport // Child viewports
}

func New(x, y, w, h int, parent *Viewport) (*Viewport, error) {
	vp := &Viewport{
		H:        h,
		W:        w,
		X:        x,
		Y:        y,
		title:    "",
		bordered: false,
		active:   false,
		content:  nil,
		parent:   parent,
		children: make([]*Viewport, 0),
	}

	if parent != nil {
		parent.AddChild(vp)
	}

	return vp, nil
}

func (vp *Viewport) SetBordered(bordered bool) {
	vp.bordered = bordered
}
func (vp *Viewport) SetActive(active bool) {
	vp.active = active
}
func (vp *Viewport) AddChild(child *Viewport) {
	vp.children = append(vp.children, child)
}

func (vp *Viewport) Draw(s *Screen) {

	if vp.bordered {
		st := tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite)

		borders := singleBorderRunes
		if vp.active {
			st = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorRed)
		}

		var sx, sy int

		for x := 0; x != vp.W; x++ {
			// top
			sx, sy, _ = vp.GetScreenCoordinates(x, 0)
			s.Scr.SetCell(sx, sy, st, borders[1])
			// bottom
			sx, sy, _ = vp.GetScreenCoordinates(x, vp.H-1)
			s.Scr.SetCell(sx, sy, st, borders[5])
		}
		for y := 0; y != vp.H; y++ {
			// left
			sx, sy, _ = vp.GetScreenCoordinates(0, y)
			s.Scr.SetCell(sx, sy, st, borders[3])
			// right
			sx, sy, _ = vp.GetScreenCoordinates(vp.W-1, y)
			s.Scr.SetCell(sx, sy, st, borders[7])

		}

		sx, sy, _ = vp.GetScreenCoordinates(0, 0)
		s.Scr.SetCell(sx, sy, st, borders[0])
		sx, sy, _ = vp.GetScreenCoordinates(vp.W-1, 0)
		s.Scr.SetCell(sx, sy, st, borders[2])
		sx, sy, _ = vp.GetScreenCoordinates(vp.W-1, vp.H-1)
		s.Scr.SetCell(sx, sy, st, borders[4])
		sx, sy, _ = vp.GetScreenCoordinates(0, vp.H-1)
		s.Scr.SetCell(sx, sy, st, borders[6])

		if vp.title != "" {
			tmp := fmt.Sprintf("[ %s ]", vp.title)

			sx, sy, _ = vp.GetScreenCoordinates(2, 0)
			s.DrawText(sx, sy, st, tmp)
		}
	}

	// Draw child viewports
	for _, child := range vp.children {
		child.Draw(s)
	}
}

// GetScreenCoordinates gets the absolute screen coordinates based on the viewport and parent viewports
func (vp *Viewport) GetScreenCoordinates(x, y int) (int, int, error) {
	rx, ry := x, y

	cvp := vp
	for cvp != nil {
		rx += cvp.X
		ry += cvp.Y
		if cvp.parent == nil {
			break
		}

		cvp = cvp.parent
	}

	if rx < 0 || ry < 0 {
		return -1, -1, errOutOfViewportBounds
	}
	if rx >= cvp.W || ry >= cvp.H {
		return -1, -1, errOutOfViewportBounds
	}

	return rx, ry, nil
}

func (vp *Viewport) SetTitle(title string) {
	vp.title = title
}
