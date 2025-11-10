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
	rustUpInitFilePath = "/tmp/rustup-init.sh"
	cargoHomeDir       = "/opt/rust"
	cargoHomeBin       = "/opt/rust/bin"
)

func (g *generalGraph) installRust(root llb.State, version *string) llb.State {
	base := llb.Image(curlImage)
	builder := base.Run(
		llb.Shlexf(`sh -c "curl --proto '=https' --tlsv1.2 -sSf -o %s https://sh.rustup.rs"`, rustUpInitFilePath),
		llb.WithCustomName("[internal] download rustup-init.sh"),
	).Root()
	root = root.File(
		llb.Copy(builder, rustUpInitFilePath, rustUpInitFilePath),
		llb.WithCustomName("[internal] copy the rustup-init.sh"),
	)
	if version != nil {
		root = root.AddEnv("RUSTUP_VERSION", *version)
	}
	root = root.AddEnv("CARGO_HOME", cargoHomeDir).File(
		llb.Mkdir(cargoHomeDir, 0755, llb.WithParents(true), llb.WithUIDGID(g.uid, g.gid)),
		llb.WithCustomNamef("[internal] create cargo dir: %s", cargoHomeDir),
	).Run(
		llb.Shlexf(`sh -c "sh %[1]s -y -q --no-modify-path && rm %[1]s"`, rustUpInitFilePath),
		llb.WithCustomName("[internal] install rust"),
	).Root()
	g.RuntimeEnvPaths = append(g.RuntimeEnvPaths, cargoHomeBin)
	return root
}
