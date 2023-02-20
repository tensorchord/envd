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

package docker

import (
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("docker", func() {
	When("given the a lowercase tag", func() {
		It("should return the tag identically", func() {
			tag := "test:test"
			newTag, err := NormalizeName(tag)
			Expect(err).NotTo(HaveOccurred())
			Expect(newTag).To(Equal(tag))
		})
	})
	When("given the a uppercase tag", func() {
		It("should return the tag lowcased", func() {
			tag := "Test:test"
			newTag, err := NormalizeName(tag)
			Expect(err).NotTo(HaveOccurred())
			Expect(newTag).NotTo(Equal(tag))
			Expect(newTag).To(Equal(strings.ToLower(tag)))
		})
	})
})
