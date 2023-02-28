package docs

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	e2e "github.com/tensorchord/envd/e2e/v1"
)

var _ = Describe("rlang_mnist", Ordered, func() {
	exampleName := "rlang_mnist"
	testcase := "e2e"
	e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
	BeforeAll(e.BuildImage(true))
	BeforeEach(e.RunContainer())
	FIt("execute runtime command `Rscript`", func() {
		res, err := e.ExecRuntimeCommand("rlang-mnist")
		Expect(err).To(BeNil())
		isNumeric := "TRUE"
		Expect(res).To(BeEquivalentTo(isNumeric))
	})
	AfterEach(e.DestroyContainer())
})
