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
	"context"

	"github.com/tensorchord/envd/pkg/lang/lsp"
	cli "github.com/urfave/cli/v2"
)

var CommandLSP = &cli.Command{
	Hidden: true,
	Name:   "lsp",
	Usage:  "Start envd language server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "Address (hostname:port) to listen on",
		},
	},
	Action: startLSP,
}

func startLSP(clicontext *cli.Context) error {
	s := lsp.New()
	err := s.Start(context.Background(), clicontext.String("address"))
	return err
}
