package driver

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
)

type Client interface {
	// Load loads the image from the reader to the docker host.
	Load(ctx context.Context, r io.ReadCloser, quiet bool) error
	StartBuildkitd(ctx context.Context, tag, name, mirror string) (string, error)

	Exec(ctx context.Context, cname string, cmd []string) error

	GetImageWithCacheHashLabel(ctx context.Context, image string, hash string) (types.ImageSummary, error)
	RemoveImage(ctx context.Context, image string) error

	PruneImage(ctx context.Context) (types.ImagesPruneReport, error)

	Stats(ctx context.Context, cname string, statChan chan<- *Stats, done <-chan bool) error
}
