package syncthing

import (
	"encoding/xml"
	"fmt"

	"github.com/syncthing/syncthing/lib/config"
	"github.com/tensorchord/envd/pkg/util/fileutil"
)

const (
	DefaultLocalPort = "83845"
)

// @source: https://docs.syncthing.net/users/config.html
func InitConfig() *config.Configuration {
	return &config.Configuration{
		Version: 37,
		GUI: config.GUIConfiguration{
			Enabled:    true,
			RawAddress: fmt.Sprintf("0.0.0.0:%s", DefaultLocalPort),
			APIKey:     "envd",
			Theme:      "default",
		},
		Options: config.OptionsConfiguration{
			GlobalAnnEnabled:     false,
			LocalAnnEnabled:      false,
			ReconnectIntervalS:   1,
			StartBrowser:         true, // TODO: disable later
			NATEnabled:           false,
			URAccepted:           1,
			URPostInsecurely:     false,
			URInitialDelayS:      1800,
			AutoUpgradeIntervalH: 0, // Disable auto upgrade
			StunKeepaliveStartS:  0, // Disable STUN keepalive
		},
	}
}

func DefaultHomeDirectory() string {
	return fmt.Sprintf("%s/syncthing", fileutil.DefaultConfigDir)
}

func GetConfigFilePath(homeDirectory string) string {
	return fmt.Sprintf("%s/config.xml", homeDirectory)
}

func GetConfigByte(cfg *config.Configuration) ([]byte, error) {
	tmp := struct {
		XMLName xml.Name `xml:"configuration"`
		*config.Configuration
	}{
		Configuration: cfg,
	}

	configByte, err := xml.MarshalIndent(tmp, "", "  ")
	if err != nil {
		return []byte{}, err
	}

	return configByte, nil
}
