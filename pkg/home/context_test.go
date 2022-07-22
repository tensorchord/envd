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

var _ = Describe("context test", func() {
	defaultContext := "default"
	testContext := "envd_home_test"
	testSocket := "0.0.0.0:12345"
	testBuilder := types.BuilderTypeTCP

	BeforeEach(func() {
		Expect(Initialize()).To(Succeed())
	})

	Describe("create with use", Ordered, func() {
		BeforeAll(func() {
			err := GetManager().ContextCreate(testContext, testBuilder, testSocket, true)
			Expect(err).NotTo(HaveOccurred())
		})

		It("current context should be the new one", func() {
			contexts, err := GetManager().ContextList()
			Expect(err).NotTo(HaveOccurred())
			Expect(contexts.Current).To(Equal(testContext))
			builder, socket, err := GetManager().ContextGetCurrent()
			Expect(err).NotTo(HaveOccurred())
			Expect(builder).To(Equal(testBuilder))
			Expect(socket).To(Equal(testSocket))
		})

		It("cannot delete the current context", func() {
			err := GetManager().ContextRemove(testContext)
			Expect(err).To(HaveOccurred())
		})

		AfterAll(func() {
			Expect(GetManager().ContextUse(defaultContext)).To(Succeed())
			contexts, err := GetManager().ContextList()
			Expect(err).NotTo(HaveOccurred())
			Expect(contexts.Current).To(Equal(defaultContext))
			Expect(GetManager().ContextRemove(testContext)).To(Succeed())
		})
	})

	Describe("create without use", Ordered, func() {
		BeforeAll(func() {
			err := GetManager().ContextCreate(testContext, testBuilder, testSocket, false)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not be able to create the same context", func() {
			err := GetManager().ContextCreate(testContext, testBuilder, testSocket, false)
			Expect(err).To(HaveOccurred())
		})

		It("should use the default context", func() {
			contexts, err := GetManager().ContextList()
			Expect(err).NotTo(HaveOccurred())
			Expect(contexts.Current).To(Equal(defaultContext))
		})

		AfterAll(func() {
			Expect(GetManager().ContextRemove(testContext)).To(Succeed())
		})
	})
})
