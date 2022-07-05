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

package builder

import (
	"context"
	"errors"
	"os"

	"github.com/golang/mock/gomock"
	"github.com/moby/buildkit/client/llb"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	mockbuildkitd "github.com/tensorchord/envd/pkg/buildkitd/mock"
	"github.com/tensorchord/envd/pkg/home"
	mockstarlark "github.com/tensorchord/envd/pkg/lang/frontend/starlark/mock"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	compileuimock "github.com/tensorchord/envd/pkg/progress/compileui/mock"
	"github.com/tensorchord/envd/pkg/progress/progresswriter"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
)

var _ = Describe("Builder", func() {
	Describe("building image", Label("buildkitd"), func() {
		var configFilePath, manifestFilePath, tag string
		BeforeEach(func() {
			configFilePath = "config.envd"
			manifestFilePath = "build.envd"
			tag = "envd-dev:test"
			Expect(home.Initialize()).NotTo(HaveOccurred())
		})
		When("building the manifest", func() {
			var b *generalBuilder
			var w compileui.Writer
			BeforeEach(func() {
				ctrl := gomock.NewController(GinkgoT())
				ctrlStarlark := gomock.NewController(GinkgoT())
				b = &generalBuilder{
					manifestFilePath: manifestFilePath,
					configFilePath:   configFilePath,
					progressMode:     "auto",
					tag:              tag,
					buildfuncname:    "build",
					logger: logrus.WithFields(logrus.Fields{
						"tag": tag,
					}),
				}
				b.Client = mockbuildkitd.NewMockClient(ctrl)
				b.Interpreter = mockstarlark.NewMockInterpreter(ctrlStarlark)

				ctrlWriter := gomock.NewController(GinkgoT())
				w = compileuimock.NewMockWriter(ctrlWriter)
				ir.DefaultGraph.Writer = w
			})

			When("failed to interpret config", func() {
				It("should get an error", func() {
					expected := errors.New("failed to interpret config")
					b.Interpreter.(*mockstarlark.MockInterpreter).EXPECT().ExecFile(
						gomock.Eq(configFilePath), "",
					).Return(nil, expected)
					pub := sshconfig.GetPublicKey()
					err := b.Build(context.TODO(), pub)
					Expect(err).To(HaveOccurred())
				})
			})

			When("failed to interpret manifest", func() {
				It("should get an error", func() {
					expected := errors.New("failed to interpret manifest")
					pub := sshconfig.GetPublicKey()
					b.Interpreter.(*mockstarlark.MockInterpreter).EXPECT().ExecFile(
						gomock.Eq(configFilePath), gomock.Eq(""),
					).Return(nil, nil)
					b.Interpreter.(*mockstarlark.MockInterpreter).EXPECT().ExecFile(
						gomock.Eq(b.manifestFilePath), gomock.Eq("build"),
					).Return(nil, expected)
					err := b.Build(context.TODO(), pub)
					Expect(err).To(HaveOccurred())
				})
			})
			It("should build successfully", func() {
				err := home.Initialize()
				Expect(err).ToNot(HaveOccurred())

				b.Client.(*mockbuildkitd.MockClient).EXPECT().Solve(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, nil).AnyTimes()

				var def *llb.Definition
				pw, err := progresswriter.NewPrinter(context.TODO(), os.Stdout, b.progressMode)
				if err != nil {
					Fail(err.Error())
				}
				close(pw.Status())
				err = b.build(context.TODO(), def, pw)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
