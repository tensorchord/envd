package ir

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/moby/buildkit/client/llb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/tensorchord/envd/pkg/flag"
)

func (g Graph) compileUbuntuAPT(root llb.State) llb.State {
	if g.UbuntuAPTSource != nil {
		logrus.WithField("source", *g.UbuntuAPTSource).Debug("using custom APT source")
		aptSource := llb.Scratch().
			File(llb.Mkdir(filepath.Dir(aptSourceFilePath), 0755, llb.WithParents(true)), llb.WithCustomName("create apt source dir")).
			File(llb.Mkfile(aptSourceFilePath, 0644, []byte(*g.UbuntuAPTSource)), llb.WithCustomName("create apt source file"))
		return llb.Merge([]llb.State{root, aptSource}, llb.WithCustomName("add apt source"))
	}
	return root
}

func (g Graph) compileRun(root llb.State) llb.State {
	if len(g.Exec) == 0 {
		return root
	} else if len(g.Exec) == 1 {
		return root.Run(llb.Shlex(g.Exec[0])).Root()
	}

	run := root.Run(llb.Shlex(g.Exec[0]))
	for _, c := range g.Exec[1:] {
		run = run.Run(llb.Shlex(c))
	}
	return run.Root()
}

func (g Graph) compileBuiltinSystemPackages(root llb.State) llb.State {
	// TODO(gaocegege): Refactor it to avoid shell configuration in built-in system packages.
	// Do not need to install bash or sh since it is built-in
	if g.Shell == shellZSH {
		g.BuiltinSystemPackages = append(g.BuiltinSystemPackages, shellZSH)
	}

	if len(g.BuiltinSystemPackages) == 0 {
		return root
	}

	// Compose the package install command.
	var sb strings.Builder
	sb.WriteString(
		"sh -c \"apt-get update && apt-get install -y --no-install-recommends")
	for _, pkg := range g.BuiltinSystemPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}
	sb.WriteString("\"")

	cacheDir := "/var/cache/apt"
	cacheLibDir := "/var/lib/apt"

	run := root.Run(llb.Shlex(sb.String()),
		llb.WithCustomNamef("(built-in packages) apt-get install %s",
			strings.Join(g.BuiltinSystemPackages, " ")))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	run.AddMount(cacheLibDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheLibDir, llb.CacheMountShared))

	// TODO(gaocegege): Refactor user to a seperate stage.
	res := run.
		Run(llb.Shlex("groupadd -g 1000 envd"), llb.WithCustomName("create user group envd")).
		Run(llb.Shlex("useradd -p \"\" -u 1000 -g envd -s /bin/sh -m envd"), llb.WithCustomName("create user envd")).
		Run(llb.Shlex("adduser envd sudo"), llb.WithCustomName("add user envd to sudoers"))
	return llb.User("envd")(res.Root())
}

func (g *Graph) compileCUDAPackages() llb.State {
	root := llb.Image(
		fmt.Sprintf("nvidia/cuda:%s.0-cudnn%s-devel-%s", *g.CUDA, *g.CUDNN, g.OS))
	g.BuiltinSystemPackages = append(g.BuiltinSystemPackages, []string{
		g.Language,
		fmt.Sprintf("%s-pip", g.Language),
	}...)
	return root
}

func (g Graph) compileSystemPackages(root llb.State) llb.State {
	if len(g.SystemPackages) == 0 {
		return root
	}

	// Compose the package install command.
	var sb strings.Builder
	sb.WriteString("sudo apt-get install -y --no-install-recommends")

	for _, pkg := range g.SystemPackages {
		sb.WriteString(fmt.Sprintf(" %s", pkg))
	}

	cacheDir := "/var/cache/apt"
	cacheLibDir := "/var/lib/apt"

	run := root.Run(llb.Shlex(sb.String()),
		llb.WithCustomNamef("(user-defined packages) apt-get install %s",
			strings.Join(g.SystemPackages, " ")))
	run.AddMount(cacheDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheDir, llb.CacheMountShared))
	run.AddMount(cacheLibDir, llb.Scratch(),
		llb.AsPersistentCacheDir("/"+cacheLibDir, llb.CacheMountShared))
	return run.Root()
}

func (g *Graph) compileBase() llb.State {
	var base llb.State
	if g.CUDA == nil && g.CUDNN == nil {
		base = llb.Image("docker.io/library/python:3.8")
	} else {
		base = g.compileCUDAPackages()
	}
	return base
}

func (g Graph) copyEnvdSSHServer() llb.State {
	// TODO(gaocegege): Remove global var ssh image.
	run := llb.Image(viper.GetString(flag.FlagSSHImage)).
		File(llb.Copy(llb.Image(viper.GetString(flag.FlagSSHImage)),
			"usr/bin/envd-ssh", "/var/envd/bin/envd-ssh",
			&llb.CopyInfo{CreateDestPath: true}), llb.WithCustomName("install envd-ssh"))
	return run
}
