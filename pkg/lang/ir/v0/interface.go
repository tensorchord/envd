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
	"github.com/cockroachdb/errors"
	"github.com/opencontainers/go-digest"

	"github.com/tensorchord/envd/pkg/editor/vscode"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/types"
)

func Base(os, language, image string) error {
	l, version, err := parseLanguage(language)
	if err != nil {
		return err
	}
	g := DefaultGraph.(*generalGraph)
	g.Language = ir.Language{
		Name:    l,
		Version: version,
	}
	if len(os) > 0 {
		g.OS = os
	}
	if image != "" {
		g.Image = &image
	}
	return nil
}

func PyPIPackage(deps []string, requirementsFile string, wheels []string) error {
	g := DefaultGraph.(*generalGraph)

	if len(deps) > 0 {
		g.PyPIPackages = append(g.PyPIPackages, deps)
	}
	g.PythonWheels = append(g.PythonWheels, wheels...)

	if requirementsFile != "" {
		g.RequirementsFile = &requirementsFile
	}

	return nil
}

func RPackage(deps []string) {
	g := DefaultGraph.(*generalGraph)

	g.RPackages = append(g.RPackages, deps...)
}

func JuliaPackage(deps []string) {
	g := DefaultGraph.(*generalGraph)

	g.JuliaPackages = append(g.JuliaPackages, deps...)
}

func SystemPackage(deps []string) {
	g := DefaultGraph.(*generalGraph)

	g.SystemPackages = append(g.SystemPackages, deps...)
}

func GPU(numGPUs int) {
	g := DefaultGraph.(*generalGraph)

	g.NumGPUs = numGPUs
}

func CUDA(version, cudnn string) {
	g := DefaultGraph.(*generalGraph)

	g.CUDA = &version
	if len(cudnn) > 0 {
		g.CUDNN = cudnn
	}
}

func VSCodePlugins(plugins []string) error {
	g := DefaultGraph.(*generalGraph)

	for _, p := range plugins {
		plugin, err := vscode.ParsePlugin(p)
		if err != nil {
			return err
		}
		g.VSCodePlugins = append(g.VSCodePlugins, *plugin)
	}
	return nil
}

// UbuntuAPT updates the Ubuntu apt source.list in the image.
func UbuntuAPT(source string) error {
	if source == "" {
		return errors.New("source is required")
	}
	g := DefaultGraph.(*generalGraph)

	g.UbuntuAPTSource = &source
	return nil
}

func PyPIIndex(url, extraURL string, trust bool) error {
	if url == "" {
		return errors.New("url is required")
	}
	g := DefaultGraph.(*generalGraph)

	g.PyPIIndexURL = &url
	g.PyPIExtraIndexURL = &extraURL
	g.PyPITrust = trust
	return nil
}

func CRANMirror(url string) error {
	g := DefaultGraph.(*generalGraph)

	g.CRANMirrorURL = &url
	return nil
}

func JuliaPackageServer(url string) error {
	g := DefaultGraph.(*generalGraph)

	g.JuliaPackageServer = &url
	return nil
}

func Shell(shell string) error {
	g := DefaultGraph.(*generalGraph)

	g.Shell = shell
	return nil
}

func Jupyter(pwd string, port int64) error {
	g := DefaultGraph.(*generalGraph)

	g.JupyterConfig = &ir.JupyterConfig{
		Token: pwd,
		Port:  port,
	}
	return nil
}

func RStudioServer() error {
	g := DefaultGraph.(*generalGraph)

	g.RStudioServerConfig = &ir.RStudioServerConfig{}
	return nil
}

func Run(commands []string, mount bool) error {
	g := DefaultGraph.(*generalGraph)

	g.Exec = append(g.Exec, ir.RunBuildCommand{
		Commands:  commands,
		MountHost: mount,
	})
	return nil
}

func Git(name, email, editor string) error {
	g := DefaultGraph.(*generalGraph)

	g.GitConfig = &ir.GitConfig{
		Name:   name,
		Email:  email,
		Editor: editor,
	}
	return nil
}

func CondaChannel(channel string, useMamba bool) error {
	g := DefaultGraph.(*generalGraph)

	g.CondaConfig.CondaChannel = &channel
	g.CondaConfig.UseMicroMamba = useMamba
	return nil
}

func CondaPackage(deps []string, channel []string, envFile string) error {
	g := DefaultGraph.(*generalGraph)

	g.CondaConfig.CondaPackages = append(
		g.CondaConfig.CondaPackages, deps...)

	g.CondaConfig.CondaEnvFileName = envFile

	if len(channel) != 0 {
		g.CondaConfig.AdditionalChannels = append(
			g.CondaConfig.AdditionalChannels, channel...)
	}
	return nil
}

func Copy(src, dest string) {
	g := DefaultGraph.(*generalGraph)

	g.Copy = append(g.Copy, ir.CopyInfo{
		Source:      src,
		Destination: dest,
	})
}

func Mount(src, dest string) {
	g := DefaultGraph.(*generalGraph)

	g.Mount = append(g.Mount, ir.MountInfo{
		Source:      src,
		Destination: dest,
	})
}

func HTTP(url, checksum, filename string) error {
	info := ir.HTTPInfo{
		URL:      url,
		Filename: filename,
	}
	if len(checksum) > 0 {
		d, err := digest.Parse(checksum)
		if err != nil {
			return err
		}
		info.Checksum = d
	}
	g := DefaultGraph.(*generalGraph)

	g.HTTP = append(g.HTTP, info)
	return nil
}

func Entrypoint(args []string) {
	g := DefaultGraph.(*generalGraph)

	g.Entrypoint = append(g.Entrypoint, args...)
}

func RuntimeCommands(commands map[string]string) {
	g := DefaultGraph.(*generalGraph)

	for k, v := range commands {
		g.RuntimeCommands[k] = v
	}
}

func RuntimeDaemon(commands [][]string) {
	g := DefaultGraph.(*generalGraph)

	g.RuntimeDaemon = append(g.RuntimeDaemon, commands...)
}

func RuntimeExpose(envdPort, hostPort int, serviceName string, listeningAddr string) error {
	g := DefaultGraph.(*generalGraph)

	g.RuntimeExpose = append(g.RuntimeExpose, ir.ExposeItem{
		EnvdPort:      envdPort,
		HostPort:      hostPort,
		ServiceName:   serviceName,
		ListeningAddr: listeningAddr,
	})
	return nil
}

func RuntimeEnviron(env map[string]string, path []string) {
	g := DefaultGraph.(*generalGraph)

	for k, v := range env {
		g.RuntimeEnviron[k] = v
	}
	g.RuntimeEnvPaths = append(g.RuntimeEnvPaths, path...)
}

func RuntimeInitScript(commands []string) {
	g := DefaultGraph.(*generalGraph)

	g.RuntimeInitScript = append(g.RuntimeInitScript, commands)
}

func Repo(url, description string) {
	g := DefaultGraph.(*generalGraph)

	g.Repo = types.RepoInfo{
		Description: description,
		URL:         url,
	}
}
