// Copyright 2022 The MIDI Authors
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
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/MIDI/pkg/flag"
	"github.com/tensorchord/MIDI/pkg/vscode"
)

func NewGraph() *Graph {
	return &Graph{
		OS:       osDefault,
		Language: languageDefault,
		CUDA:     nil,
		CUDNN:    nil,
		BuiltinSystemPackages: []string{
			"curl",
			"openssh-client",
		},

		PyPIPackages:   []string{},
		SystemPackages: []string{},
		Exec:           []llb.State{},
	}
}

var DefaultGraph = NewGraph()

func GPUEnabled() bool {
	return DefaultGraph.CUDA != nil
}

func Compile(ctx context.Context) (*llb.Definition, error) {
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

	builtinSystemStage := g.compileBuiltinSystemPackages(aptStage)
	pypiMirrorStage := g.compilePyPIMirror(builtinSystemStage)
	pypiStage := llb.Diff(aptStage, g.compilePyPIPackages(pypiMirrorStage))

	systemStage := llb.Diff(aptStage, g.compileSystemPackages(aptStage))

	sshStage := g.copyMidiSSHServer()

	vscodeStage, err := g.compileVSCode()
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get vscode plugins")
	}

	merged := llb.Merge([]llb.State{
		aptStage, systemStage, pypiStage, sshStage, vscodeStage,
	})
	return merged, nil
}

func (g *Graph) compileBase() llb.State {
	if g.CUDA == nil && g.CUDNN == nil {
		return llb.Image("docker.io/library/python:3.8")
	}
	return g.compileCUDAPackages()
}

func (g *Graph) compileCUDAPackages() llb.State {
	root := llb.Image(
		fmt.Sprintf("nvidia/cuda:%s.0-cudnn%s-devel-%s", *g.CUDA, *g.CUDNN, g.OS))
	g.BuiltinSystemPackages = append(g.BuiltinSystemPackages, []string{
		g.Language,
		fmt.Sprintf("%s-pip", g.Language),
	}...)
	return root
}

// Deprecated: Use compileCUDAPackages instead.
func (g *Graph) compileCUDAPackagesDeprecated() llb.State {
	root := llb.Image(fmt.Sprintf("nvidia/cuda:%s.0-base-%s", *g.CUDA, g.OS))
	env := root.AddEnv("DEBIAN_FRONTEND", "noninteractive").
		AddEnv("LD_LIBRARY_PATH", "/usr/local/cuda-11.0/targets/x86_64-linux/lib:/usr/local/cuda/extras/CUPTI/lib64:/usr/local/cuda/lib64:$LD_LIBRARY_PATH").
		AddEnv("LANG", "C.UTF-8")

	g.BuiltinSystemPackages = append(g.BuiltinSystemPackages, []string{
		"build-essential",
		"curl",
		"libfreetype6-dev",
		"libhdf5-serial-dev",
		"libzmq3-dev",
		"pkg-config",
		"software-properties-common",
		"unzip",
		// Python
		g.Language,
		fmt.Sprintf("%s-pip", g.Language),
		// CUDA and CUDNN
		fmt.Sprintf("cuda-command-line-tools-%s", *g.CUDA),
		fmt.Sprintf("libcublas-%s", *g.CUDA),
		fmt.Sprintf("cuda-nvrtc-%s", *g.CUDA),
		fmt.Sprintf("libcufft-%s", *g.CUDA),
		fmt.Sprintf("libcurand-%s", *g.CUDA),
		fmt.Sprintf("libcusolver-%s", *g.CUDA),
		fmt.Sprintf("libcusparse-%s", *g.CUDA),
		fmt.Sprintf("libcudnn8=%s+cuda%s", *g.CUDNN, *g.CUDA),
	}...)

	installed := g.compileBuiltinSystemPackages(env)

	run := installed.Run(llb.Shlex(`ln -s /usr/local/cuda/lib64/stubs/libcuda.so 
	/usr/local/cuda/lib64/stubs/libcuda.so.1
	&& echo "/usr/local/cuda/lib64/stubs" > /etc/ld.so.conf.d/z-cuda-stubs.conf
	&& ldconfig`))
	return run.Root()
}

func (g Graph) compilePyPIPackages(root llb.State) llb.State {
	if len(g.PyPIPackages) == 0 {
		return root
	}

	// Compose the package install command.
	var sb strings.Builder
	// TODO(gaocegege): Support per-user config to keep the mirror.
	sb.WriteString("pip install")
	for _, pkg := range g.PyPIPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}

	cacheDir := "/root/.cache/pip"

	run := root.Run(llb.Shlex(sb.String()))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	return run.Root()
}

func (g Graph) compileBuiltinSystemPackages(root llb.State) llb.State {
	if len(g.BuiltinSystemPackages) == 0 {
		return root
	}

	// Compose the package install command.
	var sb strings.Builder
	sb.WriteString(
		"sh -c \"apt-get update && apt-get install -y --no-install-recommends")
	for _, pkg := range g.BuiltinSystemPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}
	sb.WriteString("\"")

	cacheDir := "/var/cache/apt"
	cacheLibDir := "/var/lib/apt"

	run := root.Run(llb.Shlex(sb.String()))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	run.AddMount(cacheLibDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheLibDir, llb.CacheMountShared))
	return run.Root()
}

func (g Graph) compileSystemPackages(root llb.State) llb.State {
	if len(g.SystemPackages) == 0 {
		return root
	}

	// Compose the package install command.
	var sb strings.Builder
	// TODO(gaocegege): Support per-user config to keep the mirror.
	sb.WriteString("apt install")

	for _, pkg := range g.SystemPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}

	cacheDir := "/var/cache/apt"
	cacheLibDir := "/var/lib/apt"

	run := root.Run(llb.Shlex(sb.String()))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	run.AddMount(cacheLibDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheLibDir, llb.CacheMountShared))
	return run.Root()
}

func (g Graph) copyMidiSSHServer() llb.State {
	run := llb.Scratch().File(llb.Copy(llb.Local(flag.FlagContextDir),
		"examples/ssh_keypairs/public.pub", "/var/midi/remote/authorized_keys",
		&llb.CopyInfo{CreateDestPath: true})).
		File(llb.Copy(llb.Local(flag.FlagContextDir),
			"bin/midi-ssh", "/var/midi/bin/midi-ssh",
			&llb.CopyInfo{CreateDestPath: true}))
	return run
}

func (g Graph) compileVSCode() (llb.State, error) {
	inputs := []llb.State{}
	for _, p := range g.VSCodePlugins {
		vscodeClient := vscode.NewClient()
		if err := vscodeClient.DownloadOrCache(p); err != nil {
			return llb.State{}, err
		}
		ext := llb.Scratch().File(llb.Copy(llb.Local(flag.FlagCacheDir),
			vscodeClient.PluginPath(p),
			"/root/.vscode-server/extensions/"+p.String(),
			&llb.CopyInfo{CreateDestPath: true}))
		inputs = append(inputs, ext)
	}
	return llb.Merge(inputs), nil
}

func (g Graph) compileUbuntuAPT(root llb.State) llb.State {
	if g.UbuntuAPTSource != nil {
		logrus.WithField("source", *g.UbuntuAPTSource).Debug("using custom APT source")
		aptSource := llb.Scratch().
			File(llb.Mkdir(filepath.Dir(aptSourceFilePath), 0755, llb.WithParents(true))).
			File(llb.Mkfile(aptSourceFilePath, 0644, []byte(*g.UbuntuAPTSource)))
		return llb.Merge([]llb.State{root, aptSource})
	}
	return root
}

func (g Graph) compilePyPIMirror(root llb.State) llb.State {
	if g.PyPIMirror != nil {
		logrus.WithField("mirror", *g.PyPIMirror).Debug("using custom PyPI mirror")
		content := fmt.Sprintf(pypiConfigTemplate, *g.PyPIMirror)
		aptSource := llb.Scratch().
			File(llb.Mkdir(filepath.Dir(pypiMirrorFilePath), 0755, llb.WithParents(true))).
			File(llb.Mkfile(pypiMirrorFilePath, 0644, []byte(content)))
		return llb.Merge([]llb.State{root, aptSource})
	}
	return root
}
