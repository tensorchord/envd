package builder

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/spf13/viper"
	"github.com/tensorchord/MIDI/pkg/flag"
)

var _ = Describe("Builder", func() {
	Describe("building", Label("buildkitd"), func() {
		var buildkitdSocket, configFilePath, manifestFilePath, tag string
		BeforeEach(func() {
			buildkitdSocket = "docker-container://midi_buildkitd"
			configFilePath = "testdata/config.MIDI"
			manifestFilePath = "testdata/build.MIDI"
			tag = "midi-dev:test"
		})
		When("getting the wrong builtkitd address", func() {
			buildkitdSocket = "wrong"
			viper.Set(flag.FlagBuildkitdContainer, buildkitdSocket)
			It("should return an error", func() {
				_, err := New(context.TODO(), configFilePath, manifestFilePath, tag)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
