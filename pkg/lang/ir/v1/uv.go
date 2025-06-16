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

import "github.com/moby/buildkit/client/llb"

const (
	uvVersion = "0.7.10"
)

func (g generalGraph) compileUV(root llb.State) llb.State {
	if g.UVConfig == nil {
		return root
	}
	// uv configurations
	g.RuntimeEnviron["UV_LINK_MODE"] = "copy"
	g.RuntimeEnviron["UV_PYTHON_PREFERENCE"] = "only-managed"

	base := llb.Image(builderImage)
	builder := base.Run(
		llb.Shlexf(`sh -c "wget -qO- https://github.com/astral-sh/uv/releases/download/%s/uv-$(uname -m)-unknown-linux-gnu.tar.gz | tar -xz --strip-components=1 -C /tmp || exit 1"`, uvVersion),
		llb.WithCustomNamef("[internal] download uv %s", uvVersion),
	).Root()

	root = root.File(
		llb.Copy(builder, "/tmp/uv", "/usr/bin/uv"), llb.WithCustomName("[internal] install uv")).
		File(llb.Copy(builder, "/tmp/uvx", "/usr/bin/uvx"), llb.WithCustomName("[internal] install uvx"))
	return g.compileUVPython(root)
}

func (g generalGraph) compileUVPython(root llb.State) llb.State {
	if g.UVConfig == nil {
		return root
	}

	root = root.Run(
		llb.Shlexf(`uv python install %s`, g.UVConfig.PythonVersion),
		llb.WithCustomNamef("[internal] install uv Python=%s", g.UVConfig.PythonVersion),
	).Root()
	return root
}
