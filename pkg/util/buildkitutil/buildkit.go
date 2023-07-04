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
	"os"
	"strings"
	"text/template"

	"github.com/cockroachdb/errors"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

const buildkitConfigTemplate = `
[registry]
{{- range $registry := .Registries }}
  [registry."{{ if $registry.Name }}{{ $registry.Name }}{{ else }}docker.io{{ end }}"]
    {{- if $registry.UseHTTP }}
    http = true
    {{- end }}
    {{- if $registry.Mirror }}
    mirrors = ["{{ $registry.Mirror }}"]
    {{- end }}
    {{- if $registry.Ca }}
    ca=["/etc/registry/{{ $registry.Name }}_ca.pem"]
    {{- end }}
    {{- if and $registry.Cert $registry.Key }}
    [[registry."{{ if $registry.Name }}{{ $registry.Name }}{{ else }}docker.io{{ end }}".keypair]]
      key="/etc/registry/{{ $registry.Name }}_key.pem"
      cert="/etc/registry/{{ $registry.Name }}_cert.pem"
    {{- end }}
{{- end }}
`

type Registry struct {
	Name    string `json:"name"`
	Ca      string `json:"ca"`
	Cert    string `json:"cert"`
	Key     string `json:"key"`
	UseHTTP bool   `json:"use_http"`
	Mirror  string `json:"mirror"`
}

type BuildkitConfig struct {
	Registries []Registry `json:"registries"`
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

func (c *BuildkitConfig) Save() error {
	text, err := c.String()
	if err != nil {
		return err
	}

	path, err := fileutil.ConfigFile("buildkitd.toml")
	if err != nil {
		return errors.Wrap(err, "failed to get the buildkitd config file")
	}
	err = os.WriteFile(path, []byte(text), 0644)
	return err
}
