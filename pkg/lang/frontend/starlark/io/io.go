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

package io

import (
	"errors"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	"github.com/tensorchord/envd/pkg/lang/frontend/starlark/data"
	"github.com/tensorchord/envd/pkg/lang/ir"
	// envdData "github.com/tensorchord/envd/pkg/data"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

var Module = &starlarkstruct.Module{
	Name: "io",
	Members: starlark.StringDict{
		"copy":  starlark.NewBuiltin(ruleCopy, ruleFuncCopy),
		"mount": starlark.NewBuiltin(ruleMount, ruleFuncMount),
	},
}

func ruleFuncMount(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var source starlark.Value
	var destination starlark.String

	if err := starlark.UnpackArgs(ruleMount, args, kwargs,
		"src?", &source, "dest?", &destination); err != nil {
		return nil, err
	}

	var sourceStr string
	var err error

	if v, ok := source.(*data.DataSourceValue); ok {
		err = v.Init()
		if err != nil {
			return starlark.None, err
		}
		sourceStr, err = v.GetHostDir()
		if err != nil {
			return starlark.None, err
		}
	} else if vs, ok := source.(starlark.String); ok {
		sourceStr = vs.GoString()
	} else {
		return starlark.None, errors.New("invalid source")
	}

	// sourceStr := source.GoString()
	destinationStr := destination.GoString()

	logger.Debugf("rule `%s` is invoked, src=%s, dest=%s",
		ruleMount, sourceStr, destinationStr)

	// Expand source directory based on host user
	usr, _ := user.Current()
	dir := usr.HomeDir
	if sourceStr == "~" {
		sourceStr = dir
	} else if strings.HasPrefix(sourceStr, "~/") {
		sourceStr = filepath.Join(dir, sourceStr[2:])
	}
	// Expand dest directory based on container user envd
	dir = "/home/envd/"
	if destinationStr == "~" {
		destinationStr = dir
	} else if strings.HasPrefix(destinationStr, "~/") {
		destinationStr = filepath.Join(dir, destinationStr[2:])
	}
	ir.Mount(sourceStr, destinationStr)

	return starlark.None, nil
}

func ruleFuncCopy(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var source, destination starlark.String

	if err := starlark.UnpackArgs(ruleCopy, args, kwargs,
		"src?", &source, "dest?", &destination); err != nil {
		return nil, err
	}

	sourceStr := source.GoString()
	destinationStr := destination.GoString()

	logger.Debugf("rule `%s` is invoked, src=%s, dest=%s",
		ruleCopy, sourceStr, destinationStr)
	ir.Copy(sourceStr, destinationStr)

	return starlark.None, nil
}
