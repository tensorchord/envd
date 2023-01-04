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
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/types"
)

const (
	condaVersionDefault = "py39_4.11.0"
	// check the issue https://github.com/mamba-org/mamba/issues/1975
	mambaVersionDefault = "0.25.1"
	condaRootPrefix     = "/opt/conda"
	condaBinDir         = "/opt/conda/bin"
)

var (
	// this file can be used by both conda and mamba
	// https://mamba.readthedocs.io/en/latest/user_guide/configuration.html#multiple-rc-files
	condarc = "/opt/conda/.condarc"
	//go:embed install-conda.sh
	installCondaBash string
	//go:embed install-mamba.sh
	installMambaBash string
)

func (g generalGraph) compileCondaChannel(root llb.State) llb.State {
	if g.CondaConfig.CondaChannel != nil {
		logrus.WithField("conda-channel", *g.CondaChannel).Debug("using custom conda channel")
		stage := root.
			File(llb.Mkfile(condarc,
				0644, []byte(*g.CondaChannel), llb.WithUIDGID(g.uid, g.gid)), llb.WithCustomName("[internal] setting conda channel"))
		return stage
	}
	return root
}

func (g generalGraph) condaCommandPath() string {
	if g.CondaConfig.UseMicroMamba {
		return filepath.Join(condaBinDir, "micromamba")
	}
	return filepath.Join(condaBinDir, "conda")
}

func (g generalGraph) condaInitShell(shell string) string {
	path := g.condaCommandPath()
	if g.CondaConfig.UseMicroMamba {
		return fmt.Sprintf("%s shell init -p %s -s %s", path, condaRootPrefix, shell)
	}
	return fmt.Sprintf("%s init %s", path, shell)
}

func (g generalGraph) condaUpdateFromFile() string {
	args := fmt.Sprintf("update -n envd --file %s", g.CondaEnvFileName)
	if g.CondaConfig.UseMicroMamba {
		return fmt.Sprintf("%s %s", g.condaCommandPath(), args)
	}
	return fmt.Sprintf("%s env %s", g.condaCommandPath(), args)
}

func (g *generalGraph) compileCondaPackages(root llb.State) llb.State {
	if len(g.CondaConfig.CondaPackages) == 0 && len(g.CondaEnvFileName) == 0 {
		logrus.Debug("Conda packages not enabled")
		return root
	}

	cacheDir := filepath.Join(condaRootPrefix, "pkgs")
	// Refer to https://github.com/moby/buildkit/blob/31054718bf775bf32d1376fe1f3611985f837584/frontend/dockerfile/dockerfile2llb/convert_runmount.go#L46
	cacheMount := llb.Scratch().File(llb.Mkdir("/cache-conda", 0755, llb.WithParents(true)),
		llb.WithCustomName("[internal] setting conda cache mount permissions"))

	// Compose the package install command.
	var sb strings.Builder
	var run llb.ExecState

	if len(g.CondaEnvFileName) > 0 {
		sb.WriteString(g.condaUpdateFromFile())
	} else {
		if len(g.CondaConfig.AdditionalChannels) == 0 {
			sb.WriteString(fmt.Sprintf("%s install -n envd", g.condaCommandPath()))
		} else {
			sb.WriteString(fmt.Sprintf("%s install -n envd", g.condaCommandPath()))
			for _, channel := range g.CondaConfig.AdditionalChannels {
				sb.WriteString(fmt.Sprintf(" -c %s", channel))
			}
		}
		for _, pkg := range g.CondaConfig.CondaPackages {
			sb.WriteString(fmt.Sprintf(" %s", pkg))
		}
	}

	cmd := sb.String()
	run = root.Dir(g.getWorkingDir()).
		Run(llb.Shlex(cmd), llb.WithCustomNamef("[internal] %s %s",
			cmd, strings.Join(g.CondaPackages, " ")))
	run.AddMount(g.getWorkingDir(), llb.Local(flag.FlagBuildContext))
	run.AddMount(cacheDir, cacheMount,
		llb.AsPersistentCacheDir(g.CacheID(cacheDir), llb.CacheMountShared), llb.SourcePath("/cache-conda"))
	return run.Root()
}

func (g generalGraph) compileCondaEnvironment(root llb.State) (llb.State, error) {
	// Always init bash since we will use it to create jupyter notebook service.
	run := root.Run(
		llb.Shlexf(`bash -c "%s"`, g.condaInitShell("bash")),
		llb.WithCustomName("[internal] initialize conda bash environment"),
	)

	pythonVersion, err := g.getAppropriatePythonVersion()
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get python version")
	}
	// Create a conda environment.
	cmd := fmt.Sprintf("bash -c \"%s create -n envd python=%s\"", g.condaCommandPath(), pythonVersion)
	run = run.Run(llb.Shlex(cmd),
		llb.WithCustomNamef("[internal] create conda environment: %s", cmd))

	return run.Root(), nil
}

func (g *generalGraph) installConda(root llb.State) llb.State {
	root = g.updateEnvPath(root, types.DefaultCondaPath)
	// this directory is related to conda envd env meta (used by `conda env config vars set key=value`)
	g.UserDirectories = append(g.UserDirectories, fmt.Sprintf("%s/envs/envd/conda-meta", condaRootPrefix))
	if g.CondaConfig.UseMicroMamba {
		run := root.AddEnv("MAMBA_BIN_DIR", condaBinDir).
			AddEnv("MAMBA_ROOT_PREFIX", condaRootPrefix).
			AddEnv("MAMBA_VERSION", mambaVersionDefault).
			Run(llb.Shlexf("bash -c '%s'", installMambaBash),
				llb.WithCustomName("[internal] install micro mamba"))
		return run.Root()
	}
	run := root.AddEnv("CONDA_VERSION", condaVersionDefault).
		File(llb.Mkdir(condaRootPrefix, 0755, llb.WithParents(true)),
			llb.WithCustomName("[internal] create conda directory")).
		Run(llb.Shlexf("bash -c '%s'", installCondaBash),
			llb.WithCustomName("[internal] install conda"))
	return run.Root()
}
