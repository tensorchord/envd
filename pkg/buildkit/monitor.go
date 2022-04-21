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

package buildkit

import (
	"context"

	"github.com/moby/buildkit/client"
	"github.com/sirupsen/logrus"
)

type Monitor interface {
	Monitor(ctx context.Context, ch chan *client.SolveStatus) error
}

type generalMonitor struct {
}

func NewMonitor() Monitor {
	return &generalMonitor{}
}

func (g generalMonitor) Monitor(ctx context.Context, ch chan *client.SolveStatus) error {
	for {
		select {
		case ss, ok := <-ch:
			if !ok {
				return nil
			}
			for _, vs := range ss.Statuses {
				logrus.Debug(vs)
			}
		case <-ctx.Done():
			return nil
		}
	}
}
