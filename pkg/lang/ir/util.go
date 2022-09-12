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

import (
	"encoding/json"
	"os/user"
	"regexp"
	"strconv"

	"github.com/cockroachdb/errors"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func (g Graph) getWorkingDir() string {
	return fileutil.EnvdHomeDir(g.EnvironmentName)
}

func (g Graph) getExtraSourceDir() string {
	return fileutil.EnvdHomeDir("extra_source")
}

func parseLanguage(l string) (string, *string, error) {
	var language, version string
	if l == "" {
		return "", nil, errors.New("language is required")
	}

	// Get version from the string.
	re := regexp.MustCompile(`\d[\d,]*[\.]?[\d{2}]*[\.]?[\d{2}]*`)
	if !re.MatchString(l) {
		language = l
	} else {
		loc := re.FindStringIndex(l)
		language = l[:loc[0]]
		version = l[loc[0]:]
	}

	switch language {
	case "python", "r", "julia":
		return language, &version, nil
	default:
		return "", nil, errors.Newf("language %s is not supported", language)
	}
}

func getUIDGID() (int, int, error) {
	user, err := user.Current()
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get uid/gid")
	}
	// Do not support windows yet.
	uid, err := strconv.Atoi(user.Uid)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get uid")
	}
	gid, err := strconv.Atoi(user.Gid)
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get gid")
	}
	return uid, gid, nil
}

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
