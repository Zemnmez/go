//Package term provides functions for manipulating terminal
//output
package term // import "zemn.me/term"

import (
	"github.com/nsf/termbox-go"
	"image"
)

type Attribute = termbox.Attribute
type Cell = termbox.Cell

type Child struct { Canvas; Component }

type Canvas interface {
	Rect() image.Rectangle
	SetCell(pos image.Point, cell Cell)
	Buffer() [][]Cell
	Canvas(image.Rectangle) Canvas
}

type Component interface {
	Render(c Canvas) ([]Child, error)
}

type LoadingBar struct {
	Fill rune
	Empty rune
	Progress float64
}

func (l LoadingBar) Render(c Canvas) (children []Child, err error) {
	loaded := c.Canvas(image.Rect(
		0, 0,
		int(float64(c.Rect().Max.X) * l.Progress), c.Rect().Max.Y,
	))

	children = append(children, Child {
		loaded,
		Fill{ Ch: l.Fill },
	})

	unloaded := c.Canvas(image.Rect(
		int(float64(c.Rect().Max.X)*l.Progress), 0,
		c.Rect().Max.X, c.Rect().Max.Y,
	))

	children = append(children, Child {
		unloaded,
		Fill { Ch: l.Empty },
	})

	return
}

type Fill Cell
func (f Fill) Render(c Canvas) (children []Child, err error){
	rows := c.Buffer()
	for y := range rows {
		for x := range rows[y] {
			rows[y][x] = Cell(f) }}

			return
}
