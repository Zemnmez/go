//Package term provides functions for manipulating terminal
//output
package term // import "zemn.me/term"

import (
	"image"

	"github.com/nsf/termbox-go"
	"zemn.me/reactive"
	"zemn.me/reactive/tree"
)

type Attribute = termbox.Attribute
type Cell = termbox.Cell

var _ reactive.Component = &Term{}

type Term struct {
	RenderFunc func(Canvas) (components []tree.Component, err error)
	Canvas
	stop chan bool
}

func New(render func(Canvas) (components []tree.Component, err error)) (term *Term, err error) {
	err = termbox.Init()
	if err != nil {
		return
	}

	return &Term{RenderFunc: render, Canvas: newRootCanvas()}, nil
}

func (Term) Name() string { return "term" }
func (t Term) ShouldUpdate(old tree.Component) (bool, error) {
	a, b := t, *old.(*Term)
	if a.Canvas.ShouldUpdate(b.Canvas) {
		return true, nil
	}
	a.Canvas, b.Canvas = nil, nil
	return false, nil
}
func (t Term) Render() ([]tree.Component, error) { return t.RenderFunc(t.Canvas) }
func (t Term) Close()                            { close(t.stop) }
func (t *Term) Mount(s tree.StateController) {
	go func() {
		for {
			// TODO: make this a channel somehow?
			ev := termbox.PollEvent()
			switch ev.Type {
			case termbox.EventResize:
				t.Canvas = newRootCanvas()
				s.Update()
			}
		}
	}()
}

type Canvas interface {
	Rect() image.Rectangle
	SetCell(pos image.Point, cell Cell)
	Buffer() [][]Cell
	Canvas(image.Rectangle) Canvas
	ShouldUpdate(c Canvas) bool
}

type rootCanvas struct {
	Cells [][]Cell
}

func (r rootCanvas) ShouldUpdate(c Canvas) bool {
	a, b := r, c.(rootCanvas)
	return len(a.Cells) != len(b.Cells) ||
		(len(a.Cells) > 0 && (len(a.Cells[0]) != len(b.Cells[0])))
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

func (c canvas) ShouldUpdate(c2 Canvas) bool {
	a, b := c, c2.(canvas)
	return a.Width != b.Width || a.Height != b.Height
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

var _ tree.Component = LoadingBar{}

type LoadingBar struct {
	Fill     rune
	Empty    rune
	Progress float64
	Canvas
}

func (LoadingBar) Mount(tree.StateController) {}
func (LoadingBar) Close()                     {}
func (l LoadingBar) ShouldUpdate(old tree.Component) (should bool, err error) {
	a, b := l, old.(LoadingBar)
	if a.Canvas.ShouldUpdate(b.Canvas) {
		return true, nil
	}
	a.Canvas, b.Canvas = nil, nil
	should = a != b
	return
}
func (LoadingBar) Name() string { return "LoadingBar" }
func (l LoadingBar) Render() (children []tree.Component, err error) {
	termbox.Flush()
	c := l.Canvas

	/*
			 - loaded - unloaded -
		     x-----------|-------x
			 |###########|///////|
			 |###########|///////|
			 x-----------|-------|
	*/
	loaded := c.Canvas(image.Rect(
		0, 0,
		int(float64(c.Rect().Max.X)*l.Progress),
		c.Rect().Max.Y,
	))

	children = append(children, Fill{
		Canvas: loaded,
		Cell:   Cell{Ch: l.Fill},
	})

	unloaded := c.Canvas(image.Rect(
		int(float64(c.Rect().Max.X)*l.Progress), 0,
		c.Rect().Max.X, c.Rect().Max.Y,
	))

	children = append(children, Fill{
		Canvas: unloaded,
		Cell:   Cell{Ch: l.Empty},
	})

	return
}

type Fill struct {
	Cell
	Canvas
}

func (Fill) Close()                     {}
func (Fill) Mount(tree.StateController) {}
func (f Fill) ShouldUpdate(old tree.Component) (should bool, err error) {
	a, b := f, old.(Fill)
	if a.Canvas.ShouldUpdate(b.Canvas) {
		return true, nil
	}
	a.Canvas, b.Canvas = nil, nil
	return a != b, nil
}
func (Fill) Name() string { return "fill" }
func (f Fill) Render() (_ []tree.Component, err error) {
	c := f.Canvas
	rows := c.Buffer()
	for y := range rows {
		for x := range rows[y] {
			rows[y][x] = Cell(f.Cell)
		}
	}

	return
}

var _ reactive.Component = Text{}

type Text struct {
	Text string
	Canvas
}

func (Text) Name() string { return "text" }
func (t Text) ShouldUpdate(old tree.Component) (bool, error) {
	a, b := t, old.(Text)
	if a.Canvas.ShouldUpdate(b.Canvas) {
		return true, nil
	}
	a.Canvas, b.Canvas = nil, nil
	return a != b, nil
}
func (t Text) Close()                     {}
func (t Text) Mount(tree.StateController) {}
func (f Text) Render() (_ []tree.Component, err error) {
	c := f.Canvas
	runes := []rune(f.Text)
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
