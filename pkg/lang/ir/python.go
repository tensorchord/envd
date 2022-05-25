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
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
)

func (g Graph) compilePyPIPackages(root llb.State) llb.State {
	if len(g.PyPIPackages) == 0 {
		return root
	}

	cacheDir := "/home/envd/.cache/pip"

	// Compose the package install command.
	var sb strings.Builder
	// Wait until https://github.com/moby/buildkit/commit/31054718bf775bf32d1376fe1f3611985f837584 is released in v0.10.4
	sb.WriteString("sudo chown -R 1000:1000 ")
	sb.WriteString(filepath.Dir(cacheDir))
	sb.WriteString("&& pip install")
	for _, pkg := range g.PyPIPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}

	cmd := sb.String()
	// Wait until https://github.com/moby/buildkit/commit/31054718bf775bf32d1376fe1f3611985f837584 is released in v0.10.4
	// run := root.File(llb.Mkdir(cacheDir,
	// 	0755, llb.WithParents(true), llb.WithUIDGID(1000, 1000)), llb.WithCustomName("[internal] settings pip cache mount permissions")).
	run := root.
		Run(llb.Shlex(fmt.Sprintf(`sh -c "%s"`, cmd)), llb.WithCustomNamef("pip install %s",
			strings.Join(g.PyPIPackages, " ")))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	return run.Root()
}

func (g Graph) compilePyPIMirror(root llb.State) llb.State {
	if g.PyPIMirror != nil {
		logrus.WithField("mirror", *g.PyPIMirror).Debug("using custom PyPI mirror")
		content := fmt.Sprintf(pypiConfigTemplate, *g.PyPIMirror)
		pypiMirror := llb.Scratch().
			File(llb.Mkdir(filepath.Dir(pypiMirrorFilePath),
				0755, llb.WithParents(true), llb.WithUIDGID(defaultUID, defaultGID))).
			File(llb.Mkfile(pypiMirrorFilePath,
				0644, []byte(content), llb.WithUIDGID(defaultUID, defaultGID)))
		return llb.Merge([]llb.State{root, pypiMirror}, llb.WithCustomName("add PyPI mirror"))
	}
	return root
}
