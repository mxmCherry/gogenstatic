package gogenstatic_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGogenstatic(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Gogenstatic Suite")
}
