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
	"github.com/moby/buildkit/client/llb"
)

const (
	golangDefaultVersion = "1.25.3"
	golangFilePath       = "/tmp/golang.linux.tar.gz"
	golangHomeBin        = "/usr/local/go/bin"
)

func (g *generalGraph) installGolang(root llb.State, version *string) llb.State {
	goVersion := golangDefaultVersion
	if version != nil {
		goVersion = *version
	}

	base := llb.Image(curlImage)
	builder := base.Run(
		llb.Shlexf(`sh -c "wget -qO %s https://go.dev/dl/go%s.linux-$(uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/').tar.gz"`, golangFilePath, goVersion),
		llb.WithCustomNamef("[internal] download go %s", goVersion),
	).Root()

	root = root.File(
		llb.Copy(builder, golangFilePath, golangFilePath),
		llb.WithCustomNamef("[internal] prepare go %s", goVersion),
	).Run(
		llb.Shlexf(`sh -c "tar -C /usr/local -xzf %[1]s && rm %[1]s"`, golangFilePath),
		llb.WithCustomNamef("[internal] install go %s", goVersion),
	).Root()
	g.RuntimeEnvPaths = append(g.RuntimeEnvPaths, golangHomeBin)
	return root
}
