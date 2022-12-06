// Copyright 2022 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package home

import (
	"github.com/tensorchord/envd/pkg/types"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

var _ = Describe("auth test", Ordered, func() {
	defaultAuthName := "auth_name"
	deafultJWTToken := "default_token"
	ac := types.AuthConfig{
		Name:     defaultAuthName,
		JWTToken: deafultJWTToken,
	}

	BeforeEach(func() {
		Expect(Initialize()).To(Succeed())
	})

	Describe("create with use", func() {
		BeforeAll(func() {
			err := GetManager().AuthCreate(ac, true)
			Expect(err).NotTo(HaveOccurred())
		})

		It("current Auth config should be the new one", func() {
			authConf, err := GetManager().AuthGetCurrent()
			Expect(err).NotTo(HaveOccurred())
			Expect(authConf.Name).To(Equal(ac.Name))
		})

		AfterAll(func() {
			Expect(GetManager().AuthUse(defaultAuthName)).To(Succeed())
			auth, err := GetManager().AuthGetCurrent()
			Expect(err).NotTo(HaveOccurred())
			Expect(auth.Name).To(Equal(defaultAuthName))
		})
	})

	Describe("create without use", func() {
		BeforeAll(func() {
			err := GetManager().AuthCreate(ac, false)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should not be able to create the same context", func() {
			err := GetManager().AuthCreate(ac, false)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should use the default authConfig", func() {
			_, err := GetManager().AuthGetCurrent()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
