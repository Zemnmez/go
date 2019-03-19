package term_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "zemn.me/term"
	termtest "zemn.me/term/termtest"
)

var _ = Describe("Text", func() {
	It("should implement term.Component", func(done Done) {
		defer close(done)
		var _ Component = Text("hi!")
	})

	When("rendered", func() {
		It("should fill the buffer with text sequentially", func(done Done) {
			defer close(done)
			const text Text = "hello world!"
			const w = 2
			const h = (len(text) / 2) + 2
			c := termtest.NewCanvas(w, h)

			children, err := text.Render(c)
			Expect(len(children)).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())

			for i, r := range []rune(text) {
				Expect(c.Base[i].Ch).To(Equal(r),
					"%+v", c.Base)
			}
		})

		It("should truncate upon overflow", func(done Done) {
			defer close(done)
			const text Text = "hello world!"
			const w = 2
			const h = 2
			c := termtest.NewCanvas(w, h)

			children, err := text.Render(c)
			Expect(len(children)).To(Equal(0))
			Expect(err).ToNot(HaveOccurred())

			for i, r := range []rune(text) {
				if i == w*h {
					break
				}
				Expect(c.Base[i].Ch).To(Equal(r),
					"%+v", c.Base)
			}
		})
	})
})

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
