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

package json

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

func printJSON(v any) error {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	// avoid escaped from <none> into "\u003cnone\u003e"
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(v)
	if err != nil {
		return errors.Wrap(err, "failed to marshal json")
	}
	fmt.Print(buffer.String())
	return nil
}
