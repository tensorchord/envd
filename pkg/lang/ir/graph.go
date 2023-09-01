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

package ir

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	ocispecs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/tensorchord/envd/pkg/progress/compileui"
)

type Graph interface {
	Compile(ctx context.Context, envPath string, pub string, platform *ocispecs.Platform, progressMode string) (*llb.Definition, error)

	graphDebugger
	graphVisitor
	graphSerializer
}

type graphSerializer interface {
	GeneralGraphFromLabel(label []byte) (Graph, error)
}

type graphDebugger interface {
	SetWriter(w compileui.Writer)
}

type graphVisitor interface {
	GetDepsFiles(deps []string) []string
	GPUEnabled() bool
	GetNumGPUs() int
	GetShmSize() int
	IsDev() bool
	Labels() (map[string]string, error)
	ExposedPorts() (map[string]struct{}, error)
	GetEntrypoint(buildContextDir string) ([]string, error)
	GetShell() string
	GetEnvironmentName() string
	GetMount() []MountInfo
	GetJupyterConfig() *JupyterConfig
	GetRStudioServerConfig() *RStudioServerConfig
	GetExposedPorts() []ExposeItem
	DefaultCacheImporter() (*string, error)
	GetEnviron() []string
	GetHTTP() []HTTPInfo
	GetRuntimeCommands() map[string]string
	GetUser() string
	GetPlatform() *ocispecs.Platform
	GetWorkingDir() string
}
