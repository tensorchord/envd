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
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	servertypes "github.com/tensorchord/envd-server/api/types"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
	"github.com/tensorchord/envd/pkg/version"
)

func NewGraph() *Graph {
	runtimeGraph := RuntimeGraph{
		RuntimeCommands: make(map[string]string),
		RuntimeEnviron:  make(map[string]string),
	}
	langVersion := languageVersionDefault
	return &Graph{
		Image: defaultImage,
		Language: Language{
			Name:    languageDefault,
			Version: &langVersion,
		},
		CUDA:    nil,
		CUDNN:   CUDNNVersionDefault,
		NumGPUs: 0,

		PyPIPackages:    []string{},
		RPackages:       []string{},
		JuliaPackages:   []string{},
		SystemPackages:  []string{},
		Exec:            []RunBuildCommand{},
		UserDirectories: []string{},
		RuntimeEnvPaths: []string{types.DefaultPathEnv()},
		Shell:           shellBASH,
		RuntimeGraph:    runtimeGraph,
	}
}

var DefaultGraph = NewGraph()

func GPUEnabled() bool {
	return DefaultGraph.GPUEnabled()
}

func NumGPUs() int {
	return DefaultGraph.NumGPUs
}

func Compile(ctx context.Context, envName string, pub string) (*llb.Definition, error) {
	w, err := compileui.New(ctx, os.Stdout, "auto")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create compileui")
	}
	DefaultGraph.Writer = w
	DefaultGraph.EnvironmentName = envName
	DefaultGraph.PublicKeyPath = pub

	uid, gid, err := getUIDGID()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get uid/gid")
	}
	state, err := DefaultGraph.Compile(uid, gid)
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

func Labels() (map[string]string, error) {
	return DefaultGraph.Labels()
}

func ExposedPorts() (map[string]struct{}, error) {
	return DefaultGraph.ExposedPorts()
}

func CompileEntrypoint(buildContextDir string) ([]string, error) {
	return DefaultGraph.GetEntrypoint(buildContextDir)
}

func CompileEnviron() []string {
	// Add PATH and LC_ALL.
	return append(DefaultGraph.EnvString(),
		"PATH="+strings.Join(DefaultGraph.RuntimeEnvPaths, ":"),
		"LC_ALL=en_US.UTF-8",
		"LANG=C.UTF-8",
	)
}

func IsDev() bool {
	return DefaultGraph.DevTools
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
		labels[types.ImageLabelCUDNN] = g.CUDNN
	}
	labels[types.ImageLabelVendor] = types.ImageVendorEnvd
	code, err := g.RuntimeGraph.Dump()
	if err != nil {
		return labels, err
	}
	labels[types.RuntimeGraphCode] = code

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

	labels[types.ImageLabelContainerName] = string(g.EnvironmentName)
	return labels, nil
}

func (g Graph) ExposedPorts() (map[string]struct{}, error) {
	ports := make(map[string]struct{})

	// only expose ports for dev env
	if !g.DevTools {
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

func (g Graph) EnvString() []string {
	var envs []string
	for k, v := range g.RuntimeEnviron {
		envs = append(envs, fmt.Sprintf("%s=%s", k, v))
	}
	return envs
}

func (g Graph) DefaultCacheImporter() (*string, error) {
	// The base remote cache should work for all languages.
	var res string
	if g.CUDA != nil {
		res = fmt.Sprintf(
			"type=registry,ref=docker.io/%s/python-cache:envd-%s-cuda-%s-cudnn-%s",
			viper.GetString(flag.FlagDockerOrganization),
			version.GetVersionForImageTag(), *g.CUDA, g.CUDNN)
	} else {
		res = fmt.Sprintf(
			"type=registry,ref=docker.io/%s/python-cache:envd-%s",
			viper.GetString(flag.FlagDockerOrganization),
			version.GetVersionForImageTag())
	}
	return &res, nil
}

func (g *Graph) GetEntrypoint(buildContextDir string) ([]string, error) {
	if !g.DevTools {
		return g.Entrypoint, nil
	}
	g.RuntimeEnviron[types.EnvdWorkDir] = fileutil.EnvdHomeDir(filepath.Base(buildContextDir))
	return []string{"horust"}, nil
}

func (g Graph) Compile(uid, gid int) (llb.State, error) {
	g.uid = uid

	// TODO(gaocegege): Remove the hack for https://github.com/tensorchord/envd/issues/370
	g.gid = 1001
	logrus.WithFields(logrus.Fields{
		"uid": g.uid,
		"gid": g.gid,
	}).Debug("compile LLB")

	base, err := g.compileBaseImage()
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get the base image")
	}

	// prepare dev env: stable operations should be done here to make it cache friendly
	if g.DevTools {
		dev := g.compileDevPackages(base)
		sshd := g.compileSSHD(dev)
		horust := g.installHorust(sshd)
		userGroup := g.compileUserGroup(horust)
		shell, err := g.compileShell(userGroup)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to compile shell")
		}
		base = shell
	}

	lang, err := g.compileLanguage(base)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile language")
	}
	aptMirror := g.compileUbuntuAPT(base)
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
	if g.DevTools {
		prompt := g.compilePrompt(copy)
		git := g.compileGit(prompt)
		user := g.compileUserOwn(git)
		key, err := g.copySSHKey(user)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to copy ssh key")
		}
		entrypoint, err := g.compileEntrypoint(key)
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to compile entrypoint")
		}
		vscode, err := g.compileVSCode()
		if err != nil {
			return llb.State{}, errors.Wrap(err, "failed to compile VSCode extensions")
		}
		copy = llb.Merge([]llb.State{
			copy,
			vscode,
			entrypoint,
		}, llb.WithCustomName("[internal] final dev environment"))
	}

	// it's necessary to exec `run`` with the desired user
	run := g.compileRun(copy)
	mount := g.compileMountDir(run)

	g.Writer.Finish()
	return mount, nil
}
