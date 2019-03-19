//Package termtest exposes some values for testing term programs
package termtest

import (
	"image"

	"zemn.me/term"
)

func NewCanvas(width, height int) (c Canvas) {
	c = Canvas{
		Width:  width,
		Height: height,
	}

	c.Base = make([]term.Cell, width*height)

	for y := 0; y < height; y++ {
		start := y * width
		end := y*width + width
		c.Cells = append(c.Cells, c.Base[start:end])
	}

	return
}

var TestCell = term.Cell{
	Ch: 'a',
	Fg: 1,
	Bg: 1,
}

var TestBar = term.LoadingBar{
	Fill:     '#',
	Empty:    ' ',
	Progress: .5,
}

type Canvas struct {
	Base          []term.Cell
	Cells         [][]term.Cell
	Width, Height int
}

func (c Canvas) Buffer() [][]term.Cell { return c.Cells }

func (c Canvas) Rect() image.Rectangle {
	return image.Rect(
		0, 0,
		c.Width, c.Height,
	)
}

func (c Canvas) SetCell(pos image.Point, cell term.Cell) {
	c.Cells[pos.Y][pos.X] = cell
}

func (c Canvas) Canvas(r image.Rectangle) term.Canvas {
	cn := Canvas{}
	cn.Base = c.Base
	cn.Width = r.Dx()
	cn.Height = r.Dy()

	x, y := r.Min.X, r.Min.Y

	dx := r.Dx()
	for dy := 0; dy < r.Dy(); dy++ {
		cn.Cells = append(cn.Cells, c.Cells[y+dy][x:x+dx])
	}

	return cn
}
