package shell

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var _ = Describe("zsh manager", func() {
	zshManager := NewManager()
	BeforeEach(func() {
		Expect(home.Initialize()).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		// Cleanup the home cache.
		Expect(home.Initialize()).NotTo(HaveOccurred())
		m := home.GetManager()
		Expect(fileutil.RemoveAll(m.CacheDir())).NotTo(HaveOccurred())
	})
	When("cached", func() {
		It("should skip", func() {
			err := home.GetManager().MarkCache(cacheKey, true)
			Expect(err).NotTo(HaveOccurred())
			cached, err := zshManager.DownloadOrCache()
			Expect(cached).To(BeTrue())
			Expect(err).NotTo(HaveOccurred())
		})
	})
	When("not cached", func() {
		It("should download", func() {
			err := home.GetManager().MarkCache(cacheKey, false)
			Expect(err).NotTo(HaveOccurred())
			cached, err := zshManager.DownloadOrCache()
			Expect(err).NotTo(HaveOccurred())
			Expect(cached).To(BeFalse())
			exists, err := fileutil.DirExists(filepath.Join(home.GetManager().CacheDir(), "oh-my-zsh"))
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})
	})
})
