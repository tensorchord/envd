package syncthing

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"

	"github.com/syncthing/syncthing/lib/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

func GetHomeDirectory(name string) string {
	return fmt.Sprintf("%s/syncthing-%s", fileutil.DefaultConfigDir, name)
}

func GetConfigFilePath(homeDirectory string) string {
	return fmt.Sprintf("%s/config.xml", homeDirectory)
}

func (s *Syncthing) WriteLocalConfig() error {
	configBytes, err := GetConfigBytes(s.Config, XML)
	if err != nil {
		return fmt.Errorf("failed to get syncthing config bytes: %w", err)
	}

	err = fileutil.CreateDirIfNotExist(s.HomeDirectory)
	if err != nil {
		return fmt.Errorf("failed to get syncthing config file path: %w", err)
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

func (s *Syncthing) CleanLocalConfig() error {
	if err := os.RemoveAll(s.HomeDirectory); err != nil {
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
