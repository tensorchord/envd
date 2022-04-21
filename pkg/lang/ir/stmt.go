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

package ir

import (
	"github.com/moby/buildkit/client/llb"
)

var Stmt llb.State

func BaseStmt(os, language string) {
	// TODO(gaocegege): Support multi os and languages.
	base := llb.Image("docker.io/library/python:3.8")
	Stmt = base
}
