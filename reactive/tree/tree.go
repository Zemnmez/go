/*
Package tree implements functionality for rendering a tree
of stateful, encapsulated components å la React.

This package is a low-level API that describes the internal workings
of the implementation. For the consumer API, see zemn.me/term/react.



Theory

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




Practice

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

	"zemn.me/debug"
)

// A Mapper represents a method of converting
// a Component to some final representation,
// for example an HTML DOM object or an area on a screen.
//
// The Mapper is called on every compnent each time it updates.
type Mapper interface {
	// Map is called whenever a Component is rendered
	// or updated
	Map(c Component)

	// UnMap is called whenever a Component has been removed
	// i.e. closed
	UnMap(c Component)

	// Error is called when an error occurs during rendering
	Error(c Component, err error)
}

// A Node represents a state tree. The Node ultimately
// maintains state updates for a Component and determines if
// its children should decide whether to update or not.
type Node struct {
	Component
	Children []Node
	Mapper

	// used to check if the number
	// of children has changed after first render
	previouslyRendered bool
}

// NewNode constructs a new state tree rooted at the Component c,
// calling Mapper m.Map(Component) each time a state change occurs
// in a Node.
func NewNode(c Component, m Mapper) (n *Node) {
	n = new(Node)
	n.Component = c
	n.Mapper = m

	n.Update()
	return
}

func (n Node) prepareChildren(nodes ...Node) (newChildren []Node) {
	for i := range nodes {
		nodes[i].Mapper = n.Mapper
	}

	return nodes
}

// The update function re-renders the children of this Node,
// and asks them if they need to update their children via Update().
//
// Errors are handled by the Update() function, which passes them to the
// Mapper.
func (n *Node) update() (err error) {
	debug.Log(" %s performing update ", n.Component.Name())

	// Tell the mapper this Component has updated.
	n.Mapper.Map(n.Component)

	debug.Log("%s mapper updated", n.Component.Name())

	newChildren, err := n.Render()

	if err != nil {
		return
	}

	if n.previouslyRendered {
		debug.Log("%s this is not the first time this component has rendered", reflect.TypeOf(n.Component))

		if len(newChildren) != len(n.Children) {
			return fmt.Errorf(
				"had %d Children and now has %d;"+
					" the number of children a Component has is not"+
					" allowed to change",

				len(n.Children),
				len(newChildren),
			)
		}
	}

	n.previouslyRendered = true

	debug.Log("%s diffing %d children", n.Component.Name(), len(newChildren))

	if len(newChildren) > len(n.Children) {
		debug.Log(
			"%s making new Node.Children",
			n.Component.Name(),
		)
		n.Children = n.prepareChildren(make([]Node, len(newChildren))...)
	}

	for i := range newChildren {
		newChild, oldChild := newChildren[i], n.Children[i].Component

		shouldUpdate := false
		mounted := false
		unmounted := false

		switch {
		// was nil, now defined, no need to ask if update is needed
		case oldChild == nil && newChild != nil:
			shouldUpdate = true
			mounted = true

			//unmounted = false

		// both new and old were non-nil:
		// delegate to new child as to whether
		// update is needed
		case newChild != nil && oldChild != nil:
			debug.Log("[%s] ShouldUpdate?", newChild.Name())

			shouldUpdate, err = newChild.ShouldUpdate(oldChild)
			if err != nil {
				return
			}

			//mounted = false
			//unmounted = falde

		// nothing to do
		case newChild == nil && oldChild == nil:
			//shouldUpdate = false

		// removed
		case oldChild != nil && newChild == nil:
			unmounted = true
			//mounted = false
			//shouldUpdate = false
		default:
			panic(fmt.Sprintf(
				"unknown state: old %+v, new %+v",
				oldChild,
				newChild,
			))
		}

		debug.Log(
			"%s child %d was, now ",
			n.Component.Name(),
			i,
		)

		debug.Log(
			`%s child %d:
	was unmounted: %v
	was mounted: %v
	needs to be updated: %v`,
			n.Component.Name(),
			i,
			unmounted,
			mounted,
			shouldUpdate,
		)

		debug.Assert(
			!(unmounted == true && mounted == unmounted),
			"cannot be both mounted and unmounted! "+
				"mounted: %v, unmounted: %v",
			mounted, unmounted,
		)

		if unmounted {
			n.Children[i].Close()
			n.Mapper.UnMap(n.Children[i])
		}

		n.Children[i].Component = newChild

		if mounted {
			n.Children[i].Mount(&n.Children[i])
		}

		if shouldUpdate {
			n.Children[i].Update()
		}

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
func (n *Node) Update() {
	if err := n.update(); err != nil {
		debug.Log("[%s] ERROR: %s", reflect.TypeOf(n.Component), err)

		n.Mapper.Error(
			n.Component,
			fmt.Errorf(
				"Update error in Component %s: %s",
				reflect.TypeOf(n.Component),
				err,
			),
		)
	}
}

type StateController interface {
	Update()
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

	Name() string
}
