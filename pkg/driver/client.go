// Copyright 2023 The envd Authors
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

package driver

import (
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
)

type Client interface {
	// Load loads the image from the reader to the docker host.
	Load(ctx context.Context, r io.ReadCloser, quiet bool) error
	StartBuildkitd(ctx context.Context, tag, name, mirror string, timeout time.Duration) (string, error)

	Exec(ctx context.Context, cname string, cmd []string) error

	GetImageWithCacheHashLabel(ctx context.Context, image string, hash string) (types.ImageSummary, error)
	RemoveImage(ctx context.Context, image string) error

	PruneImage(ctx context.Context) (types.ImagesPruneReport, error)

	Stats(ctx context.Context, cname string, statChan chan<- *Stats, done <-chan bool) error
}
