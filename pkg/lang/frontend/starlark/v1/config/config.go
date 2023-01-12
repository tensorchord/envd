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

package config

import (
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"

	ir "github.com/tensorchord/envd/pkg/lang/ir/v1"
	"github.com/tensorchord/envd/pkg/util/starlarkutil"
)

var (
	logger = logrus.WithField("frontend", "starlark")
)

var Module = &starlarkstruct.Module{
	Name: "config",
	Members: starlark.StringDict{
		"apt_source": starlark.NewBuiltin(ruleUbuntuAptSource, ruleFuncUbuntuAptSource),
		"gpu":        starlark.NewBuiltin(ruleGPU, ruleFuncGPU),
		"jupyter":    starlark.NewBuiltin(ruleJupyter, ruleFuncJupyter),
		"cran_mirror": starlark.NewBuiltin(
			ruleCRANMirror, ruleFuncCRANMirror),
		"pip_index": starlark.NewBuiltin(
			rulePyPIIndex, ruleFuncPyPIIndex),
		"conda_channel": starlark.NewBuiltin(
			ruleCondaChannel, ruleFuncCondaChannel),
		"julia_pkg_server": starlark.NewBuiltin(
			ruleJuliaPackageServer, ruleFuncJuliaPackageServer),
		"rstudio_server": starlark.NewBuiltin(ruleRStudioServer, ruleFuncRStudioServer),
		"entrypoint":     starlark.NewBuiltin(ruleEntrypoint, ruleFuncEntrypoint),
		"repo":           starlark.NewBuiltin(ruleRepo, ruleFuncRepo),
		"owner":          starlark.NewBuiltin(ruleOwner, ruleFuncOwner),
	},
}

func ruleFuncGPU(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var numGPUs starlark.Int

	if err := starlark.UnpackArgs(ruleGPU, args, kwargs,
		"count?", &numGPUs); err != nil {
		return nil, err
	}

	numGPUsInt, ok := numGPUs.Int64()
	if ok {
		ir.GPU(int(numGPUsInt))
		logger.Debugf("Using %d GPUs", int(numGPUsInt))
	} else {
		logger.Debugf("Failed to convert gpu count to int64")
	}
	return starlark.None, nil
}

func ruleFuncJupyter(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var token starlark.String
	var port starlark.Int

	if err := starlark.UnpackArgs(ruleJupyter, args, kwargs,
		"token?", &token, "port?", &port); err != nil {
		return nil, err
	}

	pwdStr := token.GoString()

	portInt, ok := port.Int64()
	if !ok {
		return nil, errors.New("port must be an integer")
	}
	logger.Debugf("rule `%s` is invoked, password=%s, port=%d",
		ruleJupyter, pwdStr, portInt)
	if err := ir.Jupyter(pwdStr, portInt); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func ruleFuncPyPIIndex(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var url, extraURL starlark.String
	var trust bool

	if err := starlark.UnpackArgs(rulePyPIIndex, args, kwargs,
		"url", &url, "extra_url?", &extraURL, "trust?", &trust); err != nil {
		return nil, err
	}

	indexStr := url.GoString()
	extraIndexStr := extraURL.GoString()

	logger.Debugf("rule `%s` is invoked, index=%s, extraIndex=%s, trust=%t",
		rulePyPIIndex, indexStr, extraIndexStr, trust)
	if err := ir.PyPIIndex(indexStr, extraIndexStr, trust); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func ruleFuncCRANMirror(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var url starlark.String

	if err := starlark.UnpackArgs(ruleCRANMirror, args, kwargs,
		"url?", &url); err != nil {
		return nil, err
	}

	urlStr := url.GoString()

	logger.Debugf("rule `%s` is invoked, url=%s", ruleCRANMirror, urlStr)
	if err := ir.CRANMirror(urlStr); err != nil {
		return nil, err
	}
	return starlark.None, nil
}

func ruleFuncJuliaPackageServer(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var url starlark.String

	if err := starlark.UnpackArgs(ruleJuliaPackageServer, args, kwargs,
		"url?", &url); err != nil {
		return nil, err
	}

	urlStr := url.GoString()

	logger.Debugf("rule `%s` is invoked, url=%s", ruleJuliaPackageServer, urlStr)
	if err := ir.JuliaPackageServer(urlStr); err != nil {
		return nil, err
	}
	return starlark.None, nil
}

func ruleFuncUbuntuAptSource(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var source starlark.String

	if err := starlark.UnpackArgs(ruleUbuntuAptSource, args, kwargs,
		"source?", &source); err != nil {
		return nil, err
	}

	sourceStr := source.GoString()

	logger.Debugf("rule `%s` is invoked, source=%s", ruleUbuntuAptSource, sourceStr)
	if err := ir.UbuntuAPT(sourceStr); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func ruleFuncRStudioServer(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	if err := ir.RStudioServer(); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func ruleFuncCondaChannel(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var channel string

	if err := starlark.UnpackArgs(ruleCondaChannel, args, kwargs,
		"channel?", &channel); err != nil {
		return nil, err
	}

	logger.Debugf("rule `%s` is invoked, channel=%s\n",
		ruleCondaChannel, channel)
	if err := ir.CondaChannel(channel); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func ruleFuncEntrypoint(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var argv *starlark.List

	if err := starlark.UnpackArgs(ruleEntrypoint, args, kwargs, "args", &argv); err != nil {
		return nil, err
	}

	argList, err := starlarkutil.ToStringSlice(argv)
	if err != nil {
		return nil, err
	}

	logger.Debugf("user defined entrypoints: {%s}\n", argList)
	ir.Entrypoint(argList)
	return starlark.None, nil
}

func ruleFuncRepo(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var url, description string

	if err := starlark.UnpackArgs(ruleRepo, args, kwargs, "url", &url, "description?", &description); err != nil {
		return nil, err
	}

	logger.Debugf("repo info: url=%s, description=%s", url, description)
	ir.Repo(url, description)
	return starlark.None, nil
}

func ruleFuncOwner(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var (
		uid = -1
		gid = -1
	)

	if err := starlark.UnpackArgs(ruleOwner, args, kwargs, "uid", &uid, "gid", &gid); err != nil {
		return nil, err
	}

	if uid < 0 || uid > 65535 || gid < 0 || gid > 65535 {
		return nil, errors.New("get a wrong uid or gid")
	}
	logger.Debugf("owner info: uid=%d, gid=%d", uid, gid)
	ir.Owner(uid, gid)
	return starlark.None, nil
}
