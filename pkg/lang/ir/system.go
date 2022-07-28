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
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/config"
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
	root = root.AddEnv("PATH", "$PATH:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt/conda/bin:/usr/local/julia/bin:/opt/conda/envs/envd/bin")
	if len(g.Exec) == 1 {
		return root.Run(llb.Shlex(g.Exec[0])).Root()
	}

	run := root.Run(llb.Shlex(g.Exec[0]))
	for _, c := range g.Exec[1:] {
		run = run.Run(llb.Shlex(c))
	}
	return run.Root()
}

func (g *Graph) compileCUDAPackages() llb.State {
	root := llb.Image(fmt.Sprintf("docker.io/tensorchord/python:3.8-%s-cuda%s-cudnn%s", g.OS, *g.CUDA, *g.CUDNN))
	return root
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

func (g *Graph) compileBase() llb.State {
	logger := logrus.WithFields(logrus.Fields{
		"os":       g.OS,
		"language": g.Language.Name,
	})
	if g.Language.Version != nil {
		logger = logger.WithField("version", *g.Language.Version)
	}
	logger.Debug("compile base image")

	var base llb.State
	// Do not update user permission in the base image.
	if g.Image != nil {
		logger.WithField("image", *g.Image).Debugf("using custom base image")
		return llb.Image(*g.Image)
	} else if g.CUDA == nil && g.CUDNN == nil {
		switch g.Language.Name {
		case "r":
			base = llb.Image("docker.io/tensorchord/r-base:4.2")
			// r-base image already has GID 1000.
			// It is a trick, we actually use GID 1000
			if g.gid == 1000 {
				g.gid = 1001
			}
			if g.uid == 1000 {
				g.uid = 1001
			}
		case "python":
			base = llb.Image("docker.io/tensorchord/python:3.8-ubuntu20.04")
		case "julia":
			base = llb.Image("docker.io/tensorchord/julia:1.8rc1-ubuntu20.04")
		}
	} else {
		base = g.compileCUDAPackages()
	}
	var res llb.ExecState
	// TODO(gaocegege): Refactor user to a separate stage.
	if g.uid == 0 {
		res = base.
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
		res = base.
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
	return llb.User("envd")(res.Root())
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
