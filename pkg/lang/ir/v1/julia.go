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
	"github.com/sirupsen/logrus"
)

func (g generalGraph) installJulia(root llb.State) (llb.State, error) {
	return llb.State{}, errors.New("not implemented")
}

func (g generalGraph) installJuliaPackages(root llb.State) llb.State {
	if len(g.JuliaPackages) == 0 {
		return root
	}

	var sb strings.Builder

	sb.WriteString(`/usr/local/julia/bin/julia -e 'using Pkg; Pkg.add([`)
	for i, pkg := range g.JuliaPackages {
		sb.WriteString(fmt.Sprintf(`"%s"`, pkg))
		if i != len(g.JuliaPackages)-1 {
			sb.WriteString(", ")
		}
	}

	sb.WriteString(`])'`)

	// TODO(gaocegege): Support cache.
	cmd := sb.String()
	logrus.Debug("install julia packages: ", cmd)
	root = llb.User("envd")(root)
	if g.JuliaPackageServer != nil {
		root = root.AddEnv("JULIA_PKG_SERVER", *g.JuliaPackageServer)
	}
	run := root.
		Run(llb.Shlex(cmd), llb.WithCustomNamef("install julia packages"))

	return run.Root()
}
