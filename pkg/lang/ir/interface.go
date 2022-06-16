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
	"errors"

	"github.com/tensorchord/envd/pkg/editor/vscode"
)

func Base(os, language string) error {
	l, version, err := parseLanguage(language)
	if err != nil {
		return err
	}
	DefaultGraph.Language = Language{
		Name:    l,
		Version: version,
	}
	DefaultGraph.OS = os
	return nil
}

func PyPIPackage(deps []string) {
	DefaultGraph.PyPIPackages = append(DefaultGraph.PyPIPackages, deps...)
}

func RPackage(deps []string) {
	DefaultGraph.RPackages = append(DefaultGraph.RPackages, deps...)
}

func SystemPackage(deps []string) {
	DefaultGraph.SystemPackages = append(DefaultGraph.SystemPackages, deps...)
}

func GPU(numGPUs int) {
	DefaultGraph.NumGPUs = numGPUs
}

func CUDA(version, cudnn string) {
	DefaultGraph.CUDA = &version
	DefaultGraph.CUDNN = &cudnn
}

func VSCodePlugins(plugins []string) error {
	for _, p := range plugins {
		plugin, err := vscode.ParsePlugin(p)
		if err != nil {
			return err
		}
		DefaultGraph.VSCodePlugins = append(DefaultGraph.VSCodePlugins, *plugin)
	}
	return nil
}

// UbuntuAPT updates the Ubuntu apt source.list in the image.
func UbuntuAPT(mode, source string) error {
	if source == "" {
		if mode == pypiIndexModeAuto {
			// If the mode is set to `auto`, envd detects the location of the run
			// then set to the nearest mirror
			return errors.New("auto-mode not implemented")
		}
		return errors.New("source is required")
	}

	DefaultGraph.UbuntuAPTSource = &source
	return nil
}

func PyPIIndex(mode, url, extraURL string) error {
	if url == "" {
		if mode == pypiIndexModeAuto {
			// If the mode is set to `auto`, envd detects the location of the run
			// then set to the nearest index URL.
			return errors.New("auto-mode not implemented")
		}
		return errors.New("url is required")
	}

	DefaultGraph.PyPIIndexURL = &url
	DefaultGraph.PyPIExtraIndexURL = &extraURL
	return nil
}

func CRANMirror(url string) error {
	DefaultGraph.CRANMirrorURL = &url
	return nil
}

func Shell(shell string) error {
	DefaultGraph.Shell = shell
	return nil
}

func Jupyter(pwd string, port int64) error {
	DefaultGraph.JupyterConfig = &JupyterConfig{
		Password: pwd,
		Port:     port,
	}
	return nil
}

func Run(commands []string) error {
	// TODO(gaocegege): Support order-based exec.
	DefaultGraph.Exec = commands
	return nil
}

func Git(name, email, editor string) error {
	DefaultGraph.GitConfig = &GitConfig{
		Name:   name,
		Email:  email,
		Editor: editor,
	}
	return nil
}

func CondaChannel(channel string) error {
	if channel == "" {
		return errors.New("channel is required")
	}

	if !DefaultGraph.CondaEnabled() {
		DefaultGraph.CondaConfig = &CondaConfig{}
	}

	DefaultGraph.CondaConfig.CondaChannel = &channel
	return nil
}

func CondaPackage(deps []string, channel []string) {
	if !DefaultGraph.CondaEnabled() {
		DefaultGraph.CondaConfig = &CondaConfig{}
	}
	DefaultGraph.CondaConfig.CondaPackages = append(
		DefaultGraph.CondaConfig.CondaPackages, deps...)

	if len(channel) != 0 {
		DefaultGraph.CondaConfig.AdditionalChannels = append(
			DefaultGraph.CondaConfig.AdditionalChannels, channel...)
	}
}
