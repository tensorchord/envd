package app

import (
	"github.com/docker/docker/api/types"
	"os"
	"testing"
)

func Test_renderPruneReport(t *testing.T) {

	report := types.ImagesPruneReport{
		ImagesDeleted: []types.ImageDeleteResponseItem{
			{
				Deleted: "sha256:123",
			},
			{
				Untagged: "sha256:456",
			},
		},
		SpaceReclaimed: 666666,
	}
	renderPruneReport(os.Stdout, report)
}
