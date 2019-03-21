package tree_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "zemn.me/term/reactive/tree"
	"zemn.me/term/reactive/tree/treetest"
)

var _ = Describe("Tree", func() {
	It("should pass a tree of components to the mapper", func(done Done) {
		defer close(done)

		root := treetest.StaticComponent{Id: "root"}
		child_a := treetest.StaticComponent{Id: "child_a"}
		child_b := treetest.StaticComponent{Id: "child_b"}
		child_a_a := treetest.StaticComponent{Id: "child_a_a"}

		var component_list = []treetest.StaticComponent{
			root, child_a, child_b, child_a_a,
		}

		root.Children = []Component{child_a, child_b}
		child_a.Children = []Component{child_a_a}

		var r treetest.Recorder

		_ = NewNode(root, &r)

		// every component should have been passed to the mapper
		Expect(len(component_list)).To(Equal(len(r.Components)))
		var seenComponents = make(map[string]bool, len(component_list))

		for _, c := range r.Components {
			seenComponents[c.(treetest.StaticComponent).Id] =
				true
		}

		for _, staticComponent := range component_list {
			Expect(seenComponents[staticComponent.Id]).To(
				Equal(true),
				"%s", staticComponent.Id,
			)
		}
	})
})
