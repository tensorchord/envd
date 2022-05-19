package ir

import (
	"github.com/moby/buildkit/client/llb"
	"github.com/tensorchord/envd/pkg/editor/vscode"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/progress/compileui"
)

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
			&llb.CopyInfo{
				CreateDestPath: true,
			}, llb.WithUser(defaultUID)),
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
