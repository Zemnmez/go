package term_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestTerm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Term Suite")
}
