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

package v1

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os/user"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/tensorchord/envd/pkg/flag"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func (g generalGraph) getWorkingDir() string {
	return fileutil.EnvdHomeDir(g.EnvironmentName)
}

func (g generalGraph) getExtraSourceDir() string {
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

func (g generalGraph) getUIDGID() (int, int, error) {
	// Firstly check cli flag set uid and git.
	// if not set, then use the config.owner statements in envd.build
	// The current user is the last choice
	owner := viper.GetString(flag.FlagBuildOwner)
	if len(owner) > 0 {
		logrus.WithField("flag", owner).Debug("use owner")
		ids := strings.Split(owner, ":")
		if len(ids) > 2 {
			return 0, 0, errors.Newf("wrong format for owner (uid:gid): %s", owner)
		}
		uid, err := strconv.Atoi(ids[0])
		if err != nil {
			logrus.Info(err)
			return 0, 0, errors.Wrap(err, "failed to get uid")
		}
		if len(ids) == 1 {
			logrus.Debugf("gid is omitted, will use the same number as uid: %d\n", uid)
			return uid, uid, nil
		}
		gid, err := strconv.Atoi(ids[1])
		if err != nil {
			return 0, 0, errors.Wrap(err, "failed to get gid")
		}
		return uid, gid, nil
	}
	if g.uid != -1 && g.gid != -1 {
		return g.uid, g.gid, nil
	}
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
	// macOS will use staff gid=20, which is already taken in Ubuntu,
	// thus we need to hard code a new one
	// refer to https://github.com/tensorchord/envd/issues/398
	if runtime.GOOS == "darwin" && gid == 20 {
		gid = uid
	}
	return uid, gid, nil
}

// A stream of gobs is self-describing. Each data item in the stream is preceded by a specification of its type, expressed in terms of a small set of predefined types. Pointers are not transmitted, but the things they point to are transmitted; that is, the values are flattened.
// see https://pkg.go.dev/encoding/gob#hdr-Basics
// we hash the blobs to determined if the graph changed.
func GetDefaultGraphHash() string {
	var b bytes.Buffer
	err := gob.NewEncoder(&b).Encode(DefaultGraph)
	if err != nil {
		return ""
	}
	data := b.Bytes()
	hashD := md5.Sum(data)
	return hex.EncodeToString(hashD[:])
}

// GetCUDAImage finds the correct CUDA base image
// refer to https://hub.docker.com/r/nvidia/cuda/tags
func GetCUDAImage(image string, cuda *string, cudnn string, dev bool) string {
	// TODO: support CUDA 10
	target := "runtime"
	if dev {
		target = "devel"
	}
	imageTag := strings.Replace(image, ":", "", 1)

	return fmt.Sprintf("docker.io/nvidia/cuda:%s-cudnn%s-%s-%s", *cuda, cudnn, target, imageTag)
}

func (g *generalGraph) Dump() (string, error) {
	b, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	runtimeGraphCode := string(b)
	return runtimeGraphCode, nil
}

func (g *generalGraph) Load(code []byte) error {
	err := json.Unmarshal(code, g)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal")
	}
	return nil
}

func (g generalGraph) GeneralGraphFromLabel(label []byte) (ir.Graph, error) {
	newg := generalGraph{}
	err := newg.Load(label)
	if err != nil {
		return nil, err
	}
	return &newg, nil
}
