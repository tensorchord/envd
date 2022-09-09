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

package data

import (
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	envddata "github.com/tensorchord/envd/pkg/data"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

var Module = &starlarkstruct.Module{
	Name: "data",
	Members: starlark.StringDict{
		"envd": starlark.NewBuiltin(ruleEnvdManagedDataSource, ruleValueEnvdManagedDataSource),
		"path": &starlarkstruct.Module{
			Name: "path",
			Members: starlark.StringDict{
				"huggingface": starlark.String(huggingFaceDatasetPath),
				"dgl":         starlark.String(dglFaceDatasetPath),
			}},
	},
}

func ruleValueEnvdManagedDataSource(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var name starlark.String

	if err := starlark.UnpackArgs(ruleEnvdManagedDataSource, args, kwargs,
		"name?", &name); err != nil {
		return nil, err
	}
	logger.Debugf("rule `%s` is invoked, name=%s",
		ruleEnvdManagedDataSource, name)

	return NewDataSourceValue(envddata.NewEnvdManagedDataSource(name.GoString())), nil
}
