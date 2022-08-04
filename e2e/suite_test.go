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

package e2e

import (
	"os/exec"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/tensorchord/envd/pkg/version"
)

func init() {
	// Set the git tag to get the correct image in ir package.
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	tag, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	version.SetGitTagForE2ETest(string(tag))
}

func TestMain(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "envd Suite")
}
