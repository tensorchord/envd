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

package progress

import (
	"strings"

	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
)

func printVertex(vertex *client.Vertex, console consoleLogger) {
	out := []string{"-->"}
	out = append(out, vertex.Name)
	c := console
	if vertex.Cached {
		c = c.WithCached(true)
	}
	c.Printf("%s\n", strings.Join(out, " "))
}

func shortDigest(d digest.Digest) string {
	return d.Hex()[:12]
}
