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
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
)

var (
	DefaultConfigDir string
	DefaultCacheDir  string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultConfigDir = path.Join(home, ".config/envd")
	DefaultCacheDir = path.Join(home, ".config/envd")
}

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

// ConfigFile returns the location for the specified envd config file
func ConfigFile(filename string) (string, error) {
	if strings.ContainsRune(filename, os.PathSeparator) {
		return "", fmt.Errorf("filename %s should not contain any path separator", filename)
	}
	exist, err := DirExists(DefaultConfigDir)
	if err != nil {
		return "", err
	}
	if !exist {
		err = os.Mkdir(DefaultConfigDir, os.ModeDir|0700)
		if err != nil {
			return "", errors.Wrap(err, "failed to create the config dir")
		}
	}
	return path.Join(DefaultConfigDir, filename), nil
}

// CacheFile returns the location for the specified envd cache file
func CacheFile(filename string) (string, error) {
	if strings.ContainsRune(filename, os.PathSeparator) {
		return "", fmt.Errorf("filename %s should not contain any path separator", filename)
	}
	exist, err := DirExists(DefaultCacheDir)
	if err != nil {
		return "", err
	}
	if !exist {
		err = os.Mkdir(DefaultCacheDir, os.ModeDir|0700)
		if err != nil {
			return "", errors.Wrap(err, "failed to create the cache dir")
		}
	}
	return path.Join(DefaultCacheDir, filename), nil
}
