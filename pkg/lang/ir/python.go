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

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/flag"
)

const (
	pythonVersionDefault = "3.9"
)

func (g Graph) getAppropriatePythonVersion() (string, error) {
	if g.Language.Version == nil {
		return pythonVersionDefault, nil
	}

	version := *g.Language.Version
	if version == "3" || version == "" {
		return pythonVersionDefault, nil
	}
	if strings.HasPrefix(version, "3.") {
		return version, nil
	}
	return "", errors.Errorf("python version %s is not supported", version)
}

func (g Graph) compilePython(baseStage llb.State) (llb.State, error) {
	if err := g.compileJupyter(); err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile jupyter")
	}
	aptStage := g.compileUbuntuAPT(baseStage)
	systemStage := g.compileSystemPackages(aptStage)

	condaEnvStage, err := g.compileCondaEnvironment(baseStage)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile conda environment")
	}

	// Conda affects shell and python, thus we cannot do it in parallel.
	shellStage, err := g.compileShell(baseStage)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to compile shell")
	}
	condaShellStage := g.compileCondaShell(shellStage)

	diffCondaEnvStage := llb.Diff(baseStage, condaEnvStage,
		llb.WithCustomName("[internal] conda python environment"))
	diffSystemStage := llb.Diff(baseStage, systemStage,
		llb.WithCustomName("[internal] install system packages"))
	diffShellStage := llb.Diff(baseStage, condaShellStage,
		llb.WithCustomNamef("[internal] configure shell %s", g.Shell))
	prePythonStage := llb.Merge([]llb.State{
		diffSystemStage,
		diffCondaEnvStage,
		diffShellStage,
		baseStage}, llb.WithCustomName("pre-python stage"))

	condaChannelStage := g.compileCondaChannel(prePythonStage)

	condaStage := llb.Diff(prePythonStage,
		g.compileCondaPackages(condaChannelStage),
		llb.WithCustomName("[internal] install conda packages"))

	pypiMirrorStage := g.compilePyPIIndex(prePythonStage)

	pypiStage := llb.Diff(prePythonStage,
		g.compilePyPIPackages(pypiMirrorStage),
		llb.WithCustomName("[internal] install PyPI packages"))

	vscodeStage, err := g.compileVSCode()
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get vscode plugins")
	}
	sshStage, err := g.copySSHKey(prePythonStage)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to copy ssh keys")
	}
	diffSSHStage := llb.Diff(prePythonStage, sshStage,
		llb.WithCustomName("[internal] install ssh key"))

	var merged llb.State
	if vscodeStage != nil {
		merged = llb.Merge([]llb.State{
			prePythonStage, condaStage, pypiStage,
			diffSSHStage, *vscodeStage,
		}, llb.WithCustomName("[internal] generating the image"))
	} else {
		merged = llb.Merge([]llb.State{
			prePythonStage, condaStage,
			diffSSHStage, pypiStage,
		}, llb.WithCustomName("[internal] generating the image"))
	}
	merged = g.compileAlternative(merged)
	return merged, nil
}

// Set the system default python to envd's python.
func (g Graph) compileAlternative(root llb.State) llb.State {
	envdPrefix := "/opt/conda/envs/envd/bin"
	run := root.
		Run(llb.Shlexf("update-alternatives --install /usr/bin/python python %s/python 1", envdPrefix),
			llb.WithCustomName("[internal] update alternative python to envd")).
		Run(llb.Shlexf("update-alternatives --install /usr/bin/python3 python3 %s/python3 1", envdPrefix),
			llb.WithCustomName("[internal] update alternative python3 to envd")).
		Run(llb.Shlexf("update-alternatives --install /usr/bin/pip pip %s/pip 1", envdPrefix),
			llb.WithCustomName("[internal] update alternative pip to envd")).
		Run(llb.Shlexf("update-alternatives --install /usr/bin/pip3 pip3 %s/pip3 1", envdPrefix),
			llb.WithCustomName("[internal] update alternative pip3 to envd"))
	return run.Root()
}

func (g Graph) compilePyPIPackages(root llb.State) llb.State {
	if len(g.PyPIPackages) == 0 && g.RequirementsFile == nil && len(g.PythonWheels) == 0 {
		return root
	}

	// Create the envd cache directory in the container. see issue #582
	cacheDir := filepath.Join("/", "root", ".cache", "pip")
	root = g.CompileCacheDir(root, cacheDir)

	cache := root.File(llb.Mkdir("/cache/pip", 0755, llb.WithParents(true)),
		llb.WithCustomName("[internal] setting pip cache mount permissions"))

	if len(g.PyPIPackages) != 0 {
		// Compose the package install command.
		var sb strings.Builder
		// Always use the conda's pip.
		sb.WriteString("/opt/conda/envs/envd/bin/python -m pip install")
		for _, pkg := range g.PyPIPackages {
			sb.WriteString(fmt.Sprintf(" %s", pkg))
		}

		cmd := sb.String()
		logrus.WithField("command", cmd).
			Debug("Configure pip install statements")
		run := root.
			Run(llb.Shlex(sb.String()), llb.WithCustomNamef("pip install %s",
				strings.Join(g.PyPIPackages, " ")))
		// Refer to https://github.com/moby/buildkit/blob/31054718bf775bf32d1376fe1f3611985f837584/frontend/dockerfile/dockerfile2llb/convert_runmount.go#L46
		run.AddMount(cacheDir, cache,
			llb.AsPersistentCacheDir(g.CacheID(cacheDir), llb.CacheMountShared), llb.SourcePath("/cache/pip"))
		root = run.Root()
	}

	if g.RequirementsFile != nil {
		// Compose the package install command.
		var sb strings.Builder
		sb.WriteString("bash -c '")
		sb.WriteString("set -euo pipefail\n")
		sb.WriteString(fmt.Sprintf("chown -R envd:envd %s\n", g.getWorkingDir())) // Change mount dir permission
		envdCmd := strings.Builder{}
		envdCmd.WriteString(fmt.Sprintf("cd %s\n", g.getWorkingDir()))
		envdCmd.WriteString(fmt.Sprintf("/opt/conda/envs/envd/bin/python -m pip install -r  %s\n", *g.RequirementsFile))

		// Execute the command to write yaml file and conda env using envd user
		sb.WriteString(fmt.Sprintf("sudo -i -u envd bash << EOF\n%s\nEOF\n", envdCmd.String()))
		sb.WriteString("'")
		cmd := sb.String()

		logrus.WithField("command", cmd).
			Debug("Configure pip install requirements statements")
		root = root.User("root").Dir(g.getWorkingDir())
		run := root.
			Run(llb.Shlex(cmd), llb.WithCustomNamef("pip install %s", *g.RequirementsFile))
		run.AddMount(cacheDir, cache,
			llb.AsPersistentCacheDir(g.CacheID(cacheDir), llb.CacheMountShared), llb.SourcePath("/cache/pip"))
		run.AddMount(g.getWorkingDir(),
			llb.Local(flag.FlagBuildContext))
		root = run.Root()
	}

	if len(g.PythonWheels) > 0 {
		root = root.Dir(g.getWorkingDir())
		cmdTemplate := "/opt/conda/envs/envd/bin/python -m pip install %s"
		for _, wheel := range g.PythonWheels {
			run := root.Run(llb.Shlex(fmt.Sprintf(cmdTemplate, wheel)), llb.WithCustomNamef("pip install %s", wheel))
			run.AddMount(g.getWorkingDir(), llb.Local(flag.FlagBuildContext), llb.Readonly)
			run.AddMount(cacheDir, cache,
				llb.AsPersistentCacheDir(g.CacheID(cacheDir), llb.CacheMountShared), llb.SourcePath("/cache/pip"))
			root = run.Root()
		}
	}
	return root
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
		dir := filepath.Dir(pypiIndexFilePath)
		pypiMirror := root.
			File(llb.Mkdir(dir, 0755, llb.WithParents(true), llb.WithUIDGID(g.uid, g.gid)),
				llb.WithCustomNamef("[internal] setting PyPI index dir %s", dir)).
			File(llb.Mkfile(pypiIndexFilePath,
				0644, []byte(content), llb.WithUIDGID(g.uid, g.gid)),
				llb.WithCustomNamef("[internal] setting PyPI index file %s", pypiIndexFilePath))
		return pypiMirror
	}
	return root
}
