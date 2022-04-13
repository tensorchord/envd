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
	cli "github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"

	"github.com/tensorchord/MIDI/lang/frontend/starlark"
	"github.com/tensorchord/MIDI/lang/ir"
)

var CommandUp = &cli.Command{
	Name:    "up",
	Aliases: []string{"u"},
	Usage:   "build and run MIDI environment",
	UsageText: `TODO
	`,
	Action: actionUp,
}

func actionUp(clicontext *cli.Context) error {
	interpreter := starlark.NewInterpreter()
	// TODO(gaocegege): Remove func call prefix.
	if _, err := interpreter.Eval(`
midi.base("alpine3.15", "python")
`); err != nil {
		return err
	}

	bkClient, err := client.New(clicontext.Context, "unix:///run/buildkit/buildkitd.sock")
	if err != nil {
		return errors.Wrap(err, "Failed to create the buildkit client")
	}

	def, err := ir.Stmt.Marshal(clicontext.Context, llb.LinuxAmd64)
	if err != nil {
		return errors.Wrap(err, "Failed to marshal the stmt")
	}

	eg, ctx := errgroup.WithContext(clicontext.Context)
	ch := make(chan *client.SolveStatus)
	eg.Go(func() error {

		if _, err := bkClient.Solve(ctx, def, client.SolveOpt{}, ch); err != nil {
			return errors.Wrap(err, "Failed to solve the LLB definition")
		}
		return nil
	})
	return nil
}
