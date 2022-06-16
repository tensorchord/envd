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

	cacheDir := "/home/envd/.cache"

	// Compose the package install command.
	var sb strings.Builder
	if g.CondaEnabled() {
		sb.WriteString("/opt/conda/bin/conda run -n envd pip install")
	} else {
		sb.WriteString("pip install --no-warn-script-location")
	}
	for _, pkg := range g.PyPIPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}

	cmd := sb.String()
	root = llb.User("envd")(root)
	// Refer to https://github.com/moby/buildkit/blob/31054718bf775bf32d1376fe1f3611985f837584/frontend/dockerfile/dockerfile2llb/convert_runmount.go#L46
	cache := root.File(llb.Mkdir("/cache",
		0755, llb.WithParents(true), llb.WithUIDGID(g.uid, g.gid)), llb.WithCustomName("[internal] settings pip cache mount permissions"))
	run := root.
		Run(llb.Shlex(cmd), llb.WithCustomNamef("pip install %s",
			strings.Join(g.PyPIPackages, " ")))
	run.AddMount(cacheDir, cache,
		llb.AsPersistentCacheDir(g.CacheID(cacheDir), llb.CacheMountShared), llb.SourcePath("/cache"))
	return run.Root()
}

func (g Graph) compilePyPIIndex(root llb.State) llb.State {
	if g.PyPIIndexURL != nil {
		logrus.WithField("index", *g.PyPIIndexURL).Debug("using custom PyPI index")
		var extraIndex string
		if g.PyPIExtraIndexURL != nil {
			logrus.WithField("index", *g.PyPIIndexURL).Debug("using extra PyPI index")
			extraIndex = "extra-index-url=" + *g.PyPIExtraIndexURL
		}
		content := fmt.Sprintf(pypiConfigTemplate, *g.PyPIIndexURL, extraIndex)
		pypiMirror := root.
			File(llb.Mkdir(filepath.Dir(pypiIndexFilePath),
				0755, llb.WithParents(true), llb.WithUIDGID(g.uid, g.gid)),
				llb.WithCustomName("[internal] settings PyPI index")).
			File(llb.Mkfile(pypiIndexFilePath,
				0644, []byte(content), llb.WithUIDGID(g.uid, g.gid)),
				llb.WithCustomName("[internal] settings PyPI index"))
		return pypiMirror
	}
	return root
}
