// Copyright 2022 The MIDI Authors
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
	"github.com/moby/buildkit/util/progress/progresswriter"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	mockbuildkitd "github.com/tensorchord/MIDI/pkg/buildkitd/mock"
	"github.com/tensorchord/MIDI/pkg/flag"
	"github.com/tensorchord/MIDI/pkg/home"
	mockstarlark "github.com/tensorchord/MIDI/pkg/lang/frontend/starlark/mock"
)

var _ = Describe("Builder", func() {
	Describe("building image", Label("buildkitd"), func() {
		var buildkitdSocket, configFilePath, manifestFilePath, tag string
		BeforeEach(func() {
			buildkitdSocket = "docker-container://midi_buildkitd"
			configFilePath = "testdata/config.MIDI"
			manifestFilePath = "testdata/build.MIDI"
			tag = "midi-dev:test"
			viper.Set(flag.FlagBuildkitdContainer, "midi_buildkitd")
			viper.Set(flag.FlagSSHImage, "midi-ssh:latest")
			os.Setenv("DOCKER_API_VERSION", "1.41")
			DeferCleanup(func() {
				viper.Set(flag.FlagBuildkitdContainer, "")
				viper.Set(flag.FlagSSHImage, "")
			})
		})
		When("getting the wrong builtkitd address", func() {
			buildkitdSocket = "wrong"
			viper.Set(flag.FlagBuildkitdContainer, buildkitdSocket)
			It("should return an error", func() {
				_, err := New(context.TODO(), configFilePath, manifestFilePath, tag)
				Expect(err).To(HaveOccurred())
			})
		})
		When("building the manifest", func() {
			var b *generalBuilder
			BeforeEach(func() {
				ctrl := gomock.NewController(GinkgoT())
				ctrlStarlark := gomock.NewController(GinkgoT())
				b = &generalBuilder{
					manifestFilePath: manifestFilePath,
					configFilePath:   configFilePath,
					progressMode:     "auto",
					tag:              tag,
					logger: logrus.WithFields(logrus.Fields{
						"tag": tag,
					}),
				}
				b.Client = mockbuildkitd.NewMockClient(ctrl)
				b.Interpreter = mockstarlark.NewMockInterpreter(ctrlStarlark)
				pw, err := progresswriter.NewPrinter(context.TODO(), os.Stdout, b.progressMode)
				if err != nil {
					Fail(err.Error())
				}
				b.Writer = pw
			})

			When("failed to interpret config", func() {
				It("should get an error", func() {
					expected := errors.New("failed to interpret config")
					b.Interpreter.(*mockstarlark.MockInterpreter).EXPECT().ExecFile(
						gomock.Eq(configFilePath),
					).Return(nil, expected)
					err := b.Build(context.TODO())
					Expect(err).To(HaveOccurred())
				})
			})

			When("failed to interpret manifest", func() {
				It("should get an error", func() {
					expected := errors.New("failed to interpret manifest")
					b.Interpreter.(*mockstarlark.MockInterpreter).EXPECT().ExecFile(
						gomock.Eq(configFilePath),
					).Return(nil, nil)
					b.Interpreter.(*mockstarlark.MockInterpreter).EXPECT().ExecFile(
						gomock.Eq(b.manifestFilePath),
					).Return(nil, expected)
					err := b.Build(context.TODO())
					Expect(err).To(HaveOccurred())
				})
			})
			It("should build successfully", func() {
				b.Interpreter.(*mockstarlark.MockInterpreter).EXPECT().ExecFile(
					gomock.Eq(configFilePath),
				).Return(nil, nil).Times(1)
				b.Interpreter.(*mockstarlark.MockInterpreter).EXPECT().ExecFile(
					gomock.Eq(b.manifestFilePath),
				).Return(nil, nil).Times(1)
				err := home.Initialize()
				Expect(err).ToNot(HaveOccurred())
				close(b.Writer.Status())

				b.Client.(*mockbuildkitd.MockClient).EXPECT().Solve(
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
					gomock.Any(),
				).Return(nil, nil).AnyTimes()
				err = b.Build(context.TODO())
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})
})
