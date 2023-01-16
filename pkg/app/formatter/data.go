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

package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/docker/go-units"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/envd"
	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/types"
)

var FormatFlag = cli.StringFlag{
	Name:     "format",
	Usage:    `Format of output, could be "json" or "table"`,
	Aliases:  []string{"f"},
	Value:    "table",
	Required: false,
	Action:   formatterValidator,
}

func formatterValidator(clicontext *cli.Context, v string) error {
	switch v {
	case
		"table",
		"json":
		return nil
	}
	return errors.Errorf(`Argument format only allows "json" and "table", found "%v"`, v)
}

func FormatEndpoint(env types.EnvdEnvironment) string {
	var res strings.Builder
	if env.Status.JupyterAddr != nil {
		res.WriteString(fmt.Sprintf("jupyter: %s", *env.Status.JupyterAddr))
	}
	if env.Status.RStudioServerAddr != nil {
		res.WriteString(fmt.Sprintf("rstudio: %s", *env.Status.RStudioServerAddr))
	}
	return res.String()
}

func GetDetailedVersion(clicontext *cli.Context) (DetailedVersion, error) {
	context, err := home.GetManager().ContextGetCurrent()
	if err != nil {
		return DetailedVersion{}, errors.Wrap(err, "failed to get the current context")
	}
	opt := envd.Options{
		Context: context,
	}
	engine, err := envd.New(clicontext.Context, opt)
	if err != nil {
		return DetailedVersion{}, errors.Wrap(
			err, "failed to create engine for docker server",
		)
	}

	info, err := engine.GetInfo(clicontext.Context)
	if err != nil {
		return DetailedVersion{}, errors.Wrap(
			err, "failed to get detailed version info from docker server",
		)
	}

	return DetailedVersion{
		OSVersion:         info.OSVersion,
		OSType:            info.OSType,
		KernelVersion:     info.KernelVersion,
		DockerVersion:     info.ServerVersion,
		Architecture:      info.Architecture,
		DefaultRuntime:    info.DefaultRuntime,
		ContainerRuntimes: GetRuntimes(info),
	}, nil
}

type DetailedVersion struct {
	OSVersion         string
	OSType            string
	KernelVersion     string
	Architecture      string
	DockerVersion     string
	ContainerRuntimes string
	DefaultRuntime    string
}

func GetRuntimes(info *types.EnvdInfo) string {
	runtimesMap := info.Runtimes
	keys := make([]string, 0, len(runtimesMap))
	for k := range runtimesMap {
		keys = append(keys, k)
	}
	return "[" + strings.Join(keys, ",") + "]"
}

func StringOrNone(target string) string {
	if target == "" {
		return "<none>"
	}
	return target
}

func CreatedSinceString(created int64) string {
	createdAt := time.Unix(created, 0)

	if createdAt.IsZero() {
		return ""
	}

	return units.HumanDuration(time.Now().UTC().Sub(createdAt)) + " ago"
}
