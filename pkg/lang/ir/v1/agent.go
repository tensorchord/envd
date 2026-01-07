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

	"github.com/tensorchord/envd/pkg/lang/ir"
)

// https://github.com/openai/codex
const (
	codexDefaultVersion = "0.79.0"
)

func (g generalGraph) installAgentCodex(root llb.State, agent ir.CodeAgent) llb.State {
	base := llb.Image(curlImage)
	version := codexDefaultVersion
	if agent.Version != nil {
		version = *agent.Version
	}
	builder := base.Run(
		llb.Shlexf(`sh -c "wget -qO- https://github.com/openai/codex/releases/download/rust-v%s/codex-$(uname -m)-unknown-linux-musl.tar.gz | tar -xz -C /tmp || exit 1"`, version),
		llb.WithCustomNamef("[internal] download codex %s", version),
	).Run(
		llb.Shlex(`sh -c "mv /tmp/codex-$(uname -m)-unknown-linux-musl /tmp/codex"`),
		llb.WithCustomNamef("[internal] prepare codex %s", version),
	).Root()
	root = root.File(
		llb.Copy(builder, "/tmp/codex", "/usr/bin/codex"),
		llb.WithCustomName("[internal] install codex"),
	)
	return root
}
