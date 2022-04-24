package builder

import (
	"context"
	"io"
	"os"

	"github.com/moby/buildkit/client"
	"github.com/moby/buildkit/util/progress/progresswriter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/MIDI/pkg/docker"
	"github.com/tensorchord/MIDI/pkg/lang/frontend/starlark"
	"github.com/tensorchord/MIDI/pkg/lang/ir"
	"golang.org/x/sync/errgroup"
)

type Builder interface {
	Build(ctx context.Context) error
}

type generalBuilder struct {
	buildkitdSocket  string
	manifestFilePath string
	progressMode     string
	tag              string
}

func New(buildkitdSocket, manifestFilePath, tag string) Builder {
	return &generalBuilder{
		buildkitdSocket:  buildkitdSocket,
		manifestFilePath: manifestFilePath,
		// TODO(gaocegege): Support other mode?
		progressMode: "auto",
		tag:          tag,
	}
}

func (b generalBuilder) Build(ctx context.Context) error {
	interpreter := starlark.NewInterpreter()
	if _, err := interpreter.ExecFile(b.manifestFilePath); err != nil {
		return err
	}

	bkClient, err := client.New(ctx, b.buildkitdSocket, client.WithFailFast())
	if err != nil {
		return errors.Wrap(err, "failed to new buildkitd client")
	}
	defer bkClient.Close()

	def, err := ir.Compile(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to compile build.MIDI")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	pw, err := progresswriter.NewPrinter(context.TODO(), os.Stderr, b.progressMode)
	if err != nil {
		return err
	}
	mw := progresswriter.NewMultiWriter(pw)

	var writers []progresswriter.Writer
	w := mw.WithPrefix("", false)
	writers = append(writers, w)

	// Create a pipe to load the image into the docker host.
	pipeR, pipeW := io.Pipe()
	eg.Go(func() error {
		defer func() {
			for _, w := range writers {
				close(w.Status())
			}
		}()
		defer pipeW.Close()
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		logrus.Debug("building image in ", wd)
		_, err = bkClient.Solve(ctx, def, client.SolveOpt{
			Exports: []client.ExportEntry{
				{
					Type: client.ExporterDocker,
					Attrs: map[string]string{
						"name": b.tag,
					},
					Output: func(map[string]string) (io.WriteCloser, error) {
						return pipeW, nil
					},
				},
			},
			LocalDirs: map[string]string{
				"context": wd,
			},
		}, progresswriter.ResetTime(mw.WithPrefix("", false)).Status())
		if err != nil {
			err = errors.Wrap(err, "failed to solve LLB")
			logrus.Error(err)
			return err
		}
		logrus.Debug("llb def is solved successfully")
		return nil
	})

	// Watch the progress.
	eg.Go(func() error {
		// monitor := progress.NewMonitor()
		// return monitor.Monitor(ctx, ch)
		// not using shared context to not disrupt display but let is finish reporting errors
		<-pw.Done()
		return pw.Err()
	})

	// Load the image to docker host.
	eg.Go(func() error {
		defer pipeR.Close()
		dockerClient, err := docker.NewClient()
		if err != nil {
			return errors.Wrap(err, "failed to new docker client")
		}
		logrus.Debug("loading image to docker host")
		if err := dockerClient.Load(ctx, pipeR, false); err != nil {
			err = errors.Wrap(err, "failed to load docker image")
			logrus.Error(err)
			return err
		}
		logrus.Debug("loaded docker image successfully")
		return nil
	})

	go func() {
		<-ctx.Done()
		logrus.Debug("cancelling the error group")
		// Close the pipe on cancels, otherwise the whole thing hangs.
		pipeR.Close()
		pipeW.Close()
	}()

	err = eg.Wait()
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return errors.Wrap(err, "build cancelled")
		} else {
			return errors.Wrap(err, "failed to wait error group")
		}
	}

	return nil
}
