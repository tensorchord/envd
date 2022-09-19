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

package starlarkutil

import (
	"github.com/cockroachdb/errors"

	"go.starlark.net/starlark"
)

func ToStringSlice(v *starlark.List) ([]string, error) {
	if v == nil {
		return []string{}, nil
	}

	s := []string{}
	for i := 0; i < v.Len(); i++ {
		str, ok := starlark.AsString(v.Index(i))
		if !ok {
			return nil, errors.Newf("Conversion failed, expect string type, but got %s as %s", v.Index(i), v.Index(i).Type())
		}
		s = append(s, str)
	}

	return s, nil
}
