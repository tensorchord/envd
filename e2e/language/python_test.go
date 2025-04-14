// Copyright 2023 The envd Authors
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

var _ = Describe("python", Ordered, func() {
	testcase := "e2e"
	It("Should build packages successfully", func() {
		exampleName := "python/packages"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		e.BuildImage(true)()
		e.RunContainer()()
		e.DestroyContainer()()
	})
	It("Should build requirements successfully", func() {
		exampleName := "python/requirements"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		e.BuildImage(true)()
		e.RunContainer()()
		e.DestroyContainer()()
	})
	It("Should build hybrid successfully", func() {
		exampleName := "python/hybrid"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		e.BuildImage(true)()
		e.RunContainer()()
		e.DestroyContainer()()
	})

	It("Should build conda with channel successfully", func() {
		exampleName := "python/conda"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		e.BuildImage(true)()
		e.RunContainer()()
		e.DestroyContainer()()
	})

	It("Should build conda with separate channel setting successfully", func() {
		exampleName := "python/conda_channel"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		e.BuildImage(true)()
		e.RunContainer()()
		e.DestroyContainer()()
	})

	Describe("Should build uv with Python successfully", func() {
		exampleName := "python/uv"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		BeforeAll(e.BuildImage(true))
		BeforeEach(e.RunContainer())
		It("Should have Python installed", func() {
			res, err := e.ExecRuntimeCommand("uv-python")
			Expect(err).To(BeNil())
			Expect(res).To(ContainSubstring("python"))
		})
		AfterEach(e.DestroyContainer())
	})

	Describe("Should build pixi with Python successfully", func() {
		exampleName := "python/pixi"
		e := e2e.NewExample(e2e.BuildContextDirWithName(exampleName), testcase)
		BeforeAll(e.BuildImage(true))
		BeforeEach(e.RunContainer())
		It("Should have Python and dependencies installed", func() {
			res, err := e.ExecRuntimeCommand("pixi-via")
			Expect(err).To(BeNil())
			Expect(res).To(ContainSubstring("via args parser"))
		})
		AfterEach(e.DestroyContainer())
	})
})
