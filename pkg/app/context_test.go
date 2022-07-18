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

package app

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var _ = Describe("home context", func() {
	defaultContext := "default"
	BeforeEach(func() {
		Expect(home.Initialize()).NotTo(HaveOccurred())
	})
	When("check the default context", func() {
		It("should found the default context", func() {
			contexts, err := home.GetManager().ContextList()
			Expect(err).NotTo(HaveOccurred())
			Expect(contexts.Current).To(Equal(defaultContext))
		})
	})

	Describe("add a new context", Ordered, func() {
		Context("TCP context", func() {
			contextName := "tcp-test"
			builder := types.BuilderTypeTCP
			socketAddr := "0.0.0.0:12345"

			BeforeAll(func() {
				err := home.GetManager().ContextCreate(contextName, builder, socketAddr, true)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterAll(func() {
				err := home.GetManager().ContextUse(defaultContext)
				Expect(err).NotTo(HaveOccurred())
				err = home.GetManager().ContextRemove(contextName)
				Expect(err).NotTo(HaveOccurred())
				contexts, err := home.GetManager().ContextList()
				Expect(err).NotTo(HaveOccurred())
				Expect(contexts.Current).To(Equal(defaultContext))
			})

			It("should found a new context", func() {
				contexts, err := home.GetManager().ContextList()
				Expect(err).NotTo(HaveOccurred())
				Expect(contexts.Current).To(Equal(contextName))
			})

			When("connect buildkit through TCP", func() {
				name := "envd-buildkitd-tcp-test"
				BeforeEach(func() {
					cmd := exec.Command("docker", "run", "-d", "-p", "12345:8000", "--name", name, "--security-opt", "seccomp=unconfined", "--security-opt", "apparmor=unconfined", "moby/buildkit:rootless", "--addr", "tcp://0.0.0.0:8000", "--oci-worker-no-process-sandbox")
					err := cmd.Run()
					Expect(err).ToNot(HaveOccurred())
				})

				It("should be able to build image with TCP context", func() {
					args := []string{"envd.test", "--debug", "build"}
					app := New()
					err := app.Run(args)
					Expect(err).NotTo(HaveOccurred())
				})

				AfterEach(func() {
					cmd := exec.Command("docker", "stop", name)
					err := cmd.Run()
					Expect(err).ToNot(HaveOccurred())
				})
			})

			It("fail to delete the current context", func() {
				err := home.GetManager().ContextRemove(contextName)
				Expect(err).To(HaveOccurred())
			})

			It("switch to default context to delete test context", func() {
				err := home.GetManager().ContextUse(defaultContext)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
