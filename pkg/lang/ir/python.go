package ir

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
)

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
