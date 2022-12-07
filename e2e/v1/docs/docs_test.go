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

package docs

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	e2e "github.com/tensorchord/envd/e2e/v1"
	"github.com/tensorchord/envd/pkg/app"
	"github.com/tensorchord/envd/pkg/home"
)

var _ = Describe("check examples in documentation", func() {
	baseArgs := []string{
		"envd.test", "--debug",
	}

	BeforeEach(func() {
		Expect(home.Initialize()).NotTo(HaveOccurred())
		e2e.ResetEnvdApp()
		envdApp := app.New()
		err := envdApp.Run(append(baseArgs, "bootstrap"))
		Expect(err).NotTo(HaveOccurred())
	})

	It("can list envd envs", func() {
		envdApp := app.New()
		err := envdApp.Run([]string{"envd.test", "--debug", "envs", "list"})
		Expect(err).NotTo(HaveOccurred())
	})

	It("can check envd envs details", func() {
		buildContext := "testdata/minimal"
		args := append(baseArgs, []string{
			"up", "--path", buildContext, "--detach", "--force",
		}...)
		e2e.ResetEnvdApp()
		envdApp := app.New()
		err := envdApp.Run(args)
		Expect(err).NotTo(HaveOccurred())

		err = envdApp.Run([]string{"envd.test", "--debug", "envs", "describe", "--env", "minimal"})
		Expect(err).NotTo(HaveOccurred())

		destroyArgs := append(baseArgs, []string{
			"destroy", "--path", buildContext,
		}...)
		err = envdApp.Run(destroyArgs)
		Expect(err).NotTo(HaveOccurred())
	})

	upTests := []string{"testdata/minimal", "testdata/getting_started", "testdata/jupyter"}

	for _, v := range upTests {
		It(fmt.Sprintf("can up %s environment", v), func() {
			args := append(baseArgs, []string{
				"up", "--path", v, "-f", "build.envd", "--detach", "--force",
			}...)
			e2e.ResetEnvdApp()
			envdApp := app.New()
			err := envdApp.Run(args)
			Expect(err).NotTo(HaveOccurred())

			destroyArgs := append(baseArgs, []string{
				"destroy", "--path", v,
			}...)
			err = envdApp.Run(destroyArgs)
			Expect(err).NotTo(HaveOccurred())
		})
	}

	It("should be able to use envdlib", func() {
		path := "testdata/envdlib"
		args := append(baseArgs, []string{
			"up", "--path", path, "-f", "build.envd", "--detach", "--force",
		}...)
		e2e.ResetEnvdApp()
		envdApp := app.New()
		err := envdApp.Run(args)
		Expect(err).To(Succeed())

		// check the port
		time.Sleep(time.Second * 2)
		resp, err := http.Get("http://127.0.0.1:8888")
		Expect(err).To(Succeed())
		defer resp.Body.Close()
		Expect(resp.StatusCode).To(Equal(200))

		destroyArgs := append(baseArgs, []string{
			"destroy", "--path", path,
		}...)
		err = envdApp.Run(destroyArgs)
		Expect(err).To(Succeed())
	})
})
