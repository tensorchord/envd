package docs

import (
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	e2e "github.com/tensorchord/envd/e2e/v1"
)

var _ = Describe("julia_mnist", Ordered, func() {
	exampleName := "julia_mnist"
	testcase := "e2e"
	e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
	BeforeAll(e.BuildImage(true))
	BeforeEach(e.RunContainer())
	It("execute runtime command `julia-mnist`", func() {
		res, err := e.ExecRuntimeCommand("julia-mnist")
		Expect(err).To(BeNil())
		IsNumber := func(s string) bool {
			_, err = strconv.ParseFloat(s, 64)
			return err == nil
		}
		Expect(res).To(Satisfy(IsNumber))
	})
	AfterEach(e.DestroyContainer())
})
