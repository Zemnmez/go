package tree_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "zemn.me/term/reactive/tree"
)

var _ = Describe("Tree", func() {
	it("should pass a tree of components to the mapper", func(done done) {
		defer close(done)

		root := treetest.staticcomponent{id: "root"}
		child_a := treetest.staticcomponent{id: "child_a"}
		child_b := treetest.staticcomponent{id: "child_b"}
		child_a_a := treetest.staticcomponent{id: "child_a_a"}

		var component_list = []treetest.staticcomponent{
			root, child_a, child_b, child_a_a,
		}

		root.children = []component{child_a, child_b}
		child_a.children = []component{child_a_a}

		var r treetest.recorder

		_ = newnode(root, &r)

		// every component should have been passed to the mapper
		expect(len(component_list)).to(equal(len(r.components)))
		var seencomponents = make(map[string]bool, len(component_list))

		for _, c := range r.components {
			seencomponents[c.(treetest.staticcomponent).id] =
				true
		}

		for _, staticcomponent := range component_list {
			expect(seencomponents[staticcomponent.id]).to(
				equal(true),
				"%s", staticcomponent.id,
			)
		}
	})

	/*
		When("an update happens", func() {
			It("should only update children that change", func(done Done) {
				defer close(done)
				Expect(true).To(Equal(false))
			})
		}) */
})
