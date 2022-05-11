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
	"path/filepath"
	"runtime"

	"github.com/cockroachdb/errors"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

var autocompleteBASH = `
#! /bin/bash
$PROG=envd
: ${PROG:=$(basename ${BASH_SOURCE})}

_cli_bash_autocomplete() {
  if [[ "${COMP_WORDS[0]}" != "source" ]]; then
    local cur opts base
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    if [[ "$cur" == "-"* ]]; then
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} ${cur} --generate-bash-completion )
    else
      opts=$( ${COMP_WORDS[@]:0:$COMP_CWORD} --generate-bash-completion )
    fi
    COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
    return 0
  fi
}

complete -o bashdefault -o default -o nospace -F _cli_bash_autocomplete $PROG
unset PROG
`

func InsertBashCompleteEntry() error {
	var path string
	if runtime.GOOS == "darwin" {
		path = "/usr/local/etc/bash_completion.d/envd"
	} else {
		path = "/usr/share/bash-completion/completions/envd"
	}
	dirPath := filepath.Dir(path)

	dirPathExists, err := fileutil.DirExists(dirPath)
	if err != nil {
		return errors.Wrapf(err, "failed checking if %s exists", dirPath)
	}
	if !dirPathExists {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable bash-completion: %s does not exist\n", dirPath)
		return nil // bash-completion isn't available, silently fail.
	}

	pathExists, err := fileutil.FileExists(path)
	if err != nil {
		return errors.Wrapf(err, "failed checking if %s exists", path)
	}
	if pathExists {
		return nil // file already exists, don't update it.
	}

	// create the completion file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	bashEntry, err := bashCompleteEntry()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: unable to enable bash-completion: %s\n", err)
		return nil // bash-completion isn't available, silently fail.
	}

	_, err = f.Write([]byte(bashEntry))
	if err != nil {
		return errors.Wrapf(err, "failed writing to %s", path)
	}
	return nil
}

func bashCompleteEntry() (string, error) {
	return autocompleteBASH, nil
}
