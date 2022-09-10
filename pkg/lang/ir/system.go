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
	root = root.AddEnv("PATH", types.DefaultPathEnvUnix)
	logrus.Debugf("compile run: %s", strings.Join(g.Exec, " "))
	if len(g.Exec) == 1 {
		return root.Run(llb.Shlex(fmt.Sprintf("bash -c \"%s\"", g.Exec[0]))).Root()
	}

	var sb strings.Builder
	sb.WriteString("set -euo pipefail\n")
	for _, c := range g.Exec {
		sb.WriteString(c + "\n")
	}

	cmdStr := fmt.Sprintf("bash -c '%s'", sb.String())
	logrus.WithField("command", cmdStr).Debug("compile run command")
	workingDir := g.getWorkingDir()
	run := root.Dir(workingDir).
		Run(llb.Shlex(cmdStr))
	// Mount the build context into the build process.
	// TODO(gaocegege): Maybe we should make it readonly,
	// but these cases then cannot be supported:
	// run(commands=["git clone xx.git"])
	run.AddMount(workingDir, llb.Local(flag.FlagBuildContext))

	return run.Root()
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

func (g *Graph) compileCUDAPackages(org, version string) llb.State {
	return llb.Image(fmt.Sprintf(
		"docker.io/%s/python:3.9-%s-cuda%s-cudnn%s-envd-%s",
		org, g.OS, *g.CUDA, *g.CUDNN, version))
}

func (g Graph) compileSystemPackages(root llb.State) llb.State {
	if len(g.SystemPackages) == 0 {
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
	} else if g.CUDA == nil && g.CUDNN == nil {
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
			base = llb.Image(fmt.Sprintf(
				"docker.io/%s/python:3.9-ubuntu20.04-envd-%s", org, v))
		case "julia":
			base = llb.Image(fmt.Sprintf(
				"docker.io/%s/julia:1.8rc1-ubuntu20.04-envd-%s", org, v))
		}
	} else {
		base = g.compileCUDAPackages(org, v)
	}
	var res llb.ExecState

	// Install conda first.
	condaStage, err := g.installConda(base)
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to install conda")
	}

	// TODO(gaocegege): Refactor user to a separate stage.
	if g.uid == 0 {
		res = condaStage.
			Run(llb.Shlex(fmt.Sprintf("groupadd -g %d envd", 1001)),
				llb.WithCustomName("[internal] still create group envd for root context")).
			Run(llb.Shlex(fmt.Sprintf("useradd -p \"\" -u %d -g envd -s /bin/sh -m envd", 1001)),
				llb.WithCustomName("[internal] still create user envd for root context")).
			Run(llb.Shlex("usermod -s /bin/sh root"),
				llb.WithCustomName("[internal] set root default shell to /bin/sh")).
			Run(llb.Shlex("sed -i \"s/envd:x:1001:1001/envd:x:0:0/g\" /etc/passwd"),
				llb.WithCustomName("[internal] set envd uid to 0 as root")).
			Run(llb.Shlex("sed -i \"s./root./home/envd.g\" /etc/passwd"),
				llb.WithCustomName("[internal] set root home dir to /home/envd")).
			Run(llb.Shlex("sed -i \"s/envd:x:1001/envd:x:0/g\" /etc/group"),
				llb.WithCustomName("[internal] set envd group to 0 as root group")).
			Run(llb.Shlex("chown -R root:root /opt/conda"),
				llb.WithCustomName("[internal] configure user permissions"))
	} else {
		res = condaStage.
			Run(llb.Shlex(fmt.Sprintf("groupadd -g %d envd", g.gid)),
				llb.WithCustomName("[internal] create user group envd")).
			Run(llb.Shlex(fmt.Sprintf("useradd -p \"\" -u %d -g envd -s /bin/sh -m envd", g.uid)),
				llb.WithCustomName("[internal] create user envd")).
			Run(llb.Shlex("adduser envd sudo"),
				llb.WithCustomName("[internal] add user envd to sudoers")).
			Run(llb.Shlex("chown -R envd:envd /usr/local/lib"),
				llb.WithCustomName("[internal] configure user permissions")).
			Run(llb.Shlex("chown -R envd:envd /opt/conda"),
				llb.WithCustomName("[internal] configure user permissions"))
	}
	return llb.User("envd")(res.Root()), nil
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
			llb.WithUIDGID(g.uid, g.gid))).
		File(llb.Mkfile(config.ContainerAuthorizedKeysPath,
			0644, []byte(dat+" envd"), llb.WithUIDGID(g.uid, g.gid)),
			llb.WithCustomName("install ssh keys"))
	return run, nil
}
