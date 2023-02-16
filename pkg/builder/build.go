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
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/docker/cli/cli/config"
	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/session"
	"github.com/moby/buildkit/session/auth/authprovider"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"

	"github.com/tensorchord/envd/pkg/buildkitd"
	"github.com/tensorchord/envd/pkg/driver/docker"
	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/lang/version"
	"github.com/tensorchord/envd/pkg/progress/progresswriter"
	"github.com/tensorchord/envd/pkg/types"
)

func New(ctx context.Context, opt Options) (Builder, error) {
	entries, err := parseOutput(opt.OutputOpts)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse output")
	}

	c, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the current context")
	}

	logrus.WithField("entry", entries).Debug("getting exporter entry")
	// Build docker image by default
	if len(entries) == 0 {
		exportType := client.ExporterDocker
		if c.Builder == types.BuilderTypeMoby {
			exportType = "moby"
		}
		entries = []client.ExportEntry{
			{
				Type: exportType,
			},
		}
	} else if len(entries) > 1 {
		return nil, errors.New("only one output type is supported")
	}

	// Get the language version from the manifest file.
	vc, err := version.New(opt.ManifestFilePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the language version")
	}

	b := &generalBuilder{
		Options:          opt,
		manifestCodeHash: vc.GetDefaultGraphHash(),
		graph:            vc.GetDefaultGraph(),
		entries:          entries,
		logger: logrus.WithFields(logrus.Fields{
			"tag":              opt.Tag,
			"language-version": vc.GetVersion(),
		}),
		Interpreter:         vc.GetStarlarkInterpreter(opt.BuildContextDir),
		GetDepsFilesHandler: vc.GetDefaultGraph().GetDepsFiles,
	}

	var cli buildkitd.Client
	if c.Builder == types.BuilderTypeMoby {
		cli, err = buildkitd.NewMobyClient(ctx,
			c.Builder, c.BuilderAddress, "")
		if err != nil {
			return nil, errors.Wrap(err, "failed to create moby buildkit client")
		}
	} else {
		cli, err = buildkitd.NewClient(ctx,
			c.Builder, c.BuilderAddress, "")
		if err != nil {
			return nil, errors.Wrap(err, "failed to create buildkit client")
		}
	}
	b.Client = cli

	return b, nil
}

func (b generalBuilder) GetGraph() ir.Graph {
	return b.graph
}

// GPUEnabled returns true if cuda is enabled.
func (b generalBuilder) GPUEnabled() bool {
	return b.graph.GPUEnabled()
}

// NumGPUs returns the number of GPUs requested.
func (b generalBuilder) NumGPUs() int {
	return b.graph.GetNumGPUs()
}

func (b generalBuilder) Build(ctx context.Context, force bool) error {
	if !force && !b.checkIfNeedBuild(ctx) {
		return nil
	}

	def, err := b.Compile(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to compile")
	}
	b.definition = def

	pw, err := progresswriter.NewPrinter(ctx, os.Stdout, b.ProgressMode)
	if err != nil {
		return errors.Wrap(err, "failed to create progress writer")
	}

	if err = b.build(ctx, pw); err != nil {
		return errors.Wrap(err, "failed to build")
	}
	return nil
}

func (b generalBuilder) Interpret() error {
	// Evaluate config first.
	if b.ConfigFilePath != "" {
		b.logger.Debug("evaluating config file")
		if _, err := b.ExecFile(b.ConfigFilePath, ""); err != nil {
			return errors.Wrapf(err, "failed to exec starlark file %s", b.ConfigFilePath)
		}
	}

	if _, err := b.ExecFile(b.ManifestFilePath, b.BuildFuncName); err != nil {
		return errors.Wrapf(err, "failed to exec starlark file %s", b.ManifestFilePath)
	}
	return nil
}

func (b generalBuilder) Compile(ctx context.Context) (*llb.Definition, error) {
	envName := filepath.Base(b.BuildContextDir)
	def, err := b.graph.Compile(ctx, envName, b.PubKeyPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compile build.envd")
	}
	b.logger.Debug("compiled build.envd")
	return def, nil
}

func (b generalBuilder) addBuilderTag(labels *map[string]string) {
	(*labels)[types.ImageLabelCacheHash] = b.manifestCodeHash
}

// nolint:unparam
func (b generalBuilder) imageConfig(ctx context.Context) (string, error) {
	labels, err := b.graph.Labels()
	if err != nil {
		return "", errors.Wrap(err, "failed to get labels")
	}
	b.addBuilderTag(&labels)

	ports, err := b.graph.ExposedPorts()
	if err != nil {
		return "", errors.Wrap(err, "failed to get expose ports")
	}
	labels[types.ImageLabelContext] = b.BuildContextDir

	ep, err := b.graph.GetEntrypoint(b.BuildContextDir)
	if err != nil {
		return "", errors.Wrap(err, "failed to get entrypoint")
	}
	b.logger.Debugf("final entrypoint: {%s}\n", ep)

	env := b.graph.GetEnviron()
	user := b.graph.GetUser()

	data, err := ImageConfigStr(labels, ports, ep, env, user)
	if err != nil {
		return "", errors.Wrap(err, "failed to get image config")
	}
	return data, nil
}

func (b generalBuilder) defaultCacheImporter() (*string, error) {
	if b.graph != nil {
		return b.graph.DefaultCacheImporter()
	}
	return nil, nil
}

func (b generalBuilder) build(ctx context.Context, pw progresswriter.Writer) error {
	b.logger.Debug("building envd image")
	ce, err := ParseExportCache([]string{b.ExportCache}, nil)
	if err != nil {
		return errors.Wrap(err, "failed to parse export cache")
	}
	// k := platforms.Format(platforms.DefaultSpec())
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	// Create a pipe to load the image into the docker host.
	pipeR, pipeW := io.Pipe()

	for _, entry := range b.entries {
		// Set up docker config auth.
		dockerConfig := config.LoadDefaultConfigFile(os.Stderr)
		attachable := []session.Attachable{authprovider.NewDockerAuthProvider(dockerConfig)}
		b.logger.WithFields(logrus.Fields{
			"type": entry.Type,
		}).Debug("build image with buildkit")
		switch entry.Type {
		// Create default build.
		case client.ExporterDocker:
			eg.Go(func() error {
				if entry.Attrs == nil {
					entry = client.ExportEntry{
						Type: client.ExporterDocker,
						Attrs: map[string]string{
							"name": b.Tag,
						},
						Output: func(map[string]string) (io.WriteCloser, error) {
							return pipeW, nil
						},
					}
				}
				defer pipeW.Close()
				solveOpt := constructSolveOpt(ce, entry, b, attachable)
				_, err := b.Client.Build(ctx, solveOpt, "envd", b.BuildFunc(), pw.Status())
				if err != nil {
					err = errors.Wrap(&BuildkitdErr{err: err}, "Buildkit error")
					logrus.Errorf("%+v", err)
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
			func(entry client.ExportEntry) {
				eg.Go(func() error {
					solveOpt := constructSolveOpt(ce, entry, b, attachable)
					_, err := b.Client.Build(ctx, solveOpt, "envd", b.BuildFunc(), pw.Status())
					if err != nil {
						err = errors.Wrap(err, "failed to solve LLB")
						return err
					}
					b.logger.Debug("llb def is solved successfully")
					return nil
				})
			}(entry)
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
		}
		return errors.Wrap(err, "failed to wait error group")
	}
	b.logger.Debug("build successfully")
	return nil
}

func constructSolveOpt(ce []client.CacheOptionsEntry, entry client.ExportEntry,
	b generalBuilder, attachable []session.Attachable) client.SolveOpt {
	c, _ := home.GetManager().ContextGetCurrent()
	if entry.Attrs == nil && c.Builder == types.BuilderTypeMoby {
		entry = client.ExportEntry{
			Type: "moby",
			Attrs: map[string]string{
				"name": b.Tag,
			},
		}
	}
	opt := client.SolveOpt{
		CacheExports: ce,
		Exports:      []client.ExportEntry{entry},
		LocalDirs: map[string]string{
			flag.FlagCacheDir:     home.GetManager().CacheDir(),
			flag.FlagBuildContext: b.BuildContextDir,
		},
		Session: attachable,
	}
	if b.UseHTTPProxy {
		opt.FrontendAttrs = map[string]string{
			"build-arg:HTTPS_PROXY": os.Getenv("HTTPS_PROXY"),
			"build-arg:HTTP_PROXY":  os.Getenv("HTTP_PROXY"),
			"build-arg:NO_PROXY":    os.Getenv("NO_PROXY"),
		}
	}
	return opt
}
