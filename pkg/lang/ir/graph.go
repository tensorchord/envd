package ir

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/tensorchord/envd/pkg/progress/compileui"
)

type Graph interface {
	Compile(ctx context.Context, envName string, pub string) (*llb.Definition, error)

	graphDebugger
	graphVisitor
}

type graphDebugger interface {
	SetWriter(w compileui.Writer)
}

type graphVisitor interface {
	GetDepsFiles(deps []string) []string
	GPUEnabled() bool
	GetNumGPUs() int
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
}
