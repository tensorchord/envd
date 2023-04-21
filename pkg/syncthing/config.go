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
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/syncthing/syncthing/lib/config"
)

const (
	DefaultLocalPort           = "8386"
	DefaultRemotePort          = "8384"
	DefaultRemoteAPIAddress    = "127.0.0.1:8384"
	DefaultApiKey              = "envd"
	DefaultLocalDeviceAddress  = "tcp://127.0.0.1:22000"
	DefaultRemoteDeviceAddress = "tcp://127.0.0.1:22001"
)

// @source: https://docs.syncthing.net/users/config.html
func InitLocalConfig() *config.Configuration {
	return &config.Configuration{
		Version: 37,
		GUI: config.GUIConfiguration{
			Enabled:    true,
			RawAddress: fmt.Sprintf("127.0.0.1:%s", DefaultLocalPort),
			APIKey:     DefaultApiKey,
			Theme:      "default",
		},
		Options: config.OptionsConfiguration{
			GlobalAnnEnabled:     false,
			LocalAnnEnabled:      false,
			ReconnectIntervalS:   1,
			StartBrowser:         false,
			NATEnabled:           false,
			URAccepted:           -1,
			URPostInsecurely:     false,
			URInitialDelayS:      1800,
			AutoUpgradeIntervalH: 0, // Disable auto upgrade
			StunKeepaliveStartS:  0, // Disable STUN keepalive\
		},
	}
}

// Fetches the latest configuration from the syncthing rest api
func (s *Syncthing) GetConfig() (*config.Configuration, error) {
	resBody, err := s.Client.SendRequest(GET, "/rest/config", nil, nil)
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
	logrus.Debugf("Pulling latest config for: %s", s.Name)
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
			logrus.Debug("Timeout reached, config not applied")
			return fmt.Errorf("timed out waiting for configurations to apply")
		}

		events, err := s.GetConfigSavedEvents()
		if err != nil {
			return fmt.Errorf("failed to get syncthing config saved events: %w", err)
		}

		if len(events) > 0 {
			err := s.PullLatestConfig()
			if err != nil {
				return fmt.Errorf("failed to pull latest config: %w", err)
			}
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}
}

// Applies the config to the syncthing instance and waits for the config to be applied
func (s *Syncthing) ApplyConfig() error {
	configByte, err := GetConfigBytes(s.Config, JSON)
	if err != nil {
		return fmt.Errorf("failed to marshal syncthing config: %w", err)
	}

	_, err = s.Client.SendRequest(PUT, "/rest/config", nil, configByte)
	if err != nil {
		return fmt.Errorf("failed to apply syncthing config: %w", err)
	}

	err = s.WaitForConfigApply(10 * time.Second)
	if err != nil {
		return fmt.Errorf("failed to wait for syncthing config apply: %w", err)
	}

	return nil
}
