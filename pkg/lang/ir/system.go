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
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/version"
)

func (g Graph) compileUbuntuAPT(root llb.State) llb.State {
	if g.UbuntuAPTSource != nil {
		logrus.WithField("source", *g.UbuntuAPTSource).Debug("using custom APT source")
		aptSource := llb.Scratch().
			File(llb.Mkdir(filepath.Dir(aptSourceFilePath), 0755, llb.WithParents(true)),
				llb.WithCustomName("[internal] setting apt source")).
			File(llb.Mkfile(aptSourceFilePath, 0644, []byte(*g.UbuntuAPTSource)),
				llb.WithCustomName("[internal] setting apt source"))
		return llb.Merge([]llb.State{root, aptSource},
			llb.WithCustomName("[internal] setting apt source"))
	}
	return root
}

func (g Graph) compileRun(root llb.State) llb.State {
	if len(g.Exec) == 0 {
		return root
	}

	workingDir := g.getWorkingDir()
	stage := root.AddEnv("PATH", types.DefaultPathEnvUnix)
	for _, execGroup := range g.Exec {
		var sb strings.Builder
		sb.WriteString("set -euo pipefail\n")
		for _, c := range execGroup.Commands {
			sb.WriteString(c + "\n")
		}

		cmdStr := fmt.Sprintf("/usr/bin/bash -c '%s'", sb.String())
		logrus.WithField("command", cmdStr).Debug("compile run command")
		// Mount the build context into the build process.
		// TODO(gaocegege): Maybe we should make it readonly,
		// but these cases then cannot be supported:
		// run(commands=["git clone xx.git"])
		run := stage.Dir(workingDir).Run(llb.Shlex(cmdStr))
		if execGroup.MountHost {
			run.AddMount(workingDir, llb.Local(flag.FlagBuildContext))
		}
		stage = run.Root()
	}
	return stage
}

func (g Graph) compileCopy(root llb.State) llb.State {
	if len(g.Copy) == 0 {
		return root
	}

	result := root
	// Compose the copy command.
	for _, c := range g.Copy {
		result = result.File(llb.Copy(
			llb.Local(flag.FlagBuildContext), c.Source, c.Destination,
			llb.WithUIDGID(g.uid, g.gid)))
	}
	return result
}

func (g *Graph) compileCUDAPackages(org string) llb.State {
	return g.preparePythonBase(llb.Image(fmt.Sprintf(
		"docker.io/%s:%s-cudnn%s-devel-%s",
		org, *g.CUDA, g.CUDNN, g.OS)))
}

func (g Graph) compileSystemPackages(root llb.State) llb.State {
	if len(g.SystemPackages) == 0 {
		logrus.Debug("skip the apt since system package is not specified")
		return root
	}

	// Compose the package install command.
	var sb strings.Builder
	sb.WriteString("sudo apt-get update && sudo apt-get install -y --no-install-recommends")

	for _, pkg := range g.SystemPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}

	cacheDir := "/var/cache/apt"
	cacheLibDir := "/var/lib/apt"

	run := root.Run(llb.Shlex(fmt.Sprintf("bash -c \"%s\"", sb.String())),
		llb.WithCustomNamef("apt-get install %s",
			strings.Join(g.SystemPackages, " ")))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir(g.CacheID(cacheDir), llb.CacheMountShared))
	run.AddMount(cacheLibDir, llb.Scratch(),
		llb.AsPersistentCacheDir(g.CacheID(cacheLibDir), llb.CacheMountShared))
	return run.Root()
}

// nolint:unparam
func (g *Graph) compileExtraSource(root llb.State) (llb.State, error) {
	if len(g.HTTP) == 0 {
		return root, nil
	}
	inputs := []llb.State{}
	for _, httpInfo := range g.HTTP {
		src := llb.HTTP(
			httpInfo.URL,
			llb.Checksum(httpInfo.Checksum),
			llb.Filename(httpInfo.Filename),
			llb.Chown(g.uid, g.gid),
		)
		inputs = append(inputs, llb.Scratch().File(
			llb.Copy(src, "/", g.getExtraSourceDir(), &llb.CopyInfo{CreateDestPath: true}),
		))
	}
	inputs = append(inputs, root)
	return llb.Merge(inputs, llb.WithCustomName("[internal] build source layers")), nil
}

func (g *Graph) preparePythonBase(root llb.State) llb.State {
	for _, env := range types.BaseEnvironment {
		root = root.AddEnv(env.Name, env.Value)
	}

	// apt packages
	var sb strings.Builder
	sb.WriteString("apt-get update && apt-get install -y apt-utils && ")
	sb.WriteString("apt-get install -y --no-install-recommends --no-install-suggests --fix-missing ")
	sb.WriteString(strings.Join(types.BaseAptPackage, " "))
	sb.WriteString("&& rm -rf /var/lib/apt/lists/* ")
	// shell prompt
	sb.WriteString("&& curl --proto '=https' --tlsv1.2 -sSf https://starship.rs/install.sh | sh -s -- -y")

	run := root.Run(llb.Shlex(fmt.Sprintf("bash -c \"%s\"", sb.String())),
		llb.WithCustomName("[internal] install system packages"))

	return run.Root()
}

func (g Graph) compileSSHD(root llb.State) llb.State {
	sshd := root.File(llb.Copy(
		llb.Image(types.EnvdSshdImage), "/usr/bin/envd-sshd", "/var/envd/bin/envd-sshd",
		&llb.CopyInfo{CreateDestPath: true}),
		llb.WithCustomName(fmt.Sprintf("[internal] add envd-sshd from %s", types.EnvdSshdImage)))
	return sshd
}

func (g *Graph) compileBase() (llb.State, error) {
	logger := logrus.WithFields(logrus.Fields{
		"os":       g.OS,
		"language": g.Language.Name,
	})
	if g.Language.Version != nil {
		logger = logger.WithField("version", *g.Language.Version)
	}
	logger.Debug("compile base image")

	var base llb.State
	org := viper.GetString(flag.FlagDockerOrganization)
	v := version.GetVersionForImageTag()
	// Do not update user permission in the base image.
	if g.Image != nil {
		logger.WithField("image", *g.Image).Debugf("using custom base image")
		return llb.Image(*g.Image), nil
	} else if g.CUDA == nil {
		switch g.Language.Name {
		case "r":
			base = llb.Image(fmt.Sprintf("docker.io/%s/r-base:4.2-envd-%s", org, v))
			// r-base image already has GID 1000.
			// It is a trick, we actually use GID 1000
			if g.gid == 1000 {
				g.gid = 1001
			}
			if g.uid == 1000 {
				g.uid = 1001
			}
		case "python":
			// TODO(keming) use user input `base(os="")`
			base = g.preparePythonBase(llb.Image(types.PythonBaseImage))
		case "julia":
			base = llb.Image(fmt.Sprintf(
				"docker.io/%s/julia:1.8rc1-ubuntu20.04-envd-%s", org, v))
		}
	} else {
		base = g.compileCUDAPackages("nvidia/cuda")
	}

	// Install conda first.
	condaStage, err := g.installConda(base)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to install conda")
	}
	supervisor := g.installHorust(condaStage)
	sshdStage := g.compileSSHD(supervisor)
	source, err := g.compileExtraSource(sshdStage)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get extra sources")
	}
	aptStage := g.compileUbuntuAPT(source)
	final := g.compileUserGroup(aptStage)
	return final, nil
}

func (g Graph) copySSHKey(root llb.State) (llb.State, error) {
	public := DefaultGraph.PublicKeyPath
	bdat, err := os.ReadFile(public)
	dat := strings.TrimSuffix(string(bdat), "\n")
	if err != nil {
		return llb.State{}, errors.Wrap(err, "Cannot read public SSH key")
	}
	run := root.
		File(llb.Mkdir("/var/envd", 0755, llb.WithParents(true),
			llb.WithUIDGID(g.uid, g.gid)),
			llb.WithCustomName("[internal] create dir for ssh key")).
		File(llb.Mkfile(config.ContainerAuthorizedKeysPath,
			0644, []byte(dat+" envd"), llb.WithUIDGID(g.uid, g.gid)),
			llb.WithCustomName("[internal] install ssh keys"))
	return run, nil
}
