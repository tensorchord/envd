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

package v0

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
)

func (g generalGraph) compileRLang(baseStage llb.State) (llb.State, error) {
	if err := g.compileJupyter(); err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile jupyter")
	}
	aptStage := g.compileUbuntuAPT(baseStage)
	builtinSystemStage := aptStage

	sshStage, err := g.copySSHKey(builtinSystemStage)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to copy ssh keys")
	}
	diffSSHStage := llb.Diff(builtinSystemStage, sshStage, llb.WithCustomName("install ssh keys"))

	// Conda affects shell and python, thus we cannot do it in parallel.
	shellStage, err := g.compileShell(builtinSystemStage)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile shell")
	}
	diffShellStage := llb.Diff(builtinSystemStage, shellStage, llb.WithCustomName("install shell"))

	systemStage := llb.Diff(builtinSystemStage, g.compileSystemPackages(builtinSystemStage),
		llb.WithCustomName("install system packages"))

	// TODO(terrytangyuan): Support RStudio local server
	rPackageInstallStage := llb.Diff(builtinSystemStage,
		g.installRPackages(builtinSystemStage), llb.WithCustomName("install R packages"))

	vscodeStage, err := g.compileVSCode()
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get vscode plugins")
	}

	var merged llb.State
	if vscodeStage != nil {
		merged = llb.Merge([]llb.State{
			builtinSystemStage, systemStage, diffShellStage,
			diffSSHStage, rPackageInstallStage, *vscodeStage,
		}, llb.WithCustomName("[internal] generating the image"))
	} else {
		merged = llb.Merge([]llb.State{
			builtinSystemStage, systemStage, diffShellStage,
			diffSSHStage, rPackageInstallStage,
		}, llb.WithCustomName("[internal] generating the image"))
	}
	return merged, nil
}

func (g generalGraph) installRPackages(root llb.State) llb.State {
	if len(g.RPackages) == 0 {
		return root
	}
	// TODO(terrytangyuan): Support different CRAN mirrors
	var sb strings.Builder
	mirrorURL := "https://cran.rstudio.com"
	if g.CRANMirrorURL != nil {
		mirrorURL = *g.CRANMirrorURL
	}
	sb.WriteString(fmt.Sprintf(`R -e 'options(repos = c(CRAN = "%s")); install.packages(c(`, mirrorURL))
	for i, pkg := range g.RPackages {
		sb.WriteString(fmt.Sprintf(`"%s"`, pkg))
		if i != len(g.RPackages)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(`))'`)

	// TODO(terrytangyuan): Support cache.
	cmd := sb.String()
	root = llb.User("envd")(root)
	run := root.Run(llb.Shlex(cmd), llb.WithCustomNamef("install R packages"))
	return run.Root()
}
