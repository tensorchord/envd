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
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/builder"
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var CommandBuild = &cli.Command{
	Name:     "build",
	Category: CategoryBasic,
	Aliases:  []string{"b"},
	Usage:    "Build the envd environment",
	Description: `
To build an image using build.envd:
	$ envd build
To build and push the image to a registry:
	$ envd build --output type=image,name=docker.io/username/image,push=true
`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "tag",
			Usage:       "Name and optionally a tag in the 'name:tag' format",
			Aliases:     []string{"t"},
			DefaultText: "PROJECT:dev",
		},
		&cli.PathFlag{
			Name:    "from",
			Usage:   "Function to execute, format `file:func`",
			Aliases: []string{"f"},
			Value:   "build.envd:build",
		},
		&cli.BoolFlag{
			Name:    "use-proxy",
			Usage:   "Use HTTPS_PROXY/HTTP_PROXY/NO_PROXY in the build process",
			Aliases: []string{"proxy"},
			Value:   false,
		},
		&cli.PathFlag{
			Name:    "path",
			Usage:   "Path to the directory containing the build.envd",
			Aliases: []string{"p"},
			Value:   ".",
		},
		&cli.PathFlag{
			Name:    "public-key",
			Usage:   "Path to the public key",
			Aliases: []string{"pubk"},
			Value:   sshconfig.GetPublicKeyOrPanic(),
			Hidden:  true,
		},
		&cli.StringFlag{
			Name:    "output",
			Usage:   "Output destination (e.g. type=tar,dest=path,push=true)",
			Aliases: []string{"o"},
		},
		&cli.BoolFlag{
			Name:  "force",
			Usage: "Force rebuild the image",
			Value: false,
		},
		// https://github.com/urfave/cli/issues/1134#issuecomment-1191407527
		&cli.StringFlag{
			Name:    "export-cache",
			Usage:   "Export the cache (e.g. type=registry,ref=<image>)",
			Aliases: []string{"ec"},
		},
		&cli.StringFlag{
			Name:    "import-cache",
			Usage:   "Import the cache (e.g. type=registry,ref=<image>)",
			Aliases: []string{"ic"},
		},
	},
	Action: build,
}

func build(clicontext *cli.Context) error {
	telemetry.GetReporter().Telemetry("build", nil)
	opt, err := ParseBuildOpt(clicontext)
	if err != nil {
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		"build-context": opt.BuildContextDir,
		"build-file":    opt.ManifestFilePath,
		"config":        opt.ConfigFilePath,
		"tag":           opt.Tag,
	})
	logger.WithFields(logrus.Fields{
		"builder-options": opt,
	}).Debug("starting build command")

	builder, err := GetBuilder(clicontext, opt)
	if err != nil {
		return err
	}
	if err = InterpretEnvdDef(builder); err != nil {
		return err
	}
	return BuildImage(clicontext, builder)
}

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
	tag, err = docker.NormalizeNamed(tag)
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
