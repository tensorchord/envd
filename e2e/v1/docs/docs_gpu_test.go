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

package docs

import (
	. "github.com/onsi/ginkgo/v2"

	e2e "github.com/tensorchord/envd/e2e/v1"
)

var _ = Describe("check GPU examples in documentation", Ordered, func() {
	e := e2e.NewExample(e2e.BuildContextDirWithName("complex"), "e2e-doc")
	It("should be able to build the GPU example", e.BuildImage(true))
	AfterAll(e.DestroyContainer())
})
