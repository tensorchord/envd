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
	"context"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"

	"github.com/tensorchord/envd/pkg/progress/compileui"
)

func NewGraph() *Graph {
	return &Graph{
		OS:       osDefault,
		Language: languageDefault,
		CUDA:     nil,
		CUDNN:    nil,
		BuiltinSystemPackages: []string{
			// TODO(gaocegege): Move them into the base image.
			"curl",
			"openssh-client",
			"git",
			"sudo",
			"tini",
		},

		PyPIPackages:   []string{},
		SystemPackages: []string{},
		Exec:           []string{},
		Shell:          shellBASH,
	}
}

var DefaultGraph = NewGraph()

func GPUEnabled() bool {
	return DefaultGraph.CUDA != nil
}

func Compile(ctx context.Context) (*llb.Definition, error) {
	w, err := compileui.New(ctx, os.Stdout, "auto")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compileui")
	}
	DefaultGraph.Writer = w
	state, err := DefaultGraph.Compile()
	if err != nil {
		return nil, err
	}
	// TODO(gaocegege): Support multi platform.
	def, err := state.Marshal(ctx, llb.LinuxAmd64)
	if err != nil {
		return nil, err
	}
	return def, nil
}

func (g Graph) Compile() (llb.State, error) {
	// TODO(gaocegege): Support more OS and langs.
	base := g.compileBase()
	aptStage := g.compileUbuntuAPT(base)
	pypiMirrorStage := g.compilePyPIMirror(aptStage)

	g.compileJupyter()
	// TODO(gaocegege): Make apt update a seperate stage to
	// parallel system and user-defined package installation.
	builtinSystemStage := g.compileBuiltinSystemPackages(pypiMirrorStage)
	shellStage, err := g.compileShell(builtinSystemStage)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile shell")
	}
	diffShellStage := llb.Diff(builtinSystemStage, shellStage, llb.WithCustomName("install shell"))
	pypiStage := llb.Diff(builtinSystemStage, g.compilePyPIPackages(builtinSystemStage), llb.WithCustomName("install PyPI packages"))
	systemStage := llb.Diff(builtinSystemStage, g.compileSystemPackages(builtinSystemStage), llb.WithCustomName("install system packages"))
	sshStage := g.copyEnvdSSHServer()

	vscodeStage, err := g.compileVSCode()
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get vscode plugins")
	}

	var merged llb.State
	if vscodeStage != nil {
		merged = llb.Merge([]llb.State{
			builtinSystemStage, systemStage, pypiStage, sshStage, *vscodeStage, diffShellStage,
		}, llb.WithCustomName("merging all components into one"))
	} else {
		merged = llb.Merge([]llb.State{
			builtinSystemStage, systemStage, pypiStage, sshStage, diffShellStage,
		}, llb.WithCustomName("merging all components into one"))
	}

	// TODO(gaocegege): Support order-based exec.
	run := g.compileRun(merged)
	g.Writer.Finish()
	return run, nil
}
