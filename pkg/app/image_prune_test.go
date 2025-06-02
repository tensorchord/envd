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

package app

import (
	"os"
	"testing"

	dockerimage "github.com/docker/docker/api/types/image"
)

func Test_renderPruneReport(t *testing.T) {

	report := dockerimage.PruneReport{
		ImagesDeleted: []dockerimage.DeleteResponse{
			{
				Deleted: "sha256:123",
			},
			{
				Untagged: "sha256:456",
			},
		},
		SpaceReclaimed: 666666,
	}
	err := renderPruneReport(os.Stdout, report)
	if err != nil {
		t.Errorf("renderPruneReport() error = %v", err)
	}
}
