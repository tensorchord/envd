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

package config

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ssh config", func() {
	When("giving a empty ssh config", func() {
		It("Should add/remove the config successfully", func() {
			env := "test-ssh-config"
			iface := "localhost"
			port := 8888
			keyPath := "key"
			eo := EntryOptions{
				Name:               BuildHostname(env),
				IFace:              iface,
				Port:               port,
				PrivateKeyPath:     keyPath,
				EnableHostKeyCheck: false,
				EnableAgentForward: true,
			}
			err := add(getSSHConfigPath(), eo)
			Expect(err).NotTo(HaveOccurred())

			actual, err := GetPort(env)
			Expect(err).NotTo(HaveOccurred())
			Expect(actual).To(Equal(port))

			err = remove(getSSHConfigPath(), env)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
