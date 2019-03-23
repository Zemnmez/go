package tree_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"zemn.me/debug"
	. "zemn.me/reactive/tree"
	"zemn.me/reactive/tree/treetest"
)

var _ = Describe("Tree", func() {
	root := &treetest.StaticComponent{Id: "root"}
	child_a := &treetest.StaticComponent{Id: "child_a"}
	child_b := &treetest.StaticComponent{Id: "child_b"}
	child_a_a := &treetest.StaticComponent{Id: "child_a_a"}
	child_a_b := &treetest.StaticComponent{Id: "child_a_b"}

	var Component_list = []*treetest.StaticComponent{
		root, child_a, child_b, child_a_a, child_a_b,
	}

	root.Children = []Component{child_a, child_b}
	child_a.Children = []Component{child_a_a, child_a_b}

	var rec treetest.Recorder

	rootNode := NewNode(root, &rec)

	When("a tree of Component is constructed", func() {
		// every Component should have been passed to the mapper
		Context("the mapper", func() {
			It("should have been passed all components", func() {
				Expect(len(Component_list)).To(Equal(len(rec.Components)),
					"num components: %d\nnum rendered: %d",
					len(Component_list),
					len(rec.Components),
				)

				Expect(rec.Components).To(HaveLen(len(Component_list)))
			})
		})

		Context("the root node", func() {
			It("should have children", func() {
				Expect(len(rootNode.Children)).To(
					BeNumerically(">", 0),
				)
			})
		})

		Context("the Component", func() {
			for _, StaticComponent := range Component_list {
				It("should have been passed to the mapper", func() {
					var seenComponents = make(map[string]bool, len(Component_list))

					for _, c := range rec.Components {
						seenComponents[c.(*treetest.StaticComponent).Id] =
							true
					}

					Expect(seenComponents[StaticComponent.Id]).To(
						Equal(true),
						"%s", StaticComponent.Id,
					)
				})

				It("should have been told it was mounted", func() {
					Expect(len(StaticComponent.MountCalls)).To(
						BeNumerically(">", 0),
					)
				})

				It("should not have been asked whether to update", func() {
					Expect(StaticComponent.ShouldUpdateCalls).To(HaveLen(0))
				})
			}
		})

	})
})

var _ = Describe("Tree Update", func() {
	root := &treetest.StaticComponent{Id: "root"}
	child_a := &treetest.StaticComponent{Id: "child_a"}
	child_b := &treetest.StaticComponent{Id: "child_b"}
	child_a_a := &treetest.StaticComponent{Id: "child_a_a"}
	child_a_b := &treetest.StaticComponent{Id: "child_a_b"}

	var Component_list = []*treetest.StaticComponent{
		root, child_a, child_b, child_a_a, child_a_b,
	}

	_ = Component_list

	root.Children = []Component{child_a, child_b}
	child_a.Children = []Component{child_a_a, child_a_b}

	var rec treetest.Recorder

	rootNode := NewNode(root, &rec)

	_ = rootNode

	When("a component updates", func() {
		// NB: root -> child a
		// root -> child b
		// child a -> child a_a
		// child a -> child a_b

		prevRenderCalls := len(child_a.RenderCalls)

		var oldChildren []treetest.StaticComponent
		for _, child := range child_a.Children {
			oldChildren = append(oldChildren, *child.(*treetest.StaticComponent))
		}

		debug.Log("forcing an update of %s", child_a.Name())
		child_a.ForceUpdate()

		It("should re-render", func() {
			Expect(child_a.RenderCalls).To(HaveLen(prevRenderCalls + 1))
		})

		for i, newChild := range child_a.Children {
			oldChild := oldChildren[i]

			Context(fmt.Sprintf("its child, %s", newChild.(*treetest.StaticComponent).Id), func() {
				It("should be asked whether to update", func() {
					Expect(newChild.(*treetest.StaticComponent).ShouldUpdateCalls).To(
						HaveLen(len(oldChild.ShouldUpdateCalls) + 1),
					)
				})

			})
		}
	})
})
