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

package autocomplete

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func InsertFishCompleteEntry(clicontext *cli.Context) error {
	_, err := exec.LookPath("fish")
	if err != nil {
		return errors.Errorf("can't find fish in this system, stop settings the fish-completion")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.Wrapf(err, "unable obtain user directory", err)
	}
	path := filepath.Join(homeDir, ".config/fish/completions/envd.fish")
	dirPath := filepath.Dir(path)

	dirPathExists, err := fileutil.DirExists(dirPath)
	if err != nil {
		return errors.Wrapf(err, "failed checking if %s exists", dirPath)
	}
	if !dirPathExists {
		return errors.Errorf("unable to enable fish-completion: %s does not exists", dirPath)
	}

	fishEntry, err := FishCompleteEntry(clicontext)
	if err != nil {
		return errors.Wrapf(err, "unable to enable fish-completion")
	}
	if err = os.WriteFile(path, []byte(fishEntry), 0644); err != nil {
		return errors.Wrapf(err, "failed writing to %s", path)
	}
	return nil
}

func FishCompleteEntry(clicontext *cli.Context) (string, error) {
	return clicontext.App.ToFishCompletion()
}
