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

package ir

import "encoding/json"

func (rg *RuntimeGraph) Dump() (string, error) {
	b, err := json.Marshal(rg)
	if err != nil {
		return "", nil
	}
	runtimeGraphCode := string(b)
	return runtimeGraphCode, nil
}

func (rg *RuntimeGraph) Load(code []byte) error {
	var newrg *RuntimeGraph
	err := json.Unmarshal(code, newrg)
	if err != nil {
		return err
	}
	rg.RuntimeCommands = newrg.RuntimeCommands
	rg.RuntimeDaemon = newrg.RuntimeDaemon
	rg.RuntimeEnviron = newrg.RuntimeEnviron
	rg.RuntimeExpose = newrg.RuntimeExpose
	return nil
}
