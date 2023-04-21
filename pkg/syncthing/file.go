// Copyright 2023 The envd Authors
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

package syncthing

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"

	"github.com/syncthing/syncthing/lib/config"

	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func GetHomeDirectory() (string, error) {
	return fileutil.ConfigFile("syncthing")
}

func GetConfigFilePath(homeDirectory string) string {
	return filepath.Join(homeDirectory, "config.xml")
}

func (s *Syncthing) WriteLocalConfig() error {
	configBytes, err := GetConfigBytes(s.Config, XML)
	if err != nil {
		return fmt.Errorf("failed to get syncthing config bytes: %w", err)
	}

	err = os.MkdirAll(s.HomeDirectory, 0777)
	if err != nil {
		return fmt.Errorf("failed to : %w", err)
	}

	configFilePath := GetConfigFilePath(s.HomeDirectory)
	if err = fileutil.CreateIfNotExist(configFilePath); err != nil {
		return fmt.Errorf("failed to get syncthing config file path: %w", err)
	}

	if err = os.WriteFile(configFilePath, configBytes, 0666); err != nil {
		return fmt.Errorf("failed to write syncthing config file: %w", err)
	}
	return nil
}

func CleanLocalConfig(name string) error {
	configPath, err := GetHomeDirectory()
	if err != nil {
		return fmt.Errorf("failed to get syncthing config file path: %w", err)
	}

	if err := os.RemoveAll(configPath); err != nil {
		return fmt.Errorf("failed to remove syncthing config file: %w", err)
	}
	return nil
}

const (
	XML  = "xml"
	JSON = "json"
)

// Get syncthing configuration in bytes with format XML or JSON
func GetConfigBytes(cfg *config.Configuration, outputType string) (configByte []byte, err error) {
	xmlStruct := struct {
		XMLName xml.Name `xml:"configuration"`
		*config.Configuration
	}{
		Configuration: cfg,
	}

	jsonStruct := struct {
		*config.Configuration
	}{
		Configuration: cfg,
	}

	switch outputType {
	case XML:

		configByte, err = xml.MarshalIndent(xmlStruct, "", "  ")
		if err != nil {
			return []byte{}, err
		}
	case JSON:
		configByte, err = json.MarshalIndent(jsonStruct, "", "  ")
		if err != nil {
			return []byte{}, err
		}
	default:
		return []byte{}, fmt.Errorf("invalid output type: %s", outputType)
	}

	return configByte, nil
}
