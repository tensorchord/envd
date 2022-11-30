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

package ir

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
)

// compileUserOwn chown related directories
func (g *Graph) compileUserOwn(root llb.State) llb.State {
	if g.uid == 0 {
		g.RuntimeEnviron["USER"] = "root"
		return root
	}
	g.RuntimeEnviron["USER"] = "envd"
	g.User = "envd"
	if len(g.UserDirectories) == 0 {
		return root.User("envd")
	}
	for _, dir := range g.UserDirectories {
		root = root.Run(llb.Shlex(fmt.Sprintf("chown -R envd:envd %s", dir)),
			llb.WithCustomNamef("[internal] configure user permissions for %s", dir)).Root()
	}
	return root.User("envd")
}

// compileUserGroup creates user `envd`
func (g *Graph) compileUserGroup(root llb.State) llb.State {
	if g.Language.Name == "r" {
		// r-base image already has GID 1000.
		// It is a trick, we actually use GID 1000
		if g.gid == 1000 {
			g.gid = 1001
		}
		if g.uid == 1000 {
			g.uid = 1001
		}
	}

	var res llb.ExecState
	if g.uid == 0 {
		res = root.
			Run(llb.Shlex(fmt.Sprintf("groupadd -g %d envd", 1001)),
				llb.WithCustomName("[internal] still create group envd for root context")).
			Run(llb.Shlex(fmt.Sprintf("useradd -p \"\" -u %d -g envd -s /bin/sh -m envd", 1001)),
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
			Run(llb.Shlex(fmt.Sprintf("groupadd -g %d envd", g.gid)),
				llb.WithCustomName("[internal] create user group envd")).
			Run(llb.Shlex(fmt.Sprintf("useradd -p \"\" -u %d -g envd -s /bin/sh -m envd", g.uid)),
				llb.WithCustomName("[internal] create user envd")).
			Run(llb.Shlex("usermod -a -G sudo envd"),
				llb.WithCustomName("[internal] add user envd to sudoers")).
			Run(llb.Shlex(fmt.Sprintf("install -d -o envd -g %d -m 0700 /home/envd/.config /home/envd/.cache", g.gid)),
				llb.WithCustomName("[internal] mkdir config and cache dir"))
	}
	return res.Root()
}
