// Copyright 2022 The MIDI Authors
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

package compileui

import (
	"context"
	"fmt"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/containerd/console"
	"github.com/morikuni/aec"
	"github.com/sirupsen/logrus"
)

type Writer interface {
	Print(s string)
	Finish()
}

type generalWriter struct {
	console console.Console
	phase   string
	trace   *trace
	doneCh  chan bool
	repeatd bool
}

func New(ctx context.Context, out console.File, mode string) (Writer, error) {
	var c console.Console
	switch mode {
	case "auto":
		if cons, err := console.ConsoleFromFile(out); err == nil {
			c = cons
		} else {
			return nil, errors.Wrap(err, "failed to get console")
		}
	case "plain":
	default:
		return nil, errors.Errorf("invalid progress mode %s", mode)
	}

	t := newTrace(out, true)
	t.init()

	w := &generalWriter{
		console: c,
		phase:   "parse build.MIDI and download/cache dependencies",
		trace:   t,
		doneCh:  make(chan bool),
		repeatd: false,
	}
	// TODO(gaocegege): Have a result chan
	//nolint
	go w.run(ctx)
	return w, nil
}

func (w generalWriter) Print(s string) {
	fmt.Fprintln(w.console, s)
}

func (w generalWriter) Finish() {
	w.doneCh <- true
}

func (w *generalWriter) run(ctx context.Context) error {
	displayTimeout := 100 * time.Millisecond
	ticker := time.NewTicker(displayTimeout)
	width, height := w.getSize()
	logger := logrus.WithFields(logrus.Fields{
		"console-height": height,
		"console-width":  width,
	})
	logger.Debug("print compile progress")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-w.doneCh:
			return nil
		case <-ticker.C:
			b := aec.EmptyBuilder.Up(1)
			if !w.repeatd {
				b = b.Down(1)
			}
			w.repeatd = true
			fmt.Fprint(w.console, b.Column(0).ANSI)
			fmt.Fprint(w.console, aec.Hide)
			defer fmt.Fprint(w.console, aec.Show)
			s := fmt.Sprintf("[+] ⌚ %s %.1fs\n", w.phase, time.Since(*w.trace.startTime).Seconds())
			fmt.Fprint(w.console, s)
		}
	}
}

func (w generalWriter) getSize() (int, int) {
	width := 80
	height := 10
	if w.console != nil {
		size, err := w.console.Size()
		if err == nil && size.Width > 0 && size.Height > 0 {
			width = int(size.Width)
			height = int(size.Height)
		}
	}
	return width, height
}
