//Package term provides functions for manipulating terminal
//output
package term // import "zemn.me/term"

import (
	"github.com/nsf/termbox-go"
)

type Vec2d [2]int
type Pos Vec2d
type Size Vec2d

type Attribute = termbox.Attribute
type Cell = termbox.Cell

type Canvas interface {
	Size() Size
	SetCell(pos Pos, chr rune, fg, bg Attribute)
	Buffer() []Cell
	Area(Pos, Size)
}

type Component interface {
	Render(c Canvas) (err error)
}

var _ Component = LoadingBar{}
type LoadingBar struct {
	Fill rune
	Empty rune
}

func (l LoadingBar) Render(c Canvas) {
	size := c.Size()

}

var _ Component = Fill{}
type Fill Cell
func (f Fill) Render(c Canvas) {
	buf := c.Buffer()
	for i := range buf { buf[i] = Cell(f) }
}
