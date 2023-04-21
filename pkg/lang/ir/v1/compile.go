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
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
	servertypes "github.com/tensorchord/envd-server/api/types"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func NewGraph() ir.Graph {
	runtimeGraph := ir.RuntimeGraph{
		RuntimeCommands: make(map[string]string),
		RuntimeEnviron:  make(map[string]string),
		RuntimeEnvPaths: []string{types.DefaultSystemPath},
	}
	return &generalGraph{
		uid:               -1,
		gid:               -1,
		Image:             defaultImage,
		CUDA:              nil,
		CUDNN:             CUDNNVersionDefault,
		NumGPUs:           0,
		EnvdSyntaxVersion: "v1",

		PyPIPackages:    [][]string{},
		RPackages:       [][]string{},
		JuliaPackages:   [][]string{},
		SystemPackages:  []string{},
		Exec:            []ir.RunBuildCommand{},
		UserDirectories: []string{},
		Shell:           shellBASH,
		RuntimeGraph:    runtimeGraph,
	}
}

var DefaultGraph = NewGraph()

func (g *generalGraph) SetWriter(w compileui.Writer) {
	g.Writer = w
}

func (g generalGraph) GetHTTP() []ir.HTTPInfo {
	return g.HTTP
}

func (g generalGraph) GetNumGPUs() int {
	return g.NumGPUs
}

func (g generalGraph) GetShell() string {
	return g.Shell
}

func (g generalGraph) GetMount() []ir.MountInfo {
	return g.Mount
}

func (g generalGraph) GetEnvironmentName() string {
	return g.EnvironmentName
}

func (g generalGraph) GetJupyterConfig() *ir.JupyterConfig {
	return g.JupyterConfig
}

func (g generalGraph) GetRStudioServerConfig() *ir.RStudioServerConfig {
	return g.RStudioServerConfig
}

func (g generalGraph) GetExposedPorts() []ir.ExposeItem {
	return g.RuntimeExpose
}

func (g generalGraph) GetRuntimeCommands() map[string]string {
	return g.RuntimeCommands
}

func (g *generalGraph) Compile(ctx context.Context, envName string, pub string) (*llb.Definition, error) {
	w, err := compileui.New(ctx, os.Stdout, "auto")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compileui")
	}
	g.Writer = w
	g.EnvironmentName = envName
	g.PublicKeyPath = pub

	uid, gid, err := g.getUIDGID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get uid/gid")
	}
	state, err := g.CompileLLB(uid, gid)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile the graph")
	}
	// TODO(gaocegege): Support multi platform.
	def, err := state.Marshal(ctx, llb.LinuxAmd64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal the llb definition")
	}
	return def, nil
}

func (g generalGraph) GetEnviron() []string {
	return append(g.EnvString(),
		"LC_ALL=en_US.UTF-8",
		"LANG=C.UTF-8",
	)
}

func (g generalGraph) GetUser() string {
	return g.User
}

func (g generalGraph) GPUEnabled() bool {
	return g.CUDA != nil
}

func (g generalGraph) Labels() (map[string]string, error) {
	labels := make(map[string]string)

	labels[types.ImageLabelSyntaxVer] = g.EnvdSyntaxVersion

	str, err := json.Marshal(g.SystemPackages)
	if err != nil {
		return nil, err
	}
	labels[types.ImageLabelAPT] = string(str)
	pyPackages := []string{}
	for _, pkg := range g.PyPIPackages {
		pyPackages = append(pyPackages, pkg...)
	}
	str, err = json.Marshal(pyPackages)
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
		labels[types.ImageLabelCUDNN] = g.CUDNN
	}
	labels[types.ImageLabelVendor] = types.ImageVendorEnvd
	code, err := g.RuntimeGraph.Dump()
	if err != nil {
		return labels, err
	}
	labels[types.RuntimeGraphCode] = code

	code, err = g.Dump()
	if err != nil {
		return labels, err
	}
	labels[types.GeneralGraphCode] = code

	ports := []servertypes.EnvironmentPort{}
	ports = append(ports, servertypes.EnvironmentPort{
		Name: "ssh",
		Port: config.SSHPortInContainer,
	})
	if g.JupyterConfig != nil {
		ports = append(ports, servertypes.EnvironmentPort{
			Name: "jupyter",
			Port: config.JupyterPortInContainer,
		})
	}
	if g.RStudioServerConfig != nil {
		ports = append(ports, servertypes.EnvironmentPort{
			Name: "rstudio-server",
			Port: config.RStudioServerPortInContainer,
		})
	}

	if g.RuntimeExpose != nil && len(g.RuntimeExpose) > 0 {
		for _, item := range g.RuntimeExpose {
			ports = append(ports, servertypes.EnvironmentPort{
				Name: item.ServiceName,
				Port: int32(item.EnvdPort),
			})
		}
	}

	portsData, err := json.Marshal(ports)
	if err != nil {
		return labels, err
	}
	labels[types.ImageLabelPorts] = string(portsData)

	repoInfo, err := json.Marshal(g.Repo)
	if err != nil {
		return labels, err
	}
	labels[types.ImageLabelRepo] = string(repoInfo)

	labels[types.ImageLabelContainerName] = g.EnvironmentName
	return labels, nil
}

func (g generalGraph) ExposedPorts() (map[string]struct{}, error) {
	ports := make(map[string]struct{})

	// only expose ports for dev env
	if !g.Dev {
		return ports, nil
	}

	ports[fmt.Sprintf("%d/tcp", config.SSHPortInContainer)] = struct{}{}
	if g.JupyterConfig != nil {
		ports[fmt.Sprintf("%d/tcp", config.JupyterPortInContainer)] = struct{}{}
	}
	if g.RStudioServerConfig != nil {
		ports[fmt.Sprintf("%d/tcp", config.RStudioServerPortInContainer)] = struct{}{}
	}

	if g.RuntimeExpose != nil && len(g.RuntimeExpose) > 0 {
		for _, item := range g.RuntimeExpose {
			ports[fmt.Sprintf("%d/tcp", item.EnvdPort)] = struct{}{}
		}
	}

	return ports, nil
}

func (g generalGraph) EnvString() []string {
	var envs []string
	for k, v := range g.RuntimeEnviron {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	envs = append(envs, fmt.Sprintf("PATH=%s", strings.Join(g.RuntimeEnvPaths, ":")))
	return envs
}

func (g generalGraph) DefaultCacheImporter() (*string, error) {
	// We don't have remote cache for v1
	return nil, nil
}

func (g *generalGraph) GetEntrypoint(buildContextDir string) ([]string, error) {
	if !g.Dev {
		return g.Entrypoint, nil
	}
	g.RuntimeEnviron[types.EnvdWorkDir] = fileutil.EnvdHomeDir(filepath.Base(buildContextDir))
	return []string{"horust"}, nil
}

func (g *generalGraph) CompileLLB(uid, gid int) (llb.State, error) {
	g.uid = uid
	g.gid = gid
	logrus.WithFields(logrus.Fields{
		"uid": g.uid,
		"gid": g.gid,
	}).Debug("compile LLB")

	base, err := g.compileBaseImage()
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get the base image")
	}
	aptMirror := g.compileUbuntuAPT(base)

	// prepare dev env: stable operations should be done here to make it cache friendly
	if g.Dev {
		dev := g.compileDevPackages(aptMirror)
		sshd := g.compileSSHD(dev)
		horust := g.installHorust(sshd)
		starship := g.compileStarship(horust)
		userGroup := g.compileUserGroup(starship)
		aptMirror = userGroup
	}

	lang, err := g.compileLanguage(aptMirror)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile language")
	}
	systemPackages := g.compileSystemPackages(aptMirror)
	merge := llb.Merge([]llb.State{
		base,
		llb.Diff(base, lang, llb.WithCustomName("[internal] prepare language")),
		llb.Diff(base, systemPackages, llb.WithCustomName("[internal] install system packages")),
	}, llb.WithCustomName("[internal] language environment and system packages"))
	packages := g.compileLanguagePackages(merge)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile language")
	}

	source, err := g.compileExtraSource(packages)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile extra source")
	}
	copy := g.compileCopy(source)

	// dev postprocessing: related to UID, which may not be cached
	if g.Dev {
		git := g.compileGit(copy)
		user := g.compileUserOwn(git)
		key, err := g.copySSHKey(user)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to copy ssh key")
		}
		shell, err := g.compileShell(key)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to compile shell")
		}
		prompt := g.compilePrompt(shell)
		entrypoint, err := g.compileEntrypoint(prompt)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to compile entrypoint")
		}
		vscode, err := g.compileVSCode(entrypoint)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to compile VSCode extensions")
		}
		copy = vscode
	}

	// it's necessary to exec `run` with the desired user
	run := g.compileRun(copy)
	mount := g.compileMountDir(run)

	g.Writer.Finish()
	return mount, nil
}
