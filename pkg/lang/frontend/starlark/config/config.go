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
	"github.com/tensorchord/envd/pkg/lang/ir"
	"go.starlark.net/starlark"
	"go.starlark.net/starlarkstruct"
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
	var password starlark.String
	var port starlark.Int

	if err := starlark.UnpackArgs(ruleJupyter, args, kwargs,
		"password?", &password, "port?", &port); err != nil {
		return nil, err
	}

	pwdStr := ""
	if password != starlark.String("") {
		pwdStr = password.GoString()
	}

	portInt, ok := port.Int64()
	if !ok {
		return nil, errors.New("port must be an integer")
	}
	logger.Debugf("rule `%s` is invoked, password=%s, port=%d", ruleJupyter,
		pwdStr, portInt)
	if err := ir.Jupyter(pwdStr, portInt); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func ruleFuncPyPIIndex(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var mode, url, extraURL starlark.String

	if err := starlark.UnpackArgs(rulePyPIIndex, args, kwargs,
		"mode?", &mode, "url?", &url, "extra_url?", &extraURL); err != nil {
		return nil, err
	}

	modeStr := ""
	if mode != starlark.String("") {
		modeStr = mode.GoString()
	}
	indexStr := ""
	if url != starlark.String("") {
		indexStr = url.GoString()
	}
	extraIndexStr := ""
	if extraURL != starlark.String("") {
		extraIndexStr = extraURL.GoString()
	}

	logger.Debugf("rule `%s` is invoked, mode=%s, index=%s, extraIndex=%s", rulePyPIIndex,
		modeStr, indexStr, extraIndexStr)
	if err := ir.PyPIIndex(modeStr, indexStr, extraIndexStr); err != nil {
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

	urlStr := ""
	if url != starlark.String("") {
		urlStr = url.GoString()
	}

	logger.Debugf("rule `%s` is invoked, url=%s", ruleCRANMirror, urlStr)
	if err := ir.CRANMirror(urlStr); err != nil {
		return nil, err
	}
	return starlark.None, nil
}

func ruleFuncUbuntuAptSource(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var mode, source starlark.String

	if err := starlark.UnpackArgs(ruleUbuntuAptSource, args, kwargs,
		"mode?", &mode, "source?", &source); err != nil {
		return nil, err
	}

	modeStr := ""
	if mode != starlark.String("") {
		modeStr = mode.GoString()
	}
	sourceStr := ""
	if source != starlark.String("") {
		sourceStr = source.GoString()
	}

	logger.Debugf("rule `%s` is invoked, mode=%s, source=%s", ruleUbuntuAptSource,
		modeStr, sourceStr)
	if err := ir.UbuntuAPT(modeStr, sourceStr); err != nil {
		return nil, err
	}

	return starlark.None, nil
}

func ruleFuncCondaChannel(thread *starlark.Thread, _ *starlark.Builtin,
	args starlark.Tuple, kwargs []starlark.Tuple) (starlark.Value, error) {
	var channel starlark.String

	if err := starlark.UnpackArgs(ruleCondaChannel, args, kwargs,
		"channel?", &channel); err != nil {
		return nil, err
	}

	channelStr := ""
	if channel != starlark.String("") {
		channelStr = channel.GoString()
	}

	logger.Debugf("rule `%s` is invoked, channel=%s", ruleCondaChannel,
		channelStr)
	if err := ir.CondaChannel(channelStr); err != nil {
		return nil, err
	}

	return starlark.None, nil
}
