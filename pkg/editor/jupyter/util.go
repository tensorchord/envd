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

package jupyter

import (
	"strconv"

	"github.com/tensorchord/envd/pkg/lang/ir"
)

func GenerateCommand(g ir.Graph, notebookDir string) []string {
	if g.JupyterConfig == nil {
		return nil
	}

	var cmd []string
	// Use python in conda env.
	if g.CondaEnabled() {
		cmd = append(cmd, "/opt/conda/bin/python3")
	} else {
		cmd = append(cmd, "python3")
	}

	cmd = append(cmd, []string{
		"-m", "notebook",
		"--ip", "0.0.0.0", "--notebook-dir", notebookDir,
	}...)

	if g.JupyterConfig.Password != "" {
		cmd = append(cmd, "--NotebookApp.password", g.JupyterConfig.Password,
			"--NotebookApp.token", "''")
	} else {
		cmd = append(cmd, "--NotebookApp.password", "''",
			"--NotebookApp.token", "''")
	}
	if g.JupyterConfig.Port != 0 {
		p := strconv.Itoa(int(g.JupyterConfig.Port))
		cmd = append(cmd, "--port", p)
	}
	return cmd
}
