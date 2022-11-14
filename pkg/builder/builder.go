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
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

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
)

type Builder interface {
	Build(ctx context.Context, force bool) error
	Interpret() error
	// Compile compiles envd IR to LLB.
	Compile(ctx context.Context) (*llb.Definition, error)
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
	// ExportCache is the option to export cache.
	// e.g. type=registry,ref=docker.io/username/image
	ExportCache string
	// ImportCache is the option to import cache.
	// e.g. type=registry,ref=docker.io/username/image
	ImportCache string
	// UseHTTPProxy uses HTTPS_PROXY/HTTP_PROXY/NO_PROXY in the build process.
	UseHTTPProxy bool
}

type BuildkitdErr struct {
	err error
}

func (e *BuildkitdErr) Error() string {
	return e.err.Error()
}
func (e *BuildkitdErr) Format(s fmt.State, verb rune) { errors.FormatError(e, s, verb) }

type generalBuilder struct {
	Options
	manifestCodeHash string
	entries          []client.ExportEntry

	definition *llb.Definition

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

	manifestHash := ir.GetDefaultGraphHash()

	b := &generalBuilder{
		Options:          opt,
		manifestCodeHash: manifestHash,
		entries:          entries,
		logger: logrus.WithFields(logrus.Fields{
			"tag": opt.Tag,
		}),
	}

	c, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get the current context")
	}
	cli, err := buildkitd.NewClient(ctx, c.Builder, c.BuilderAddress, "")
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
	def, err := ir.Compile(ctx, envName, b.PubKeyPath)
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
	labels, err := ir.Labels()
	if err != nil {
		return "", errors.Wrap(err, "failed to get labels")
	}
	b.addBuilderTag(&labels)

	ports, err := ir.ExposedPorts()
	if err != nil {
		return "", errors.Wrap(err, "failed to get expose ports")
	}
	labels[types.ImageLabelContext] = b.BuildContextDir

	ep, err := ir.CompileEntrypoint(b.BuildContextDir)
	if err != nil {
		return "", errors.Wrap(err, "failed to get entrypoint")
	}
	b.logger.Debugf("final entrypoint: {%s}\n", ep)

	env := ir.CompileEnviron()

	data, err := ImageConfigStr(labels, ports, ep, env)
	if err != nil {
		return "", errors.Wrap(err, "failed to get image config")
	}
	return data, nil
}

func (b generalBuilder) defaultCacheImporter() (*string, error) {
	if ir.DefaultGraph != nil {
		return ir.DefaultGraph.DefaultCacheImporter()
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
		attachable := []session.Attachable{authprovider.NewDockerAuthProvider(os.Stderr)}
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
				solveOpt := client.SolveOpt{
					CacheExports: ce,
					Exports:      []client.ExportEntry{entry},
					LocalDirs: map[string]string{
						flag.FlagCacheDir:     home.GetManager().CacheDir(),
						flag.FlagBuildContext: b.BuildContextDir,
					},
					Session: attachable,
				}
				if b.UseHTTPProxy {
					solveOpt.FrontendAttrs = map[string]string{
						"build-arg:HTTPS_PROXY": os.Getenv("HTTPS_PROXY"),
						"build-arg:HTTP_PROXY":  os.Getenv("HTTP_PROXY"),
						"build-arg:NO_PROXY":    os.Getenv("NO_PROXY"),
					}
				}
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
			eg.Go(func() error {
				solveOpt := client.SolveOpt{
					CacheExports: ce,
					Exports:      []client.ExportEntry{entry},
					LocalDirs: map[string]string{
						flag.FlagCacheDir:     home.GetManager().CacheDir(),
						flag.FlagBuildContext: b.BuildContextDir,
					},
					Session: attachable,
				}
				if b.UseHTTPProxy {
					solveOpt.FrontendAttrs = map[string]string{
						"build-arg:HTTPS_PROXY": os.Getenv("HTTPS_PROXY"),
						"build-arg:HTTP_PROXY":  os.Getenv("HTTP_PROXY"),
						"build-arg:NO_PROXY":    os.Getenv("NO_PROXY"),
					}
				}
				_, err := b.Client.Build(ctx, solveOpt, "envd", b.BuildFunc(), pw.Status())
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
	b.logger.Debug("build successfully")
	return nil
}

func (b generalBuilder) checkIfNeedBuild(ctx context.Context) bool {
	if ir.DefaultGraph.HTTP != nil {
		return true
	}
	depsFiles := []string{
		b.ConfigFilePath,
	}
	depsFiles = getDepsFiles(depsFiles)
	isUpdated, err := b.checkDepsFileUpdate(ctx, b.Tag, b.ManifestFilePath, depsFiles)
	if err != nil {
		b.logger.Debugf("failed to check manifest update: %s", err)
	}
	if !isUpdated {
		b.logger.Infof("manifest is not updated, skip building")
		return false
	}
	return true
}

func getDepsFiles(deps []string) []string {
	tHandle := reflect.TypeOf(*ir.DefaultGraph)
	vHandle := reflect.ValueOf(*ir.DefaultGraph)
	deps = searchFileInGraph(tHandle, vHandle, deps)
	return deps
}

// Match all filed in ir.Graph with the given keyword
func likeFileFiled(str string) bool {
	nameKeyword := []string{
		"File",
		"Path",
		"Wheels",
	}
	if len(nameKeyword) == 0 {
		return true
	}
	re := regexp.MustCompile(strings.Join(nameKeyword, "|"))
	return re.MatchString(str)
}

// search all files in Graph
func searchFileInGraph(tHandle reflect.Type, vHandle reflect.Value, deps []string) []string {
	for i := 0; i < vHandle.NumField(); i++ {
		v := vHandle.Field(i)
		if v.Type().Kind() == reflect.Struct {
			t := v.Type()
			deps = searchFileInGraph(t, v, deps)
		} else if v.Type().Kind() == reflect.Ptr {
			if v.Type().Elem().Kind() == reflect.Struct {
				if v.Elem().CanAddr() {
					t := v.Type().Elem()
					deps = searchFileInGraph(t, v.Elem(), deps)
				}
			}
		} else {
			t := tHandle.Field(i)
			fieldName := t.Name
			if likeFileFiled(fieldName) {
				typeName := t.Type.String()
				if v.Interface() != nil {
					if typeName == "string" {
						deps = append(deps, v.Interface().(string))
					}
					if typeName == "*string" {
						deps = append(deps, *(v.Interface().(*string)))
					}
					if typeName == "[]string" {
						filesList := v.Interface().([]string)
						deps = append(deps, filesList...)
					}
				}
			}
		}
	}
	return deps
}

// nolint:unparam
func (b generalBuilder) checkDepsFileUpdate(ctx context.Context, tag string, manifest string, deps []string) (bool, error) {
	dockerClient, err := docker.NewClient(ctx)
	if err != nil {
		return true, err
	}

	image, err := dockerClient.GetImageWithCacheHashLabel(ctx, tag, b.manifestCodeHash)
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
		// Only need to use the latest modified time
		if modifiedtime > latestTimestamp {
			latestTimestamp = modifiedtime
		}
	}
	if latestTimestamp > imageCreatedTime {
		return true, nil
	}
	return false, nil
}
