package tree_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "zemn.me/term/reactive/tree"
	"zemn.me/term/reactive/tree/treetest"
)

var _ = Describe("Tree", func() {
	When("a tree of Component is constructed", func() {
		It("should pass the Components to the mapper", func() {

			root := &treetest.StaticComponent{Id: "root"}
			child_a := &treetest.StaticComponent{Id: "child_a"}
			child_b := &treetest.StaticComponent{Id: "child_b"}
			child_a_a := &treetest.StaticComponent{Id: "child_a_a"}

			var Component_list = []*treetest.StaticComponent{
				root, child_a, child_b, child_a_a,
			}

			root.Children = []Component{child_a, child_b}
			child_a.Children = []Component{child_a_a}

			var r treetest.Recorder

			_ = NewNode(root, &r)

			// every Component should have been passed to the mapper
			Expect(len(Component_list)).To(Equal(len(r.Components)))
			var seenComponents = make(map[string]bool, len(Component_list))

			for _, c := range r.Components {
				seenComponents[c.(*treetest.StaticComponent).Id] =
					true
			}
			for _, StaticComponent := range Component_list {
				Expect(seenComponents[StaticComponent.Id]).To(
					Equal(true),
					"%s", StaticComponent.Id,
				)

				// component should have mounted
				Expect(len(StaticComponent.MountCalls)).To(Equal(1))

				// should not have asked whether to update
				Expect(len(StaticComponent.ShouldUpdateCalls)).To(Equal(0))
			}
		})

	})

	/*
		When("an update happens", func() {
			It("should only update.Children that change", func(done Done) {
				defer close(done)
				Expect(true).To(Equal(false))
			})
		}) */
})
