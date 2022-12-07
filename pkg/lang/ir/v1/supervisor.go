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

package v1

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/types"
)

const (
	horustTemplate = `
name = "%[1]s"
command = """
%[2]s
"""
stdout = "/var/log/horust/%[1]s_stdout.log"
stderr = "/var/log/horust/%[1]s_stderr.log"
user = "${USER}"
working-directory = "${%[3]s}"
%[4]s

[environment]
keep-env = true

[restart]
strategy = "on-failure"
backoff = "1s"
attempts = 2

[termination]
wait = "5s"
`
)

func (g generalGraph) installHorust(root llb.State) llb.State {
	horust := root.
		File(llb.Copy(llb.Image(types.HorustImage), "/", "/usr/local/bin"),
			llb.WithCustomName("[internal] install horust")).
		File(llb.Mkdir(types.HorustServiceDir, 0755, llb.WithParents(true)),
			llb.WithCustomNamef("[internal] mkdir for horust service: %s", types.HorustServiceDir)).
		File(llb.Mkdir(types.HorustLogDir, 0777, llb.WithParents(true)),
			llb.WithCustomNamef("[internal] mkdir for horust log: %s", types.HorustLogDir))
	return horust
}

func (g generalGraph) addNewProcess(root llb.State, name, command string, depends []string) llb.State {
	var sb strings.Builder
	if len(depends) != 0 {
		sb.WriteString("start-after = [")
		for _, d := range depends {
			sb.WriteString("\"")
			sb.WriteString(d)
			sb.WriteString("\",")
		}
		sb.WriteString("]\n")
	}
	template := fmt.Sprintf(horustTemplate, name, command, types.EnvdWorkDir, sb.String())

	filename := filepath.Join(types.HorustServiceDir, fmt.Sprintf("%s.toml", name))
	supervisor := root.File(llb.Mkfile(filename, 0644, []byte(template), llb.WithUIDGID(g.uid, g.gid)), llb.WithCustomNamef("[internal] create file %s", filename))
	return supervisor
}

func (g generalGraph) compileEntrypoint(root llb.State) (llb.State, error) {
	if len(g.Entrypoint) > 0 {
		return root, errors.New("`config.entrypoint` is only for custom image, maybe you need `runtime.init`")
	}
	cmd := fmt.Sprintf("/var/envd/bin/envd-sshd --port %d --shell %s", config.SSHPortInContainer, g.Shell)
	entrypoint := g.addNewProcess(root, "sshd", cmd, nil)
	var deps []string
	if g.RuntimeInitScript != nil {
		for i, command := range g.RuntimeInitScript {
			entrypoint = g.addNewProcess(entrypoint, fmt.Sprintf("init_%d", i), fmt.Sprintf("/bin/bash -c 'set -euo pipefail\n%s'", strings.Join(command, "\n")), nil)
			deps = append(deps, fmt.Sprintf("init_%d", i))
		}
	}

	if g.RuntimeDaemon != nil {
		for i, command := range g.RuntimeDaemon {
			entrypoint = g.addNewProcess(entrypoint, fmt.Sprintf("daemon_%d", i), strings.Join(command, " "), deps)
		}
	}

	if g.JupyterConfig != nil {
		jupyterCmd := g.generateJupyterCommand("")
		entrypoint = g.addNewProcess(entrypoint, "jupyter", strings.Join(jupyterCmd, " "), deps)
	}

	if g.RStudioServerConfig != nil {
		rstudioCmd := g.generateRStudioCommand("")
		entrypoint = g.addNewProcess(entrypoint, "rstudio", strings.Join(rstudioCmd, " "), deps)
	}

	return entrypoint, nil
}
