//Package term provides functions for manipulating terminal
//output
package term // import "zemn.me/term"

import (
	"image"

	"github.com/nsf/termbox-go"
	"zemn.me/reactive"
)

type Component interface {
	Render(c Canvas) (children []Component, err error)
	Mount(s reactive.StateController)
	Close()
	ShouldUpdate(old Component) (bool, error)
}

type WithCanvas struct {
	Component
	Canvas
}

func (w WithCanvas) Mount(s reactive.StateController) { w.Component.Mount(s) }
func (w WithCanvas) Close() { w.Component.Close() }
func (w WithCanvas) ShouldUpdate(old reactive.Component) (should bool, err error) { return w.Component.ShouldUpdate(old) }
func (w WithCanvas) Render() (children []reactive.Component, err error) {
	c, err := w.Component.Render(w.Canvas)
	if err != nil { return }

	children = make([]reactive.Component, len(c))
	for i := range c {
		children[i] =
	}
}

type Attribute = termbox.Attribute
type Cell = termbox.Cell

type Canvas interface {
	Rect() image.Rectangle
	SetCell(pos image.Point, cell Cell)
	Buffer() [][]Cell
	Canvas(image.Rectangle) Canvas
}

type rootCanvas struct {
	Cells [][]Cell
}

func (c rootCanvas) Rect() image.Rectangle {
	w, h := termbox.Size()
	return image.Rect(
		0, 0,
		w, h,
	)
}

func (c rootCanvas) SetCell(pos image.Point, cell Cell) {
	termbox.SetCell(pos.X, pos.Y, cell.Ch, cell.Fg, cell.Bg)
}
func (c rootCanvas) Buffer() [][]Cell { return c.Cells }
func (c rootCanvas) Canvas(r image.Rectangle) Canvas {
	cn := canvas{}
	cn.Width = r.Dx()
	cn.Height = r.Dy()

	x, y := r.Min.X, r.Min.Y

	dx := r.Dx()
	for dy := 0; dy < r.Dy(); dy++ {
		cn.Cells = append(cn.Cells, c.Cells[y+dy][x:x+dx])
	}

	return cn
}

type canvas struct {
	Cells         [][]Cell
	Width, Height int
}

func (c canvas) Buffer() [][]Cell { return c.Cells }

func (c canvas) Rect() image.Rectangle {
	return image.Rect(
		0, 0,
		c.Width, c.Height,
	)
}

func (c canvas) SetCell(pos image.Point, cell Cell) {
	c.Cells[pos.Y][pos.X] = cell
}

func (c canvas) Canvas(r image.Rectangle) Canvas {
	cn := canvas{}
	cn.Width = r.Dx()
	cn.Height = r.Dy()

	x, y := r.Min.X, r.Min.Y

	dx := r.Dx()
	for dy := 0; dy < r.Dy(); dy++ {
		cn.Cells = append(cn.Cells, c.Cells[y+dy][x:x+dx])
	}

	return cn
}

func newRootCanvas() (c rootCanvas) {
	cells := termbox.CellBuffer()
	width, height := termbox.Size()
	for y := 0; y < height; y++ {
		start := y * width
		end := y*width + width
		c.Cells = append(c.Cells, cells[start:end])
	}

	return
}

func NewCanvas() (c Canvas, done func(), err error) {
	err = termbox.Init()
	if err != nil {
		return
	}

	c = newRootCanvas()
	done = termbox.Close
	return
}

type LoadingBar struct {
	Fill     rune
	Empty    rune
	Progress float64
	Canvas
}

func (l LoadingBar) Render() (children []reactive.Component, err error) {
	loaded := c.Canvas.Canvas(image.Rect(
		0, 0,
		int(float64(c.Rect().Max.X)*l.Progress), c.Rect().Max.Y,
	))

	children = append(children, Child{
		loaded,
		Fill{Ch: l.Fill},
	})

	unloaded := c.Canvas.Canvas(image.Rect(
		int(float64(c.Rect().Max.X)*l.Progress), 0,
		c.Rect().Max.X, c.Rect().Max.Y,
	))

	children = append(children, Child{
		unloaded,
		Fill{Ch: l.Empty},
	})

	return
}

type Fill struct {
	Cell
}

func (Fill) Close()                                                         {}
func (Fill) Mount(reactive.StateController)                                 {}
func (f Fill) ShouldUpdate(old reactive.Component) (should bool, err error) { return old.Cell != f.Cell, nil }
func (Fill) Name() string { return "fill" }
func (f Fill) Render(c Canvas) (_ []reactive.Component, err error) {
	rows := c.Buffer()
	for y := range rows {
		for x := range rows[y] {
			rows[y][x] = Cell(f)
		}
	}

	return
}

type Text string

func (f Text) Render(c Canvas) (_ []reactive.Component, err error) {
	runes := []rune(f)
	for i, r := range runes {
		x := i % c.Rect().Dx()
		y := i / c.Rect().Dx()
		if x == c.Rect().Dx() ||
			y == c.Rect().Dy() {
			break
		}
		c.Buffer()[y][x] = Cell{Ch: r}
	}
	return
}
