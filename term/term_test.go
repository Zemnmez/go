package term_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "zemn.me/term"
	termtest "zemn.me/term/termtest"
)

var _ = Describe("LoadingBar", func() {
	It("should implement term.Component", func(done Done) {
		defer close(done)
		var _ Component = LoadingBar{}
	})

	When("rendered as half done", func() {
		It("should produce two children", func(done Done) {
			defer close(done)
			c := termtest.NewCanvas(2, 2)
			children, err := termtest.TestBar.Render(c)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(children)).To(Equal(2))
		})
	})
})

var _ = Describe("Fill", func() {
	It("should implement interfaces correctly", func(done Done) {
		defer close(done)
		var _ Component = Fill{}
	})

	When("rendered", func() {
		It("should have no children", func(done Done) {
			defer close(done)
			c := termtest.NewCanvas(20, 30)
			children, err := Fill(termtest.TestCell).Render(c)
			Expect(err).ToNot(HaveOccurred())
			Expect(len(children)).To(Equal(0))
		})

		It("should fill a canvas", func(done Done) {
			defer close(done)
			c := termtest.NewCanvas(20, 30)
			Fill(termtest.TestCell).Render(c)
			for _, cell := range c.Base {
				Expect(cell).To(Equal(termtest.TestCell))
			}
		})

	})

})
