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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	e2e "github.com/tensorchord/envd/e2e/v1"
	"github.com/tensorchord/envd/pkg/app"
	"github.com/tensorchord/envd/pkg/home"
)

var _ = Describe("get env command", func() {
	args := []string{
		"envd.test", "--debug", "envs", "list",
	}
	BeforeEach(func() {
		Expect(home.Initialize()).NotTo(HaveOccurred())
		envdApp := app.New()
		err := envdApp.Run([]string{"envd.test", "--debug", "bootstrap"})
		Expect(err).NotTo(HaveOccurred())
	})
	When("given the right arguments", func() {
		It("should get the environments successfully", func() {
			envdApp := app.New()
			e2e.ResetEnvdApp()
			err := envdApp.Run(args)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
