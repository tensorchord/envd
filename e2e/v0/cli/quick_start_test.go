// Copyright 2022 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cli

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	e2e "github.com/tensorchord/envd/e2e/v0"
)

var _ = Describe("e2e quickstart", Ordered, func() {
	exampleName := "quick-start"
	testcase := "e2e"
	e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
	BeforeAll(e.BuildImage(true))
	BeforeEach(e.RunContainer())
	It("execute python demo.py", func() {
		res, err := e.Exec("python demo.py")
		Expect(err).To(BeNil())
		Expect(res).To(Equal("[2 3 4]"))
	})
	AfterEach(e.DestroyContainer())
})
