//Package reactive exposes APIs allowing React-like trees of
//components.
//
// For more in-depth information on the design, see ./tree.
package reactive

import "zemn.me/reactive/tree"

func Render(c Component, m Mapper) {
	tree.NewNode(c, m)
}

type StateController interface {
	Update()
}

type Component interface {
	// Mount is called when the Component is
	// mapped to some representation, like an HTML element,
	// or drawing area.
	Mount(s tree.StateController)

	// Close is called when the Component is removed
	// from the representation by its parent replacing it
	// with `nil`.
	Close()

	// ShouldUpdate is called each time a parent re-renders.
	// ShouldUpdate is used to determine if this Component should
	// itself re-render.
	ShouldUpdate(old tree.Component) (bool, error)

	// Render produces a list of child components, or nothing.
	// a []Component mus *not* change length or shuffle the
	// defintions of its child components around, as the order
	// has to be used to compare changes in state.
	Render() ([]tree.Component, error)

	Name() string
}

// A Mapper represents a method of converting
// a Component to some final representation,
// for example an HTML DOM object or an area on a screen.
//
// The Mapper is called on every compnent each time it updates.
type Mapper interface {
	// Map is called whenever a Component is rendered
	// or updated
	Map(c tree.Component)

	// UnMap is called whenever a Component has been removed
	// i.e. closed
	UnMap(c tree.Component)

	// Error is called when an error occurs during rendering
	Error(c tree.Component, err error)
}
