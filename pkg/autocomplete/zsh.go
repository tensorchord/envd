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

package autocomplete

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/containerd/containerd/log"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var autocompleteZSH = `
#compdef envd

_cli_zsh_autocomplete() {

  local -a opts
  local cur
  cur=${words[-1]}
  if [[ "$cur" == "-"* ]]; then
    opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} ${cur} --generate-bash-completion)}")
  else
    opts=("${(@f)$(_CLI_ZSH_AUTOCOMPLETE_HACK=1 ${words[@]:0:#words[@]-1} --generate-bash-completion)}")
  fi

  if [[ "${opts[1]}" != "" ]]; then
    _describe 'values' opts
  else
    _files
  fi

  return
}

compdef _cli_zsh_autocomplete envd`

var zshConfig = `
# envd zsh-completion
[ -f ~/.config/envd/envd.zsh ] && source ~/.config/envd/envd.zsh
`

// If debugging this, it might be required to run `rm ~/.zcompdump*` to remove the cache
func InsertZSHCompleteEntry() error {
	// check the system has zsh
	_, err := exec.LookPath("zsh")
	if err != nil {
		return errors.Errorf("can't find zsh in this system, stop setting the zsh-completion.")
	}

	// should be the same on linux and macOS
	filename := "envd.zsh"
	homeDir := os.Getenv("HOME")
	dirs := []string{
		"/usr/share/zsh/site-functions",
		"/usr/local/share/zsh/site-functions",
		fileutil.DefaultConfigDir,
	}

	var f *os.File
	var lastErr error
	path := ""
	for _, dir := range dirs {
		dirPathExists, err := fileutil.DirExists(dir)
		if err != nil {
			return errors.Wrapf(err, "failed to check if %s exists", dir)
		}
		if dirPathExists {
			path = fmt.Sprintf("%s/%s", dir, filename)
			log.L.Debugf("use the zsh-completion path for envd: %s", path)

			pathExists, err := fileutil.FileExists(path)
			if err != nil {
				lastErr = errors.Wrapf(err, "failed to check if %s exists", path)
			}
			if pathExists {
				return nil // file already exists, don't update it.
			}

			// create the completion file
			f, err = os.Create(path)
			if err != nil {
				lastErr = err
				continue
			}

			break
		}
	}

	if f == nil {
		return lastErr
	}
	defer f.Close()

	compEntry, err := ZshCompleteEntry()
	if err != nil {
		return errors.Wrapf(err, "Warning: unable to enable zsh-completion")
	}

	_, err = f.Write([]byte(compEntry))
	if err != nil {
		return errors.Wrapf(err, "failed writing to %s", path)
	}

	if strings.HasPrefix(path, homeDir) {
		zshFile, err := os.OpenFile(fmt.Sprintf("%s/.zshrc", homeDir), os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
		if err != nil {
			log.L.Warnf("unable to open the `~/.zshrc`, please add the following lines into `~/.zshrc` to get the envd zsh completion:\n"+
				"    %s\n", zshConfig)
			return err
		}
		defer zshFile.Close()

		_, err = fmt.Fprintf(zshFile, "%s\n", zshConfig)
		if err != nil {
			log.L.Warnf("unable to write the `~/.zshrc`, please add the following lines into `~/.zshrc` to get the envd zsh completion:\n"+
				"    %s\n", zshConfig)
			return err
		}
	}

	return deleteZcompdump()
}

func ZshCompleteEntry() (string, error) {
	return autocompleteZSH, nil
}

func deleteZcompdump() error {
	var homeDir string
	sudoUser, found := os.LookupEnv("SUDO_USER")
	if !found {
		var err error
		homeDir, err = os.UserHomeDir()
		if err != nil {
			return errors.Wrapf(err, "failed to lookup current user home dir")
		}
	} else {
		currentUser, err := user.Lookup(sudoUser)
		if err != nil {
			return errors.Wrapf(err, "failed to lookup user %s", sudoUser)
		}
		homeDir = currentUser.HomeDir
	}
	files, err := os.ReadDir(homeDir)
	if err != nil {
		return errors.Wrapf(err, "failed to read dir %s", homeDir)
	}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), ".zcompdump") {
			path := filepath.Join(homeDir, f.Name())
			err := os.Remove(path)
			if err != nil {
				return errors.Wrapf(err, "failed to remove %s", path)
			}
		}
	}
	return nil
}
