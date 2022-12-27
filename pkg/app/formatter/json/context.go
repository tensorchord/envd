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
	"fmt"

	"github.com/tensorchord/envd/pkg/types"
)

type contextInfo struct {
	Context     string `json:"context"`
	Builder     string `json:"builder"`
	BuilderAddr string `json:"builder_addr"`
	Runner      string `json:"runner"`
	RunnerAddr  string `json:"runner_addr,omitempty"`
	Current     bool   `json:"current"`
}

func PrintContext(contexts types.EnvdContext) error {
	output := []contextInfo{}
	for _, p := range contexts.Contexts {
		item := contextInfo{
			Context:     p.Name,
			Builder:     string(p.Builder),
			BuilderAddr: fmt.Sprintf("%s://%s", p.Builder, p.BuilderAddress),
			Runner:      string(p.Runner),
			Current:     p.Name == contexts.Current,
		}
		if p.RunnerAddress != nil {
			item.RunnerAddr = *p.RunnerAddress
		}
		output = append(output, item)
	}
	return printJSON(output)
}
