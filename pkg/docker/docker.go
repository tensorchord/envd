// Copyright 2022 The MIDI Authors
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

package docker

import (
	"context"
	"io"

	"github.com/docker/docker/client"
)

type Client interface {
	Load(ctx context.Context, r io.ReadCloser, quiet bool) error
}

type generalClient struct {
	*client.Client
}

func NewClient() (Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return generalClient{cli}, nil
}

// Load loads the docker image from the reader into the docker host.
// It's up to the caller to close the io.ReadCloser.
func (g generalClient) Load(ctx context.Context, r io.ReadCloser, quiet bool) error {
	resp, err := g.ImageLoad(ctx, r, quiet)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return nil
}
