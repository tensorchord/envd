// Copyright 2022 The MIDI Authors
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
	"context"

	"github.com/moby/buildkit/client/llb"
)

// A Graph contains the state,
// such as its call stack and thread-local storage.
type Graph struct {
	OS             string
	Language       string
	PyPIPackages   []string
	SystemPackages []string
}

func NewGraph() *Graph {
	return &Graph{
		OS:             osDefault,
		Language:       languageDefault,
		PyPIPackages:   []string{},
		SystemPackages: []string{},
	}
}

var DefaultGraph = NewGraph()

func Base(os, language string) {
	DefaultGraph.Language = language
	DefaultGraph.OS = os
}

func PyPIPackage(deps []string) {
	DefaultGraph.PyPIPackages = append(DefaultGraph.PyPIPackages, deps...)
}

func Compile(ctx context.Context) (*llb.Definition, error) {
	state := DefaultGraph.Compile()
	// TODO(gaocegege): Support multi platform.
	def, err := state.Marshal(ctx, llb.LinuxAmd64)
	if err != nil {
		return nil, err
	}
	return def, nil
}

func (g Graph) Compile() llb.State {
	// TODO(gaocegege): Support more OS and langs.
	base := llb.Image("docker.io/library/python:3.8")
	return g.compilePyPIPackages(base)
}

func (g Graph) compilePyPIPackages(root llb.State) llb.State {
	// TODO(gaocegege): Support per-user config to keep the mirror.
	cmd := "pip install -i https://mirror.sjtu.edu.cn/pypi/web/simple"
	cacheDir := "/root/.cache/pip"
	for _, pkg := range g.PyPIPackages {
		cmd = cmd + " " + pkg
	}
	run := root.Run(llb.Shlex(cmd))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	return run.Root()
}
