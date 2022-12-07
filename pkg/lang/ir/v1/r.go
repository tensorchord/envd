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

package v1

import (
	"fmt"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
)

func (g generalGraph) installRLang(root llb.State) (llb.State, error) {
	return llb.State{}, errors.New("not implemented")
}

func (g generalGraph) installRPackages(root llb.State) llb.State {
	if len(g.RPackages) == 0 {
		return root
	}
	// TODO(terrytangyuan): Support different CRAN mirrors
	var sb strings.Builder
	mirrorURL := "https://cran.rstudio.com"
	if g.CRANMirrorURL != nil {
		mirrorURL = *g.CRANMirrorURL
	}
	sb.WriteString(fmt.Sprintf(`R -e 'options(repos = c(CRAN = "%s")); install.packages(c(`, mirrorURL))
	for i, pkg := range g.RPackages {
		sb.WriteString(fmt.Sprintf(`"%s"`, pkg))
		if i != len(g.RPackages)-1 {
			sb.WriteString(", ")
		}
	}
	sb.WriteString(`))'`)

	// TODO(terrytangyuan): Support cache.
	cmd := sb.String()
	root = llb.User("envd")(root)
	run := root.Run(llb.Shlex(cmd), llb.WithCustomNamef("install R packages"))
	return run.Root()
}
