package builder

import (
	"context"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/solver/pb"
)

const (
	keyFilename = "filename"
)

func (b generalBuilder) BuildFunc() func(ctx context.Context, c client.Client) (*client.Result, error) {
	return func(ctx context.Context, c client.Client) (*client.Result, error) {
		opts := c.BuildOpts().Opts
		filename := opts[keyFilename]
		if filename == "" {
			filename = "build.envd"
		}

		depsFiles := []string{
			b.pubKeyPath,
			b.configFilePath,
			b.manifestFilePath,
		}
		isUpdated, err := b.CheckDepsFileUpdate(ctx, b.tag, depsFiles)
		if err != nil {
			b.logger.Debugf("failed to check manifest update: %s", err)
		}
		if !isUpdated {
			b.logger.Infof("manifest is not updated, skip building")
			return nil, nil
		}

		def, err := b.compile(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to compile")
		}

		res, err := c.Solve(ctx, client.SolveRequest{
			Definition: def.ToPB(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to solve in BuildFunc")
		}
		ctr, err := c.NewContainer(ctx, client.NewContainerRequest{
			Mounts: []client.Mount{
				{
					Dest:      "/",
					MountType: pb.MountType_BIND,
					Ref:       res.Ref,
				},
			},
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create new container")
		}

		defer ctr.Release(ctx)
		proc, err := ctr.Start(ctx, client.StartRequest{
			// Args:   cfg.Args,
			// Env:    cfg.Env,
			// User:   cfg.User,
			// Cwd:    cfg.Cwd,
			// Tty:    cfg.Tty,
			// Stdin:  cfg.Stdin,
			// Stdout: cfg.Stdout,
			// Stderr: cfg.Stderr,
			Stdout: os.Stdout,
			Args:   []string{"/bin/sh", "-c", "ls"},
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to start the container")
		}
		errCh := make(chan error)
		doneCh := make(chan struct{})
		go func() {
			if err := proc.Wait(); err != nil {
				errCh <- err
				return
			}
			close(doneCh)
		}()
		select {
		case <-doneCh:
		case <-ctx.Done():
			return nil, ctx.Err()
		case err := <-errCh:
			return nil, errors.Wrap(err, "failed to wait the container")
		}
		return res, nil
	}
}
