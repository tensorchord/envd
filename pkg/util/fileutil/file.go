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
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/go-git/go-git/v5"
	"github.com/sirupsen/logrus"
)

var (
	DefaultConfigDir  string
	DefaultCacheDir   string
	DefaultEnvdLibDir string
)

func init() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	DefaultConfigDir = filepath.Join(home, ".config", "envd")
	DefaultCacheDir = filepath.Join(home, ".cache", "envd")
	DefaultEnvdLibDir = filepath.Join(DefaultCacheDir, "envdlib")
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

// FindFileAbsPath returns the absolute path for the given path and file
func FindFileAbsPath(path, fileName string) (string, error) {
	if len(fileName) <= 0 {
		return "", errors.New("file name is empty")
	}
	manifest := filepath.Join(path, fileName)
	exist, err := FileExists(manifest)
	if err != nil {
		return "", err
	}
	var absPath string
	if exist {
		absPath, err = filepath.Abs(manifest)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}
	// check if ${PWD}/fileName exists
	absPath, err = filepath.Abs(fileName)
	if err != nil {
		return "", err
	}
	return absPath, nil
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
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to stat file")
	}

	logrus.WithField("filename", f).Debug("Creating file")
	_, err = os.Create(f)
	if err != nil {
		return errors.Wrap(err, "failed to create file")
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

// ConfigFile returns the location for the specified envd config file
func ConfigFile(filename string) (string, error) {
	return validateAndJoin(DefaultConfigDir, filename)
}

// CacheFile returns the location for the specified envd cache file
func CacheFile(filename string) (string, error) {
	return validateAndJoin(DefaultCacheDir, filename)
}

func validateAndJoin(dir, file string) (string, error) {
	if strings.ContainsRune(file, os.PathSeparator) {
		return "", errors.Newf("filename %s should not contain any path separator", file)
	}
	if err := os.MkdirAll(dir, os.ModeDir|0700); err != nil {
		return "", errors.Wrap(err, "failed to create the dir")
	}
	return filepath.Join(dir, file), nil
}

// DownloadOrUpdateGitRepo downloads (if not exist) or update (if exist)
func DownloadOrUpdateGitRepo(url string) (path string, err error) {
	logger := logrus.WithField("git", url)
	path = filepath.Join(DefaultEnvdLibDir, strings.ReplaceAll(url, "/", "_"))
	var repo *git.Repository
	exist, err := DirExists(path)
	if err != nil {
		return
	}
	if !exist {
		logger.Debugf("clone repo to %s", path)
		// check https://github.com/go-git/go-git/issues/305
		_, err = git.PlainClone(path, false, &git.CloneOptions{
			URL: url,
		})
		if err != nil {
			return
		}
	} else {
		logger.Debugf("repo already exists in %s", path)
		repo, err = git.PlainOpen(path)
		if err != nil {
			return
		}
		var wt *git.Worktree
		wt, err = repo.Worktree()
		if err != nil {
			return
		}
		logger.Debug("try to pull latest")
		err = wt.Pull(&git.PullOptions{})
		if err != nil && errors.Is(err, git.NoErrAlreadyUpToDate) {
			return path, nil
		}
	}

	return path, nil
}

// EnvdHomeDir returns the envd user path inside the environment
func EnvdHomeDir(path ...string) string {
	return filepath.Join(append([]string{"/", "home", "envd"}, path...)...)
}

// DefaultHomeDir returns the default user path inside the environment
func DefaultHomeDir(path ...string) string {
	return filepath.Join(append([]string{"~"}, path...)...)
}
