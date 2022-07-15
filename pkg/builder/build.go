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
	"github.com/moby/buildkit/frontend/gateway/client"
)

func (b generalBuilder) BuildFunc() func(ctx context.Context, c client.Client) (*client.Result, error) {
	return func(ctx context.Context, c client.Client) (*client.Result, error) {
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
