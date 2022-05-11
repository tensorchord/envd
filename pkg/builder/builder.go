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
)

type Builder interface {
	Build(ctx context.Context) error
	GPUEnabled() bool
}

type generalBuilder struct {
	manifestFilePath string
	configFilePath   string
	progressMode     string
	tag              string
	buildContextDir  string

	logger *logrus.Entry
	starlark.Interpreter
	buildkitd.Client
}

func New(ctx context.Context,
	configFilePath, manifestFilePath, buildContextDir, tag string) (Builder, error) {
	b := &generalBuilder{
		manifestFilePath: manifestFilePath,
		configFilePath:   configFilePath,
		buildContextDir:  buildContextDir,
		// TODO(gaocegege): Support other mode?
		progressMode: "auto",
		tag:          tag,
		logger: logrus.WithFields(logrus.Fields{
			"tag": tag,
		}),
	}

	cli, err := buildkitd.NewClient(ctx)
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

func (b generalBuilder) Build(ctx context.Context) error {
	def, err := b.compile(ctx)
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
	if _, err := b.ExecFile(b.configFilePath); err != nil {
		return errors.Wrap(err, "failed to exec starlark file")
	}

	if _, err := b.ExecFile(b.manifestFilePath); err != nil {
		return errors.Wrap(err, "failed to exec starlark file")
	}
	return nil
}

func (b generalBuilder) compile(ctx context.Context) (*llb.Definition, error) {
	if err := b.interpret(); err != nil {
		return nil, errors.Wrap(err, "failed to interpret")
	}
	def, err := ir.Compile(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile build.envd")
	}
	b.logger.Debug("compiled build.envd")
	return def, nil
}

func (b generalBuilder) build(ctx context.Context, def *llb.Definition, pw progresswriter.Writer) error {
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

	err := eg.Wait()
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
