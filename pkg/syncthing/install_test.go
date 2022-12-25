package syncthing

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestSyncThing(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Books Suite")
}

var _ = Describe("Syncthing", func() {
	BeforeEach(func() {
		if _, err := os.Stat(getSyncthingBinPath()); err == nil {
			err := os.Remove(getSyncthingBinPath())
			Expect(err).To(BeNil())
		}
	})

	Describe("Install", func() {
		It("Installs binary in cache directory", func() {
			err := Install()
			Expect(err).To(BeNil())

			Expect(IsInstalled()).To(BeTrue())
		})
	})
})
