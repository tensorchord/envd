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
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/home"
)

var _ = Describe("up command", func() {
	buildContext := "testdata/up-test"
	env := "up-test"
	baseArgs := []string{
		"envd.test", "--debug",
	}
	BeforeEach(func() {
		Expect(home.Initialize()).NotTo(HaveOccurred())
		app := New()
		err := app.Run(append(baseArgs, "bootstrap"))
		Expect(err).NotTo(HaveOccurred())
		cli, err := docker.NewClient(context.TODO())
		Expect(err).NotTo(HaveOccurred())
		_, err = cli.Destroy(context.TODO(), env)
		Expect(err).NotTo(HaveOccurred())
	})
	When("given the right arguments", func() {
		It("should up and destroy successfully", func() {
			args := append(baseArgs, []string{
				"up", "--path", buildContext, "--detach",
			}...)
			app := New()
			err := app.Run(args)
			Expect(err).NotTo(HaveOccurred())

			depsArgs := append(baseArgs, []string{
				"get", "envs", "deps", "--env", env,
			}...)

			err = app.Run(depsArgs)
			Expect(err).NotTo(HaveOccurred())

			destroyArgs := append(baseArgs, []string{
				"destroy", "--path", buildContext,
			}...)
			err = app.Run(destroyArgs)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
