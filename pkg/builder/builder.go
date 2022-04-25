// Copyright 2022 The MIDI Authors
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

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/progress/progresswriter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tensorchord/MIDI/pkg/docker"
	"github.com/tensorchord/MIDI/pkg/flag"
	"github.com/tensorchord/MIDI/pkg/lang/frontend/starlark"
	"github.com/tensorchord/MIDI/pkg/lang/ir"
	"golang.org/x/sync/errgroup"
)

type Builder interface {
	Build(ctx context.Context) error
	GPUEnabled() bool
}

type generalBuilder struct {
	buildkitdSocket  string
	manifestFilePath string
	progressMode     string
	tag              string

	logger *logrus.Entry
}

func New(buildkitdSocket, manifestFilePath, tag string) Builder {
	return &generalBuilder{
		buildkitdSocket:  buildkitdSocket,
		manifestFilePath: manifestFilePath,
		// TODO(gaocegege): Support other mode?
		progressMode: "auto",
		tag:          tag,
		logger:       logrus.WithField("tag", tag),
	}
}

// GPUEnabled returns true if cuda is enabled.
// It
func (b generalBuilder) GPUEnabled() bool {
	return ir.GPUEnabled()
}

func (b generalBuilder) Build(ctx context.Context) error {
	interpreter := starlark.NewInterpreter()
	if _, err := interpreter.ExecFile(b.manifestFilePath); err != nil {
		return err
	}

	bkClient, err := client.New(ctx, b.buildkitdSocket, client.WithFailFast())
	if err != nil {
		return errors.Wrap(err, "failed to new buildkitd client")
	}
	defer bkClient.Close()

	def, err := ir.Compile(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to compile build.MIDI")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	pw, err := progresswriter.NewPrinter(context.TODO(), os.Stderr, b.progressMode)
	if err != nil {
		return err
	}
	mw := progresswriter.NewMultiWriter(pw)

	var writers []progresswriter.Writer
	w := mw.WithPrefix("", false)
	writers = append(writers, w)

	// Create a pipe to load the image into the docker host.
	pipeR, pipeW := io.Pipe()
	eg.Go(func() error {
		defer func() {
			for _, w := range writers {
				close(w.Status())
			}
		}()
		defer pipeW.Close()
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		b.logger.Debug("building image in ", wd)
		_, err = bkClient.Solve(ctx, def, client.SolveOpt{
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
				flag.FlagContextDir: wd,
				flag.FlagCacheDir:   viper.GetString(flag.FlagCacheDir),
			},
		}, progresswriter.ResetTime(mw.WithPrefix("", false)).Status())
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
		// monitor := progress.NewMonitor()
		// return monitor.Monitor(ctx, ch)
		// not using shared context to not disrupt display but let is finish reporting errors
		<-pw.Done()
		return pw.Err()
	})

	// Load the image to docker host.
	eg.Go(func() error {
		defer pipeR.Close()
		dockerClient, err := docker.NewClient()
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
