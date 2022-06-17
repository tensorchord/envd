// Copyright 2022 The envd Authors
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

package docker

import (
	"fmt"
	"strings"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/editor/jupyter"
	"github.com/tensorchord/envd/pkg/lang/ir"
)

const (
	template = `set -e
/var/envd/bin/envd-ssh --authorized-keys %s --port %d --shell %s &
%s
wait -n`
)

func entrypointSH(g ir.Graph, workingDir string, sshPort int) string {
	if g.JupyterConfig != nil {
		cmds := jupyter.GenerateCommand(g, workingDir)
		return fmt.Sprintf(template,
			config.ContainerauthorizedKeysPath, sshPort, g.Shell,
			strings.Join(cmds, " "))
	}
	return fmt.Sprintf(template,
		config.ContainerauthorizedKeysPath, sshPort, g.Shell, "")
}
