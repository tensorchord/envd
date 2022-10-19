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

package ir

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/client/llb"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/types"
)

const (
	horustTemplate = `
name = "%[1]s"
command = "%[2]s"
stdout = "/var/logs/%[1]s_stdout.log"
stderr = "/var/logs/%[1]s_stderr.log"
user = "${USER}"
working-directory = "${ENVD_WORKDIR}"

[environment]
keep-env = true
re-export = [ "PATH", "SHELL", "USER", "ENVD_WORKDIR" ]

[restart]
strategy = "on-failure"
backoff = "1s"
attempts = 5

[termination]
wait = "5s"
`
)

func (g Graph) addNewProcess(root llb.State, name, command string) llb.State {
	template := fmt.Sprintf(horustTemplate, name, command)
	filename := filepath.Join(types.HorustServiceDir, fmt.Sprintf("%s.toml", name))
	supervisor := root.File(llb.Mkfile(filename, 0644, []byte(template), llb.WithUIDGID(g.uid, g.gid)))
	return supervisor
}

func (g Graph) compileEntrypoint(root llb.State) llb.State {
	if g.Image != nil {
		return root
	}
	cmd := fmt.Sprintf("/var/envd/bin/envd-sshd --port %d --shell %s", config.SSHPortInContainer, g.Shell)
	entrypoint := g.addNewProcess(root, "sshd", cmd)

	if g.RuntimeDaemon != nil {
		for i, command := range g.RuntimeDaemon {
			entrypoint = g.addNewProcess(entrypoint, fmt.Sprintf("daemon_%d", i), fmt.Sprintf("%s &\n", strings.Join(command, " ")))
		}
	}

	if g.JupyterConfig != nil {
		jupyterCmd := g.generateJupyterCommand("")
		entrypoint = g.addNewProcess(entrypoint, "jupyter", strings.Join(jupyterCmd, " "))
	}

	if g.RStudioServerConfig != nil {
		rstudioCmd := g.generateRStudioCommand("")
		entrypoint = g.addNewProcess(entrypoint, "rstudio", strings.Join(rstudioCmd, " "))
	}

	return entrypoint
}
