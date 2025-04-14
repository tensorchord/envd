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
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/tensorchord/envd/pkg/editor/vscode"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/types"
)

// A Graph contains the state,
// such as its call stack and thread-local storage.
// TODO(gaocegeg): Refactor it to support order.
type generalGraph struct {
	uid int `default:"-1"`
	gid int `default:"-1"`

	Languages         []ir.Language
	EnvdSyntaxVersion string
	Image             string
	User              string

	Shell   string
	Dev     bool
	CUDA    *string
	CUDNN   string
	NumGPUs int
	ShmSize int

	UbuntuAPTSource    *string
	CRANMirrorURL      *string
	JuliaPackageServer *string
	PyPIIndexURL       *string
	PyPIExtraIndexURL  *string
	PyPITrust          bool

	PublicKeyPath string

	PyPIPackages     [][]string
	RequirementsFile *string
	PythonWheels     []string
	RPackages        [][]string
	JuliaPackages    [][]string
	SystemPackages   []string

	VSCodePlugins   []vscode.Plugin
	UserDirectories []string

	Exec       []ir.RunBuildCommand
	Copy       []ir.CopyInfo
	Mount      []ir.MountInfo
	HTTP       []ir.HTTPInfo
	Entrypoint []string

	Repo types.RepoInfo

	*ir.JupyterConfig
	*ir.GitConfig
	*ir.CondaConfig
	*ir.UVConfig
	*ir.PixiConfig
	*ir.RStudioServerConfig

	Writer compileui.Writer `json:"-"`
	// EnvironmentName is the base name of the environment.
	// It is the BaseDir(BuildContextDir)
	// e.g. mnist, streamlit-mnist
	EnvironmentName string
	// EnvironmentPath is the full path of this environment.
	EnvironmentPath string
	// WorkingDir is the working directory of this environment.
	// This only affect the `WorkingDir` in the image config.
	WorkingDir string

	// (v1) disable `merge` op for `moby` builder
	// check https://github.com/tensorchord/envd/issues/1693
	DisableMergeOp bool

	ir.RuntimeGraph

	Platform *ocispecs.Platform
}

const (
	shellBASH = "bash"
	shellZSH  = "zsh"
	shellFish = "fish"
)
