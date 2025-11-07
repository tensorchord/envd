// Copyright 2025 The envd Authors
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

	"github.com/tensorchord/envd/e2e"
)

var _ = Describe("rust", Ordered, func() {
	testcase := "e2e"

	Describe("Should install rust lang successfully", func() {
		exampleName := "rust"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		BeforeAll(e.BuildImage(true))
		BeforeEach(e.RunContainer())
		It("Should have cargo installed", func() {
			res, err := e.ExecRuntimeCommand("show")
			Expect(err).To(BeNil())
			Expect(res).To(ContainSubstring("toolchain"))
		})
		AfterEach(e.DestroyContainer())
	})

	Describe("Should install golang successfully", func() {
		exampleName := "golang"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		BeforeAll(e.BuildImage(true))
		BeforeEach(e.RunContainer())
		It("Should have go installed", func() {
			res, err := e.ExecRuntimeCommand("version")
			Expect(err).To(BeNil())
			Expect(res).To(ContainSubstring("go version"))
		})
		AfterEach(e.DestroyContainer())
	})
})
