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
	_ "embed"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/client/llb"
)

const (
	juliaRootDir = "/opt/julia"               // Location of downloaded Julia binary and other files
	juliaBinDir  = "/opt/julia/bin"           // Location of Julia executable binary file
	juliaPkgDir  = "/opt/julia/user_packages" // Location of additional packages installed via Julia
	juliaBinName = "julia.tar.gz"             // Julia archive name
)

//go:embed julia.sh
var downloadJuliaBashScript string

// getJuliaBinary returns the llb.State only after setting up Julia environment
// A successful run of getJuliaBinary should set up the Julia environment
func (g generalGraph) getJuliaBinary(root llb.State) llb.State {

	base := llb.Image(builderImage)
	builder := base.
		Run(llb.Shlexf("sh -c '%s'", downloadJuliaBashScript),
			llb.WithCustomName("[internal] downloading julia binary")).Root()

	var path = filepath.Join("/tmp", juliaBinName)
	setJulia := root.
		File(llb.Copy(builder, path, path),
			llb.WithCustomNamef("[internal] copying %s to /tmp", juliaBinName)).
		File(llb.Mkdir(juliaRootDir, 0755, llb.WithParents(true)),
			llb.WithCustomNamef("[internal] creating %s folder for julia binary", juliaRootDir)).
		Run(llb.Shlexf(`bash -c "tar zxvf %s --strip 1 -C %s && rm %s"`, path, juliaRootDir, path),
			llb.WithCustomNamef("[internal] unpack julia archive under %s", juliaRootDir))

	return setJulia.Root()
}

// installJulia returns the llb.State only after adding the Julia environment to $PATH
// A successful run of installJulia should add Julia to global environment path
func (g *generalGraph) installJulia(root llb.State) llb.State {

	confJulia := g.getJuliaBinary(root)
	confJulia = g.updateEnvPath(confJulia, juliaBinDir)

	return confJulia
}

// installJuliaPackages returns the llb.State only after installing required Julia packages
// A successful run of installJuliaPackages should install Julia packages under "/opt/julia/user_packages" and export the path
func (g *generalGraph) installJuliaPackages(root llb.State) llb.State {

	if len(g.JuliaPackages) == 0 {
		return root
	}

	root = root.File(llb.Mkdir(juliaPkgDir, 0755, llb.WithParents(true)),
		llb.WithCustomName("[internal] creating folder for julia packages"))

	// Allow root to utilize the installed Julia environment
	root = g.updateEnvPath(root, juliaBinDir)

	// Export "/opt/julia/user_packages" as the additional library path for root
	root = root.AddEnv("JULIA_DEPOT_PATH", juliaPkgDir)

	// Export "/opt/julia/user_packages" as the additional library path for users
	g.RuntimeEnviron["JULIA_DEPOT_PATH"] = juliaPkgDir

	// Change owner of the "/opt/julia/user_packages" to users
	g.UserDirectories = append(g.UserDirectories, juliaPkgDir)

	for _, packages := range g.JuliaPackages {
		command := fmt.Sprintf(`julia -e 'using Pkg; Pkg.add(["%s"])'`, strings.Join(packages, `","`))
		run := root.
			Run(llb.Shlex(command), llb.WithCustomNamef("[internal] installing Julia pacakges: %s", strings.Join(packages, " ")))
		root = run.Root()
	}

	return root
}
