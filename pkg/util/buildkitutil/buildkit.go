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
[registry]
{{- range $index, $value := .Registry }}
  [registry."{{ if $value }}{{ $value }}{{ else }}docker.io{{ end }}"]
  	{{- if $.UseHTTP }}
    http = true
  	{{- end }}
  	{{- if $.Mirror }}
    mirrors = ["{{ $.Mirror }}"]
  	{{- end }}
    {{- if $.Ca }}
    ca=["{{ index $.Ca $index }}"]
    key=["{{ index $.Key $index }}"]
    cert=["{{ index $.Cert $index }}"]
    {{- end }}
{{- end }}
`

type BuildkitConfig struct {
	Mirror   string
	UseHTTP  bool
	Registry []string
	Ca       []string
	Cert     []string
	Key      []string
	Bindings []string
}

func (c *BuildkitConfig) String() (string, error) {
	tmpl, err := template.New("buildkitConfig").Parse(buildkitConfigTemplate)
	if err != nil {
		return "", err
	}
	var config strings.Builder
	err = tmpl.Execute(&config, c)
	if err != nil {
		return "", err
	}

	return config.String(), nil
}
