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
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/types"
)

func (g *generalGraph) compileJulia(baseStage llb.State) (llb.State, error) {
	baseStage = g.updateEnvPath(baseStage, types.DefaultJuliaPath)
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

	shellStage, err := g.compileShell(builtinSystemStage)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile shell")
	}
	diffShellStage := llb.Diff(builtinSystemStage, shellStage, llb.WithCustomName("install shell"))

	systemStage := llb.Diff(builtinSystemStage, g.compileSystemPackages(builtinSystemStage),
		llb.WithCustomName("install system packages"))

	juliaStage := llb.Diff(builtinSystemStage,
		g.installJuliaPackages(builtinSystemStage), llb.WithCustomName("install julia packages"))

	vscodeStage, err := g.compileVSCode()
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get vscode plugins")
	}

	var merged llb.State
	if vscodeStage != nil {
		merged = llb.Merge([]llb.State{
			builtinSystemStage, systemStage, diffShellStage,
			diffSSHStage, juliaStage, *vscodeStage,
		}, llb.WithCustomName("[internal] generating the image"))
	} else {
		merged = llb.Merge([]llb.State{
			builtinSystemStage, systemStage, diffShellStage,
			diffSSHStage, juliaStage,
		}, llb.WithCustomName("[internal] generating the image"))
	}
	return merged, nil
}

func (g generalGraph) installJuliaPackages(root llb.State) llb.State {
	if len(g.JuliaPackages) == 0 {
		return root
	}

	var sb strings.Builder

	sb.WriteString(`/usr/local/julia/bin/julia -e 'using Pkg; Pkg.add([`)
	for i, pkg := range g.JuliaPackages {
		sb.WriteString(fmt.Sprintf(`"%s"`, pkg))
		if i != len(g.JuliaPackages)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString(`])'`)

	// TODO(gaocegege): Support cache.
	cmd := sb.String()
	logrus.Debug("install julia packages: ", cmd)
	root = llb.User("envd")(root)
	if g.JuliaPackageServer != nil {
		root = root.AddEnv("JULIA_PKG_SERVER", *g.JuliaPackageServer)
	}
	root = root.AddEnv("PATH", "/usr/local/julia/bin")
	run := root.
		Run(llb.Shlex(cmd), llb.WithCustomNamef("install julia packages"))

	return run.Root()
}
