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
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	e2e "github.com/tensorchord/envd/e2e/v1"
	"github.com/tensorchord/envd/pkg/app"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var _ = Describe("init project", Ordered, func() {
	var path string
	BeforeAll(func() {
		Expect(home.Initialize()).To(Succeed())
		envdApp := app.New()
		err := envdApp.Run([]string{"envd.test", "--debug", "bootstrap"})
		Expect(err).To(Succeed())
		e2e.ResetEnvdApp()
		path, err = os.MkdirTemp("", "envd_init_test_*")
		Expect(err).To(Succeed())
		err = os.WriteFile(filepath.Join(path, "requirements.txt"), []byte("via"), 0666)
		Expect(err).To(Succeed())
	})

	It("init python env", func() {
		envdApp := app.New()
		err := envdApp.Run([]string{"envd.test", "--debug", "init", "-p", path})
		Expect(err).To(Succeed())
		exist, err := fileutil.FileExists(filepath.Join(path, "build.envd"))
		Expect(err).To(Succeed())
		Expect(exist).To(BeTrue())
	})

	Describe("run init env", Ordered, func() {
		var e *e2e.Example
		BeforeAll(func() {
			// have to use `path` inside ginkgo closure
			e = e2e.NewExample(path, "init_test")
			e.RunContainer()()
		})
		It("exec installed command inside container", func() {
			_, err := e.Exec("via --help")
			Expect(err).To(Succeed())
			e.DestroyContainer()
		})
	})

	AfterAll(func() {
		os.RemoveAll(path)
	})
})
