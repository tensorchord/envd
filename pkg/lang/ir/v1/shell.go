// Copyright 2023 The envd Authors
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

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"

	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/shell"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

const (
	starshipConfig = `
[container]
format = "[$symbol \\[envd\\]]($style)"

[sudo]
disabled = false
symbol = "sudo "

[python]
symbol = "Py "

[status]
format = '[\[$status:$common_meaning$signal_name\]]($style) '
disabled = false

[git_branch]
symbol = "git "

[git_commit]
tag_symbol = " tag "

[git_status]
ahead = ">"
behind = "<"
diverged = "<>"
renamed = "r"
deleted = "x"
`

	fishVersion  = "4.0.1"
	fishAssetURL = "https://github.com/fish-shell/fish-shell/releases/download/%[1]s/fish-static-$(uname -m)-%[1]s.tar.xz"
)

func (g *generalGraph) compileShell(root llb.State) (_ llb.State, err error) {
	g.RuntimeEnviron["SHELL"] = "bash"
	switch g.Shell {
	case shellZSH:
		g.RuntimeEnviron["SHELL"] = "/usr/bin/zsh"
		if root, err = g.compileZSH(root); err != nil {
			return llb.State{}, err
		}
	case shellFish:
		g.RuntimeEnviron["SHELL"] = "/usr/bin/fish"
		root = g.compileFish(root)
	}
	root = g.compileCondaShell(root)
	root = g.compilePixiShell(root)
	return root, nil
}

func (g generalGraph) compilePixiShell(root llb.State) llb.State {
	if g.PixiConfig == nil {
		return root
	}

	switch g.Shell {
	case shellZSH:
		root = root.Run(
			llb.Shlex(`sh -c 'echo "eval \"\$(pixi completion --shell zsh)\"" >> ~/.zshrc'`),
			llb.WithCustomName("[internal] setting pixi zsh config"),
		).Root()
	case shellFish:
		root = root.Run(
			llb.Shlex(`sh -c 'echo "pixi completion --shell fish | source" >> ~/.config/fish/config.fish'`),
			llb.WithCustomName("[internal] setting pixi fish config"),
		).Root()
	case shellBASH:
		root = root.Run(
			llb.Shlex(`sh -c 'echo "eval \"\$(pixi completion --shell bash)\"" >> ~/.bashrc'`),
			llb.WithCustomName("[internal] setting pixi bash config"),
		).Root()
	}
	return root
}

func (g *generalGraph) compileCondaShell(root llb.State) llb.State {
	if g.CondaConfig == nil {
		return root
	}

	findDir := fileutil.DefaultHomeDir
	if g.Dev {
		findDir = fileutil.EnvdHomeDir
	}
	rcPath := findDir(".bashrc")
	activateFile := "activate"
	switch g.Shell {
	case shellZSH:
		rcPath = findDir(".zshrc")
	case shellFish:
		rcPath = findDir(".config/fish/config.fish")
		activateFile = "activate.fish"
	}
	run := root.
		Run(llb.Shlexf("bash -c \"%s\"", g.condaInitShell(g.Shell)),
			llb.WithCustomNamef("[internal] init conda %s env", g.Shell)).
		Run(llb.Shlexf(`bash -c 'echo "source %s/%s envd" >> %s'`, condaBinDir, activateFile, rcPath),
			llb.WithCustomNamef("[internal] add conda environment to %s", rcPath))
	return run.Root()
}

func (g *generalGraph) compilePrompt(root llb.State) llb.State {
	// starship config
	config := root.
		File(llb.Mkdir(defaultConfigDir, 0755, llb.WithParents(true)),
			llb.WithCustomName("[internal] creating config dir")).
		File(llb.Mkfile(starshipConfigPath, 0644, []byte(starshipConfig), llb.WithUIDGID(g.uid, g.gid)),
			llb.WithCustomName("[internal] setting prompt starship config"))

	run := config.Run(llb.Shlexf(`bash -c 'echo "eval \"\$(starship init bash)\"" >> %s'`, fileutil.EnvdHomeDir(".bashrc")),
		llb.WithCustomName("[internal] setting prompt bash config")).Root()

	switch g.Shell {
	case shellZSH:
		run = run.Run(
			llb.Shlexf(`bash -c 'echo "eval \"\$(starship init zsh)\"" >> %s'`, fileutil.EnvdHomeDir(".zshrc")),
			llb.WithCustomName("[internal] setting prompt zsh config")).Root()
	case shellFish:
		run = run.Run(
			llb.Shlexf(`bash -c 'echo "starship init fish | source" >> %s'`, fileutil.EnvdHomeDir(".config/fish/config.fish")),
			llb.WithCustomName("[internal] setting prompt fish config")).Root()
	}
	return run
}

func (g generalGraph) compileZSH(root llb.State) (llb.State, error) {
	installPath := fileutil.EnvdHomeDir("install.sh")
	zshrcPath := fileutil.EnvdHomeDir(".zshrc")
	ohMyZSHPath := fileutil.EnvdHomeDir(".oh-my-zsh")
	m := shell.NewManager()
	g.Writer.LogZSH(compileui.ActionStart, false)
	cached, err := m.DownloadOrCache()
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to download oh-my-zsh")
	}
	g.Writer.LogZSH(compileui.ActionEnd, cached)
	zshStage := root.
		File(llb.Copy(llb.Local(flag.FlagCacheDir), "oh-my-zsh", ohMyZSHPath,
			&llb.CopyInfo{CreateDestPath: true})).
		File(llb.Mkfile(installPath, 0666, []byte(m.InstallScript())))
	zshrc := zshStage.Run(llb.Shlexf("bash %s", installPath),
		llb.WithCustomName("[internal] install oh-my-zsh")).
		File(llb.Mkfile(zshrcPath, 0666, []byte(m.ZSHRC())))
	return zshrc, nil
}

func (g generalGraph) compileFish(root llb.State) llb.State {
	base := llb.Image(builderImage)
	url := fmt.Sprintf(fishAssetURL, fishVersion)
	builder := base.Run(
		llb.Shlexf(`sh -c "wget -qO- %s | tar -xJf - -C /tmp || exit 1"`, url),
		llb.WithCustomName("[internal] download fish shell"),
	).Root()
	root = root.File(
		llb.Copy(builder, "/tmp/fish", "/usr/bin/fish"),
		llb.WithCustomName("[internal] copy fish shell from the builder image")).
		Run(llb.Shlex(`sh -c "echo yes | fish --install"`),
			llb.WithCustomName("[internal] install fish shell")).Root()

	return root
}
