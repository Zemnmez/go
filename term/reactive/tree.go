package reactive

type node struct {
	Component
	Children []node
}

func (n *node) Update() {
	children := n.Component.Render()
	if len(c.Children) != len(children) {
		panic("different number of children returned." +
			"number of children must remain the same.")
	}


	for i, oldChild := range n.Children {
		shouldUpdate, err := oldChild.Component.ShouldUpdate(children[i])
		if err != nil { panic(err) }

		// no change
		if !shouldUpdate {
			continue
		}

		n.Children[i].Component = children[i]
		n.Children[i].Update()
	}
}

func newTree(c Component) (n node) {
	n.Component = c
	children := c.Render()
	if len(children) <= 0 {
		return
	}
	c.Children = make([]node, len(children))
	for i, comp := range c.Children {
		if comp == nil {
			// empty slot
			continue
		}
		c.Children[i] = NewTree(comp)
	}
	return
}

