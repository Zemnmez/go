package main

import (
	"time"

	"github.com/nsf/termbox-go"

	"zemn.me/reactive"
	"zemn.me/reactive/tree"
	"zemn.me/term"
)

func main() {
	if err := do(); err != nil {
		panic(err)
	}
}

var _ tree.Component = &FakeProcess{}

type FakeProcess struct {
	term.LoadingBar
	done chan bool
}

func (FakeProcess) Name() string { return "fakeprocess" }
func (f FakeProcess) Close()     { close(f.done) }
func (f FakeProcess) ShouldUpdate(old tree.Component) (bool, error) {
	return f != *old.(*FakeProcess), nil
}
func (f FakeProcess) Render() (components []tree.Component, err error) {
	return []tree.Component{f.LoadingBar}, nil
}

func (f *FakeProcess) Mount(s tree.StateController) {
	go func() {
		for {
			select {
			case <-f.done:
				return
			case <-time.After(1 * time.Second):
				f.LoadingBar.Progress = (f.LoadingBar.Progress + 0.1)
				if f.LoadingBar.Progress > 1 {
					f.LoadingBar.Progress = 0
				}
				s.Update()
			}
		}
	}()
}

func do() (err error) {
	term, err := term.New(func(c term.Canvas) (components []tree.Component, err error) {
		return []tree.Component{
			&FakeProcess{
				LoadingBar: term.LoadingBar{
					Canvas:   c,
					Empty:    ' ',
					Fill:     '#',
					Progress: .7,
				},
			},
		}, nil
	})

	defer termbox.Close()
	reactive.Render(term, mapper{})
	<-time.After(10 * time.Second)
	return nil
}

type mapper struct{}

func (mapper) Map(tree.Component) {
	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}
func (mapper) UnMap(tree.Component) {
	if err := termbox.Flush(); err != nil {
		panic(err)
	}
}
func (mapper) Error(c tree.Component, err error) {
	panic(err)
}
