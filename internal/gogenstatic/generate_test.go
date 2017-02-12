package gogenstatic_test

import (
	"bytes"
	"io/ioutil"
	"regexp"

	. "github.com/mxmCherry/gogenstatic/internal/gogenstatic"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generate", func() {
	It("should generate static", func() {
		buf := bytes.NewBuffer(nil)

		err := Generate(buf, "testdata/src")
		Expect(err).NotTo(HaveOccurred())

		Expect(normalise(buf.String())).To(Equal(readFile("testdata/static.go")))
	})
})

func normalise(code string) string {
	return regexp.
		MustCompile("var modTime = time\\.Unix\\(\\d+, 0\\)").
		ReplaceAllString(code, "var modTime = time.Unix(0000000000, 0)")
}

func readFile(p string) string {
	b, err := ioutil.ReadFile(p)
	Expect(err).NotTo(HaveOccurred())
	return string(b)
}
