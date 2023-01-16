package docs

import (
	. "github.com/onsi/ginkgo/v2"

	e2e "github.com/tensorchord/envd/e2e/v1"
)

var _ = Describe("lang", Ordered, func() {
	It("Should build rlang environment successfully", func() {
		exampleName := "rlang"
		testcase := "e2e-extra"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		e.BuildImage(true)()
		e.RunContainer()()
		e.DestroyContainer()()
	})
	It("Should build Julia environment successfully", func() {
		exampleName := "julia"
		testcase := "e2e-extra"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		e.BuildImage(true)()
		e.RunContainer()()
		e.DestroyContainer()()
	})
})
