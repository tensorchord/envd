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
	"fmt"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"

	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/shell"
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
map_symbol = true
disabled = false
`
)

func (g *Graph) compileShell(root llb.State) (llb.State, error) {
	// starship config
	config := llb.Scratch().
		File(llb.Mkdir(filepath.Dir(defaultConfigDir), 0755, llb.WithParents(true)),
			llb.WithCustomName("[internal] creating config dir")).
		File(llb.Mkfile(starshipConfigPath, 0644, []byte(starshipConfig)),
			llb.WithCustomName("[internal] setting prompt config"))
	config.Run(
		llb.Shlex(`bash -c 'echo "eval \"\$(starship init bash)\"" >> /home/envd/.bashrc'`),
		llb.WithCustomName("[internal] setting prompt config"))
	root = llb.Merge([]llb.State{root, config}, llb.WithCustomName("[internal] setting prompt config"))

	if g.Shell == shellZSH {
		return g.compileZSH(root)
	}
	return root, nil
}

func (g Graph) compileZSH(root llb.State) (llb.State, error) {
	installPath := "/home/envd/install.sh"
	zshrcPath := "/home/envd/.zshrc"
	ohMyZSHPath := "/home/envd/.oh-my-zsh"
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
	zshrc := zshStage.Run(llb.Shlex(fmt.Sprintf("bash %s", installPath)),
		llb.WithCustomName("install oh-my-zsh")).
		File(llb.Mkfile(zshrcPath,
			0644, []byte(m.ZSHRC()), llb.WithUIDGID(g.uid, g.gid)))
	config := zshrc.Run(
		llb.Shlex(`bash -c 'echo "eval \"\$(starship init zsh)\"" >> /home/envd/.zshrc'`),
		llb.WithCustomName("[interna] setting prompt config"))
	return config.Root(), nil
}
