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
	"context"
	"io"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/tensorchord/envd/pkg/buildkitd"
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/progress/progresswriter"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

type Builder interface {
	Build(ctx context.Context, pub string) error
	GPUEnabled() bool
	NumGPUs() int
}

type generalBuilder struct {
	manifestFilePath string
	configFilePath   string
	progressMode     string
	tag              string
	buildContextDir  string
	outputType       string
	outputDest       string

	logger *logrus.Entry
	starlark.Interpreter
	buildkitd.Client
}

func New(ctx context.Context, configFilePath, manifestFilePath, buildContextDir, tag, output string, debug bool) (Builder, error) {
	outputType, outputDest, err := parseOutput(output)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse output")
	}
	var mode string = "auto"
	if debug {
		mode = "plain"
	}

	b := &generalBuilder{
		manifestFilePath: manifestFilePath,
		configFilePath:   configFilePath,
		outputType:       outputType,
		outputDest:       outputDest,
		buildContextDir:  buildContextDir,
		// TODO(gaocegege): Support other mode?
		progressMode: mode,
		tag:          tag,
		logger: logrus.WithFields(logrus.Fields{
			"tag": tag,
		}),
	}

	cli, err := buildkitd.NewClient(ctx, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create buildkit client")
	}
	b.Client = cli

	b.Interpreter = starlark.NewInterpreter()
	return b, nil
}

// GPUEnabled returns true if cuda is enabled.
func (b generalBuilder) GPUEnabled() bool {
	return ir.GPUEnabled()
}

// NumGPUs returns the number of GPUs requested.
func (b generalBuilder) NumGPUs() int {
	return ir.NumGPUs()
}

func (b generalBuilder) Build(ctx context.Context, pub string) error {
	def, err := b.compile(ctx, pub)
	if err != nil {
		return errors.Wrap(err, "failed to compile")
	}

	pw, err := progresswriter.NewPrinter(ctx, os.Stdout, b.progressMode)
	if err != nil {
		return errors.Wrap(err, "failed to create progress writer")
	}

	if err = b.build(ctx, def, pw); err != nil {
		return errors.Wrap(err, "failed to build")
	}
	return nil
}

func (b generalBuilder) interpret() error {
	// Evaluate config first.
	if _, err := b.ExecFile(b.configFilePath, ""); err != nil {
		return errors.Wrap(err, "failed to exec starlark file")
	}

	if _, err := b.ExecFile(b.manifestFilePath, "build"); err != nil {
		return errors.Wrap(err, "failed to exec starlark file")
	}
	return nil
}

func (b generalBuilder) compile(ctx context.Context, pub string) (*llb.Definition, error) {
	if err := b.interpret(); err != nil {
		return nil, errors.Wrap(err, "failed to interpret")
	}
	def, err := ir.Compile(ctx, fileutil.Base(b.buildContextDir), pub)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile build.envd")
	}
	b.logger.Debug("compiled build.envd")
	return def, nil
}

func (b generalBuilder) labels(ctx context.Context) (string, error) {
	labels, err := ir.Labels()
	if err != nil {
		return "", errors.Wrap(err, "failed to get labels")
	}
	labels[types.ImageLabelContext] = b.buildContextDir
	data, err := ImageConfigStr(labels)
	if err != nil {
		return "", errors.Wrap(err, "failed to get image config")
	}
	return data, nil
}

func (b generalBuilder) build(ctx context.Context, def *llb.Definition, pw progresswriter.Writer) error {
	labels, err := b.labels(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get labels")
	}
	// k := platforms.Format(platforms.DefaultSpec())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	// Create a pipe to load the image into the docker host.
	pipeR, pipeW := io.Pipe()

	eg.Go(func() error {
		defer pipeW.Close()
		_, err := b.Solve(ctx, def, client.SolveOpt{
			Exports: []client.ExportEntry{
				{
					Type: client.ExporterDocker,
					Attrs: map[string]string{
						"name": b.tag,
						// Ref https://github.com/r2d4/mockerfile/blob/140c6a912bbfdae220febe59ab535ef0acba0e1f/pkg/build/build.go#L65
						"containerimage.config": labels,
					},
					Output: func(map[string]string) (io.WriteCloser, error) {
						return pipeW, nil
					},
				},
			},
			LocalDirs: map[string]string{
				flag.FlagContextDir: b.buildContextDir,
				flag.FlagCacheDir:   home.GetManager().CacheDir(),
			},
			// TODO(gaocegege): Use llb.WithProxy to implement it.
			FrontendAttrs: map[string]string{
				"build-arg:HTTPS_PROXY": os.Getenv("HTTPS_PROXY"),
			},
		}, pw.Status())
		if err != nil {
			err = errors.Wrap(err, "failed to solve LLB")
			b.logger.Error(err)
			return err
		}
		b.logger.Debug("llb def is solved successfully")
		return nil
	})

	// Watch the progress.
	eg.Go(func() error {
		// not using shared context to not disrupt display but let is finish reporting errors
		<-pw.Done()
		return pw.Err()
	})

	if b.outputDest != "" {
		// Save the image to the output file.
		eg.Go(func() error {
			defer pipeR.Close()
			f, err := os.Create(b.outputDest)
			if err != nil {
				return err
			}

			defer f.Close()
			_, err = io.Copy(f, pipeR)
			if err != nil {
				return err
			}

			b.logger.Debug("export the image successfully")
			return nil
		})
	}

	if b.outputDest == "" {
		// Load the image to docker host.
		eg.Go(func() error {
			defer pipeR.Close()
			dockerClient, err := docker.NewClient(ctx)
			if err != nil {
				return errors.Wrap(err, "failed to new docker client")
			}
			b.logger.Debug("loading image to docker host")
			if err := dockerClient.Load(ctx, pipeR, true); err != nil {
				err = errors.Wrap(err, "failed to load docker image")
				b.logger.Error(err)
				return err
			}
			b.logger.Debug("loaded docker image successfully")
			return nil
		})
	}

	err = eg.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			b.logger.Debug("cancelling the error group")
			// Close the pipe on cancels, otherwise the whole thing hangs.
			pipeR.Close()
			return errors.Wrap(err, "build cancelled")
		} else {
			return errors.Wrap(err, "failed to wait error group")
		}
	}
	return nil
}
