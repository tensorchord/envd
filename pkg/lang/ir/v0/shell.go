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
)

func (g *generalGraph) compileShell(root llb.State) (llb.State, error) {
	if g.Shell == shellZSH {
		g.RuntimeEnviron["SHELL"] = "/usr/bin/zsh"
		return g.compileZSH(root)
	}
	g.RuntimeEnviron["SHELL"] = "/usr/bin/bash"
	return root, nil
}

func (g *generalGraph) compileCondaShell(root llb.State) llb.State {
	rcPath := fileutil.EnvdHomeDir(".bashrc")
	if g.Shell == shellZSH {
		rcPath = fileutil.EnvdHomeDir(".zshrc")
	}
	run := root.
		Run(llb.Shlexf("bash -c \"%s\"", g.condaInitShell(g.Shell)),
			llb.WithCustomNamef("[internal] init conda %s env", g.Shell)).
		Run(llb.Shlexf(`bash -c 'echo "source %s/activate envd" >> %s'`, condaBinDir, rcPath),
			llb.WithCustomNamef("[internal] add conda environment to %s", rcPath))
	return run.Root()
}

func (g *generalGraph) compilePrompt(root llb.State) llb.State {
	// skip this for customized image
	if g.Image != nil {
		return root
	}
	// starship config
	config := root.
		File(llb.Mkdir(defaultConfigDir, 0755, llb.WithParents(true)),
			llb.WithCustomName("[internal] creating config dir")).
		File(llb.Mkfile(starshipConfigPath, 0644, []byte(starshipConfig), llb.WithUIDGID(g.uid, g.gid)),
			llb.WithCustomName("[internal] setting prompt starship config"))

	run := config.Run(llb.Shlexf(`bash -c 'echo "eval \"\$(starship init bash)\"" >> %s'`, fileutil.EnvdHomeDir(".bashrc")),
		llb.WithCustomName("[internal] setting prompt bash config")).Root()

	if g.Shell == shellZSH {
		run = run.Run(
			llb.Shlexf(`bash -c 'echo "eval \"\$(starship init zsh)\"" >> %s'`, fileutil.EnvdHomeDir(".zshrc")),
			llb.WithCustomName("[internal] setting prompt zsh config")).Root()
	}
	return run
}

func (g generalGraph) compileZSH(root llb.State) (llb.State, error) {
	installPath := fileutil.EnvdHomeDir("install.sh")
	zshrcPath := fileutil.EnvdHomeDir(".zshrc")
	ohMyZSHPath := fileutil.EnvdHomeDir(".oh-my-zsh")
	m := shell.NewManager()
	g.Writer.LogZSH(compileui.ActionStart, false)
	if cached, err := m.DownloadOrCache(); err != nil {
		return llb.State{}, errors.Wrap(err, "failed to download oh-my-zsh")
	} else {
		g.Writer.LogZSH(compileui.ActionEnd, cached)
	}
	zshStage := root.
		File(llb.Copy(llb.Local(flag.FlagCacheDir), "oh-my-zsh", ohMyZSHPath,
			&llb.CopyInfo{CreateDestPath: true}, llb.WithUIDGID(g.uid, g.gid))).
		File(llb.Mkfile(installPath,
			0644, []byte(m.InstallScript()), llb.WithUIDGID(g.uid, g.gid)))
	zshrc := zshStage.Run(llb.Shlexf("bash %s", installPath),
		llb.WithCustomName("[internal] install oh-my-zsh")).
		File(llb.Mkfile(zshrcPath,
			0644, []byte(m.ZSHRC()), llb.WithUIDGID(g.uid, g.gid)))
	return zshrc, nil
}
