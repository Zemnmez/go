//Package reactive exposes functionality for rendering
// and re-rendering a tree of stateful components a la React.
package reactive

type Element struct {
	value interface{}
}

func (e Element) Render(s StateController) []Component {
	panic("an Element should not be rendered; instead its value should be used directly")

}

type StateController interface {
	Update(func())
}

type Component interface {
	Mount(s StateController)
	Close()

	Render(s StateController) ([]Render, err error)
}

type Render func() Component

/*

type Counter struct {
	count int

	close chan-> bool
}

func (c *Counter) Close() {
	close(c.close)
}

func (c *Counter) Mount(s StateController) {
	c.close = make(chan bool)

	go func() {
		for {
			select {
			case <-time.After(1*time.Second):
				s.Update(func() { c.count++ })
			case <-c.close:
				return
			}
		}
	}
}

type Text struct {
	text string
}

func (t Text) Render() []Render {
	return []Render{
		reactive.Element{ t.text }
	}
}

func (f Fill) Render(s StateController) (children []Component, err error){
	return []Component{
		reactive.Element{ f.count },
	}, nil
}

type App struct { }
func (a App) Render() {
	return Counter{}
}

*/
