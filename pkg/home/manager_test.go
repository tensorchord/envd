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

package home

import (
	"os"
	"path/filepath"

	"github.com/adrg/xdg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("home manager", func() {
	When("initialized", func() {
		It("should initialized successfully", func() {
			Expect(Initialize()).NotTo(HaveOccurred())
			m := GetManager()
			Expect(m.CacheDir()).To(Equal(filepath.Join(xdg.CacheHome, "envd")))
			Expect(m.ConfigFile()).To(Equal(filepath.Join(xdg.ConfigHome, "envd/config.envd")))
		})
		It("should return the cache status", func() {
			Expect(os.RemoveAll(filepath.Join(xdg.CacheHome, "envd/cache.status"))).NotTo(HaveOccurred())
			Expect(Initialize()).NotTo(HaveOccurred())
			m := GetManager()
			m.(*generalManager).cacheMap = make(map[string]bool)
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
