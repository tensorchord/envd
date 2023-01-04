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
	"strings"

	"github.com/moby/buildkit/client/llb"

	"github.com/tensorchord/envd/pkg/types"
)

// compileUserOwn chown related directories
func (g *generalGraph) compileUserOwn(root llb.State) llb.State {
	if g.uid == 0 {
		g.RuntimeEnviron["USER"] = "root"
		return root
	}
	g.RuntimeEnviron["USER"] = "envd"
	g.User = "envd"
	for _, dir := range g.UserDirectories {
		root = root.Run(llb.Shlexf("chown -R envd:envd %s", dir),
			llb.WithCustomNamef("[internal] configure user permissions for %s", dir)).Root()
	}
	user := root.User("envd")
	// re-add the env since it's a different user
	for _, env := range types.BaseEnvironment {
		user = user.AddEnv(env.Name, env.Value)
	}
	user = user.AddEnv("PATH", strings.Join(g.RuntimeEnvPaths, ":"))
	return user
}

// compileUserGroup creates user `envd`
func (g *generalGraph) compileUserGroup(root llb.State) llb.State {
	var res llb.ExecState
	if g.uid == 0 {
		res = root.
			Run(llb.Shlexf("groupadd -g %d envd", 1001),
				llb.WithCustomName("[internal] still create group envd for root context")).
			Run(llb.Shlexf(`useradd -p "" -u %d -g envd -s /bin/sh -m envd`, 1001),
				llb.WithCustomName("[internal] still create user envd for root context")).
			Run(llb.Shlex("usermod -s /bin/sh root"),
				llb.WithCustomName("[internal] set root default shell to /bin/sh")).
			Run(llb.Shlex("sed -i \"s/envd:x:1001:1001/envd:x:0:0/g\" /etc/passwd"),
				llb.WithCustomName("[internal] set envd uid to 0 as root")).
			Run(llb.Shlex("sed -i \"s./root./home/envd.g\" /etc/passwd"),
				llb.WithCustomName("[internal] set root home dir to /home/envd")).
			Run(llb.Shlex("sed -i \"s/envd:x:1001/envd:x:0/g\" /etc/group"),
				llb.WithCustomName("[internal] set envd group to 0 as root group"))
	} else {
		res = root.
			Run(llb.Shlexf("groupadd -g %d envd", g.gid),
				llb.WithCustomNamef("[internal] create user group envd(g:%d)", g.gid)).
			Run(llb.Shlexf(`useradd -p "" -u %d -g envd -s /bin/sh -m envd`, g.uid),
				llb.WithCustomNamef("[internal] create user envd(u:%d)", g.uid)).
			Run(llb.Shlex("usermod -a -G sudo envd"),
				llb.WithCustomName("[internal] add user envd to sudoers")).
			Run(llb.Shlexf("install -d -o envd -g %d -m 0700 /home/envd/.config /home/envd/.cache", g.gid),
				llb.WithCustomName("[internal] mkdir config and cache"))
	}
	return res.Root()
}
