package home

import (
	"path/filepath"

	"github.com/adrg/xdg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var _ = Describe("home manager", func() {
	BeforeEach(func() {
		// Cleanup the home cache.
		Expect(Initialize()).NotTo(HaveOccurred())
		m := GetManager()
		Expect(fileutil.RemoveAll(m.CacheDir())).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		// Cleanup the home cache.
		Expect(Initialize()).NotTo(HaveOccurred())
		m := GetManager()
		Expect(fileutil.RemoveAll(m.CacheDir())).NotTo(HaveOccurred())
	})
	When("initialized", func() {
		It("should initialized successfully", func() {
			Expect(Initialize()).NotTo(HaveOccurred())
			m := GetManager()
			Expect(m.CacheDir()).To(Equal(filepath.Join(xdg.CacheHome, "envd")))
			Expect(m.ConfigFile()).To(Equal(filepath.Join(xdg.ConfigHome, "envd/config.envd")))
		})
		It("should return the cache status", func() {
			Expect(Initialize()).NotTo(HaveOccurred())
			m := GetManager()
			Expect(m.Cached("test")).To(BeFalse())
			Expect(m.MarkCache("test", true)).To(Succeed())
			Expect(m.Cached("test")).To(BeTrue())
			// Restart the init process, the cache should be persistent.
			Expect(Initialize()).NotTo(HaveOccurred())
			m = GetManager()
			Expect(m.Cached("test")).To(BeTrue())
		})
	})
})
