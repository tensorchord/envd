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
	"github.com/sirupsen/logrus"
)

const (
	condarc = "/home/envd/.condarc"
)

func (g Graph) CondaEnabled() bool {
	return g.CondaConfig != nil
}

func (g Graph) compileCondaChannel(root llb.State) llb.State {
	if g.CondaConfig != nil && g.CondaConfig.CondaChannel != nil {
		logrus.WithField("conda-channel", *g.CondaChannel).Debug("using custom connda channel")
		stage := root.
			File(llb.Mkfile(condarc,
				0644, []byte(*g.CondaChannel), llb.WithUIDGID(g.uid, g.gid)), llb.WithCustomName("[internal] settings conda channel"))
		return stage
	}
	return root
}

func (g Graph) compileCondaPackages(root llb.State) llb.State {
	if !g.CondaEnabled() || len(g.CondaConfig.CondaPackages) == 0 {
		return root
	}

	cacheDir := "/opt/conda/pkgs"

	// Compose the package install command.
	var sb strings.Builder
	if len(g.CondaConfig.AdditionalChannels) == 0 {
		sb.WriteString("/opt/conda/bin/conda install -n envd")

	} else {
		sb.WriteString("/opt/conda/bin/conda install -n envd -c")
		for _, channel := range g.CondaConfig.AdditionalChannels {
			sb.WriteString(fmt.Sprintf(" %s", channel))
		}
	}

	for _, pkg := range g.CondaConfig.CondaPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}

	cmd := sb.String()
	root = llb.User("envd")(root)
	// Refer to https://github.com/moby/buildkit/blob/31054718bf775bf32d1376fe1f3611985f837584/frontend/dockerfile/dockerfile2llb/convert_runmount.go#L46
	cache := root.File(llb.Mkdir("/cache",
		0755, llb.WithParents(true), llb.WithUIDGID(g.uid, g.gid)),
		llb.WithCustomName("[internal] settings conda cache mount permissions"))
	run := root.
		Run(llb.Shlex(cmd), llb.WithCustomNamef("conda install %s",
			strings.Join(g.CondaPackages, " ")))
	run.AddMount(cacheDir, cache,
		llb.AsPersistentCacheDir(g.CacheID(cacheDir), llb.CacheMountShared), llb.SourcePath("/cache"))
	return run.Root()
}

func (g Graph) setCondaENV(root llb.State) llb.State {
	if !g.CondaEnabled() {
		return root
	}

	root = llb.User("envd")(root)
	// Always init bash since we will use it to create jupyter notebook service.
	run := root.Run(llb.Shlex("bash -c \"/opt/conda/bin/conda init bash\""), llb.WithCustomName("[internal] initialize conda bash environment"))

	pythonVersion := "3.9"
	if g.Language.Version != nil {
		pythonVersion = *g.Language.Version
	}
	cmd := fmt.Sprintf(
		"bash -c \"/opt/conda/bin/conda create -n envd python=%s\"", pythonVersion)
	// Create a conda environment.
	run = run.Run(llb.Shlex(cmd), llb.WithCustomName("[internal] create conda environment"))

	switch g.Shell {
	case shellBASH:
		run = run.Run(
			llb.Shlex(`bash -c 'echo "source /opt/conda/bin/activate envd" >> /home/envd/.bashrc'`),
			llb.WithCustomName("[internal] add conda environment to bashrc"))
	case shellZSH:
		run = run.Run(
			llb.Shlex(fmt.Sprintf("bash -c \"/opt/conda/bin/conda init %s\"", g.Shell)),
			llb.WithCustomNamef("[internal] initialize conda %s environment", g.Shell)).Run(
			llb.Shlex(`bash -c 'echo "source /opt/conda/bin/activate envd" >> /home/envd/.zshrc'`),
			llb.WithCustomName("[internal] add conda environment to zshrc"))
	}
	return run.Root()
}
