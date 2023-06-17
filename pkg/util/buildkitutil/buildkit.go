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
	"text/template"
)

const buildkitConfigTemplate = `
[registry."{{ if .Registry }}{{ .Registry }}{{ else }}docker.io{{ end }}"]{{ if .Mirror }}
  mirrors = ["{{ .Mirror }}"]{{ end }}
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
