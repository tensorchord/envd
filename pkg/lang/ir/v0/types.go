package v0

import (
	"github.com/tensorchord/envd/pkg/editor/vscode"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/progress/compileui"
	"github.com/tensorchord/envd/pkg/types"
)

// A Graph contains the state,
// such as its call stack and thread-local storage.
// TODO(gaocegeg): Refactor it to support order.
type generalGraph struct {
	uid int
	gid int

	OS string
	ir.Language
	Image *string

	Shell   string
	CUDA    *string
	CUDNN   string
	NumGPUs int

	UbuntuAPTSource    *string
	CRANMirrorURL      *string
	JuliaPackageServer *string
	PyPIIndexURL       *string
	PyPIExtraIndexURL  *string

	PublicKeyPath string

	PyPIPackages     []string
	RequirementsFile *string
	PythonWheels     []string
	RPackages        []string
	JuliaPackages    []string
	SystemPackages   []string

	VSCodePlugins   []vscode.Plugin
	UserDirectories []string
	RuntimeEnvPaths []string

	Exec       []ir.RunBuildCommand
	Copy       []ir.CopyInfo
	Mount      []ir.MountInfo
	HTTP       []ir.HTTPInfo
	Entrypoint []string

	Repo types.RepoInfo

	*ir.JupyterConfig
	*ir.GitConfig
	*ir.CondaConfig
	*ir.RStudioServerConfig

	Writer compileui.Writer
	// EnvironmentName is the base name of the environment.
	// It is the BaseDir(BuildContextDir)
	// e.g. mnist, streamlit-mnist
	EnvironmentName string

	ir.RuntimeGraph
}

const (
	shellBASH = "bash"
	shellZSH  = "zsh"
)
