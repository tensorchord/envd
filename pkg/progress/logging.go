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

package progress

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var currentConsoleMutex sync.Mutex

// consoleLogger is a writer for consoles.
type consoleLogger struct {
	prefix        string
	disableColors bool
	isCached      bool

	// The following are shared between instances and are protected by the mutex.
	mu             *sync.Mutex
	prefixColors   map[string]*color.Color
	nextColorIndex *int
	w              io.Writer
}

// Current returns the current console.
func Current(disableColors bool) consoleLogger {
	return consoleLogger{
		w:              os.Stdout,
		disableColors:  disableColors || color.NoColor,
		prefixColors:   make(map[string]*color.Color),
		nextColorIndex: new(int),
		mu:             &currentConsoleMutex,
	}
}

func (cl consoleLogger) clone() consoleLogger {
	return consoleLogger{
		w:              cl.w,
		prefix:         cl.prefix,
		isCached:       cl.isCached,
		prefixColors:   cl.prefixColors,
		disableColors:  cl.disableColors,
		nextColorIndex: cl.nextColorIndex,
		mu:             cl.mu,
	}
}

// WithPrefix returns a ConsoleLogger with a prefix added.
func (cl consoleLogger) WithPrefix(prefix string) consoleLogger {
	ret := cl.clone()
	ret.prefix = prefix
	return ret
}

// Prefix returns the console's prefix.
func (cl consoleLogger) Prefix() string {
	return cl.prefix
}

// WithCached returns a ConsoleLogger with isCached flag set accordingly.
func (cl consoleLogger) WithCached(isCached bool) consoleLogger {
	ret := cl.clone()
	ret.isCached = isCached
	return ret
}

// PrintSuccess prints the success message.
func (cl consoleLogger) PrintSuccess() {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	successColor.Fprintf(cl.w, "=========================== SUCCESS ===========================\n")
}

// Printf prints formatted text to the console.
func (cl consoleLogger) Printf(format string, args ...interface{}) {
	cl.mu.Lock()
	defer cl.mu.Unlock()
	text := fmt.Sprintf(format, args...)
	text = strings.TrimSuffix(text, "\n")
	for _, line := range strings.Split(text, "\n") {
		cl.printPrefix()
		cl.w.Write([]byte(line))
		cl.w.Write([]byte("\n"))
	}
}

// PrintBytes prints bytes directly to the console.
func (cl consoleLogger) PrintBytes(data []byte) {
	// TODO: This does not deal well with control characters, because of the prefix.
	cl.mu.Lock()
	defer cl.mu.Unlock()
	if !bytes.Contains(data, []byte("\n")) {
		// No prefix when it's not a complete line.
		cl.w.Write(data)
	} else {
		adjustedData := bytes.TrimSuffix(data, []byte("\n"))
		for _, line := range bytes.Split(adjustedData, []byte("\n")) {
			cl.printPrefix()
			cl.w.Write(line)
			cl.w.Write([]byte("\n"))
		}
	}
}

func (cl consoleLogger) printPrefix() {
	// Assumes mu locked.
	if cl.prefix == "" {
		return
	}
	c := noColor
	if !cl.disableColors {
		var found bool
		c, found = cl.prefixColors[cl.prefix]
		if !found {
			c = availablePrefixColors[*cl.nextColorIndex]
			cl.prefixColors[cl.prefix] = c
			*cl.nextColorIndex = (*cl.nextColorIndex + 1) % len(availablePrefixColors)
		}
	}
	c.Fprintf(cl.w, "%s", cl.prefix)
	cl.w.Write([]byte(" | "))
	if cl.isCached {
		cl.w.Write([]byte("*"))
		cachedColor.Fprintf(cl.w, "cached")
		cl.w.Write([]byte("* "))
	}
}
