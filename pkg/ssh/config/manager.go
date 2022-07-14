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

// https://gist.github.com/stefanprodan/2d20d0c6fdab6f14ce8219464e8b4b9a
// Refer to okteto/pkg/ssh/exec.go

package config

import (
	"io"
	"os"
)

type SshManager struct {
	isWsl bool
}

func (*SshManager) GenerateKeys() error {
	publicKeyPath, privateKeyPath, err := getDefaultKeyPaths()
	if err != nil {
		return err
	}
	err = generateKeys(publicKeyPath, privateKeyPath, bitSize)

	if err != nil {
		return err
	}

}

func copy(src, dst string, permission int) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	err = os.Chmod(dst, 0700)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
