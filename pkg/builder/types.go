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
