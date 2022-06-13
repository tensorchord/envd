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

// package lsp is mainly copied from https://github.com/tilt-dev/starlark-lsp/blob/main/pkg/cli/start.go
package lsp

import (
	"context"
	"io"
	"io/fs"
	"net"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/tilt-dev/starlark-lsp/pkg/analysis"
	"github.com/tilt-dev/starlark-lsp/pkg/document"
	"github.com/tilt-dev/starlark-lsp/pkg/server"
	"go.lsp.dev/jsonrpc2"
	"go.lsp.dev/protocol"
	"go.uber.org/zap"

	"github.com/tensorchord/envd/envd"
)

type BuiltinAnalyzerOptionProvider = func() analysis.AnalyzerOption
type BuiltinFSProvider = func() fs.FS

var builtinAnalyzerOption BuiltinAnalyzerOptionProvider = nil
var providedManagerOptions []document.ManagerOpt

type Server interface {
	Start(ctx context.Context, addr string) error
}

type generalServer struct {
	fsProvider func() fs.FS
}

func New() Server {
	return &generalServer{
		fsProvider: envd.ApiStubs,
	}
}

func (s generalServer) Start(ctx context.Context, addr string) error {
	var err error

	builtinAnalyzerOption = func() analysis.AnalyzerOption {
		return analysis.WithBuiltins(s.fsProvider())
	}

	analyzer, err := createAnalyzer(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to create analyzer")
	}
	if addr != "" {
		err = runSocketServer(ctx, addr, analyzer)
	} else {
		err = runStdioServer(ctx, analyzer)
	}
	if err == context.Canceled {
		err = nil
	}
	return err
}

func runStdioServer(ctx context.Context, analyzer *analysis.Analyzer) error {
	ctx, cancel := context.WithCancel(ctx)
	logger := protocol.LoggerFromContext(ctx)
	logger.Debug("running in stdio mode")
	stdio := struct {
		io.ReadCloser
		io.Writer
	}{
		os.Stdin,
		os.Stdout,
	}

	return launchHandler(ctx, cancel, stdio, analyzer)
}

func runSocketServer(ctx context.Context, addr string, analyzer *analysis.Analyzer) error {
	ctx, cancel := context.WithCancel(ctx)
	var lc net.ListenConfig
	listener, err := lc.Listen(ctx, "tcp4", addr)
	if err != nil {
		cancel()
		return err
	}
	defer func() {
		_ = listener.Close()
	}()

	logger := protocol.LoggerFromContext(ctx).
		With(zap.String("local_addr", listener.Addr().String()))
	ctx = protocol.WithLogger(ctx, logger)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				cancel()
				return nil
			}
			logger.Warn("failed to accept connection", zap.Error(err))
		}
		logger.Debug("accepted connection",
			zap.String("remote_addr", conn.RemoteAddr().String()))

		err = launchHandler(ctx, cancel, conn, analyzer)
		if err != nil {
			cancel()
			return err
		}
	}
}

func initializeConn(conn io.ReadWriteCloser, logger *zap.Logger) (jsonrpc2.Conn, protocol.Client) {
	stream := jsonrpc2.NewStream(conn)
	jsonConn := jsonrpc2.NewConn(stream)
	notifier := protocol.ClientDispatcher(jsonConn, logger.Named("notify"))

	return jsonConn, notifier
}

func createHandler(cancel context.CancelFunc, notifier protocol.Client, analyzer *analysis.Analyzer) jsonrpc2.Handler {
	docManager := document.NewDocumentManager(providedManagerOptions...)
	s := server.NewServer(cancel, notifier, docManager, analyzer)
	h := s.Handler(server.StandardMiddleware...)
	return h
}

func launchHandler(ctx context.Context, cancel context.CancelFunc, conn io.ReadWriteCloser, analyzer *analysis.Analyzer) error {
	logger := protocol.LoggerFromContext(ctx)
	jsonConn, notifier := initializeConn(conn, logger)
	h := createHandler(cancel, notifier, analyzer)
	jsonConn.Go(ctx, h)

	select {
	case <-ctx.Done():
		_ = jsonConn.Close()
		return ctx.Err()
	case <-jsonConn.Done():
		if ctx.Err() == nil {
			if errors.Unwrap(jsonConn.Err()) != io.EOF {
				// only propagate connection error if context is still valid
				return jsonConn.Err()
			}
		}
	}

	return nil
}

func createAnalyzer(ctx context.Context) (*analysis.Analyzer, error) {
	opts := []analysis.AnalyzerOption{
		analysis.WithStarlarkBuiltins(), builtinAnalyzerOption(),
	}

	return analysis.NewAnalyzer(ctx, opts...)
}
