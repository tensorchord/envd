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

package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/solver/pb"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	buildutil "github.com/tensorchord/envd/pkg/app/build"
	"github.com/tensorchord/envd/pkg/app/telemetry"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
)

var CommandDebugLLB = &cli.Command{
	Name:        "llb",
	Category:    CategoryOther,
	Aliases:     []string{"b"},
	Usage:       "dump buildkit LLB in human-readable format.",
	Description: ``,

	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "dot",
			Usage: "Output dot format",
		},
		&cli.PathFlag{
			Name:    "path",
			Usage:   "Path to the directory containing the build.envd",
			Aliases: []string{"p"},
			Value:   ".",
		},
		&cli.PathFlag{
			Name:    "from",
			Usage:   "Function to execute, format `file:func`",
			Aliases: []string{"f"},
			Value:   "build.envd:build",
		},
		&cli.PathFlag{
			Name:    "public-key",
			Usage:   "Path to the public key",
			Aliases: []string{"pubk"},
			Value:   sshconfig.GetPublicKeyOrPanic(),
			Hidden:  true,
		},
	},

	Action: debugLLB,
}

func debugLLB(clicontext *cli.Context) error {
	telemetry.GetReporter().Telemetry("debug-llb")
	opt, err := buildutil.ParseBuildOpt(clicontext)
	if err != nil {
		return err
	}

	logger := logrus.WithFields(logrus.Fields{
		"build-context": opt.BuildContextDir,
		"build-file":    opt.ManifestFilePath,
		"config":        opt.ConfigFilePath,
		"tag":           opt.Tag,
	})
	logger.WithFields(logrus.Fields{
		"builder-options": opt,
	}).Debug("starting debug llb command")

	builder, err := buildutil.GetBuilder(clicontext, opt)
	if err != nil {
		return err
	}
	if err = buildutil.InterpretEnvdDef(builder); err != nil {
		return err
	}

	def, err := builder.Compile(clicontext.Context)
	if err != nil {
		return errors.Wrap(err, "failed to compile envd IR to LLB")
	}

	ops, err := loadLLB(def)
	if err != nil {
		return errors.Wrap(err, "failed to load LLB")
	}

	if clicontext.Bool("dot") {
		writeDot(ops, os.Stdout)
	} else {
		enc := json.NewEncoder(os.Stdout)
		for _, op := range ops {
			if err := enc.Encode(op); err != nil {
				return errors.Wrap(err, "failed to encode LLB op")
			}
		}
	}
	return nil
}

type llbOp struct {
	Op         pb.Op
	Digest     digest.Digest
	OpMetadata pb.OpMetadata
}

// Refer to https://github.com/moby/buildkit/blob/master/cmd/buildctl/debug/dumpllb.go#L17:5
func loadLLB(def *llb.Definition) ([]llbOp, error) {
	var ops []llbOp
	for _, dt := range def.Def {
		var op pb.Op
		if err := (&op).Unmarshal(dt); err != nil {
			return nil, errors.Wrap(err, "failed to parse op")
		}
		dgst := digest.FromBytes(dt)
		ent := llbOp{Op: op, Digest: dgst, OpMetadata: def.Metadata[dgst]}
		ops = append(ops, ent)
	}
	return ops, nil
}

func writeDot(ops []llbOp, w io.Writer) {
	// TODO: print OpMetadata
	fmt.Fprintln(w, "digraph {")
	defer fmt.Fprintln(w, "}")
	for _, op := range ops {
		name, shape := attr(op.Digest, op.Op)
		fmt.Fprintf(w, "  %q [label=%q shape=%q];\n", op.Digest, name, shape)
	}
	for _, op := range ops {
		for i, inp := range op.Op.Inputs {
			label := ""
			if eo, ok := op.Op.Op.(*pb.Op_Exec); ok {
				for _, m := range eo.Exec.Mounts {
					if int(m.Input) == i && m.Dest != "/" {
						label = m.Dest
					}
				}
			}
			fmt.Fprintf(w, "  %q -> %q [label=%q];\n", inp.Digest, op.Digest, label)
		}
	}
}

func attr(dgst digest.Digest, op pb.Op) (string, string) {
	switch op := op.Op.(type) {
	case *pb.Op_Source:
		return op.Source.Identifier, "ellipse"
	case *pb.Op_Exec:
		return generateExecNode(op.Exec)
	case *pb.Op_Build:
		return "build", "box3d"
	case *pb.Op_Merge:
		return "merge", "invtriangle"
	case *pb.Op_Diff:
		return "diff", "doublecircle"
	case *pb.Op_File:
		names := []string{}

		for _, action := range op.File.Actions {
			var name string

			switch act := action.Action.(type) {
			case *pb.FileAction_Copy:
				name = fmt.Sprintf("copy{src=%s, dest=%s}", act.Copy.Src, act.Copy.Dest)
			case *pb.FileAction_Mkfile:
				name = fmt.Sprintf("mkfile{path=%s}", act.Mkfile.Path)
			case *pb.FileAction_Mkdir:
				name = fmt.Sprintf("mkdir{path=%s}", act.Mkdir.Path)
			case *pb.FileAction_Rm:
				name = fmt.Sprintf("rm{path=%s}", act.Rm.Path)
			}

			names = append(names, name)
		}
		return strings.Join(names, ","), "note"
	default:
		return dgst.String(), "plaintext"
	}
}

func generateExecNode(op *pb.ExecOp) (string, string) {
	mounts := []string{}
	for _, m := range op.Mounts {
		mstr := fmt.Sprintf("selector=%s, target=%s, mount-type=%s", m.Selector,
			m.Dest, m.MountType)
		if m.CacheOpt != nil {
			mstr = mstr + fmt.Sprintf(" cache-id=%s, cache-share-mode = %s",
				m.CacheOpt.ID, m.CacheOpt.Sharing)
		}
		mounts = append(mounts, mstr)
	}

	name := fmt.Sprintf("user=%s, cwd=%s, args={%s}, mounts={%s}, env={%s}",
		op.Meta.User,
		op.Meta.Cwd,
		strings.Join(op.Meta.Args, " "),
		strings.Join(mounts, " "),
		strings.Join(op.Meta.Env, " "),
	)

	return name, "box"
}
