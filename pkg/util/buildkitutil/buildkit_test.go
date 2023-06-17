package buildkitutil

import (
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
				SetCA:   true,
			},
			`
[registry."registry.example.com"]
  ca=["/etc/registry/ca.pem"]
  [[registry."registry.example.com".keypair]]
	key="/etc/registry/key.pem"
	cert="/etc/registry/cert.pem"
`,
		},
	}

	for _, tc := range testCases {
		config, err := tc.bc.String()
		assert.NoError(t, err)
		assert.Equal(t, tc.expected, config)
	}
}
