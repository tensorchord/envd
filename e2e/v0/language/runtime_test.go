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

package language

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	e2e "github.com/tensorchord/envd/e2e/v0"
)

var _ = Describe("runtime", Ordered, func() {
	exampleName := "runtime"
	testcase := "e2e"
	e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
	BeforeAll(e.BuildImage(true))
	BeforeEach(e.RunContainer())
	It("execute runtime command `numpy`", func() {
		res, err := e.ExecRuntimeCommand("numpy")
		Expect(err).To(BeNil())
		Expect(res).To(Equal("[2 3 4]"))
	})
	It("execute runtime command `root`", func() {
		_, err := e.ExecRuntimeCommand("root")
		Expect(err).To(BeNil())
	})
	AfterEach(e.DestroyContainer())
})
