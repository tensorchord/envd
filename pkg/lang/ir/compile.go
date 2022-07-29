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
	"fmt"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
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

func ExposedPorts() (map[string]struct{}, error) {
	return DefaultGraph.ExposedPorts()
}

func Entrypoint(buildContextDir string) ([]string, error) {
	return DefaultGraph.Entrypoint(buildContextDir)
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

func (g Graph) ExposedPorts() (map[string]struct{}, error) {
	ports := make(map[string]struct{})
	ports[fmt.Sprintf("%d/tcp", config.SSHPortInContainer)] = struct{}{}
	if g.JupyterConfig != nil {
		ports[fmt.Sprintf("%d/tcp", config.JupyterPortInContainer)] = struct{}{}
	}
	if g.RStudioServerConfig != nil {
		ports[fmt.Sprintf("%d/tcp", config.RStudioServerPortInContainer)] = struct{}{}
	}

	return ports, nil
}

func (g Graph) Entrypoint(buildContextDir string) ([]string, error) {
	// Do not set entrypoint if the image is customized.
	if g.Image != nil {
		logrus.Debug("skip entrypoint because the image is customized")
		return []string{}, nil
	}

	ep := []string{
		"tini",
		"--",
		"bash",
		"-c",
	}

	template := `set -e
/var/envd/bin/envd-ssh --authorized-keys %s --port %d --shell %s &
%s
wait -n`

	// Generate jupyter and rstudio server commands.
	var customCmd strings.Builder
	if g.JupyterConfig != nil {
		workingDir := fmt.Sprintf("/home/envd/%s", fileutil.Base(buildContextDir))
		jupyterCmd := g.generateJupyterCommand(workingDir)
		customCmd.WriteString(strings.Join(jupyterCmd, " "))
		customCmd.WriteString("\n")
	}
	if g.RStudioServerConfig != nil {
		workingDir := fmt.Sprintf("/home/envd/%s", fileutil.Base(buildContextDir))
		rstudioCmd := g.generateRStudioCommand(workingDir)
		customCmd.WriteString(strings.Join(rstudioCmd, " "))
		customCmd.WriteString("\n")
	}

	cmd := fmt.Sprintf(template,
		config.ContainerAuthorizedKeysPath,
		config.SSHPortInContainer, g.Shell, customCmd.String())
	ep = append(ep, cmd)

	logrus.WithField("entrypoint", ep).Debug("generate entrypoint")
	return ep, nil
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
	// Use custom logic when image is specified.
	if g.Image != nil {
		merged, err = g.compileCustomPython(aptStage)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to compile custom python image")
		}
	} else {
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
	}

	prompt := g.compilePrompt(merged)
	copy := g.compileCopy(prompt)
	// TODO(gaocegege): Support order-based exec.
	run := g.compileRun(copy)
	finalStage, err := g.compileGit(run)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile git")
	}
	g.Writer.Finish()
	return finalStage, nil
}
