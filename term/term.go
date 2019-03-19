//Package term provides functions for manipulating terminal
//output
package term // import "zemn.me/term"

import (
	"github.com/nsf/termbox-go"
)

type DrawableAPI interface {
	SetCell(pos Pos, ch rune, fg, bg termbox.Attribute)
	Size() (w, h int)
}

type Drawable interface {
	Draw(api DrawableAPI) (err error)
	Update() (<-chan bool)
}

type Canvas struct {
	initOnce sync.Once
}

func (c Canvas) init() (err error) {
	err = termbox.Init()
	if err != nil { return }
	return
}

type Pos [2]int
func (v Pos) X() int { return v[0] }
func (v Pos) Y() int { return v[1] }
func (v *Pos) SetX(x int) { v[0] = x }
func (v *Pos) SetY(y int) { v[1] = y }
func (p Pos) Add(pn ...Pos) Pos {
	for _, vec := range pn {
		for i, v := range vec {
			p[i] += v
		}
	}

	return p
}

func (c *Canvas) Init() (err error) {
	c.initOnce.Do(func() {
		err = c.init()
	})

	return
}

func (c Canvas) Close() (err error) {
	termbox.Close()
	return
}

type drawApiImpl struct { }
func (drawApiImpl) SetCell(p Pos, ch rune, fg, bg termbox.Attribute) {
	termbox.SetCell(p.X(), p.Y(), ch, termbox.ColorBlue, bg)
}
func (drawApiImpl) Size() (w, h int) { return termbox.Size() }

var drawApi = drawApiImpl{}

func (c *Canvas) Draw(d Drawable) (err error) {
	if err = c.Init(); err != nil { return }
	return d.Draw(drawApi)
}

type SimpleText struct {
	Pos
	Text string
}

func (s SimpleText) Draw(api DrawableAPI) (err error) {
	for i, r := range s.Text {
		api.SetCell(s.Pos.Add(Pos{ i, 0 }), r, 0, 0)
	}

	return
}

func (Canvas) Flush() { termbox.Flush() }

func (s SimpleText) Update() (<-chan bool) { return nil }
