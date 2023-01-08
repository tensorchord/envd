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
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	e2e "github.com/tensorchord/envd/e2e/v0"
	"github.com/tensorchord/envd/pkg/app"
	"github.com/tensorchord/envd/pkg/driver/docker"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var _ = Describe("build command", Ordered, func() {
	buildTestName := "testdata/build-test"
	customImageTestName := "testdata/custom-image-test"
	When("given the right arguments", func() {
		It("should build successfully", func() {
			envdApp := app.New()
			e2e.ResetEnvdApp()
			err := envdApp.Run([]string{"envd.test", "--debug", "bootstrap"})
			Expect(err).NotTo(HaveOccurred())
			_, err = docker.NewClient(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			c := types.Context{Runner: types.RunnerTypeDocker}
			opt := envd.Options{Context: &c}
			envdEngine, err := envd.New(context.TODO(), opt)
			Expect(err).NotTo(HaveOccurred())
			_, err = envdEngine.Destroy(context.TODO(), buildTestName)
			Expect(err).NotTo(HaveOccurred())
			envdApp = app.New()
			args := []string{
				"envd.test", "--debug", "build", "--path", buildTestName,
			}
			err = envdApp.Run(args)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	When("given the custom image", func() {
		It("should build successfully", func() {
			envdApp := app.New()
			e2e.ResetEnvdApp()
			err := envdApp.Run([]string{"envd.test", "--debug", "bootstrap"})
			Expect(err).NotTo(HaveOccurred())
			_, err = docker.NewClient(context.TODO())
			Expect(err).NotTo(HaveOccurred())
			c := types.Context{Runner: types.RunnerTypeDocker}
			opt := envd.Options{Context: &c}
			envdEngine, err := envd.New(context.TODO(), opt)
			Expect(err).NotTo(HaveOccurred())
			_, err = envdEngine.Destroy(context.TODO(), buildTestName)
			Expect(err).NotTo(HaveOccurred())
			envdApp = app.New()
			args := []string{
				"envd.test", "--debug", "build", "--path", customImageTestName,
			}
			err = envdApp.Run(args)
			Expect(err).NotTo(HaveOccurred())
		})
	})
	AfterAll(func() {
		Expect(home.Initialize()).NotTo(HaveOccurred())
		envdApp := app.New()
		e2e.ResetEnvdApp()
		err := envdApp.Run([]string{"envd.test", "--debug", "bootstrap"})
		Expect(err).NotTo(HaveOccurred())
		_, err = docker.NewClient(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		c := types.Context{Runner: types.RunnerTypeDocker}
		opt := envd.Options{Context: &c}
		envdEngine, err := envd.New(context.TODO(), opt)
		Expect(err).NotTo(HaveOccurred())
		_, err = envdEngine.Destroy(context.TODO(), buildTestName)
		Expect(err).NotTo(HaveOccurred())
		// Init DefaultGraph.
		// ir.DefaultGraph = ir.NewGraph()
	})
})
