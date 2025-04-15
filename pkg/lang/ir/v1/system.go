// Copyright 2023 The envd Authors
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

package v1

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func (g generalGraph) compileUbuntuAPT(root llb.State) llb.State {
	if g.UbuntuAPTSource != nil {
		logrus.WithField("source", *g.UbuntuAPTSource).Debug("using custom APT source")
		aptSource := root.
			File(llb.Mkdir(filepath.Dir(aptSourceFilePath), 0755, llb.WithParents(true)),
				llb.WithCustomName("[internal] setting apt source")).
			File(llb.Mkfile(aptSourceFilePath, 0644, []byte(*g.UbuntuAPTSource)),
				llb.WithCustomName("[internal] setting apt source"))
		return aptSource
	}
	return root
}

func (g generalGraph) compileRun(root llb.State) llb.State {
	if len(g.Exec) == 0 {
		return root
	}

	workingDir := g.getWorkingDir()
	for _, execGroup := range g.Exec {
		var sb strings.Builder
		sb.WriteString("set -euo pipefail\n")
		for _, c := range execGroup.Commands {
			sb.WriteString(c + "\n")
		}

		cmdStr := fmt.Sprintf("bash -c '%s'", sb.String())
		logrus.WithField("command", cmdStr).Debug("compile run command")
		// mount host here is read-only
		run := root.Dir(workingDir).Run(llb.Shlex(cmdStr))
		if execGroup.MountHost {
			run.AddMount(workingDir, llb.Local(flag.FlagBuildContext))
		}
		root = run.Root()
	}
	return root
}

func (g generalGraph) compileCopy(root llb.State) llb.State {
	if len(g.Copy) == 0 {
		return root
	}

	result := root
	// Compose the copy command.
	for _, c := range g.Copy {
		var from llb.State
		if c.Image == "" {
			from = llb.Local(flag.FlagBuildContext)
		} else {
			from = llb.Image(c.Image)
		}
		result = result.File(llb.Copy(
			from, c.Source, c.Destination,
			&llb.CopyInfo{CreateDestPath: true},
			llb.WithUIDGID(g.uid, g.gid)))
	}
	return result
}

func (g generalGraph) compileSystemPackages(root llb.State) llb.State {
	if len(g.SystemPackages) == 0 {
		logrus.Debug("skip the apt since system package is not specified")
		return root
	}

	// Compose the package install command.
	var sb strings.Builder
	sb.WriteString("apt-get update && apt-get install -y --no-install-recommends")

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

// nolint:unparam
func (g *generalGraph) compileExtraSource(root llb.State) llb.State {
	if len(g.HTTP) == 0 {
		return root
	}
	if g.DisableMergeOp {
		for _, httpInfo := range g.HTTP {
			src := llb.HTTP(
				httpInfo.URL,
				llb.Checksum(httpInfo.Checksum),
				llb.Filename(httpInfo.Filename),
				llb.Chown(g.uid, g.gid),
			)
			root = root.File(llb.Copy(
				src, "/", g.getExtraSourceDir(), &llb.CopyInfo{CreateDestPath: true},
			))
		}
		return root
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
	return llb.Merge(inputs, llb.WithCustomName("[internal] build source layers"))
}

func (g *generalGraph) compileLanguage(root llb.State) (llb.State, error) {
	langs := []llb.State{}
	lang := root
	var err error

	if g.DisableMergeOp {
		for _, language := range g.Languages {
			switch language.Name {
			case "python":
				root, err = g.installPython(root)
			case "r":
				root = g.installRLang(root)
			case "julia":
				root = g.installJulia(root)
			}
		}
		return root, err
	}

	for _, language := range g.Languages {
		switch language.Name {
		case "python":
			lang, err = g.installPython(root)
		case "r":
			lang = g.installRLang(root)
		case "julia":
			lang = g.installJulia(root)
		}
		langs = append(langs, lang)
	}
	if len(langs) <= 1 {
		return lang, err
	}
	for i, lang := range g.Languages {
		langs[i] = llb.Diff(root, langs[i], llb.WithCustomNamef("[internal] build %s env", lang.Name))
	}
	return llb.Merge(append([]llb.State{root}, langs...),
		llb.WithCustomName("[internal] merge all the language environments")), err
}

func (g *generalGraph) compileLanguagePackages(root llb.State) llb.State {
	// Use default python in the base image if install.python() is not specified.
	g.compileJupyter()
	index := g.compilePyPIIndex(root)
	pack := g.compilePyPIPackages(index)
	if g.CondaConfig != nil {
		channel := g.compileCondaChannel(pack)
		pack = g.compileCondaPackages(channel)
	}
	if g.UVConfig != nil {
		pack = g.compileUV(pack)
	}
	if g.PixiConfig != nil {
		pack = g.compilePixi(pack)
	}

	for _, language := range g.Languages {
		switch language.Name {
		case "r":
			pack = g.installRPackages(pack)
		case "julia":
			pack = g.installJuliaPackages(pack)
		}
	}
	return pack
}

func (g *generalGraph) compileDevPackages(root llb.State) llb.State {
	// apt packages
	var sb strings.Builder
	sb.WriteString("apt-get update && apt-get install -y apt-utils && ")
	sb.WriteString("apt-get install -y --no-install-recommends --no-install-suggests --fix-missing ")
	sb.WriteString(strings.Join(types.BaseAptPackage, " "))
	sb.WriteString("&& rm -rf /var/lib/apt/lists/*")

	run := root.Run(llb.Shlexf(`bash -c "%s"`, sb.String()),
		llb.WithCustomName("[internal] install built-in packages"))

	return run.Root()
}

func (g generalGraph) compileStarship(root llb.State) llb.State {
	starship := root.File(llb.Copy(
		llb.Image(types.EnvdStarshipImage), "/usr/local/bin/starship", "/usr/local/bin/starship",
		&llb.CopyInfo{CreateDestPath: true}),
		llb.WithCustomName(fmt.Sprintf("[internal] add envd-starship from %s", types.EnvdStarshipImage)))
	return starship
}

func (g generalGraph) compileSSHD(root llb.State) llb.State {
	sshd := root.File(llb.Copy(
		llb.Image(types.EnvdSshdImage), "/usr/bin/envd-sshd", "/var/envd/bin/envd-sshd",
		&llb.CopyInfo{CreateDestPath: true}),
		llb.WithCustomName(fmt.Sprintf("[internal] add envd-sshd from %s", types.EnvdSshdImage)))
	return sshd
}

func (g *generalGraph) compileBaseImage() (llb.State, error) {
	// TODO: find another way to install CUDA
	if g.CUDA != nil {
		g.Image = GetCUDAImage(g.Image, g.CUDA, g.CUDNN, g.Dev)
	}

	logger := logrus.WithFields(logrus.Fields{
		"image":    g.Image,
		"language": g.Languages,
	})
	logger.Debug("compile base image")

	// Fix https://github.com/tensorchord/envd/issues/1147.
	// Fetch the image metadata from base image.
	base := llb.Image(g.Image)
	// fetching the image config may take some time
	config, err := ir.FetchImageConfig(context.Background(), g.Image, g.Platform)
	if err != nil {
		return llb.State{}, errors.Wrapf(err, "failed to get the image config, check if the image(%s) exists", g.Image)
	}

	// Set the environment variables to RuntimeEnviron to keep it in the resulting image.
	logger.Logger.Debugf("inherit envs from base image: %s", config.Env)
	for _, e := range config.Env {
		// in case the env value also contains `=`
		kv := strings.SplitN(e, "=", 2)
		g.RuntimeEnviron[kv[0]] = kv[1]
		if kv[0] == "PATH" {
			// deduplicate the PATH but keep the order as:
			// 0. default Unix PATH
			// 1. configured paths in the Starlark frontend `runtime.environ(extra_path=[...])`
			// 2. paths in the base image
			// 3. others added during the image building (Python paths, etc.)

			// iterate over the original paths and add them to the map
			pathMap := make(map[string]bool)
			for _, path := range g.RuntimeEnvPaths {
				pathMap[path] = true
			}
			// split the PATH into different paths
			newPaths := strings.Split(kv[1], ":")
			// iterate over the new paths
			for _, path := range newPaths {
				// check if the path is already in the map
				if _, ok := pathMap[path]; !ok {
					// if not, add the path to the map and slice
					pathMap[path] = true
					g.RuntimeEnvPaths = append(g.RuntimeEnvPaths, path)
				}
			}
		}
	}

	// add necessary envs
	for _, env := range types.BaseEnvironment {
		base = base.AddEnv(env.Name, env.Value)
	}
	for k, v := range g.RuntimeEnviron {
		base = base.AddEnv(k, v)
	}

	if !g.Dev {
		if len(g.Entrypoint) == 0 {
			g.Entrypoint = config.Entrypoint
		}
		g.User = config.User
		g.WorkingDir = config.WorkingDir
	} else {
		// for dev mode, we will create an `envd` user
		g.User = ""
		g.WorkingDir = g.getWorkingDir()
	}
	return base, nil
}

func (g generalGraph) copySSHKey(root llb.State) (llb.State, error) {
	public := g.PublicKeyPath
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

func (g generalGraph) compileMountDir(root llb.State) llb.State {
	mount := root
	if g.Dev {
		// create the ENVD_WORKDIR as a placeholder (envd-server may not mount this dir)
		workDir := fileutil.EnvdHomeDir(g.EnvironmentName)
		mount = root.File(llb.Mkdir(workDir, 0755, llb.WithParents(true), llb.WithUIDGID(g.uid, g.gid)),
			llb.WithCustomNamef("[internal] create work dir: %s", workDir))
	}

	for _, m := range g.Mount {
		mount = mount.File(llb.Mkdir(m.Destination, 0755, llb.WithParents(true),
			llb.WithUIDGID(g.uid, g.gid)),
			llb.WithCustomNamef("[internal] create dir for runtime.mount %s", m.Destination),
		)
	}
	return mount
}

func (g *generalGraph) updateEnvPath(root llb.State, path string) llb.State {
	g.RuntimeEnvPaths = append(g.RuntimeEnvPaths, path)
	return root.AddEnv("PATH", strings.Join(g.RuntimeEnvPaths, ":"))
}
