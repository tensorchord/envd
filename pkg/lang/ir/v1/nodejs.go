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

// from https://nodejs.org/download/release/
const (
	nodejsDefaultVersion = "25.1.0"
	nodejsTempDir        = "/tmp/nodejs"
	nodejsHomeDir        = "/opt/nodejs"
	nodejsHomeBin        = "/opt/nodejs/bin"
)

func (g *generalGraph) installNodeJS(root llb.State, version *string) llb.State {
	nodejsVersion := nodejsDefaultVersion
	if version != nil {
		nodejsVersion = *version
	}

	base := llb.Image(curlImage)
	builder := base.Run(
		llb.Shlexf(`sh -c "mkdir %[1]s && wget -qO- https://nodejs.org/download/release/v%[2]s/node-v%[2]s-linux-$(uname -m | sed -e 's/x86_64/x64/').tar.xz | tar -xJ --strip-components=1 -C %[1]s || exit 1"`, nodejsTempDir, nodejsVersion),
		llb.WithCustomNamef("[internal] download nodejs %s", nodejsVersion),
	).Root()

	root = root.File(
		llb.Copy(builder, nodejsTempDir, nodejsHomeDir),
		llb.WithCustomNamef("[internal] prepare nodejs %s", nodejsVersion),
	)
	g.RuntimeEnvPaths = append(g.RuntimeEnvPaths, nodejsHomeBin)
	return root
}
