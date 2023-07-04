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
				Registries: []Registry{
					{
						Name:    "registry.example.com",
						Ca:      "/etc/registry/ca.pem",
						Cert:    "/etc/registry/cert.pem",
						Key:     "/etc/registry/key.pem",
						UseHTTP: false,
						Mirror:  "https://mirror.example.com",
					},
				},
			},
			`
[registry]
  [registry."registry.example.com"]
    mirrors = ["https://mirror.example.com"]
    ca=["/etc/registry/registry.example.com_ca.pem"]
    [[registry."registry.example.com".keypair]]
      key="/etc/registry/registry.example.com_key.pem"
      cert="/etc/registry/registry.example.com_cert.pem"
`,
		},
		{
			BuildkitConfig{
				Registries: []Registry{
					{
						Name:    "registry.example.com",
						UseHTTP: true,
						Mirror:  "https://mirror.example.com",
					},
					{
						Name:   "docker.io",
						Mirror: "https://mirror.example.com",
					},
				},
			},
			`
[registry]
  [registry."registry.example.com"]
    http = true
    mirrors = ["https://mirror.example.com"]
  [registry."docker.io"]
    mirrors = ["https://mirror.example.com"]
`,
		},
		{
			BuildkitConfig{
				Registries: []Registry{},
			},
			`
[registry]
`,
		},
		{
			BuildkitConfig{
				Registries: []Registry{
					{
						Name:    "registry1.example.com",
						Ca:      "/etc/registry/ca1.pem",
						Cert:    "/etc/registry/cert1.pem",
						Key:     "/etc/registry/key1.pem",
						UseHTTP: true,
						Mirror:  "https://mirror.example.com",
					},
					{
						Name:   "registry2.example.com",
						Ca:     "/etc/registry/ca2.pem",
						Cert:   "/etc/registry/cert2.pem",
						Key:    "/etc/registry/key2.pem",
						Mirror: "https://mirror.example.com",
					},
				},
			},
			`
[registry]
  [registry."registry1.example.com"]
    http = true
    mirrors = ["https://mirror.example.com"]
    ca=["/etc/registry/registry1.example.com_ca.pem"]
    [[registry."registry1.example.com".keypair]]
      key="/etc/registry/registry1.example.com_key.pem"
      cert="/etc/registry/registry1.example.com_cert.pem"
  [registry."registry2.example.com"]
    mirrors = ["https://mirror.example.com"]
    ca=["/etc/registry/registry2.example.com_ca.pem"]
    [[registry."registry2.example.com".keypair]]
      key="/etc/registry/registry2.example.com_key.pem"
      cert="/etc/registry/registry2.example.com_cert.pem"
`,
		},
	}

	for _, tc := range testCases {
		config, err := tc.bc.String()
		assert.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(tc.expected), strings.TrimSpace(config))
	}
}
