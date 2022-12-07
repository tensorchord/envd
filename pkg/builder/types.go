package builder

import (
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd/pkg/buildkitd"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark"
	"github.com/tensorchord/envd/pkg/lang/ir"
)

type Options struct {
	// ManifestFilePath is the path to the manifest file `build.envd`.
	ManifestFilePath string
	// ConfigFilePath is the path to the config file `config.envd`.
	ConfigFilePath string
	// ProgressMode is the output mode (auto, plain).
	ProgressMode string
	// Tag is the name of the image.
	Tag string
	// BuildContextDir is the directory of the build context.
	BuildContextDir string
	// BuildFuncName is the name of the build func.
	BuildFuncName string
	// PubKeyPath is the path to the ssh public key.
	PubKeyPath string
	// OutputOpts is the output options.
	OutputOpts string
	// ExportCache is the option to export cache.
	// e.g. type=registry,ref=docker.io/username/image
	ExportCache string
	// ImportCache is the option to import cache.
	// e.g. type=registry,ref=docker.io/username/image
	ImportCache string
	// UseHTTPProxy uses HTTPS_PROXY/HTTP_PROXY/NO_PROXY in the build process.
	UseHTTPProxy bool
}

type generalBuilder struct {
	Options
	manifestCodeHash string
	entries          []client.ExportEntry

	definition *llb.Definition

	logger *logrus.Entry
	starlark.Interpreter
	buildkitd.Client

	graph ir.Graph

	GetDepsFilesHandler func([]string) []string
}
