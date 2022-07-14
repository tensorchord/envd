package builder

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/frontend/gateway/client"
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
			return nil, err
		}

		return res, nil
	}
}
