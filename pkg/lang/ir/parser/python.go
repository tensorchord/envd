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

package parser

import (
	"bufio"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
)

func ParsePythonRequirements(requirementsFile string) ([]string, error) {
	parsed := []string{}
	file, err := os.Open(requirementsFile)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read the file %s", requirementsFile)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if line[0] == '#' {
			continue
		}
		parsed = append(parsed, line)
	}

	logrus.Debug("parsed python requirements: ", parsed)

	if err := scanner.Err(); err != nil {
		return nil, errors.Wrapf(err, "failed to scan the file %s", requirementsFile)
	}
	return parsed, nil
}
