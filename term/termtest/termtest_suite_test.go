package termtest_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTermtest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Termtest Suite")
}
