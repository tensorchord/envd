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

package build

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/builder"
	"github.com/tensorchord/envd/pkg/driver/docker"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

// refer to https://github.com/moby/moby/blob/b139a7636f3b6a3d9ad0e2d6dc8bcb687ba2f2cc/daemon/names/names.go#L6
var containerNamePattern = regexp.MustCompile(`[a-zA-Z0-9][a-zA-Z0-9_.-]*`)

func DetectEnvironment(clicontext *cli.Context, buildOpt builder.Options) error {
	context, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the current context")
	}
	opt := envd.Options{
		Context: context,
	}
	engine, err := envd.New(clicontext.Context, opt)
	if err != nil {
		return errors.Wrap(err, "failed to create the docker client")
	}
	// detect if the current environment is running before building
	ctr := filepath.Base(buildOpt.BuildContextDir)
	running, err := engine.IsRunning(clicontext.Context, ctr)
	if err != nil {
		return err
	}
	force := clicontext.Bool("force")
	if running && !force {
		logrus.Errorf("detect container %s is running, please save your data and stop the running container if you need to envd up again.", ctr)
		return errors.Newf("\"%s\" is stil running, please run `envd destroy --name %s` stop it first", ctr, ctr)
	}
	return nil
}

func GetBuilder(clicontext *cli.Context, opt builder.Options) (builder.Builder, error) {
	builder, err := builder.New(clicontext.Context, opt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create the builder")
	}
	return builder, nil
}

func InterpretEnvdDef(builder builder.Builder) error {
	if err := builder.Interpret(); err != nil {
		return errors.Wrap(err, "failed to interpret")
	}
	return nil
}

func BuildImage(clicontext *cli.Context, builder builder.Builder) error {
	force := clicontext.Bool("force")
	if err := builder.Build(clicontext.Context, force); err != nil {
		return errors.Wrap(err, "failed to build the image")
	}
	return nil
}

func CreateEnvNameFromDir(absDir string) (string, error) {
	curDir := filepath.Base(absDir)
	matches := containerNamePattern.FindAllString(curDir, -1)
	if len(matches) == 0 {
		return "", errors.Newf("cannot create a legal container name from %s", curDir)
	}
	name := strings.Join(matches, "")
	// align with docker image name length
	if len(name) > 30 {
		name = name[:30]
	}
	if name != curDir {
		logrus.Debugf("dir %s is not a legal container name, normalize it to %s\n", curDir, name)
	}
	return name, nil
}

func ParseBuildOpt(clicontext *cli.Context) (builder.Options, error) {
	buildContext, err := filepath.Abs(clicontext.Path("path"))
	if err != nil {
		return builder.Options{}, errors.Wrap(err, "failed to get absolute path of the build context")
	}
	fileName, funcName, err := builder.ParseFromStr(clicontext.String("from"))
	if err != nil {
		return builder.Options{}, err
	}

	manifest, err := fileutil.FindFileAbsPath(buildContext, fileName)
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
	// The current container engine is only Docker. It should be expanded to support other container engines.
	tag, err = docker.NormalizeName(tag)
	if err != nil {
		return builder.Options{}, err
	}
	output := clicontext.String("output")
	exportCache := clicontext.String("export-cache")
	importCache := clicontext.String("import-cache")
	useProxy := clicontext.Bool("use-proxy")

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
		UseHTTPProxy:     useProxy,
	}

	debug := clicontext.Bool("debug")
	if debug {
		opt.ProgressMode = "plain"
	}
	return opt, nil
}
