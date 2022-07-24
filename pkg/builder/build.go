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

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend/gateway/client"
)

func (b generalBuilder) BuildFunc() func(ctx context.Context, c client.Client) (*client.Result, error) {
	return func(ctx context.Context, c client.Client) (*client.Result, error) {
		def, err := b.compile(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to compile")
		}

		imageConfig, err := b.imageConfig(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to get labels")
		}

		sreq := client.SolveRequest{
			Definition: def.ToPB(),
		}
		if b.Options.ImportCache != "" {
			ci, err := ParseImportCache([]string{b.Options.ImportCache})
			if err != nil {
				return nil, errors.Wrap(err, "failed to get the import cache")
			}
			sreq.CacheImports = ci
		}
		res, err := c.Solve(ctx, sreq)
		if err != nil {
			return nil, err
		}

		res.AddMeta(exptypes.ExporterImageConfigKey, []byte(imageConfig))
		b.logger.Debugf("setting image config: %s", imageConfig)

		return res, nil
	}
}
