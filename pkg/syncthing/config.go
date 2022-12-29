package syncthing

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/facebookgo/subset"
	"github.com/r3labs/diff/v3"
	"github.com/sirupsen/logrus"
	"github.com/syncthing/syncthing/lib/config"
)

const (
	DefaultLocalPort  = "8386"
	DefaultRemotePort = "8384"
	DefaultApiKey     = "envd"
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

func (s *Syncthing) GetConfig() (*config.Configuration, error) {
	resBody, err := s.ApiCall(GET, "/rest/config", nil, []byte{})
	if err != nil {
		return nil, fmt.Errorf("failed to fetch syncthing config: %w", err)
	}

	cfg := &config.Configuration{}
	err = json.Unmarshal(resBody, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal syncthing config: %w", err)
	}

	return cfg, nil
}

func (s *Syncthing) PullLatestConfig() error {
	cfg, err := s.GetConfig()
	if err != nil {
		return fmt.Errorf("failed to fetch syncthing config: %w", err)
	}

	s.Config = cfg
	s.PrevConfig = cfg.Copy()
	return nil
}

func (s *Syncthing) WaitForConfigApply(timeout time.Duration) error {
	start := time.Now()
	for {
		if time.Since(start) > timeout {
			return fmt.Errorf("timed out waiting for configurations to apply")
		}

		events, err := s.GetConfigSavedEvents()
		if err != nil {
			return fmt.Errorf("failed to get syncthing config saved events: %w", err)
		}

		// Check if the applied config is the most recent config
		for _, event := range events {
			if s.ConfigChangesApplied(event.Data) {
				err := s.PullLatestConfig()
				if err != nil {
					return fmt.Errorf("failed to pull latest config: %w", err)
				}

				return nil
			}
		}

		time.Sleep(500 * time.Millisecond)
	}
}

func (s *Syncthing) ConfigChangesApplied(newConfig config.Configuration) bool {
	var configDiff config.Configuration
	_, err := diff.Merge(s.PrevConfig, s.Config.Copy(), configDiff)
	if err != nil {
		logrus.Error(fmt.Errorf("failed to check for changes in config: %w", err))
		return false
	}

	return subset.Check(configDiff, newConfig)

}

func (s *Syncthing) ApplyConfig() error {
	configByte, err := GetConfigBytes(s.Config, JSON)
	if err != nil {
		return fmt.Errorf("failed to marshal syncthing config: %w", err)
	}

	_, err = s.ApiCall(PUT, "/rest/config", nil, configByte)
	if err != nil {
		return fmt.Errorf("failed to apply syncthing config: %w", err)
	}

	err = s.WaitForConfigApply(10 * time.Second)
	if err != nil {
		return fmt.Errorf("failed to wait for syncthing config apply: %w", err)
	}

	return nil
}
