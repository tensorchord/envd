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

	"github.com/moby/buildkit/client/llb"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

const (
	templateGitConfig = `
[user]
	email = %s
	name = %s
[core]
	editor = %s

`
)

func (g *generalGraph) compileGit(root llb.State) llb.State {
	if g.GitConfig == nil {
		return root
	}
	content := fmt.Sprintf(templateGitConfig, g.GitConfig.Email, g.GitConfig.Name, g.GitConfig.Editor)
	installPath := fileutil.EnvdHomeDir(".gitconfig")
	gitStage := root.File(llb.Mkfile(installPath,
		0644, []byte(content), llb.WithUIDGID(g.uid, g.gid)))
	return gitStage
}
