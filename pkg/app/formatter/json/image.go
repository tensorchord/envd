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

package json

import (
	"time"

	"github.com/docker/docker/pkg/stringid"
	"github.com/docker/go-units"

	"github.com/tensorchord/envd/pkg/types"
)

type imgInfo struct {
	Name    string `json:"name"`
	Context string `json:"endpoint,omitempty"`
	GPU     bool   `json:"gpu"`
	CUDA    string `json:"cuda,omitempty"`
	CUDNN   string `json:"cudnn,omitempty"`
	ImageID string `json:"image_id"`
	Created string `json:"created"`
	Size    string `json:"size"`
}

func PrintImages(imgs []types.EnvdImage) error {
	output := []imgInfo{}
	for _, img := range imgs {
		CreatedAt := time.Unix(img.Created, 0)
		item := imgInfo{
			Name:    img.Name,
			Context: img.BuildContext,
			GPU:     img.GPU,
			CUDA:    img.CUDA,
			CUDNN:   img.CUDNN,
			ImageID: stringid.TruncateID(img.Digest),
			Created: units.HumanDuration(time.Now().UTC().Sub(CreatedAt)),
			Size:    units.HumanSizeWithPrecision(float64(img.Size), 3),
		}
		output = append(output, item)
	}
	return printJSON(output)
}
