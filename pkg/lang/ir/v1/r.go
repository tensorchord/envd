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

package v1

import (
	"fmt"
	"strings"

	"github.com/moby/buildkit/client/llb"
)

const rPath = "/usr/local/lib/R/site-library"

func (g *generalGraph) installRLang(root llb.State) llb.State {
	g.UserDirectories = append(g.UserDirectories, rPath)
	prepare := root.Run(llb.Shlex(`sh -c "
wget -qO- https://cloud.r-project.org/bin/linux/ubuntu/marutter_pubkey.asc | gpg --dearmor -o /usr/share/keyrings/r-project.gpg &&
echo "deb [signed-by=/usr/share/keyrings/r-project.gpg] https://cloud.r-project.org/bin/linux/ubuntu jammy-cran40/" | tee -a /etc/apt/sources.list.d/r-project.list
"`), llb.WithCustomName("add R public GPG key")).Root()
	run := prepare.Run(
		llb.Shlex(`sh -c "apt-get update && apt-get install -y --no-install-recommends r-base"`),
		llb.WithCustomNamef("[internal] apt install R environment from CRAN repository"))
	return run.Root()
}

func (g generalGraph) installRPackages(root llb.State) llb.State {
	if len(g.RPackages) == 0 {
		return root
	}

	mirrorURL := "https://cran.rstudio.com"
	if g.CRANMirrorURL != nil {
		mirrorURL = *g.CRANMirrorURL
	}

	for _, packages := range g.RPackages {
		command := fmt.Sprintf(`R -e 'options(repos = "%s"); install.packages(c("%s"), lib = "%s")'`, mirrorURL, strings.Join(packages, `","`), rPath)
		run := root.
			Run(llb.Shlex(command), llb.WithCustomNamef("[internal] installing R packages: %s", strings.Join(packages, " ")))
		root = run.Root()
	}
	return root
}
