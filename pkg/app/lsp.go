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

	"github.com/tensorchord/envd/pkg/lang/lsp"
	cli "github.com/urfave/cli/v2"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// go lsp server uses zap, thus we have to keep two loggers (logrus and zap) in envd.
// zap is only used in lsp subcommand.
var logLevel = zap.NewAtomicLevelAt(zapcore.WarnLevel)

var CommandLSP = &cli.Command{
	// Hide the command since users are not expected to use it directly.
	// it is only used in vscode-envd extension.
	Hidden: true,
	Name:   "lsp",
	Usage:  "Start envd language server",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "address",
			Usage: "Address (hostname:port) to listen on",
		},
		// It is not possible to use the global debug flag because we cannot get it here.
		// The UX is not good enough since we provide two flags about debug:
		// envd --debug lsp --debug
		// Users are not expected to use this command directly, thus it should be fine.
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug logging",
		},
	},
	Action: startLSP,
}

func startLSP(clicontext *cli.Context) error {
	if clicontext.Bool("debug") {
		logLevel.SetLevel(zapcore.DebugLevel)
	}

	logger, cleanup := newzapLogger()
	defer cleanup()
	ctx := protocol.WithLogger(clicontext.Context, logger)

	s := lsp.New()
	err := s.Start(ctx, clicontext.String("address"))
	return err
}

func newzapLogger() (logger *zap.Logger, cleanup func()) {
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = logLevel
	cfg.Development = false
	logger, err := cfg.Build()
	if err != nil {
		panic(fmt.Errorf("failed to initialize logger: %v", err))
	}

	cleanup = func() {
		_ = logger.Sync()
	}
	return logger, cleanup
}
