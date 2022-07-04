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

package compileui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/containerd/console"
	"github.com/morikuni/aec"
	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/editor/vscode"
)

type Writer interface {
	LogVSCodePlugin(p vscode.Plugin, action Action, cached bool)
	LogZSH(action Action, cached bool)
	Finish()
}

type generalWriter struct {
	console     console.Console
	modeConsole bool
	phase       string
	trace       *trace
	doneCh      chan bool
	repeatd     bool
	result      *Result
	lineCount   int
}

func New(ctx context.Context, out console.File, mode string) (Writer, error) {
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

	modeConsole := c != nil
	t := newTrace(out, modeConsole)
	t.init()

	w := &generalWriter{
		console:     c,
		modeConsole: modeConsole,
		phase:       "parse build.envd and download/cache dependencies",
		trace:       t,
		doneCh:      make(chan bool),
		repeatd:     false,
		result: &Result{
			plugins: make([]*PluginInfo, 0),
		},
		lineCount: 0,
	}
	go func() {
		// TODO(gaocegege): Print in text.
		if modeConsole {
			// TODO(gaocegege): Have a result chan
			// nolint
			w.run(ctx)
		} else {
			<-ctx.Done()
		}
	}()
	return w, nil
}

func (w *generalWriter) LogVSCodePlugin(p vscode.Plugin, action Action, cached bool) {
	switch action {
	case ActionStart:
		c := time.Now()
		w.result.plugins = append(w.result.plugins, &PluginInfo{
			Plugin:    p,
			startTime: &c,
			cached:    cached,
		})
	case ActionEnd:
		c := time.Now()
		for i, plugin := range w.result.plugins {
			if plugin.String() == p.String() {
				w.result.plugins[i].endTime = &c
				w.result.plugins[i].cached = cached
			}
		}
	}
}

func (w *generalWriter) LogZSH(action Action, cached bool) {
	switch action {
	case ActionStart:
		c := time.Now()
		w.result.ZSHInfo = &ZSHInfo{
			OHMYZSH:   "oh-my-zsh",
			startTime: &c,
			cached:    cached,
		}
	case ActionEnd:
		c := time.Now()
		w.result.ZSHInfo.endTime = &c
		w.result.ZSHInfo.cached = cached
	}
}

func (w generalWriter) Finish() {
	if w.modeConsole {
		w.doneCh <- true
	}
}

func (w *generalWriter) run(ctx context.Context) error {
	if w.modeConsole {
		displayTimeout := 100 * time.Millisecond
		ticker := time.NewTicker(displayTimeout)
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-w.doneCh:
				w.output(true)
				return nil
			case <-ticker.C:
				w.output(false)
			}
		}
	}
	return nil
}

func (w *generalWriter) output(finished bool) {
	width, _ := w.getSize()
	b := aec.EmptyBuilder.Up(uint(1 + w.lineCount))
	if !w.repeatd {
		b = b.Down(1)
	}
	w.repeatd = true
	if finished {
		fmt.Fprint(w.console, colorRun)
	}
	fmt.Fprint(w.console, b.Column(0).ANSI)
	fmt.Fprint(w.console, aec.Hide)
	defer fmt.Fprint(w.console, aec.Show)

	statusStr := ""
	if finished {
		statusStr = "✅ (finished)"
	}
	s := fmt.Sprintf("[+] ⌚ %s %.1fs %s \n",
		w.phase, time.Since(*w.trace.startTime).Seconds(), statusStr)
	fmt.Fprint(w.console, s)
	loc := 0

	// output shell info.
	if w.result.ZSHInfo != nil {
		timer := time.Since(*w.result.ZSHInfo.startTime).Seconds()
		if w.result.ZSHInfo.endTime != nil {
			timer = w.result.ZSHInfo.endTime.Sub(*w.result.ZSHInfo.startTime).Seconds()
		}
		template := " => download %s"
		if w.result.ZSHInfo.cached {
			template = " => 💽 (cached) download %s"
		}
		timerStr := fmt.Sprintf(" %3.1fs\n", timer)
		out := fmt.Sprintf(template, w.result.ZSHInfo.OHMYZSH)
		out = align(out, timerStr, width)
		fmt.Fprint(w.console, out)
		loc++
	}

	// output vscode plugins.
	for _, p := range w.result.plugins {
		if p.startTime == nil {
			continue
		}
		timer := time.Since(*p.startTime).Seconds()
		if p.endTime != nil {
			timer = p.endTime.Sub(*p.startTime).Seconds()
		}
		timerStr := fmt.Sprintf(" %3.1fs\n", timer)
		template := " => download %s"
		if p.cached {
			template = " => 💽 (cached) download %s"
		}
		out := fmt.Sprintf(template, p.Plugin)
		out = align(out, timerStr, width)
		fmt.Fprint(w.console, out)
		loc++
	}

	// override previous content
	if diff := w.lineCount - loc; diff > 0 {
		logrus.WithFields(logrus.Fields{
			"diff":    diff,
			"plugins": len(w.result.plugins),
		}).Debug("override previous content", diff)
		for i := 0; i < diff; i++ {
			fmt.Fprintln(w.console, strings.Repeat(" ", width))
		}
		fmt.Fprint(w.console, aec.EmptyBuilder.Up(uint(diff)).Column(0).ANSI)
	}
	w.lineCount = loc
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

func align(l, r string, w int) string {
	return fmt.Sprintf("%-[2]*[1]s %[3]s", l, w-len(r)-1, r)
}
