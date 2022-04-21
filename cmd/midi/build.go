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

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/tensorchord/MIDI/pkg/buildkit"
	"github.com/tensorchord/MIDI/pkg/docker"
	"github.com/tensorchord/MIDI/pkg/lang/frontend/starlark"
	"github.com/tensorchord/MIDI/pkg/lang/ir"
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
	},
	Action: actionBuild,
}

func actionBuild(clicontext *cli.Context) error {
	interpreter := starlark.NewInterpreter()
	if _, err := interpreter.ExecFile("./examples/basics/os.midi"); err != nil {
		return err
	}

	bkClient, err := client.New(clicontext.Context, "unix:///run/buildkit/buildkitd.sock")
	if err != nil {
		return errors.Wrap(err, "failed to new buildkitd client")
	}

	def, err := ir.Stmt.Marshal(clicontext.Context, llb.LinuxAmd64)
	if err != nil {
		return errors.Wrap(err, "failed to marshal LLB")
	}

	ctx, cancel := context.WithCancel(clicontext.Context)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)
	ch := make(chan *client.SolveStatus)

	// Create a pipe to load the image into the docker host.
	pipeR, pipeW := io.Pipe()
	eg.Go(func() error {
		if _, err := bkClient.Solve(ctx, def, client.SolveOpt{
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
		}, ch); err != nil {
			err = errors.Wrap(err, "failed to solve LLB")
			logrus.Error(err)
			return err
		}
		logrus.Debug("LLB Def is solved successfully")
		return nil
	})

	// Watch the progress.
	eg.Go(func() error {
		monitor := buildkit.NewMonitor()
		return monitor.Monitor(ctx, ch)
	})

	// Load the image to docker host.
	eg.Go(func() error {
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
		defer pipeR.Close()
		logrus.Debug("Loaded docker image successfully")
		return nil
	})

	go func() {
		for range ctx.Done() {
			logrus.Debug("cancelling the reader group")
			// Close read pipe on cancels, otherwise the whole thing hangs.
			pipeR.Close()
		}
	}()

	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait error group")
	}
	return nil
}
