package ir

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"
)

var Stmt llb.State

func BaseStmt(os, language string) {
	base := llb.Image(fmt.Sprintf("docker.io/library/%s:%s", language, os))
	run := base.Run(llb.Shlex("ls -la"))
	Stmt = run.State
}
