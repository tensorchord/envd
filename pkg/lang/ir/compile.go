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
	"encoding/json"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/types"
)

func NewGraph() *Graph {
	return &Graph{
		OS: osDefault,
		Language: Language{
			Name: languageDefault,
		},
		CUDA:    nil,
		CUDNN:   nil,
		NumGPUs: -1,

		PyPIPackages:   []string{},
		RPackages:      []string{},
		JuliaPackages:  []string{},
		SystemPackages: []string{},
		Exec:           []string{},
		Shell:          shellBASH,
	}
}

var DefaultGraph = NewGraph()

func GPUEnabled() bool {
	return DefaultGraph.GPUEnabled()
}

func NumGPUs() int {
	return DefaultGraph.NumGPUs
}

func Compile(ctx context.Context, cachePrefix string, pub string) (*llb.Definition, error) {
	w, err := compileui.New(ctx, os.Stdout, "auto")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compileui")
	}
	DefaultGraph.Writer = w
	DefaultGraph.CachePrefix = cachePrefix
	DefaultGraph.PublicKeyPath = pub

	uid, gid, err := getUIDGID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get uid/gid")
	}
	state, err := DefaultGraph.Compile(uid, gid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile")
	}
	// TODO(gaocegege): Support multi platform.
	return state.Marshal(ctx, llb.LinuxAmd64)
}

func Labels() (map[string]string, error) {
	return DefaultGraph.Labels()
}

func (g Graph) GPUEnabled() bool {
	return g.CUDA != nil
}

func (g Graph) Labels() (map[string]string, error) {
	labels := make(map[string]string)
	str, err := json.Marshal(g.SystemPackages)
	if err != nil {
		return nil, err
	}
	labels[types.ImageLabelAPT] = string(str)
	str, err = json.Marshal(g.PyPIPackages)
	if err != nil {
		return nil, err
	}
	labels[types.ImageLabelPyPI] = string(str)
	str, err = json.Marshal(g.RPackages)
	if err != nil {
		return nil, err
	}
	labels[types.ImageLabelR] = string(str)
	if g.GPUEnabled() {
		labels[types.ImageLabelGPU] = "true"
		labels[types.ImageLabelCUDA] = *g.CUDA
		if g.CUDNN != nil {
			labels[types.ImageLabelCUDNN] = *g.CUDNN
		}
	}
	labels[types.ImageLabelVendor] = types.ImageVendorEnvd

	return labels, nil
}

func (g Graph) Compile(uid, gid int) (llb.State, error) {
	g.uid = uid

	// TODO(gaocegege): Remove the hack for https://github.com/tensorchord/envd/issues/370
	g.gid = 1001
	logrus.WithFields(logrus.Fields{
		"uid": g.uid,
		"gid": g.gid,
	}).Debug("compile LLB")

	// TODO(gaocegege): Support more OS and langs.
	base := g.compileBase()
	aptStage := g.compileUbuntuAPT(base)
	var merged llb.State
	var err error
	switch g.Language.Name {
	case "r":
		merged, err = g.compileRLang(aptStage)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to compile r language")
		}
	case "python":
		merged, err = g.compilePython(aptStage)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to compile python")
		}
	case "julia":
		merged, err = g.compileJulia(aptStage)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to compile julia")
		}
	}

	// TODO(gaocegege): Support order-based exec.
	run := g.compileRun(merged)
	finalStage, err := g.compileGit(run)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile git")
	}
	g.Writer.Finish()
	return finalStage, nil
}

func (g Graph) compileJulia(aptStage llb.State) (llb.State, error) {
	g.compileJupyter()
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
		}, llb.WithCustomName("merging all components into one"))
	} else {
		merged = llb.Merge([]llb.State{
			builtinSystemStage, systemStage, diffShellStage,
			diffSSHStage, juliaStage,
		}, llb.WithCustomName("merging all components into one"))
	}
	return merged, nil
}

func (g Graph) compileRLang(aptStage llb.State) (llb.State, error) {
	g.compileJupyter()
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
		}, llb.WithCustomName("merging all components into one"))
	} else {
		merged = llb.Merge([]llb.State{
			builtinSystemStage, systemStage, diffShellStage,
			diffSSHStage, rPackageInstallStage,
		}, llb.WithCustomName("merging all components into one"))
	}
	return merged, nil
}

func (g Graph) compilePython(aptStage llb.State) (llb.State, error) {
	condaChanelStage := g.compileCondaChannel(aptStage)
	pypiMirrorStage := g.compilePyPIIndex(condaChanelStage)

	g.compileJupyter()
	builtinSystemStage := pypiMirrorStage

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

	condaEnvStage := g.setCondaENV(shellStage)

	condaStage := llb.Diff(builtinSystemStage,
		g.compileCondaPackages(condaEnvStage),
		llb.WithCustomName("install conda packages"))

	pypiStage := llb.Diff(condaEnvStage,
		g.compilePyPIPackages(condaEnvStage),
		llb.WithCustomName("install PyPI packages"))
	systemStage := llb.Diff(builtinSystemStage, g.compileSystemPackages(builtinSystemStage),
		llb.WithCustomName("install system packages"))

	vscodeStage, err := g.compileVSCode()
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get vscode plugins")
	}

	var merged llb.State
	if vscodeStage != nil {
		merged = llb.Merge([]llb.State{
			builtinSystemStage, systemStage, condaStage,
			diffSSHStage, pypiStage, *vscodeStage,
		}, llb.WithCustomName("merging all components into one"))
	} else {
		merged = llb.Merge([]llb.State{
			builtinSystemStage, systemStage, condaStage,
			diffSSHStage, pypiStage,
		}, llb.WithCustomName("merging all components into one"))
	}
	return merged, nil
}
