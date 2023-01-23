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
	"context"
	"time"

	"github.com/cockroachdb/errors"
	ui "github.com/gizak/termui/v3"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/app/telemetry"
	"github.com/tensorchord/envd/pkg/driver/docker"
	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/metrics"
	"github.com/tensorchord/envd/pkg/types"
)

var CommandTop = &cli.Command{
	Name:     "top",
	Category: CategoryManagement,
	Usage:    "Show statistics about the containers managed by the environment.",
	Flags:    []cli.Flag{},
	Action:   top,
}

func top(clicontext *cli.Context) error {
	defer ui.Close()
	if err := ui.Init(); err != nil {
		return err
	}

	context, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return errors.Wrap(err, "failed to get the current context")
	}
	opt := envd.Options{
		Context: context,
	}
	telemetry.GetReporter().Telemetry("top", telemetry.AddField("runner", context.Runner))

	envdEngine, err := envd.New(clicontext.Context, opt)
	if err != nil {
		return err
	}
	envs, err := envdEngine.ListEnvironment(clicontext.Context)
	if err != nil {
		return err
	}

	dockerClient, err := docker.NewClient(clicontext.Context)
	if err != nil {
		return errors.Wrap(err, "failed to create the docker client")
	}

	collector, err := metrics.GetCollector("docker", dockerClient)
	if err != nil {
		return err
	}

	rows := initGrid(clicontext.Context, envs, collector)
	if err != nil {
		return err
	}
	tickerCount := 1
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return nil
			}
		case <-ticker:
			if tickerCount == 100 {
				return nil
			}
			tickerCount++
			ui.Render(rows...)
		}
	}
}

func initGrid(ctx context.Context, envs []types.EnvdEnvironment, collector metrics.Collector) []ui.Drawable {
	// There will be a header
	rowNumber := len(envs) + 1
	rows := make([]*metrics.WidgetRow, rowNumber)
	header := metrics.NewWidgetRow(0)
	header.Add(metrics.NewNameCol("Name"))
	header.Add(metrics.NewNameCol("CPU"))
	header.Add(metrics.NewNameCol("Memory"))
	rows[0] = header
	for i, env := range envs {
		row := metrics.NewWidgetRow(i + 1)
		metricsChan := collector.Watch(ctx, env.Name)
		row.Add(metrics.NewNameCol(env.Name))
		row.Add(metrics.NewCPUCol(metricsChan))
		row.Add(metrics.NewMEMCol(metricsChan))
		rows[i+1] = row
	}
	rrows := make([]ui.Drawable, len(rows))
	for i, v := range rows {
		rrows[i] = v
	}
	return rrows
}
