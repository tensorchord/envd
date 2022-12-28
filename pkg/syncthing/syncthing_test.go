package syncthing_test

import (
	"fmt"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tensorchord/envd/pkg/syncthing"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func TestSyncthing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Syncthing Suite")
}

var _ = Describe("Syncthing", func() {
	BeforeEach(func() {
		if _, err := os.Stat(syncthing.GetSyncthingBinPath()); err == nil {
			err := os.Remove(syncthing.GetSyncthingBinPath())
			Expect(err).To(BeNil())
		}
	})

	Describe("Syncthing", func() {
		It("Starts and stops syncthing", func() {
			s := syncthing.Syncthing{}
			err := s.Start()
			Expect(err).To(BeNil())

			Expect(s.IsRunning()).To(BeTrue())

			err = s.Stop()
			Expect(err).To(BeNil())
		})
	})

	Describe("Syncthing config", func() {
		BeforeEach(func() {
			os.RemoveAll(syncthing.DefaultHomeDirectory())
		})

		AfterEach(func() {
			os.RemoveAll(syncthing.DefaultHomeDirectory())
		})

		It("Initializes syncthing configuration", func() {
			s, err := syncthing.InitializeLocalSyncthing()
			Expect(err).To(BeNil())

			Expect(s.Port).To(Equal(syncthing.DefaultLocalPort))
			Expect(s.Config.GUI.Address()).To(Equal(fmt.Sprintf("0.0.0.0:%s", syncthing.DefaultLocalPort)))

			dirExists, err := fileutil.DirExists(s.HomeDirectory)
			Expect(err).To(BeNil())
			Expect(dirExists).To(BeTrue())

			configFilePath := syncthing.GetConfigFilePath(s.HomeDirectory)
			fileExists, err := fileutil.FileExists(configFilePath)
			Expect(err).To(BeNil())
			Expect(fileExists).To(BeTrue())

		})
	})
})

var _ = Describe("Syncthing Install", func() {
	BeforeEach(func() {
		if _, err := os.Stat(syncthing.GetSyncthingBinPath()); err == nil {
			err := os.Remove(syncthing.GetSyncthingBinPath())
			Expect(err).To(BeNil())
		}
	})
	Describe("Install", func() {

		It("Installs binary in cache directory", func() {
			err := syncthing.InstallSyncthing()
			Expect(err).To(BeNil())

			Expect(syncthing.IsInstalled()).To(BeTrue())
		})
	})

})
