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

	e2e "github.com/tensorchord/envd/e2e/v1"
	"github.com/tensorchord/envd/pkg/app"
	"github.com/tensorchord/envd/pkg/driver/docker"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var _ = Describe("up command", Ordered, func() {
	buildContext := "testdata/up-test"
	env := "up-test"
	baseArgs := []string{
		"envd.test", "--debug",
	}
	BeforeAll(func() {
		Expect(home.Initialize()).NotTo(HaveOccurred())
		envdApp := app.New()
		err := envdApp.Run(append(baseArgs, "bootstrap"))
		Expect(err).NotTo(HaveOccurred())
		_, err = docker.NewClient(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		c := types.Context{Runner: types.RunnerTypeDocker}
		opt := envd.Options{Context: &c}
		envdEngine, err := envd.New(context.TODO(), opt)
		Expect(err).NotTo(HaveOccurred())
		_, err = envdEngine.Destroy(context.TODO(), env)
		Expect(err).NotTo(HaveOccurred())

	})
	When("given the right arguments", func() {
		It("should up and destroy successfully", func() {
			args := append(baseArgs, []string{
				"up", "--path", buildContext, "--detach", "--force",
			}...)
			e2e.ResetEnvdApp()
			envdApp := app.New()
			err := envdApp.Run(args)
			Expect(err).NotTo(HaveOccurred())

			depsArgs := append(baseArgs, []string{
				"envs", "describe", "--env", env,
			}...)

			err = envdApp.Run(depsArgs)
			Expect(err).NotTo(HaveOccurred())

			destroyArgs := append(baseArgs, []string{
				"destroy", "--path", buildContext,
			}...)
			err = envdApp.Run(destroyArgs)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
