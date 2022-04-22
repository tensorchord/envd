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

package main

import (
	"context"
	"io"
	"os"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/progress/progresswriter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/tensorchord/MIDI/pkg/docker"
	"github.com/tensorchord/MIDI/pkg/lang/frontend/starlark"
	"github.com/tensorchord/MIDI/pkg/lang/ir"
	"github.com/tensorchord/MIDI/pkg/progress"
)

var CommandBuild = &cli.Command{
	Name:      "build",
	Aliases:   []string{"b"},
	Usage:     "build MIDI environment",
	UsageText: `TODO`,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:    "tag",
			Usage:   "Name and optionally a tag in the 'name:tag' format",
			Aliases: []string{"t"},
		},
		&cli.PathFlag{
			Name:    "file",
			Usage:   "Name of the build.MIDI (Default is 'PATH/build.MIDI')",
			Aliases: []string{"f"},
			Value:   "./build.MIDI",
		},
	},

	Action: actionBuild,
}

func actionBuild(clicontext *cli.Context) error {
	path := clicontext.Path("file")
	if path == "" {
		return errors.New("file does not exist")
	}

	interpreter := starlark.NewInterpreter()
	if _, err := interpreter.ExecFile(path); err != nil {
		return err
	}

	bkClient, err := client.New(clicontext.Context, "unix:///run/buildkit/buildkitd.sock", client.WithFailFast())
	if err != nil {
		return errors.Wrap(err, "failed to new buildkitd client")
	}
	defer bkClient.Close()

	def, err := ir.Compile(clicontext.Context)
	if err != nil {
		return errors.Wrap(err, "failed to compile build.MIDI")
	}

	ctx, cancel := context.WithCancel(clicontext.Context)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	// ch := make(chan *client.SolveStatus)

	pw, err := progresswriter.NewPrinter(context.TODO(), os.Stderr, clicontext.String("progress"))
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
		_, err := bkClient.Solve(ctx, def, client.SolveOpt{
			Exports: []client.ExportEntry{
				{
					Type: client.ExporterDocker,
					Attrs: map[string]string{
						"name": clicontext.String("tag"),
					},
					Output: func(map[string]string) (io.WriteCloser, error) {
						return pipeW, nil
					},
				},
			},
		}, progresswriter.ResetTime(mw.WithPrefix("", false)).Status())
		if err != nil {
			err = errors.Wrap(err, "failed to solve LLB")
			logrus.Error(err)
			return err
		}
		logrus.Debug("LLB Def is solved successfully")
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
		logrus.Debug("Loading image to docker host")
		if err := dockerClient.Load(ctx, pipeR, false); err != nil {
			err = errors.Wrap(err, "failed to load docker image")
			logrus.Error(err)
			return err
		}
		logrus.Debug("Loaded docker image successfully")
		return nil
	})

	go func() {
		<-ctx.Done()
		logrus.Debug("cancelling the error group")
		// Close the pipe on cancels, otherwise the whole thing hangs.
		pipeR.Close()
		pipeW.Close()
	}()

	err = eg.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.Wrap(err, "build cancelled")
		} else {
			return errors.Wrap(err, "failed to wait error group")
		}
	}
	monitor := progress.NewMonitor()
	monitor.Success()
	return nil
}
