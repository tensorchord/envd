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
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	ir "github.com/tensorchord/envd/pkg/lang/ir/v1"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

var Module = &starlarkstruct.Module{
	Name: "io",
	Members: starlark.StringDict{
		"copy": starlark.NewBuiltin(ruleCopy, ruleFuncCopy),
		"http": starlark.NewBuiltin(ruleHTTP, ruleFuncHTTP),
	},
}

func ruleFuncCopy(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var source, destination starlark.String

	if err := starlark.UnpackArgs(ruleCopy, args, kwargs,
		"host_path?", &source, "envd_path?", &destination); err != nil {
		return nil, err
	}

	sourceStr := source.GoString()
	destinationStr := destination.GoString()

	logger.Debugf("rule `%s` is invoked, src=%s, dest=%s\n",
		ruleCopy, sourceStr, destinationStr)
	ir.Copy(sourceStr, destinationStr)

	return starlark.None, nil
}

func ruleFuncHTTP(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var url, checksum, filename string
	if err := starlark.UnpackArgs(ruleHTTP, args, kwargs,
		"url", &url, "checksum?", &checksum, "filename?", &filename); err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, ruleHTTP, url=%s, checksum=%s, filename=%s\n",
		ruleHTTP, url, checksum, filename)
	if err := ir.HTTP(url, checksum, filename); err != nil {
		return nil, err
	}
	return starlark.None, nil
}
