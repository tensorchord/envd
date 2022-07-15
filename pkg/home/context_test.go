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

package home

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tensorchord/envd/pkg/types"
)

var _ = Describe("home context", func() {
	defaultContext := "default"
	BeforeEach(func() {
		Expect(Initialize()).NotTo(HaveOccurred())
	})
	When("check the default context", func() {
		It("should found the default context", func() {
			contexts, err := GetManager().ContextList()
			Expect(err).NotTo(HaveOccurred())
			Expect(contexts.Current).To(Equal(defaultContext))
		})
	})

	Describe("add a new context", Ordered, func() {
		Context("TCP context", func() {
			contextName := "tcp-test"
			builder := types.BuilderTypeTCP
			socketAddr := "0.0.0.0:1234"

			It("can create a new context", func() {
				err := GetManager().ContextCreate(contextName, builder, socketAddr, true)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should found a new context", func() {
				contexts, err := GetManager().ContextList()
				Expect(err).NotTo(HaveOccurred())
				Expect(contexts.Current).To(Equal(contextName))
				Expect(len(contexts.Contexts)).To(Equal(2))
			})

			It("fail to delete the current context", func() {
				err := GetManager().ContextRemove(contextName)
				Expect(err).To(HaveOccurred())
			})

			It("switch to default context to delete test context", func() {
				err := GetManager().ContextUse(defaultContext)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should be the default context", func() {
				contexts, err := GetManager().ContextList()
				Expect(err).NotTo(HaveOccurred())
				Expect(contexts.Current).To(Equal(defaultContext))
			})

			It("can delete another context", func() {
				err := GetManager().ContextRemove(contextName)
				Expect(err).NotTo(HaveOccurred())
			})
			
			It("should have only one context", func() {
				contexts, err := GetManager().ContextList()
				Expect(err).NotTo(HaveOccurred())
				Expect(len(contexts.Contexts)).To(Equal(1))
			})
		})
	})
})
