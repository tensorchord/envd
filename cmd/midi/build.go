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
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/tensorchord/MIDI/pkg/buildkit"
	"github.com/tensorchord/MIDI/pkg/lang/frontend/starlark"
	"github.com/tensorchord/MIDI/pkg/lang/ir"
)

var CommandBuild = &cli.Command{
	Name:    "build",
	Aliases: []string{"b"},
	Usage:   "build MIDI environment",
	UsageText: `TODO
	`,
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

	eg, ctx := errgroup.WithContext(clicontext.Context)
	ch := make(chan *client.SolveStatus)
	eg.Go(func() error {
		if _, err := bkClient.Solve(ctx, def, client.SolveOpt{
			Exports: []client.ExportEntry{
				{
					Type: client.ExporterImage,
					Attrs: map[string]string{
						"name": "docker.io/username/image",
						"push": "1",
					},
				},
			},
		}, ch); err != nil {
			return errors.Wrap(err, "failed to solve LLB")
		}
		logrus.Debug("LLB Def is solved successfully")
		return nil
	})

	eg.Go(func() error {
		monitor := buildkit.NewMonitor()
		return monitor.Monitor(ctx, ch)
	})

	err = eg.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait error group")
	}
	return nil
}
