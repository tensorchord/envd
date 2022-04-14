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

package starlark

import (
	"github.com/sirupsen/logrus"
	"github.com/tensorchord/MIDI/pkg/lang/ir"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
)

var Module = &starlarkstruct.Module{
	Name: "midi",
	Members: starlark.StringDict{
		ruleBase: starlark.NewBuiltin(ruleBase, ruleFuncBase),
	},
}

func ruleFuncBase(thread *starlark.Thread, _ *starlark.Builtin, args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var os, language starlark.String

	if err := starlark.UnpackPositionalArgs(ruleBase, args, kwargs, 2, &os, &language); err != nil {
		return nil, err
	}

	logrus.Debugf("rule `base` is invoked, os=%s, language=%s", os.GoString(), language.GoString())
	ir.BaseStmt(os.GoString(), language.GoString())

	return starlark.None, nil
}
