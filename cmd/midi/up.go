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
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	cli "github.com/urfave/cli/v2"

	"github.com/tensorchord/MIDI/lang/frontend/starlark"
	"github.com/tensorchord/MIDI/lang/llb"
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
	if _, err := interpreter.Eval("midi.base(\"ubuntu\", \"go\")"); err != nil {
		return err
	}
	if err := grpcclient.RunFromEnvironment(appcontext.Context(), llb.Build); err != nil {
		return err
	}
	return nil
}
