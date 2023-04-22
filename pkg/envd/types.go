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

package envd

import (
	"fmt"
	"time"

	"github.com/schollz/progressbar/v3"
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/envd-server/api/types"

	"github.com/tensorchord/envd/pkg/lang/ir"
)

const (
	Localhost = "127.0.0.1"
)

var (
	waitingInterval = 1 * time.Second
)

type StartOptions struct {
	Image           string
	EnvironmentName string
	BuildContext    string
	NumGPU          int
	NumCPU          string
	CPUSet          string
	NumMem          string
	Timeout         time.Duration
	ShmSize         int
	Forced          bool
	SshdHost        string

	EngineSource
}

type EngineSource struct {
	DockerSource     *DockerSource
	EnvdServerSource *EnvdServerSource
}

type DockerSource struct {
	Graph        ir.Graph
	MountOptions []string
}

type EnvdServerSource struct {
	Sync bool
}

type StartResult struct {
	// TODO(gaocegege): Make result a chan, to send running status to the receiver.
	SSHPort int
	Address string
	Name    string

	Ports []types.EnvironmentPort
}

type ProgressBar struct {
	bar        *progressbar.ProgressBar
	currStage  int
	totalStage int
	notify     chan struct{}
}

func InitProgressBar(stage int) *ProgressBar {
	done := make(chan struct{})
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSpinnerType(11),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)

	go func() {
		timer := time.NewTicker(time.Millisecond * 100)
		for {
			select {
			case <-done:
				return
			case <-timer.C:
				_ = bar.Add(1)
			}
		}
	}()

	b := ProgressBar{
		notify:     done,
		bar:        bar,
		totalStage: stage,
	}
	return &b
}

func (b *ProgressBar) updateTitle(title string) {
	b.currStage += 1
	b.bar.Describe(fmt.Sprintf("[cyan][%d/%d][reset] %s",
		b.currStage,
		b.totalStage,
		title,
	))
}

func (b *ProgressBar) finish() {
	b.notify <- struct{}{}
	if err := b.bar.Finish(); err != nil {
		logrus.Infof("stop progress bar err: %v\n", err)
	}
}
