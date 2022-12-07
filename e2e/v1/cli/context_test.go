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

package cli

import (
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	e2e "github.com/tensorchord/envd/e2e/v1"
	"github.com/tensorchord/envd/pkg/app"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var _ = Describe("home context", func() {
	defaultContext := "default"
	BeforeEach(func() {
		Expect(home.Initialize()).To(Succeed())
	})
	When("check the default context", func() {
		It("should found the default context", func() {
			contexts, err := home.GetManager().ContextList()
			Expect(err).NotTo(HaveOccurred())
			Expect(contexts.Current).To(Equal(defaultContext))
		})
	})

	Describe("add a new context", Ordered, func() {
		testContext := "envd_home_test"
		testBuilderAddress := "0.0.0.0:12345"
		testBuilder := types.BuilderTypeTCP
		testRunner := types.RunnerTypeEnvdServer
		testRunnerAddress := "http://localhost"
		c := types.Context{
			Name:           testContext,
			Builder:        testBuilder,
			BuilderAddress: testBuilderAddress,
			Runner:         testRunner,
			RunnerAddress:  &testRunnerAddress,
		}

		BeforeAll(func() {
			err := home.GetManager().ContextCreate(c, true)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should find a new context", func() {
			contexts, err := home.GetManager().ContextList()
			Expect(err).NotTo(HaveOccurred())
			Expect(contexts.Current).To(Equal(testContext))
		})

		Describe("connect buildkit through TCP", Ordered, func() {
			name := "envd-buildkitd-tcp-test"
			buildContext := "testdata/build-test"
			dockerArgs := []string{
				"run", "-d", "-p", "12345:8000", "--rm", "--name", name,
				"--security-opt", "seccomp=unconfined", "--security-opt", "apparmor=unconfined",
				"moby/buildkit:rootless", "--addr", "tcp://0.0.0.0:8000", "--oci-worker-no-process-sandbox",
			}
			BeforeAll(func() {
				cmd := exec.Command("docker", dockerArgs...)
				err := cmd.Run()
				Expect(err).ToNot(HaveOccurred())
			})

			It("should be able to build image with TCP context", func() {
				args := []string{"envd.test", "--debug", "build", "--path", buildContext}
				envdApp := app.New()
				e2e.ResetEnvdApp()
				err := envdApp.Run(args)
				Expect(err).NotTo(HaveOccurred())
			})

			AfterAll(func() {
				cmd := exec.Command("docker", "stop", name)
				err := cmd.Run()
				Expect(err).ToNot(HaveOccurred())
			})
		})

		It("fail to delete the current context", func() {
			err := home.GetManager().ContextRemove(testContext)
			Expect(err).To(HaveOccurred())
		})

		AfterAll(func() {
			err := home.GetManager().ContextUse(defaultContext)
			Expect(err).NotTo(HaveOccurred())
			err = home.GetManager().ContextRemove(testContext)
			Expect(err).NotTo(HaveOccurred())
			contexts, err := home.GetManager().ContextList()
			Expect(err).NotTo(HaveOccurred())
			Expect(contexts.Current).To(Equal(defaultContext))
		})
	})
})
