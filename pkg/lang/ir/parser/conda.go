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

package parser

import (
	"github.com/cockroachdb/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type CondaEnv struct {
	EnvName       string
	CondaPackages []string
	PipPackages   []string
}

func ParseCondaEnvYaml(condaEnvFile string) (*CondaEnv, error) {
	env := &CondaEnv{}
	config := viper.New()
	config.SetConfigFile(condaEnvFile)
	if err := config.ReadInConfig(); err != nil {
		return nil, errors.Wrapf(err, "failed to read the Conda Env File %s", condaEnvFile)
	}
	env.EnvName = config.GetString("name")
	for _, dependencies := range config.Get("dependencies").([]interface{}) {
		switch dependencies.(type) {
		case string:
			env.CondaPackages = append(env.CondaPackages, dependencies.(string))
		default:
			for _, pkgs := range dependencies.(map[string]interface{}) {
				for _, p := range pkgs.([]interface{}) {
					env.PipPackages = append(env.CondaPackages, p.(string))
				}
			}
		}
	}
	logrus.Debug("parsed python requirements: ", env)
	return env, nil
}
