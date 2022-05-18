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
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/tensorchord/envd/pkg/editor/vscode"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/shell"
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

func (g *Graph) compileBase() llb.State {
	var base llb.State
	if g.CUDA == nil && g.CUDNN == nil {
		base = llb.Image("docker.io/library/python:3.8")
	} else {
		base = g.compileCUDAPackages()
	}
	return base
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

func (g *Graph) compileShell(root llb.State) (llb.State, error) {
	if g.Shell == shellZSH {
		return g.compileZSH(root)
	}
	return root, nil
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

	cacheDir := "/home/envd/.cache/pip"
	cmd := sb.String()
	run := root.Run(llb.Shlex(cmd), llb.WithCustomNamef("pip install %s",
		strings.Join(g.PyPIPackages, " ")))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	return run.Root()
}

func (g Graph) compileBuiltinSystemPackages(root llb.State) llb.State {
	// TODO(gaocegege): Refactor it to avoid shell configuration in built-in system packages.
	// Do not need to install bash or sh since it is built-in
	if g.Shell == shellZSH {
		g.BuiltinSystemPackages = append(g.BuiltinSystemPackages, shellZSH)
	}

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

	run := root.Run(llb.Shlex(sb.String()),
		llb.WithCustomNamef("(built-in packages) apt-get install %s",
			strings.Join(g.BuiltinSystemPackages, " ")))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	run.AddMount(cacheLibDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheLibDir, llb.CacheMountShared))

	// TODO(gaocegege): Refactor user to a seperate stage.
	run = run.
		Run(llb.Shlex("groupadd -g 1000 envd"), llb.WithCustomName("create user group envd")).
		Run(llb.Shlex("useradd -p \"\" -u 1000 -g envd -s /bin/sh -m envd"), llb.WithCustomName("create user envd")).
		Run(llb.Shlex("adduser envd sudo"), llb.WithCustomName("add user envd to sudoers"))
	return llb.User("envd")(run.Root())
}

func (g Graph) compileSystemPackages(root llb.State) llb.State {
	if len(g.SystemPackages) == 0 {
		return root
	}

	// Compose the package install command.
	var sb strings.Builder
	sb.WriteString("sudo apt-get install -y --no-install-recommends")

	for _, pkg := range g.SystemPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}

	cacheDir := "/var/cache/apt"
	cacheLibDir := "/var/lib/apt"

	run := root.Run(llb.Shlex(sb.String()),
		llb.WithCustomNamef("(user-defined packages) apt-get install %s",
			strings.Join(g.SystemPackages, " ")))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	run.AddMount(cacheLibDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheLibDir, llb.CacheMountShared))
	return run.Root()
}

func (g Graph) copyEnvdSSHServer() llb.State {
	// TODO(gaocegege): Remove global var ssh image.
	run := llb.Image(viper.GetString(flag.FlagSSHImage)).
		File(llb.Copy(llb.Image(viper.GetString(flag.FlagSSHImage)),
			"usr/bin/envd-ssh", "/var/envd/bin/envd-ssh",
			&llb.CopyInfo{CreateDestPath: true}), llb.WithCustomName("install envd-ssh"))
	return run
}

func (g Graph) compileVSCode() (*llb.State, error) {
	if len(g.VSCodePlugins) == 0 {
		return nil, nil
	}
	inputs := []llb.State{}
	for _, p := range g.VSCodePlugins {
		vscodeClient := vscode.NewClient()
		g.Writer.LogVSCodePlugin(p, compileui.ActionStart, false)
		if cached, err := vscodeClient.DownloadOrCache(p); err != nil {
			return nil, err
		} else {
			g.Writer.LogVSCodePlugin(p, compileui.ActionEnd, cached)
		}
		ext := llb.Scratch().File(llb.Copy(llb.Local(flag.FlagCacheDir),
			vscodeClient.PluginPath(p),
			"/home/envd/.vscode-server/extensions/"+p.String(),
			&llb.CopyInfo{CreateDestPath: true}),
			llb.WithCustomNamef("install vscode plugin %s", p.String()))
		inputs = append(inputs, ext)
	}
	layer := llb.Merge(inputs, llb.WithCustomName("merging plugins for vscode"))
	return &layer, nil
}

func (g *Graph) compileJupyter() {
	if g.JupyterConfig != nil {
		g.PyPIPackages = append(g.PyPIPackages, "jupyter")
	}
}

func (g Graph) compileUbuntuAPT(root llb.State) llb.State {
	if g.UbuntuAPTSource != nil {
		logrus.WithField("source", *g.UbuntuAPTSource).Debug("using custom APT source")
		aptSource := llb.Scratch().
			File(llb.Mkdir(filepath.Dir(aptSourceFilePath), 0755, llb.WithParents(true)), llb.WithCustomName("create apt source dir")).
			File(llb.Mkfile(aptSourceFilePath, 0644, []byte(*g.UbuntuAPTSource)), llb.WithCustomName("create apt source file"))
		return llb.Merge([]llb.State{root, aptSource}, llb.WithCustomName("add apt source"))
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
		return llb.Merge([]llb.State{root, aptSource}, llb.WithCustomName("add PyPI mirror"))
	}
	return root
}

func (g Graph) compileZSH(root llb.State) (llb.State, error) {
	installPath := "/home/envd/install.sh"
	m := shell.NewManager()
	g.Writer.LogZSH(compileui.ActionStart, false)
	if cached, err := m.DownloadOrCache(); err != nil {
		return llb.State{}, errors.Wrap(err, "failed to download oh-my-zsh")
	} else {
		g.Writer.LogZSH(compileui.ActionEnd, cached)
	}
	zshStage := root.
		File(llb.Copy(llb.Local(flag.FlagCacheDir), "oh-my-zsh", "/home/envd/.oh-my-zsh",
			&llb.CopyInfo{CreateDestPath: true})).
		File(llb.Mkfile(installPath, 0644, []byte(m.InstallScript())))
	run := zshStage.Run(llb.Shlex(fmt.Sprintf("bash %s", installPath)),
		llb.WithCustomName("install oh-my-zsh"))
	return run.Root(), nil
}

func (g Graph) compileRun(root llb.State) llb.State {
	if len(g.Exec) == 0 {
		return root
	} else if len(g.Exec) == 1 {
		return root.Run(llb.Shlex(g.Exec[0])).Root()
	}

	run := root.Run(llb.Shlex(g.Exec[0]))
	for _, c := range g.Exec[1:] {
		run = run.Run(llb.Shlex(c))
	}
	return run.Root()
}
