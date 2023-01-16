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
	"fmt"
	"strings"

	"github.com/moby/buildkit/client/llb"
)

// nolint:unparam
func (g generalGraph) compileCustomPython(baseStage llb.State) (llb.State, error) {
	aptStage := g.compileUbuntuAPT(baseStage)
	pypiMirrorStage := g.compilePyPIIndex(aptStage)
	systemStage := g.compileCustomSystemPackages(pypiMirrorStage)
	pypiStage := g.compileCustomPyPIPackages(systemStage)

	return pypiStage, nil
}

func (g generalGraph) compileCustomPyPIPackages(root llb.State) llb.State {
	if len(g.PyPIPackages) == 0 {
		return root
	}

	cacheDir := "/home/root/.cache"
	// Refer to https://github.com/moby/buildkit/blob/31054718bf775bf32d1376fe1f3611985f837584/frontend/dockerfile/dockerfile2llb/convert_runmount.go#L46
	cache := llb.Scratch().File(llb.Mkdir("/cache", 0755, llb.WithParents(true)),
		llb.WithCustomName("[internal] settings pip cache mount permissions"))

	for _, packages := range g.PyPIPackages {
		cmd := fmt.Sprintf("pip install %s", strings.Join(packages, " "))
		run := root.
			Run(llb.Shlex(cmd), llb.WithCustomNamef("pip install %s",
				strings.Join(packages, " ")))
		run.AddMount(cacheDir, cache,
			llb.AsPersistentCacheDir(g.CacheID(cacheDir), llb.CacheMountShared), llb.SourcePath("/cache"))
		root = run.Root()
	}
	return root
}

func (g generalGraph) compileCustomSystemPackages(root llb.State) llb.State {
	if len(g.SystemPackages) == 0 {
		return root
	}

	// Compose the package install command.
	var sb strings.Builder
	sb.WriteString("apt-get update && apt-get install -y --no-install-recommends --no-install-suggests --fix-missing")

	for _, pkg := range g.SystemPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}

	cacheDir := "/var/cache/apt"
	cacheLibDir := "/var/lib/apt"

	run := root.Run(llb.Shlexf(`bash -c "%s"`, sb.String()),
		llb.WithCustomNamef("apt-get install %s",
			strings.Join(g.SystemPackages, " ")))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir(g.CacheID(cacheDir), llb.CacheMountShared))
	run.AddMount(cacheLibDir, llb.Scratch(),
		llb.AsPersistentCacheDir(g.CacheID(cacheLibDir), llb.CacheMountShared))
	return run.Root()
}
