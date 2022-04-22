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
	"fmt"

	"github.com/moby/buildkit/client/llb"
)

// A Graph contains the state,
// such as its call stack and thread-local storage.
type Graph struct {
	OS       string
	Language string
	CUDA     *string
	CUDNN    *string

	BuiltinSystemPackages []string
	PyPIPackages          []string
	SystemPackages        []string

	Exec []llb.State
}

func NewGraph() *Graph {
	return &Graph{
		OS:                    osDefault,
		Language:              languageDefault,
		CUDA:                  nil,
		CUDNN:                 nil,
		BuiltinSystemPackages: []string{},

		PyPIPackages:   []string{},
		SystemPackages: []string{},
		Exec:           []llb.State{},
	}
}

var DefaultGraph = NewGraph()

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
	base := g.compileBase()
	system := g.compileSystemPackages(base)
	pypi := g.compilePyPIPackages(base)
	merged := llb.Merge([]llb.State{
		system, pypi,
	})
	return merged
}

func (g *Graph) compileBase() llb.State {
	if g.CUDA == nil && g.CUDNN == nil {
		return llb.Image("docker.io/library/python:3.8")
	}
	return g.compileCUDAPackages()
}

func (g Graph) compileCUDAPackages() llb.State {
	root := llb.Image(
		fmt.Sprintf("nvidia/cuda:%s.0-cudnn%s-devel-%s", *g.CUDA, *g.CUDNN, g.OS))
	g.BuiltinSystemPackages = append(g.BuiltinSystemPackages, []string{
		g.Language,
		fmt.Sprintf("%s-pip", g.Language),
	}...)
	installed := g.compileBuiltinSystemPackages(root)
	return installed
}

// Deprecated: Use compileCUDAPackages instead.
func (g *Graph) compileCUDAPackagesDeprecated() llb.State {
	root := llb.Image(fmt.Sprintf("nvidia/cuda:%s.0-base-%s", *g.CUDA, g.OS))
	env := root.AddEnv("DEBIAN_FRONTEND", "noninteractive").
		AddEnv("LD_LIBRARY_PATH", "/usr/local/cuda-11.0/targets/x86_64-linux/lib:/usr/local/cuda/extras/CUPTI/lib64:/usr/local/cuda/lib64:$LD_LIBRARY_PATH").
		AddEnv("LANG", "C.UTF-8")

	g.BuiltinSystemPackages = append(g.BuiltinSystemPackages, []string{
		"build-essential",
		"curl",
		"libfreetype6-dev",
		"libhdf5-serial-dev",
		"libzmq3-dev",
		"pkg-config",
		"software-properties-common",
		"unzip",
		// Python
		g.Language,
		fmt.Sprintf("%s-pip", g.Language),
		// CUDA and CUDNN
		fmt.Sprintf("cuda-command-line-tools-%s", *g.CUDA),
		fmt.Sprintf("libcublas-%s", *g.CUDA),
		fmt.Sprintf("cuda-nvrtc-%s", *g.CUDA),
		fmt.Sprintf("libcufft-%s", *g.CUDA),
		fmt.Sprintf("libcurand-%s", *g.CUDA),
		fmt.Sprintf("libcusolver-%s", *g.CUDA),
		fmt.Sprintf("libcusparse-%s", *g.CUDA),
		fmt.Sprintf("libcudnn8=%s+cuda%s", *g.CUDNN, *g.CUDA),
	}...)

	installed := g.compileBuiltinSystemPackages(env)

	run := installed.Run(llb.Shlex(`ln -s /usr/local/cuda/lib64/stubs/libcuda.so 
	/usr/local/cuda/lib64/stubs/libcuda.so.1
	&& echo "/usr/local/cuda/lib64/stubs" > /etc/ld.so.conf.d/z-cuda-stubs.conf
	&& ldconfig`))
	return run.Root()
}

func (g Graph) compilePyPIPackages(root llb.State) llb.State {
	if len(g.PyPIPackages) == 0 {
		return root
	}
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

func (g Graph) compileBuiltinSystemPackages(root llb.State) llb.State {
	if len(g.BuiltinSystemPackages) == 0 {
		return root
	}
	cmd := "sh -c \"apt-get update && apt-get install -y --no-install-recommends"
	cacheDir := "/var/cache/apt"
	cacheLibDir := "/var/lib/apt"
	for _, pkg := range g.BuiltinSystemPackages {
		cmd = cmd + " " + pkg
	}
	cmd += "\""
	run := root.Run(llb.Shlex(cmd))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	run.AddMount(cacheLibDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheLibDir, llb.CacheMountShared))
	return run.Root()
}

func (g Graph) compileSystemPackages(root llb.State) llb.State {
	if len(g.SystemPackages) == 0 {
		return root
	}
	// TODO(gaocegege): Support per-user config to keep the mirror.
	cmd := "apt install"
	cacheDir := "/var/cache/apt"
	cacheLibDir := "/var/lib/apt"
	for _, pkg := range g.SystemPackages {
		cmd = cmd + " " + pkg
	}
	run := root.Run(llb.Shlex(cmd))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	run.AddMount(cacheLibDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheLibDir, llb.CacheMountShared))
	return run.Root()
}
