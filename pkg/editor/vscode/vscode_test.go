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

package vscode

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Visual Studio Code", func() {
	Describe("Plugin", func() {
		It("should get the latest version successfully", func() {
			url, err := GetLatestVersionURL(Plugin{
				Publisher: "redhat",
				Extension: "java",
			})
			Expect(err).To(BeNil())
			Expect(url).NotTo(Equal(""))
		})
		It("should be able to parse", func() {
			tcs := []struct {
				name              string
				expectedExtension string
				expectedVersion   string
				expectedPublisher string
				expectedErr       bool
			}{
				{
					name:              "ms-python.python-2021.12.1559732655",
					expectedPublisher: "ms-python",
					expectedExtension: "python",
					expectedVersion:   "2021.12.1559732655",
					expectedErr:       false,
				},
				{
					name:              "ms-vscode.cpptools-1.7.1",
					expectedPublisher: "ms-vscode",
					expectedExtension: "cpptools",
					expectedVersion:   "1.7.1",
					expectedErr:       false,
				},
				{
					name:              "github.copilot-1.12.5517",
					expectedPublisher: "github",
					expectedExtension: "copilot",
					expectedVersion:   "1.12.5517",
					expectedErr:       false,
				},
				{
					name:              "dbaeumer.vscode-eslint-1.1.1",
					expectedPublisher: "dbaeumer",
					expectedExtension: "vscode-eslint",
					expectedVersion:   "1.1.1",
					expectedErr:       false,
				},
				{
					name:              "dbaeumer.vscode-eslint",
					expectedPublisher: "dbaeumer",
					expectedExtension: "vscode-eslint",
					expectedVersion:   "",
					expectedErr:       false,
				},
				{
					name:        "test",
					expectedErr: true,
				},
			}
			for _, tc := range tcs {
				p, err := ParsePlugin(tc.name)
				if tc.expectedErr {
					Expect(err).To(HaveOccurred())
				} else {
					Expect(err).ToNot(HaveOccurred())
					Expect(p.Publisher).To(Equal(tc.expectedPublisher))
					Expect(p.Extension).To(Equal(tc.expectedExtension))
					if tc.expectedVersion != "" {
						Expect(p.Version).NotTo(BeNil())
						Expect(*p.Version).To(Equal(tc.expectedVersion))
					}
				}
			}
		})
	})
})
