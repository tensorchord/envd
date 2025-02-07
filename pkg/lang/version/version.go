// Copyright 2023 The envd Authors
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

package version

import (
	"bufio"
	"log"
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/tensorchord/envd/pkg/lang/frontend/starlark"
	starlarkv1 "github.com/tensorchord/envd/pkg/lang/frontend/starlark/v1"
	"github.com/tensorchord/envd/pkg/lang/ir"
	v1 "github.com/tensorchord/envd/pkg/lang/ir/v1"
)

// Checker gets the version from the comment.
// # syntax=v0
type Getter interface {
	GetVersion() Version
	GetDefaultGraph() ir.Graph
	GetDefaultGraphHash() string
	GetStarlarkInterpreter(buildContextDir string) starlark.Interpreter
}

type Version string

const (
	// V0 is no longer supported in envd v1.
	V0 Version = "v0"
	// V1 is the v1 version of the starlark frontend language.
	// V1 is the default version.
	V1 Version = "v1"
	// VersionUnknown is the unknown version of the starlark frontend language.
	VersionUnknown Version = "unknown"
)

type generalGetter struct {
	v Version
}

func NewByVersion(ver string) Getter {
	g := &generalGetter{}
	switch ver {
	case string(V1):
		logrus.Debug("explicit using v1")
	case string(V0):
		logrus.Fatal("v0 is no longer supported in envd v1, try to use v1")
	default:
		logrus.Debug("unknown version, using v1 by default")
	}
	g.v = V1
	return g
}

// New returns a new version getter.
func New(file string) (Getter, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fscanner := bufio.NewScanner(f)
	fscanner.Scan()
	comment := fscanner.Text()

	g := &generalGetter{}
	if strings.Contains(comment, "# syntax=v1") {
		logrus.Debug("explicit using v1")
	} else if strings.Contains(comment, "# syntax=v0") {
		logrus.Fatal("v0 is no longer supported in envd v1, try to use v1")
	} else {
		logrus.Debug("unknown version, using v1 by default")
	}
	g.v = V1
	return g, nil
}

func (g generalGetter) GetVersion() Version {
	return g.v
}

func (g generalGetter) GetDefaultGraph() ir.Graph {
	switch g.v {
	case V1:
		return v1.DefaultGraph
	case V0:
		logrus.Fatal("v0 is no longer supported in envd v1, try to use v1")
		return v1.DefaultGraph
	default:
		return nil
	}
}

func (g generalGetter) GetDefaultGraphHash() string {
	switch g.v {
	case V1:
		return v1.GetDefaultGraphHash()
	case V0:
		logrus.Fatal("v0 is no longer supported in envd v1, try to use v1")
		return v1.GetDefaultGraphHash()
	default:
		return ""
	}
}

func (g generalGetter) GetStarlarkInterpreter(buildContextDir string) starlark.Interpreter {
	switch g.v {
	case V1:
		return starlarkv1.NewInterpreter(buildContextDir)
	case V0:
		logrus.Fatal("v0 is no longer supported in envd v1, try to use v1")
		return starlarkv1.NewInterpreter(buildContextDir)
	default:
		return nil
	}
}
