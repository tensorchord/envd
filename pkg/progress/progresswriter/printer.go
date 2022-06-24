// Copyright 2022 The envd Authors
// Copyright 2022 The buildkit Authors
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

package progresswriter

import (
	"context"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/containerd/console"
	"github.com/moby/buildkit/client"

	"github.com/tensorchord/envd/pkg/progress/progressui"
)

type printer struct {
	status chan *client.SolveStatus
	done   <-chan struct{}
	err    error
}

func (p *printer) Done() <-chan struct{} {
	return p.done
}

func (p *printer) Err() error {
	return p.err
}

func (p *printer) Status() chan *client.SolveStatus {
	if p == nil {
		return nil
	}
	return p.status
}

type tee struct {
	Writer
	status chan *client.SolveStatus
}

func (t *tee) Status() chan *client.SolveStatus {
	return t.status
}

func Tee(w Writer, ch chan *client.SolveStatus) Writer {
	st := make(chan *client.SolveStatus)
	t := &tee{
		status: st,
		Writer: w,
	}
	go func() {
		for v := range st {
			w.Status() <- v
			ch <- v
		}
		close(w.Status())
		close(ch)
	}()
	return t
}

func NewPrinter(ctx context.Context, out console.File, mode string) (Writer, error) {
	statusCh := make(chan *client.SolveStatus)
	doneCh := make(chan struct{})

	pw := &printer{
		status: statusCh,
		done:   doneCh,
	}

	if v := os.Getenv("BUILDKIT_PROGRESS"); v != "" && mode == "auto" {
		mode = v
	}

	var c console.Console
	switch mode {
	case "auto", "tty", "":
		if cons, err := console.ConsoleFromFile(out); err == nil {
			c = cons
		} else {
			if mode == "tty" {
				return nil, errors.Wrap(err, "failed to get console")
			}
		}
	case "plain":
	default:
		return nil, errors.Errorf("invalid progress mode %s", mode)
	}

	go func() {
		// not using shared context to not disrupt display but let it finish reporting errors
		_, pw.err = progressui.DisplaySolveStatus(ctx, "build envd environment", c, out, statusCh)
		close(doneCh)
	}()
	return pw, nil
}
