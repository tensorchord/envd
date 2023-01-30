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

package v0

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"os/user"
	"regexp"
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

func getUIDGID() (int, int, error) {
	owner := viper.GetString(flag.FlagBuildOwner)
	if len(owner) > 0 {
		logrus.WithField("flag", owner).Info("use owner")
		ids := strings.Split(owner, ":")
		if len(ids) > 2 {
			return 0, 0, errors.Newf("wrong format for owner (uid:gid): %s", owner)
		}
		uid, err := strconv.Atoi(ids[0])
		if err != nil {
			logrus.Info(err)
			return 0, 0, errors.Wrap(err, "failed to get uid")
		}
		// if omit gid, will use the uid as gid
		if len(ids) == 1 {
			return uid, uid, nil
		}
		gid, err := strconv.Atoi(ids[1])
		if err != nil {
			return 0, 0, errors.Wrap(err, "failed to get gid")
		}
		return uid, gid, nil
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
