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

package shell

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var _ = Describe("zsh manager", Serial, func() {
	zshManager := NewManager()
	BeforeEach(func() {
		Expect(home.Initialize()).NotTo(HaveOccurred())
	})
	AfterEach(func() {
		// Cleanup the home cache.
		Expect(home.Initialize()).NotTo(HaveOccurred())
		Expect(os.RemoveAll(filepath.Join(fileutil.DefaultCacheDir, "cache.status"))).NotTo(HaveOccurred())
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
