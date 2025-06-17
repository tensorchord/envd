// Copyright 2025 The envd Authors
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

package v1

import (
	"bytes"
	"text/template"

	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
)

const (
	pixiVersion        = "0.48.0"
	pixiConfigTemplate = `
{{- if .UsePixiMirror -}}
[mirrors]
"https://conda.anaconda.org" = ["https://prefix.dev/"]
# prefix.dev doesn't mirror conda-forge's label channels
# pixi uses the longest matching prefix for the mirror
"https://conda.anaconda.org/conda-forge/label" = [
  "https://conda.anaconda.org/conda-forge/label",
]
{{- end -}}
{{- if .PyPIIndex -}}
[pypi-config]
index-url = "{{ .PyPIIndex }}"
{{- end -}}
`
)

func (g generalGraph) compilePixi(root llb.State) llb.State {
	if g.PixiConfig == nil {
		return root
	}

	base := llb.Image(builderImage)
	builder := base.Run(
		llb.Shlexf(`sh -c "wget -qO- https://github.com/prefix-dev/pixi/releases/download/v%s/pixi-$(uname -m)-unknown-linux-musl.tar.gz | tar -xz -C /tmp || exit 1"`, pixiVersion),
		llb.WithCustomNamef("[internal] download pixi %s", pixiVersion),
	).Root()

	root = root.File(
		llb.Copy(builder, "/tmp/pixi", "/usr/bin/pixi"), llb.WithCustomName("[internal] install pixi"),
	)

	return g.compilePixiConfig(root)
}

func (g generalGraph) compilePixiConfig(root llb.State) llb.State {
	if !g.PixiConfig.UsePixiMirror && g.PixiConfig.PyPIIndex == nil {
		return root
	}

	root = root.File(
		llb.Mkdir("/etc/pixi", 0777, llb.WithParents(true)),
		llb.WithCustomName("[internal] create pixi config directory"),
	)

	tmpl, err := template.New("pixi-config").Parse(pixiConfigTemplate)
	if err != nil {
		logrus.Errorf("failed to parse pixi config template: %v", err)
		return root
	}
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, g.PixiConfig)
	if err != nil {
		logrus.Errorf("failed to execute pixi config template: %v", err)
		return root
	}
	root = root.File(
		llb.Mkfile("/etc/pixi/config.toml", 0755, buf.Bytes(), llb.WithUIDGID(g.uid, g.gid)),
		llb.WithCustomName("[internal] create pixi config file"),
	)
	return root
}
