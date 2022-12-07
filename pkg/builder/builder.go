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

	"github.com/moby/buildkit/client/llb"

	"github.com/tensorchord/envd/pkg/lang/ir"
)

type Builder interface {
	Build(ctx context.Context, force bool) error
	Interpret() error
	// Compile compiles envd IR to LLB.
	Compile(ctx context.Context) (*llb.Definition, error)
	GPUEnabled() bool
	NumGPUs() int
	GetGraph() ir.Graph
}
