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
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/editor/vscode"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func (g generalGraph) compileVSCode() (llb.State, error) {
	if len(g.VSCodePlugins) == 0 {
		return llb.Scratch(), nil
	}
	inputs := []llb.State{}
	for _, p := range g.VSCodePlugins {
		vscodeClient, err := vscode.NewClient(vscode.MarketplaceVendorOpenVSX)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to create vscode client")
		}
		g.Writer.LogVSCodePlugin(p, compileui.ActionStart, false)
		if cached, err := vscodeClient.DownloadOrCache(p); err != nil {
			return llb.State{}, err
		} else {
			g.Writer.LogVSCodePlugin(p, compileui.ActionEnd, cached)
		}
		ext := llb.Scratch().File(llb.Copy(llb.Local(flag.FlagCacheDir),
			vscodeClient.PluginPath(p),
			fileutil.EnvdHomeDir(".vscode-server", "extensions", p.String()),
			&llb.CopyInfo{
				CreateDestPath: true,
			}, llb.WithUIDGID(g.uid, g.gid)),
			llb.WithCustomNamef("install vscode plugin %s", p.String()))
		inputs = append(inputs, ext)
	}
	layer := llb.Merge(inputs, llb.WithCustomName("merging plugins for vscode"))
	return layer, nil
}

// nolint:unused
func (g *generalGraph) compileJupyter() error {
	if g.JupyterConfig == nil {
		return nil
	}

	g.PyPIPackages = append(g.PyPIPackages, []string{"jupyter"})
	switch g.Language.Name {
	case "python":
		return nil
	default:
		return errors.Newf("Jupyter is not supported in %s yet", g.Language.Name)
	}
}

func (g generalGraph) generateJupyterCommand(workingDir string) []string {
	if g.JupyterConfig == nil {
		return nil
	}

	if g.JupyterConfig.Token == "" {
		g.JupyterConfig.Token = "''"
	}

	// get from env if not set
	if len(workingDir) == 0 {
		workingDir = fmt.Sprintf("${%s}", types.EnvdWorkDir)
	}

	cmd := []string{
		"python3", "-m", "notebook",
		"--ip", "0.0.0.0", "--notebook-dir", workingDir,
		"--NotebookApp.token", g.JupyterConfig.Token,
		"--port", strconv.Itoa(config.JupyterPortInContainer),
	}

	if g.uid == 0 {
		cmd = append(cmd, "--allow-root")
	}

	return cmd
}

// nolint:unparam
func (g generalGraph) generateRStudioCommand(workingDir string) []string {
	if g.RStudioServerConfig == nil {
		return nil
	}

	// get from env if not set
	// if len(workingDir) == 0 {
	// 	workingDir = fmt.Sprintf("${%s}", types.EnvdWorkDir)
	// }

	return []string{
		// TODO(gaocegege): Remove root permission here.
		"sudo",
		"/usr/lib/rstudio-server/bin/rserver",
		// TODO(gaocegege): Support working dir.
	}
}
