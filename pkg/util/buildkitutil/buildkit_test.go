// Copyright 2023 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package buildkitutil

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildkitWithRegistry(t *testing.T) {
	testCases := []struct {
		bc       BuildkitConfig
		expected string
	}{
		{
			BuildkitConfig{
				Registry: "registry.example.com",
				Mirror:   "https://mirror.example.com",
				UseHTTP:  true,
			},
			`
[registry."registry.example.com"]
  mirrors = ["https://mirror.example.com"]
  http = true
`,
		},
		{
			BuildkitConfig{
				Registry: "registry.example.com",
				SetCA:    true,
			},
			`
[registry."registry.example.com"]
  http = false
  ca=["/etc/registry/ca.pem"]
  [[registry."registry.example.com".keypair]]
	key="/etc/registry/key.pem"
	cert="/etc/registry/cert.pem"
`,
		},
		{
			BuildkitConfig{},
			`
[registry."docker.io"]
  http = false
			`,
		},
	}

	for _, tc := range testCases {
		config, err := tc.bc.String()
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(tc.expected), strings.TrimSpace(config))
	}
}
