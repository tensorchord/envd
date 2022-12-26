package language

import (
	. "github.com/onsi/ginkgo/v2"

	e2e "github.com/tensorchord/envd/e2e/v1"
)

var _ = Describe("rlang", Ordered, func() {
	It("Should build rlang environment successfully", func() {
		exampleName := "rlang"
		testcase := "e2e"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		e.BuildImage(true)()
		e.RunContainer()()
		e.DestroyContainer()()
	})
})
