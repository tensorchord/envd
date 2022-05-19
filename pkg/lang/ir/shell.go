package ir

import (
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/shell"
)

func (g *Graph) compileShell(root llb.State) (llb.State, error) {
	if g.Shell == shellZSH {
		return g.compileZSH(root)
	}
	return root, nil
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
			&llb.CopyInfo{CreateDestPath: true}, llb.WithUser(defaultUID))).
		File(llb.Mkfile(installPath, 0644, []byte(m.InstallScript()), llb.WithUser(defaultUID)))
	run := zshStage.Run(llb.Shlex(fmt.Sprintf("bash %s", installPath)),
		llb.WithCustomName("install oh-my-zsh"))
	return run.Root(), nil
}
