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
	"os"

	"github.com/cockroachdb/errors"
	"github.com/golang/mock/gomock"
	"github.com/moby/buildkit/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	mockbuildkitd "github.com/tensorchord/envd/pkg/buildkitd/mock"
	"github.com/tensorchord/envd/pkg/home"
	mockstarlark "github.com/tensorchord/envd/pkg/lang/frontend/starlark/mock"
	v0 "github.com/tensorchord/envd/pkg/lang/ir/v0"
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
				pub, err := sshconfig.GetPublicKey()
				Expect(err).NotTo(HaveOccurred())
				b = &generalBuilder{
					Options: Options{
						ManifestFilePath: manifestFilePath,
						ConfigFilePath:   configFilePath,
						ProgressMode:     "plain",
						Tag:              tag,
						BuildFuncName:    "build",
						PubKeyPath:       pub,
					},
					logger: logrus.WithFields(logrus.Fields{
						"tag": tag,
					}),
				}
				b.Client = mockbuildkitd.NewMockClient(ctrl)
				b.Interpreter = mockstarlark.NewMockInterpreter(ctrlStarlark)

				ctrlWriter := gomock.NewController(GinkgoT())
				w = compileuimock.NewMockWriter(ctrlWriter)
				v0.DefaultGraph.SetWriter(w)
			})

			When("build error", func() {
				It("should get an error", func() {
					b.entries = []client.ExportEntry{
						{
							Type: client.ExporterDocker,
						},
					}

					b.Client.(*mockbuildkitd.MockClient).EXPECT().Build(gomock.Any(),
						gomock.Any(), gomock.Eq("envd"), gomock.Any(), gomock.Any()).
						Return(nil, errors.New("build error"))

					pw, err := progresswriter.NewPrinter(context.TODO(), os.Stdout, b.ProgressMode)
					Expect(err).NotTo(HaveOccurred())

					close(pw.Status())
					err = b.build(context.TODO(), pw)
					Expect(err).To(HaveOccurred())
				})
			})

			It("should build successfully", func() {
				err := home.Initialize()
				Expect(err).ToNot(HaveOccurred())

				b.entries = []client.ExportEntry{
					{
						Type: client.ExporterDocker,
					},
				}

				b.Client.(*mockbuildkitd.MockClient).EXPECT().Build(gomock.Any(),
					gomock.Any(), gomock.Eq("envd"), gomock.Any(), gomock.Any()).
					Return(nil, nil)

				pw, err := progresswriter.NewPrinter(context.TODO(), os.Stdout, b.ProgressMode)
				Expect(err).NotTo(HaveOccurred())

				close(pw.Status())
				err = b.build(context.TODO(), pw)
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
