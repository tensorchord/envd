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
	"strings"

	"github.com/moby/buildkit/client/llb"
)

func (g Graph) compileCustomPython(aptStage llb.State) (llb.State, error) {
	pypiMirrorStage := g.compilePyPIIndex(aptStage)

	builtinSystemStage := pypiMirrorStage

	pypiStage := llb.Diff(builtinSystemStage,
		g.compileCustomPyPIPackages(builtinSystemStage),
		llb.WithCustomName("install PyPI packages"))
	systemStage := llb.Diff(builtinSystemStage, g.compileSystemPackages(builtinSystemStage),
		llb.WithCustomName("install system packages"))

	merged := llb.Merge([]llb.State{
		builtinSystemStage, systemStage, pypiStage,
	}, llb.WithCustomName("merging all components into one"))

	return merged, nil
}

func (g Graph) compileCustomPyPIPackages(root llb.State) llb.State {
	if len(g.PyPIPackages) == 0 {
		return root
	}

	cacheDir := "/home/root/.cache"

	// Compose the package install command.
	var sb strings.Builder
	// Always use the conda's pip.
	sb.WriteString("pip install")
	for _, pkg := range g.PyPIPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}

	cmd := sb.String()
	// Refer to https://github.com/moby/buildkit/blob/31054718bf775bf32d1376fe1f3611985f837584/frontend/dockerfile/dockerfile2llb/convert_runmount.go#L46
	cache := root.File(llb.Mkdir("/cache", 0755, llb.WithParents(true)),
		llb.WithCustomName("[internal] settings pip cache mount permissions"))
	run := root.
		Run(llb.Shlex(cmd), llb.WithCustomNamef("pip install %s",
			strings.Join(g.PyPIPackages, " ")))
	run.AddMount(cacheDir, cache,
		llb.AsPersistentCacheDir(g.CacheID(cacheDir), llb.CacheMountShared),
		llb.SourcePath("/cache"))
	return run.Root()
}
