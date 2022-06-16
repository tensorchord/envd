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
	"github.com/tensorchord/envd/pkg/editor/vscode"
	"github.com/tensorchord/envd/pkg/progress/compileui"
)

// A Graph contains the state,
// such as its call stack and thread-local storage.
type Graph struct {
	uid int
	gid int

	OS string
	Language

	Shell   string
	CUDA    *string
	CUDNN   *string
	NumGPUs int

	UbuntuAPTSource   *string
	CRANMirrorURL     *string
	PyPIIndexURL      *string
	PyPIExtraIndexURL *string

	PublicKeyPath string

	PyPIPackages   []string
	RPackages      []string
	SystemPackages []string

	VSCodePlugins []vscode.Plugin

	Exec []string
	*JupyterConfig
	*GitConfig
	*CondaConfig

	Writer      compileui.Writer
	CachePrefix string
}

type Language struct {
	Name    string
	Version *string
}

type CondaConfig struct {
	CondaPackages      []string
	AdditionalChannels []string
	CondaChannel       *string
}

type GitConfig struct {
	Name   string
	Email  string
	Editor string
}

type JupyterConfig struct {
	Password string
	Port     int64
}

const (
	shellBASH = "bash"
	shellZSH  = "zsh"
)
