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
	"context"
	"errors"

	"github.com/moby/buildkit/client"
	"github.com/opencontainers/go-digest"
	"github.com/sirupsen/logrus"
)

// Monitor monitors the progress and print the log to the console.
type Monitor interface {
	Monitor(ctx context.Context, ch chan *client.SolveStatus) error
	Success()
}

var defaultMonitor Monitor

type generalMonitor struct {
	console consoleLogger
}

func NewMonitor() Monitor {
	if defaultMonitor == nil {
		defaultMonitor = generalMonitor{
			console: Current(false),
		}
	}
	return defaultMonitor
}

func (g generalMonitor) Monitor(ctx context.Context, ch chan *client.SolveStatus) error {
	vertexLoggers := make(map[digest.Digest]*logrus.Entry)
	vertexConsoles := make(map[digest.Digest]consoleLogger)
	vertices := make(map[digest.Digest]*client.Vertex)
	introducedVertex := make(map[digest.Digest]bool)

	for {
		select {
		case ss, ok := <-ch:
			if !ok {
				return nil
			}
			for _, vertex := range ss.Vertexes {
				logger := logrus.WithContext(ctx).
					WithField("name", vertex.Name).
					WithField("vertex", shortDigest(vertex.Digest)).
					WithField("cached", vertex.Cached).
					WithField("error", vertex.Error)

				vertexLoggers[vertex.Digest] = logger
				targetConsole := g.console.WithPrefix(vertex.Name)
				vertexConsoles[vertex.Digest] = targetConsole
				vertices[vertex.Digest] = vertex
				if !introducedVertex[vertex.Digest] && (vertex.Cached || vertex.Started != nil) {
					introducedVertex[vertex.Digest] = true
					printVertex(vertex, targetConsole)
					logger.Debug("Vertex started or cached")
				}
				if vertex.Error != "" {
					if !introducedVertex[vertex.Digest] {
						introducedVertex[vertex.Digest] = true
						printVertex(vertex, targetConsole)
					}
					targetConsole.Printf("ERROR: (%s) %s\n", vertex.Name, vertex.Error)
					logger.Error(errors.New(vertex.Error))
				}
				for _, vs := range ss.Statuses {
					vertex, found := vertices[vs.Vertex]
					if !found {
						// No logging for internal operations.
						continue
					}
					logger := vertexLoggers[vs.Vertex]
					targetConsole := vertexConsoles[vs.Vertex]
					progress := int32(0)
					if vs.Total != 0 {
						progress = int32(100.0 * float32(vs.Current) / float32(vs.Total))
					}
					if vs.Completed != nil {
						progress = 100
					}
					logger = logger.WithField("progress", progress).WithField("name", vs.Name)
					if !introducedVertex[vertex.Digest] {
						introducedVertex[vertex.Digest] = true
						printVertex(vertex, targetConsole)
					}
					logger.Debug(vs.ID)
					targetConsole.Printf("%s %d%%\n", vs.ID, progress)
				}
				for _, logLine := range ss.Logs {
					vertex, found := vertices[logLine.Vertex]
					if !found {
						// No logging for internal operations.
						continue
					}
					logger := vertexLoggers[logLine.Vertex]
					targetConsole := vertexConsoles[logLine.Vertex]
					if !introducedVertex[logLine.Vertex] {
						introducedVertex[logLine.Vertex] = true
						printVertex(vertex, targetConsole)
					}
					targetConsole.PrintBytes(logLine.Data)
					logger.Debug(string(logLine.Data))
				}
			}
		case <-ctx.Done():
			return nil
		}
	}
}

func (g generalMonitor) Success() {
	g.console.PrintSuccess()
}
