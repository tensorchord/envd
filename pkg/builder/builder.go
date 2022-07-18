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
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/tensorchord/envd/pkg/buildkitd"
	"github.com/tensorchord/envd/pkg/docker"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/lang/frontend/starlark"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/progress/progresswriter"
	"github.com/tensorchord/envd/pkg/types"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

type Builder interface {
	Build(ctx context.Context) error
	GPUEnabled() bool
	NumGPUs() int
}

type Options struct {
	// ManifestFilePath is the path to the manifest file `build.envd`.
	ManifestFilePath string
	// ConfigFilePath is the path to the config file `config.envd`.
	ConfigFilePath string
	// ProgressMode is the output mode (auto, plain).
	ProgressMode string
	// Tag is the name of the image.
	Tag string
	// BuildContextDir is the directory of the build context.
	BuildContextDir string
	// BuildFuncName is the name of the build func.
	BuildFuncName string
	// PubKeyPath is the path to the ssh public key.
	PubKeyPath string
	// OutputOpts is the output options.
	OutputOpts string
}

type generalBuilder struct {
	Options
	entries []client.ExportEntry

	logger *logrus.Entry
	starlark.Interpreter
	buildkitd.Client
}

func New(ctx context.Context, opt Options) (Builder, error) {
	entries, err := parseOutput(opt.OutputOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse output")
	}

	logrus.WithField("entry", entries).Debug("getting exporter entry")
	// Build docker image by default
	if len(entries) == 0 {
		entries = []client.ExportEntry{
			{
				Type: client.ExporterDocker,
			},
		}
	} else if len(entries) > 1 {
		return nil, errors.New("only one output type is supported")
	}

	b := &generalBuilder{
		Options: opt,
		entries: entries,
		logger: logrus.WithFields(logrus.Fields{
			"tag": opt.Tag,
		}),
	}

	currentDriver, currentSocket, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the current context")
	}
	cli, err := buildkitd.NewClient(ctx, currentDriver, currentSocket, "")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create buildkit client")
	}
	b.Client = cli

	b.Interpreter = starlark.NewInterpreter(opt.BuildContextDir)
	return b, nil
}

// GPUEnabled returns true if cuda is enabled.
func (b generalBuilder) GPUEnabled() bool {
	return ir.GPUEnabled()
}

// NumGPUs returns the number of GPUs requested.
func (b generalBuilder) NumGPUs() int {
	return ir.NumGPUs()
}

// Always return updated when met error
func (b generalBuilder) CheckDepsFileUpdate(ctx context.Context, tag string, deps []string) (bool, error) {
	dockerClient, err := docker.NewClient(ctx)
	if err != nil {
		return true, err
	}
	image, err := dockerClient.GetImage(ctx, tag)
	if err != nil {
		return true, err
	}
	imageCreatedTime := image.Created

	latestTimestamp := int64(0)
	for _, dep := range deps {
		file, err := os.Stat(dep)
		if err != nil {
			return true, err
		}
		modifiedtime := file.ModTime().Unix()
		// Only needt o use the latest modified time
		if modifiedtime > latestTimestamp {
			latestTimestamp = modifiedtime
		}
	}
	if latestTimestamp > imageCreatedTime {
		return true, nil
	}
	return false, nil
}

func (b generalBuilder) Build(ctx context.Context) error {
	pw, err := progresswriter.NewPrinter(ctx, os.Stdout, b.ProgressMode)
	if err != nil {
		return errors.Wrap(err, "failed to create progress writer")
	}

	if err = b.build(ctx, pw); err != nil {
		return errors.Wrap(err, "failed to build")
	}
	return nil
}

func (b generalBuilder) interpret() error {
	// Evaluate config first.
	if _, err := b.ExecFile(b.ConfigFilePath, ""); err != nil {
		return errors.Wrap(err, "failed to exec starlark file")
	}

	if _, err := b.ExecFile(b.ManifestFilePath, b.BuildFuncName); err != nil {
		return errors.Wrap(err, "failed to exec starlark file")
	}
	return nil
}

func (b generalBuilder) compile(ctx context.Context) (*llb.Definition, error) {
	if err := b.interpret(); err != nil {
		return nil, errors.Wrap(err, "failed to interpret")
	}
	def, err := ir.Compile(ctx, fileutil.Base(b.BuildContextDir), b.PubKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile build.envd")
	}
	b.logger.Debug("compiled build.envd")
	return def, nil
}

func (b generalBuilder) imageConfig(ctx context.Context) (string, error) {
	labels, err := ir.Labels()
	if err != nil {
		return "", errors.Wrap(err, "failed to get labels")
	}
	ports, err := ir.ExposedPorts()
	if err != nil {
		return "", errors.Wrap(err, "failed to get expose ports")
	}
	labels[types.ImageLabelContext] = b.BuildContextDir

	ep, err := ir.Entrypoint(b.BuildContextDir)
	if err != nil {
		return "", errors.Wrap(err, "failed to get entrypoint")
	}
	data, err := ImageConfigStr(labels, ports, ep)
	if err != nil {
		return "", errors.Wrap(err, "failed to get image config")
	}
	return data, nil
}

func (b generalBuilder) build(ctx context.Context, pw progresswriter.Writer) error {
	imageConfig, err := b.imageConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to get labels")
	}
	// k := platforms.Format(platforms.DefaultSpec())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	// Create a pipe to load the image into the docker host.
	pipeR, pipeW := io.Pipe()

	for _, entry := range b.entries {
		// Set up docker config auth.
		attachable := []session.Attachable{authprovider.NewDockerAuthProvider(os.Stderr)}
		switch entry.Type {
		// Create default build.
		case client.ExporterDocker:
			eg.Go(func() error {
				if entry.Attrs == nil {
					entry = client.ExportEntry{
						Type: client.ExporterDocker,
						Attrs: map[string]string{
							"name": b.Tag,
							// Ref https://github.com/r2d4/mockerfile/blob/140c6a912bbfdae220febe59ab535ef0acba0e1f/pkg/build/build.go#L65
							"containerimage.config": imageConfig,
						},
						Output: func(map[string]string) (io.WriteCloser, error) {
							return pipeW, nil
						},
					}
				}
				defer pipeW.Close()
				_, err := b.Client.Build(ctx, client.SolveOpt{
					Exports: []client.ExportEntry{entry},
					LocalDirs: map[string]string{
						// TODO(gaocegege): Move it to BuildFunc with the help
						// of llb.Local
						flag.FlagCacheDir: home.GetManager().CacheDir(),
					},
					Session: attachable,
					// TODO(gaocegege): Use llb.WithProxy to implement it.
					FrontendAttrs: map[string]string{
						"build-arg:HTTPS_PROXY": os.Getenv("HTTPS_PROXY"),
					},
				}, "envd", b.BuildFunc(), pw.Status())

				if err != nil {
					err = errors.Wrap(err, "failed to solve LLB")
					return err
				}
				b.logger.Debug("llb def is solved successfully")
				return nil
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
		default:
			eg.Go(func() error {
				_, err := b.Client.Build(ctx, client.SolveOpt{
					Exports: []client.ExportEntry{entry},
					LocalDirs: map[string]string{
						flag.FlagCacheDir: home.GetManager().CacheDir(),
					},
					Session: attachable,
					// TODO(gaocegege): Use llb.WithProxy to implement it.
					FrontendAttrs: map[string]string{
						"build-arg:HTTPS_PROXY": os.Getenv("HTTPS_PROXY"),
					},
				}, "envd", b.BuildFunc(), pw.Status())

				if err != nil {
					err = errors.Wrap(err, "failed to solve LLB")
					return err
				}
				b.logger.Debug("llb def is solved successfully")
				return nil
			})
		}
	}

	// Watch the progress.
	eg.Go(func() error {
		// not using shared context to not disrupt display but let is finish reporting errors
		<-pw.Done()
		return pw.Err()
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
