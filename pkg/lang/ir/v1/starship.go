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
	starshipDefaultVersion = "1.24.0"
)

func (g generalGraph) compileStarship(root llb.State) llb.State {
	base := llb.Image(curlImage)
	builder := base.Run(
		llb.Shlexf(`sh -c "wget -qO- https://github.com/starship/starship/releases/download/v%s/starship-$(uname -m)-unknown-linux-musl.tar.gz | tar -xz -C /tmp || exit 1"`, starshipDefaultVersion),
		llb.WithCustomNamef("[internal] download starship %s", starshipDefaultVersion),
	).Root()

	root = root.File(
		llb.Copy(builder, "/tmp/starship", "/usr/local/bin/starship"),
		llb.WithCustomNamef("[internal] install starship %s", starshipDefaultVersion),
	)
	return root
}
