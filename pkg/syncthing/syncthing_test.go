package syncthing

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSyncthing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Syncthing Suite")
}

var _ = Describe("Syncthing", func() {
	BeforeEach(func() {
		if _, err := os.Stat(getSyncthingBinPath()); err == nil {
			err := os.Remove(getSyncthingBinPath())
			Expect(err).To(BeNil())
		}
	})
	Describe("Syncthing", func() {
		It("Starts and stops syncthing", func() {
			s := Syncthing{}
			err := s.Start()
			Expect(err).To(BeNil())

			Expect(s.IsRunning()).To(BeTrue())

			err = s.Stop()
			Expect(err).To(BeNil())
		})

	})

})

var _ = Describe("Syncthing Install", func() {
	BeforeEach(func() {
		if _, err := os.Stat(getSyncthingBinPath()); err == nil {
			err := os.Remove(getSyncthingBinPath())
			Expect(err).To(BeNil())
		}
	})
	Describe("Install", func() {

		It("Installs binary in cache directory", func() {
			err := InstallSyncthing()
			Expect(err).To(BeNil())

			Expect(IsInstalled()).To(BeTrue())
		})
	})

})
