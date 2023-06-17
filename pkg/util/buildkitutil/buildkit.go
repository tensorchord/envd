package buildkitutil

import (
	"strings"
	"text/template"
)

const buildkitConfigTemplate = `
[registry."{{ if .Registry }}{{ .Registry }}{{ else }}docker.io{{ end }}"]
  mirrors = ["{{ .Mirror }}"]
  http = {{ .UseHTTP }}
  {{ if .SetCA}}ca=["/etc/registry/ca.pem"]
  [[registry."{{ if .Registry }}{{ .Registry }}{{ else }}docker.io{{ end }}".keypair]]
	key="/etc/registry/key.pem"
	cert="/etc/registry/cert.pem"
  {{ end }}
`

type BuildkitConfig struct {
	Registry string
	Mirror   string
	UseHTTP  bool
	SetCA    bool
}

func (c *BuildkitConfig) String() (string, error) {
	tmpl, err := template.New("buildkitConfig").Parse(buildkitConfigTemplate)
	if err != nil {
		return "", err
	}
	var config strings.Builder
	err = tmpl.Execute(&config, c)
	return config.String(), err
}
