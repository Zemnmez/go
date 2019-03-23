package tree_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "zemn.me/reactive/tree"
	"zemn.me/reactive/tree/treetest"
)

var _ = Describe("Tree", func() {

	root := &treetest.StaticComponent{Id: "root"}
	child_a := &treetest.StaticComponent{Id: "child_a"}
	child_b := &treetest.StaticComponent{Id: "child_b"}
	child_a_a := &treetest.StaticComponent{Id: "child_a_a"}

	var Component_list = []*treetest.StaticComponent{
		root, child_a, child_b, child_a_a,
	}

	root.Children = []Component{child_a, child_b}
	child_a.Children = []Component{child_a_a}

	When("a tree of Component is constructed", func() {

		var r treetest.Recorder

		_ = NewNode(root, &r)

		// every Component should have been passed to the mapper
		Context("the mapper", func() {
			It("should have been passed all components", func() {
				Expect(len(Component_list)).To(Equal(len(r.Components)))
			})
		})
		var seenComponents = make(map[string]bool, len(Component_list))

		for _, c := range r.Components {
			seenComponents[c.(*treetest.StaticComponent).Id] =
				true
		}

		Context("the Component", func() {
			for _, StaticComponent := range Component_list {
				It("should have been passed to the mapper", func() {
					Expect(seenComponents[StaticComponent.Id]).To(
						Equal(true),
						"%s", StaticComponent.Id,
					)
				})

				It("should have been told it was mounted", func() {
					Expect(len(StaticComponent.MountCalls)).To(Equal(1))

				})

				It("should not have been asked whether to update", func() {
					Expect(len(StaticComponent.ShouldUpdateCalls)).To(Equal(0))
				})
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
