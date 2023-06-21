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
				Registry: []string{"registry.example.com"},
				Mirror:   "https://mirror.example.com",
				UseHTTP:  true,
			},
			`
[registry]
  [registry."registry.example.com"]
    http = true
    mirrors = ["https://mirror.example.com"]
`,
		},
		{
			BuildkitConfig{
				Registry: []string{"registry.example.com"},
				Ca:       []string{"/etc/registry/ca.pem"},
				Key:      []string{"/etc/registry/key.pem"},
				Cert:     []string{"/etc/registry/cert.pem"},
			},
			`
[registry]
  [registry."registry.example.com"]
    ca=["/etc/registry/ca.pem"]
    key=["/etc/registry/key.pem"]
    cert=["/etc/registry/cert.pem"]
`,
		},
		{
			BuildkitConfig{
				Registry: []string{},
			},
			`
[registry]
`,
		},
		{
			BuildkitConfig{
				Registry: []string{"registry.example.com", "registry.example2.com"},
				Mirror:   "https://mirror.example.com",
				UseHTTP:  true,
				Ca:       []string{"/etc/registry/ca.pem", "/etc/registry2/ca.pem"},
				Key:      []string{"/etc/registry/key.pem", "/etc/registry2/key.pem"},
				Cert:     []string{"/etc/registry/cert.pem", "/etc/registry2/cert.pem"},
			},
			`
[registry]
  [registry."registry.example.com"]
    http = true
    mirrors = ["https://mirror.example.com"]
    ca=["/etc/registry/ca.pem"]
    key=["/etc/registry/key.pem"]
    cert=["/etc/registry/cert.pem"]
  [registry."registry.example2.com"]
    http = true
    mirrors = ["https://mirror.example.com"]
    ca=["/etc/registry2/ca.pem"]
    key=["/etc/registry2/key.pem"]
    cert=["/etc/registry2/cert.pem"]
`,
		},
	}

	for _, tc := range testCases {
		config, err := tc.bc.String()
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(tc.expected), strings.TrimSpace(config))
	}
}