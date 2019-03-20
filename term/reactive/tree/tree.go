/*
Package tree implements functionality for rendering a tree
of stateful, encapsulated components å la React.

This package is a low-level API that describes the internal workings
of the implementation. For the consumer API, see zemn.me/term/react.

THEORY

For a react-like component, we need components that describe
their tree of dependencies for data, and are able to tell the
library when they're changed.

Building a system like this in Go has its pros and cons.

First, Javascript has the benefit of first-class data serialization objects --
in Javascript it's easier to say 'the component has some data'
and operate on changes in that data.

Second, React's JSX makes it much harder to do things that make it difficult
to operate in a reactive fashion properly. It's vital that components have
known keys, or 'slots', so we can map each component's state to its previous
state. This is something that's very difficult, if not impossible to enforce
in vanilla Go.

However, because it's much easier and faster to compare structures in Go,
it's much less necessary to have good ShouldUpdate comparisons, and
we totally avoid any need for anything like immutable.js by just not using
pointers.

PRACTICE

We start with a root Render function which constructs a tree of Components.
This tree is passed to whatever package ultimately maps component data
to some representation, like an HTML element or an area in a console.

When a component mounts i.e. we pass it to the final representation mapper
we call Component.Mount(StateController), where StateController has an
Update() function to tell our tree that the tree has updated at this point.

Then, we can perform another Component.Render() to get new child components,
and ask those children whether they will change as a result of their new
construction by passing the old children the new state.

If a child agrees that they should re-render, we call the Update() function
of the new child, as though the child itself updated. This continues
down the tree until it reaches a node with no children.

Because we need the children to be ordered the same way for mapping old to new
state, we panic if a Component.Render() produces a different number of children.

If a Component.Render() would cause a child to be removed, it instead returns
a nil Component. When a nil Component is returned, the Component.Close() function
of the Component that used to be in that slot is called to allow it to clean
itself up.


*/
package tree

import (
	"fmt"
	"reflect"
)

func (n Node) update() {}

// A Mapper represents a method of converting
// a Component to some final representation,
// for example an HTML DOM object or an area on a screen.
//
// The Mapper is called on every compnent each time it updates.
type Mapper interface {
	Map(c Component)
	Error(c Component, err error)
}

// A Node represents a state tree. The Node ultimately
// maintains state updates for a Component and determines if
// its children should decide whether to update or not.
type Node struct {
	Component
	Children []node
	Mapper

	// used to check if the number
	// of children has changed after first render
	previouslyRendered bool
}

// NewNode constructs a new state tree rooted at the Component c,
// calling Mapper m.Map(Component) each time a state change occurs
// in a Node.
func NewNode(c Component, m Mapper) (n Node) {
	n.Component = c
	n.Mapper = m

	n.Update()
	return
}

// The update function re-renders the children of this Node,
// and asks them if they need to update their children via Update().
//
// Errors are handled by the Update() function, which passes them to the
// Mapper.
func (n *Node) update() (err error) {

	// Tell the mapper this Component has updated.
	n.Mapper.Map(n.Component)

	var err error
	children, err := n.Render()

	if err != nil {
		return
	}

	if (len(children) != len(n.Children)) && !n.previouslyRendered {
		return fmt.Errorf(
			"had %d Children and now has %d"+
				" the number of children a Component has is not"+
				" allowed to change",

			len(n.Children),
			len(children),
		)
	}

	n.previouslyRendered = true

	for i, child := range children {
		var update bool

		// we compare new to old here
		// even though it makes for uglier code,
		// because otherwise every component when updated
		// to a nil component would be an awkward case to handle.
		if child != nil {
			update, err = child.ShouldUpdate(n.Children[i].Component)
			if err != nil {
				return
			}
		}

		// if the child was previously nil, and now is not,
		// the component Mounted.
		//
		// if the child was previously non-nil, and now is,
		// the component needs to be Closed.
		if child == nil || n.Children[i].Component == nil {
			// removed
			if child == nil && n.Children[i].Component != nil {
				n.Children[i].Component.Close()
			}

			if child != nil && n.Children[i].Component == nil {
				child.Mount(n)
			}
		}

		n.Children[i].Component = child

		if !update {
			continue
		}

		n.Children[i].Update()
	}

	return

}

// The Update function is called each time a Node is constructed
// for the first time, or its state changes.
//
// The Update function constructs new Node.Children and determines if
// any have changed.
//
// If an error occurs, it is passed to the Mapper via Mapper.Error().
func (n Node) Update() {
	if err := n.update(); err != nil {
		n.Mapper.Error(
			n.Component,
			fmt.Sprintf(
				"Update error in Component %s: %s",
				reflect.TypeOf(n.Component),
				err,
			),
		)
	}
}

type Component interface {
	// Mount is called when the Component is
	// mapped to some representation, like an HTML element,
	// or drawing area.
	Mount(s StateController)

	// Close is called when the Component is removed
	// from the representation by its parent replacing it
	// with `nil`.
	Close()

	// ShouldUpdate is called each time a parent re-renders.
	// ShouldUpdate is used to determine if this Component should
	// itself re-render.
	ShouldUpdate(old Component) (bool, error)

	// Render produces a list of child components, or nothing.
	// a []Component mus *not* change length or shuffle the
	// defintions of its child components around, as the order
	// has to be used to compare changes in state.
	Render() ([]Component, error)
}
