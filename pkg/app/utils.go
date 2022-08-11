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

package app

import (
	"fmt"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd/pkg/builder"
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/urfave/cli/v2"
)

func parseBuildOpt(clicontext *cli.Context) (builder.Options, error) {
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	if err != nil {
		return builder.Options{}, errors.Wrap(err, "failed to get absolute path of the build context")
	}
	fileName, funcName, err := builder.ParseFromStr(clicontext.String("from"))
	if err != nil {
		return builder.Options{}, err
	}

	manifest, err := filepath.Abs(filepath.Join(buildContext, fileName))
	if err != nil {
		return builder.Options{}, errors.Wrap(err, "failed to get absolute path of the build file")
	}
	if manifest == "" {
		return builder.Options{}, errors.New("file does not exist")
	}

	config := home.GetManager().ConfigFile()

	tag := clicontext.String("tag")
	if tag == "" {
		logrus.Debug("tag not specified, using default")
		tag = fmt.Sprintf("%s:%s", filepath.Base(buildContext), "dev")
	}
	// The current container engine is only Docker. It should be expaned to support other container engines.
	tag, err = docker.NormalizeNamed(tag)
	if err != nil {
		return builder.Options{}, err
	}
	output := ""
	exportCache := clicontext.String("export-cache")
	importCache := clicontext.String("import-cache")

	opt := builder.Options{
		ManifestFilePath: manifest,
		ConfigFilePath:   config,
		BuildFuncName:    funcName,
		BuildContextDir:  buildContext,
		Tag:              tag,
		OutputOpts:       output,
		PubKeyPath:       clicontext.Path("public-key"),
		ProgressMode:     "auto",
		ExportCache:      exportCache,
		ImportCache:      importCache,
	}

	debug := clicontext.Bool("debug")
	if debug {
		opt.ProgressMode = "plain"
	}
	return opt, nil
}
