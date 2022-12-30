package syncthing

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"

	"github.com/facebookgo/subset"
	"github.com/r3labs/diff/v3"
	"github.com/syncthing/syncthing/lib/config"
)

const (
	DefaultLocalPort     = "8386"
	DefaultRemotePort    = "8384"
	DefaultApiKey        = "envd"
	DefaultDeviceAddress = "tcp://127.0.0.1:22001"
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
			URAccepted:           -1,
			URPostInsecurely:     false,
			URInitialDelayS:      1800,
			AutoUpgradeIntervalH: 0, // Disable auto upgrade
			StunKeepaliveStartS:  0, // Disable STUN keepalive
		},
	}
}

// Fetches the latest configuration from the syncthing rest api
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

// Fetches the latest configuration from the syncthing rest api and applies it to the syncthing struct
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
			if s.ConfigChangesApplied(event) {
				fmt.Println("Config changes applied")
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

// Checks if configuration changes are applied by checking if the config changes are a subset of the provided config
func (s *Syncthing) ConfigChangesApplied(event *ConfigSavedEvent) bool {
	fmt.Println("Checking if config changes are applied")
	newConfig := event.Data.Copy()
	if subset.Check(&s.PrevConfig, s.Config.Copy()) || subset.Check(s.Config, &newConfig) {
		return true
	}

	// Patches the changes to the latest config, if not changed, then changes applied
	// If the config changed, then there are changes that are not applied
	_, err := diff.Merge(&s.PrevConfig, s.Config, &newConfig)
	if err != nil {
		return false
	}

	res := reflect.DeepEqual(&event.Data, &newConfig)
	fmt.Print("Performed equal, result is: ", res)
	return res
}

// Applies the config to the syncthing instance and waits for the config to be applied
func (s *Syncthing) ApplyConfig() error {
	fmt.Println("Applying config for syncthing...")
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

	fmt.Println("After waiting for config apply")

	return nil
}
