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
	"fmt"
	"os/user"
	"regexp"
	"strconv"

	"github.com/cockroachdb/errors"
)

func parseLanguage(l string) (string, *string, error) {
	var language, version string
	if l == "" {
		return "", nil, errors.New("language is required")
	}

	re := regexp.MustCompile(`\d[\d,]*[\.]?[\d{2}]*[\.]?[\d{2}]*`)
	if !re.MatchString(l) {
		language = l
	} else {
		loc := re.FindStringIndex(l)
		language = l[:loc[0]]
		version = l[loc[0]:]
	}

	switch language {
	case "python":
		return "python", &version, nil
	case "r":
		return "r", &version, nil
	default:
		return "", nil, fmt.Errorf("language %s is not supported", language)
	}
}

func getUIDGID() (int, int, error) {
	user, err := user.Current()
	if err != nil {
		return 0, 0, errors.Wrap(err, "failed to get uid/gid")
	}
	// Do not support windows yet.
	if uid, err := strconv.Atoi(user.Uid); err != nil {
		return 0, 0, errors.Wrap(err, "failed to get uid")
	} else if gid, err := strconv.Atoi(user.Gid); err != nil {
		return 0, 0, errors.Wrap(err, "failed to get gid")
	} else {
		return uid, gid, nil
	}
}
