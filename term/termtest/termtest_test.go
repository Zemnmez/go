package termtest_test

import (
	"image"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"zemn.me/term"
	. "zemn.me/term/termtest"
)

var _ = Describe("Canvas", func() {
	When("generating a canvas", func() {
		var c Canvas

		const width = 10
		const height = 20

		It("should not panic", func() {
			c = NewCanvas(width, height)
		})

		It("should generate a buffer of the right size", func() {
			Expect(len(c.Base)).To(Equal(width * height))
		})

		It("should generate columns with the right size", func() {
			Expect(len(c.Cells)).To(Equal(height),
				"%+v", c.Cells)
		})

		It("should generate rows with the right size", func() {
			Expect(len(c.Cells[0])).To(Equal(width))
		})

		It("should return cells as expected", func() {
			Expect(c.Buffer()).To(Equal(c.Cells))
		})
	})

	It("should implement zemn.me/term.Canvas", func(done Done) {
		defer close(done)
		var _ term.Canvas = NewCanvas(10, 20)
	})

	It("should successfully allow writes to all cells", func(done Done) {
		defer close(done)
		const w = 300
		const h = 400
		c := NewCanvas(w, h)

		for x := 0; x < w; x++ {
			for y := 0; y < h; y++ {
				c.SetCell(image.Pt(x, y), TestCell)
			}
		}

		for _, cell := range c.Base {
			Expect(cell).To(Equal(TestCell), "%+v")
		}
	})

	When("split into a sub-canvas", func() {
		It("should have the correct rect size", func(done Done) {
			defer close(done)
			const w = 300
			const h = 400
			var c term.Canvas = NewCanvas(w, h)

			const nw = 100
			const nh = 100
			const dx = 100
			const dy = 100

			c = c.Canvas(image.Rect(
				0, 0,
				nw, nh,
			).Add(image.Pt(dx, dy)))

			r := c.Rect()

			Expect(r.Dx()).To(Equal(nw))
			Expect(r.Dy()).To(Equal(nh))
		})

		It("should successfully allow writes to all cells", func(done Done) {
			defer close(done)
			const w = 300
			const h = 400
			c := NewCanvas(w, h)

			const nw = 100
			const nh = 100
			const dx = 100
			const dy = 100

			c = c.Canvas(image.Rect(
				0, 0,
				nw, nh,
			).Add(image.Pt(dx, dy))).(Canvas)

			for x := 0; x < w; x++ {
				for y := 0; y < h; y++ {
					c.SetCell(image.Pt(x, y), TestCell)
				}
			}

			for _, cell := range c.Base {
				Expect(cell).To(Equal(TestCell), "%+v")
			}
		})
	})
})
