package compileui

import (
	"io"
	"time"
)

type trace struct {
	w           io.Writer
	startTime   *time.Time
	modeConsole bool
}

func newTrace(w io.Writer, modeConsole bool) *trace {
	return &trace{
		w:           w,
		modeConsole: modeConsole,
	}
}

func (t *trace) init() {
	current := time.Now()
	t.startTime = &current
}
