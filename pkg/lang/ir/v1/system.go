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
	"github.com/moby/buildkit/client/llb/imagemetaresolver"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/config"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

// signFolder stores the path to the apt-source signature in string format
// The value of signFolder should always be /etc/apt/keyrings
const signFolder = "/etc/apt/keyrings"

// signURI stores the third-party URI for downloading the signature of the corresponding repo
const signURI = "https://cloud.r-project.org/bin/linux/ubuntu/marutter_pubkey.asc"

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

// copyAPTSignature returns the state and the path in string format
// The returned string represents the location of the apt-source signature
// A successful run of copyAPTSignature should return path of /etc/apt/keyrings/*.asc
func (g generalGraph) copyAPTSignature(root llb.State, name string, url string) (llb.State, string) {

	var fileName = fmt.Sprintf("%s.asc", name)

	// path stores the location of the signature file
	// The value of path should be /etc/apt/keyrings/*.asc
	var path = filepath.Join(signFolder, fileName)

	base := llb.Image(builderImage)
	builder := base.
		Run(llb.Shlexf("sh -c \"curl %s >> %s\"", url, fileName),
			llb.WithCustomName("[internal] downloading apt-source signature in base image")).Root()

	aptSign := root.
		File(llb.Mkdir(signFolder, 0755, llb.WithParents(true)),
			llb.WithCustomName("[internal] setting target apt-source signature folder")).
		File(llb.Copy(builder, fileName, path),
			llb.WithCustomName("[internal] copy signature from builder"))

	return aptSign, path
}

// configRSrc returns the state and the content in DEB822 format for third-party apt-source
// The returned string contains all configuration for adding third-party into apt source list
// A successful run of configRSrc should return the content strictly follow the DEB822 format
func (g generalGraph) configRSrc(root llb.State, aptConfig ir.APTConfig, sign string) (llb.State, string) {

	var enabled = fmt.Sprintf("Enabled: %s\n", aptConfig.Enabled)
	var types = fmt.Sprintf("Types: %s\n", aptConfig.Types)
	var uris = fmt.Sprintf("URIs: %s\n", aptConfig.URIs)
	var suites = fmt.Sprintf("Suites: %s\n", aptConfig.Suites)
	var components = fmt.Sprintf("Components: %s\n", aptConfig.Components)
	var architecture = fmt.Sprintf("Architectures: %s\n", aptConfig.Arch)

	aptSign, signPath := g.copyAPTSignature(root, aptConfig.Name, sign)
	var signature = fmt.Sprintf("Signed-By: %s\n", signPath)

	var content strings.Builder
	content.WriteString(enabled)
	content.WriteString(types)
	content.WriteString(uris)
	content.WriteString(suites)
	content.WriteString(components)
	content.WriteString(signature)
	content.WriteString(architecture)

	return aptSign, content.String()
}

// compileRLang returns the llb.State only after compoiling the environment for installing R language
// A successful run of compileRLang should set up the official R apt repo into /etc/apt/sources.list.d/*.sources
func (g generalGraph) compileRLang(root llb.State) llb.State {

	var aptConfig = ir.APTConfig{
		Name:       "R-base",                                       // Name for the *.sources file in /etc/apt/sources.list.d
		Enabled:    "yes",                                          // Represents the validation of the third-party repo
		Types:      "deb",                                          // Type of the repo, binary or source code
		URIs:       "https://cloud.r-project.org/bin/linux/ubuntu", // URI repo for the package
		Suites:     "focal-cran40/",                                // Branch of the package
		Components: "",                                             // Distribution of the package. E.g. main, non-free
		Signed:     "R-base.asc",                                   // Name for the signature file
		Arch:       "",                                             // Architecture that is supported
	}

	var file = fmt.Sprintf("/etc/apt/sources.list.d/%s.sources", aptConfig.Name)

	aptSource, content := g.configRSrc(root, aptConfig, signURI)

	aptRLang := aptSource.
		File(llb.Mkdir("/etc/apt/sources.list.d/", 0755, llb.WithParents(true)),
			llb.WithCustomName("[internal] setting apt-source folder sources.list.d")).
		File(llb.Mkfile(file, 0644, []byte(content)),
			llb.WithCustomName("[internal] setting apt-source file")).
		File(llb.Mkfile("/etc/apt/apt.conf.d/DEB822.conf", 0644, []byte("APT::Sources::Use-Deb822 true;\n")),
			llb.WithCustomName("[internal] setting apt-conf file to support DEB822 format"))

	return aptRLang

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

		cmdStr := fmt.Sprintf("/usr/bin/bash -c '%s'", sb.String())
		logrus.WithField("command", cmdStr).Debug("compile run command")
		// Mount the build context into the build process.
		// TODO(gaocegege): Maybe we should make it readonly,
		// but these cases then cannot be supported:
		// run(commands=["git clone xx.git"])
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
		result = result.File(llb.Copy(
			llb.Local(flag.FlagBuildContext), c.Source, c.Destination,
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
func (g *generalGraph) compileExtraSource(root llb.State) (llb.State, error) {
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

func (g *generalGraph) compileLanguage(root llb.State) (llb.State, error) {
	lang := root
	var err error
	switch g.Language.Name {
	case "python":
		lang, err = g.installPython(root)
	case "r":
		rSrc := g.compileRLang(root)
		lang = g.installRLang(rSrc)
	case "julia":
		lang = g.installJulia(root)
	}

	return lang, err
}

func (g *generalGraph) compileLanguagePackages(root llb.State) llb.State {
	pack := root
	switch g.Language.Name {
	case "python":
		index := g.compilePyPIIndex(root)
		pypi := g.compilePyPIPackages(index)
		if g.CondaConfig == nil {
			pack = pypi
		} else {
			channel := g.compileCondaChannel(root)
			conda := g.compileCondaPackages(channel)
			pack = llb.Merge([]llb.State{
				root,
				llb.Diff(root, pypi, llb.WithCustomName("[internal] PyPI packages")),
				llb.Diff(root, conda, llb.WithCustomName("[internal] conda packages")),
			}, llb.WithCustomName("[internal] Python packages"))
		}
	case "r":
		pack = g.installRPackages(root)
	case "julia":
		pack = g.installJuliaPackages(root)
	}
	return pack
}

func (g *generalGraph) compileDevPackages(root llb.State) llb.State {
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
	sb.WriteString("&& locale-gen en_US.UTF-8")

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
		"language": g.Language.Name,
	})
	if g.Language.Version != nil {
		logger = logger.WithField("version", *g.Language.Version)
	}
	logger.Debug("compile base image")

	// Fix https://github.com/tensorchord/envd/issues/1147.
	// Fetch the image metadata from base image.
	base := llb.Image(g.Image, llb.WithMetaResolver(imagemetaresolver.Default()))
	envs, err := base.Env(context.Background())
	if err != nil {
		return llb.State{}, errors.Wrap(err, "failed to get the image metadata")
	}

	// Set the environment variables to RuntimeEnviron to keep it in the resulting image.
	for _, e := range envs {
		// in case the env value also contains `=`
		kv := strings.SplitN(e, "=", 2)
		g.RuntimeEnviron[kv[0]] = kv[1]
	}
	// TODO: inherit the USER from base
	g.User = ""
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
