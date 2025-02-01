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

package ir

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/containers/image/v5/docker"
	"github.com/containers/image/v5/image"
	"github.com/containers/image/v5/types"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
)

func FetchImageConfig(ctx context.Context, imageName string, platform *specs.Platform) (config specs.ImageConfig, err error) {
	ref, err := docker.ParseReference(fmt.Sprintf("//%s", imageName))
	if err != nil {
		return config, errors.Wrap(err, "failed to parse image reference")
	}
	sys := types.SystemContext{}
	if platform != nil {
		sys.ArchitectureChoice = platform.Architecture
		sys.OSChoice = platform.OS
	}
	src, err := ref.NewImageSource(ctx, &sys)
	if err != nil {
		return config, errors.Wrap(err, "failed to get image source from ref")
	}
	defer src.Close()
	digest, err := docker.GetDigest(ctx, &sys, ref)
	if err != nil {
		return config, errors.Wrap(err, "failed to get the image digest")
	}
	image, err := image.FromUnparsedImage(ctx, &sys, image.UnparsedInstance(src, &digest))
	if err != nil {
		return config, errors.Wrap(err, "failed to get unparsed image")
	}
	img, err := image.OCIConfig(ctx)
	if err != nil {
		return config, errors.Wrap(err, "failed to get OCI config")
	}
	return img.Config, nil
}

func (rg *RuntimeGraph) Dump() (string, error) {
	b, err := json.Marshal(rg)
	if err != nil {
		return "", err
	}
	runtimeGraphCode := string(b)
	return runtimeGraphCode, nil
}

func (rg *RuntimeGraph) Load(code []byte) error {
	err := json.Unmarshal(code, rg)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}
	return nil
}
