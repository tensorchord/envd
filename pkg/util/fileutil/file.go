/*
   Copyright The earthly Authors.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package fileutil

import (
	"os"
	"path/filepath"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
)

// FileExists returns true if the file exists
func FileExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "unable to stat %s", filename)
	}
	return !info.IsDir(), nil
}

func RemoveAll(dirname string) error {
	return os.RemoveAll(dirname)
}

// DirExists returns true if the directory exists.
func DirExists(filename string) (bool, error) {
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, errors.Wrapf(err, "unable to stat %s", filename)
	}
	return info.IsDir(), nil
}

func CreateIfNotExist(f string) error {
	_, err := os.Stat(f)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.WithField("filename", f).Debug("Creating file")
			if _, err := os.Create(f); err != nil {
				return errors.Wrap(err, "failed to create file")
			}
		} else {
			return errors.Wrap(err, "failed to stat file")
		}
	}
	return nil
}

func CWD() (string, error) {
	return os.Getwd()
}

func RootDir() (string, error) {
	cwd, err := CWD()
	if err != nil {
		return "", err
	}
	return filepath.Base(cwd), nil
}

func Base(dir string) string {
	return filepath.Base(dir)
}
