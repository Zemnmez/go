package term_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "zemn.me/term"
)

var _ = Describe("Term", func() {
	Context("in the simplest possible case", func() {
		It("should draw a simple hello world", func() {
			const Text = "hello world!"
			var c Canvas
			c.Draw(SimpleText{
				Pos: { 0, 0 },
				Text: Text,
			})

		}),
	})
})
